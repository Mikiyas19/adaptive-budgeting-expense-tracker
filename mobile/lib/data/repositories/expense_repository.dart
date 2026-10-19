import '../../domain/entities/expense.dart';
import '../datasources/database_helper.dart';
import '../../core/network/api_client.dart';
import '../../core/constants/api_constants.dart';

class ExpenseRepository {
  // Fetch all expenses from local SQLite database
  Future<List<Expense>> getLocalExpenses() async {
    final db = await DatabaseHelper.instance.database;
    final result = await db.query('expenses', orderBy: 'date DESC');
    return result.map((json) => Expense.fromMap(json)).toList();
  }

  // Insert expense locally first (Offline-First approach)
  Future<int> insertExpense(Expense expense) async {
    final db = await DatabaseHelper.instance.database;
    final id = await db.insert('expenses', expense.toMap());
    
    // Add to sync queue
    await db.insert('sync_queue', {
      'table_name': 'expenses',
      'record_id': id,
      'action': 'INSERT',
      'created_at': DateTime.now().toIso8601String(),
    });

    return id;
  }

  // Sync pending expenses with the Go backend
  Future<void> syncExpenses() async {
    try {
      final db = await DatabaseHelper.instance.database;
      final unsynced = await db.query(
        'expenses',
        where: 'is_synced = ?',
        whereArgs: [0],
      );

      for (var row in unsynced) {
        final expense = Expense.fromMap(row);
        
        // Push to Go backend REST API
        final response = await ApiClient.post(ApiConstants.expenses, {
          'category': expense.category,
          'amount': expense.amount,
          'description': expense.description,
          'date': expense.date,
        });

        // If successful, update local record as synced
        if (response.containsKey('expense') || response.containsKey('id')) {
          await db.update(
            'expenses',
            {'is_synced': 1},
            where: 'id = ?',
            whereArgs: [expense.id],
          );
        }
      }
    } catch (e) {
      // Sync failed due to offline state; will retry later
    }
  }
}
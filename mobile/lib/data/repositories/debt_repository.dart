import '../../domain/entities/debt.dart';
import '../datasources/database_helper.dart';
import '../../core/network/api_client.dart';
import '../../core/constants/api_constants.dart';

class DebtRepository {
  Future<List<Debt>> getLocalDebts() async {
    final db = await DatabaseHelper.instance.database;
    final result = await db.query('debts', orderBy: 'due_date ASC');
    return result.map((json) => Debt.fromMap(json)).toList();
  }

  Future<int> insertDebt(Debt debt) async {
    final db = await DatabaseHelper.instance.database;
    final id = await db.insert('debts', debt.toMap());
    
    await db.insert('sync_queue', {
      'table_name': 'debts',
      'record_id': id,
      'action': 'INSERT',
      'created_at': DateTime.now().toIso8601String(),
    });

    return id;
  }
}
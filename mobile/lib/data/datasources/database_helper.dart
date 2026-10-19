import 'package:path/path.dart';
import 'package:sqflite/sqflite.dart';

class DatabaseHelper {
  static final DatabaseHelper instance = DatabaseHelper._init();
  static Database? _database;

  DatabaseHelper._init();

  Future<Database> get database async {
    if (_database != null) return _database!;
    _database = await _initDB('adaptive_budgeting.db');
    return _database!;
  }

  Future<Database> _initDB(String filePath) async {
    final dbPath = await getDatabasesPath();
    final path = join(dbPath, filePath);

    return await openDatabase(
      path,
      version: 1,
      onCreate: _createDB,
    );
  }

  Future<void> _createDB(Database db, int version) async {
    const idType = 'INTEGER PRIMARY KEY AUTOINCREMENT';
    const textType = 'TEXT NOT NULL';
    const textNullable = 'TEXT';
    const realType = 'REAL NOT NULL';
    const intType = 'INTEGER NOT NULL DEFAULT 0';

    // Local Expenses Table
    await db.execute('''
      CREATE TABLE expenses (
        id $idType,
        server_id INTEGER,
        category $textType,
        amount $realType,
        description $textNullable,
        date $textType,
        is_synced $intType,
        created_at $textType
      )
    ''');

    // Local Debts Table
    await db.execute('''
      CREATE TABLE debts (
        id $idType,
        server_id INTEGER,
        creditor_debtor $textType,
        amount $realType,
        due_date $textNullable,
        status $textType,
        is_synced $intType,
        created_at $textType
      )
    ''');

    // Pending Sync Queue
    await db.execute('''
      CREATE TABLE sync_queue (
        id $idType,
        table_name $textType,
        record_id INTEGER NOT NULL,
        action $textType,
        created_at $textType
      )
    ''');
  }

  Future<void> close() async {
    final db = await instance.database;
    db.close();
  }
}
class Expense {
  final int? id;
  final int? serverId;
  final String category;
  final double amount;
  final String? description;
  final String date;
  final bool isSynced;

  Expense({
    this.id,
    this.serverId,
    required this.category,
    required this.amount,
    this.description,
    required this.date,
    this.isSynced = false,
  });

  Map<String, dynamic> toMap() {
    return {
      'id': id,
      'server_id': serverId,
      'category': category,
      'amount': amount,
      'description': description,
      'date': date,
      'is_synced': isSynced ? 1 : 0,
      'created_at': DateTime.now().toIso8601String(),
    };
  }

  factory Expense.fromMap(Map<String, dynamic> map) {
    return Expense(
      id: map['id'],
      serverId: map['server_id'],
      category: map['category'],
      amount: map['amount'],
      description: map['description'],
      date: map['date'],
      isSynced: map['is_synced'] == 1,
    );
  }
}
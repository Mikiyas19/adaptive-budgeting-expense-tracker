class Debt {
  final int? id;
  final int? serverId;
  final String creditorDebtor;
  final double amount;
  final String? dueDate;
  final String status;
  final bool isSynced;

  Debt({
    this.id,
    this.serverId,
    required this.creditorDebtor,
    required this.amount,
    this.dueDate,
    required this.status,
    this.isSynced = false,
  });

  Map<String, dynamic> toMap() {
    return {
      'id': id,
      'server_id': serverId,
      'creditor_debtor': creditorDebtor,
      'amount': amount,
      'due_date': dueDate,
      'status': status,
      'is_synced': isSynced ? 1 : 0,
      'created_at': DateTime.now().toIso8601String(),
    };
  }

  factory Debt.fromMap(Map<String, dynamic> map) {
    return Debt(
      id: map['id'],
      serverId: map['server_id'],
      creditorDebtor: map['creditor_debtor'],
      amount: map['amount'],
      dueDate: map['due_date'],
      status: map['status'],
      isSynced: map['is_synced'] == 1,
    );
  }
}
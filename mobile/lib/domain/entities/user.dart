class User {
  final int id;
  final String email;
  final String fullName;
  final double monthlyBudget;

  User({
    required this.id,
    required this.email,
    required this.fullName,
    required this.monthlyBudget,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['user_id'] ?? json['id'] ?? 0,
      email: json['email'] ?? '',
      fullName: json['full_name'] ?? '',
      monthlyBudget: (json['monthly_budget'] ?? 0.0).toDouble(),
    );
  }
}
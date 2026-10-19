class ApiConstants {
  // Use 10.0.2.2 if testing on Android Emulator, or localhost / local IP for Windows/Desktop/Web
  static const String baseUrl = 'http://10.0.2.2:8080/api/v1';


  // Endpoints
  static const String register = '$baseUrl/register';
  static const String login = '$baseUrl/login';
  static const String expenses = '$baseUrl/expenses';
  static const String debts = '$baseUrl/debts';
  static const String sync = '$baseUrl/sync';
}
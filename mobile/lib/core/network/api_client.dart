import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import '../constants/api_constants.dart';

class ApiClient {
  // Helper to get stored JWT token
  static Future<String?> _getToken() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString('jwt_token');
  }

  // Save token after successful login/register
  static Future<void> saveToken(String token) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('jwt_token', token);
  }

  // Clear token on logout
  static Future<void> logout() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('jwt_token');
  }

  // Generic POST Request
  static Future<Map<String, dynamic>> post(String endpoint, Map<String, dynamic> body) async {
    try {
      final response = await http.post(
        Uri.parse(endpoint),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode(body),
      );

      final data = jsonDecode(response.body);
      if (response.statusCode >= 200 && response.statusCode < 300) {
        return data;
      } else {
        throw Exception(data['error'] ?? 'Something went wrong');
      }
    } catch (e) {
      throw Exception('Network Error: $e');
    }
  }

  // Generic Authenticated GET Request
  static Future<dynamic> get(String endpoint) async {
    try {
      final token = await _getToken();
      final response = await http.get(
        Uri.parse(endpoint),
        headers: {
          'Content-Type': 'application/json',
          if (token != null) 'Authorization': 'Bearer $token',
        },
      );

      final data = jsonDecode(response.body);
      if (response.statusCode >= 200 && response.statusCode < 300) {
        return data;
      } else {
        throw Exception(data['error'] ?? 'Unauthorized or Error');
      }
    } catch (e) {
      throw Exception('Network Error: $e');
    }
  }
}
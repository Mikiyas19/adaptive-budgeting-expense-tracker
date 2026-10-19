import '../../core/constants/api_constants.dart';
import '../../core/network/api_client.dart';
import '../../domain/entities/user.dart';

class AuthRepository {
  Future<Map<String, dynamic>> register(String email, String password, String fullName) async {
    final response = await ApiClient.post(ApiConstants.register, {
      'email': email,
      'password': password,
      'full_name': fullName,
    });
    return response;
  }

  Future<User> login(String email, String password) async {
    final response = await ApiClient.post(ApiConstants.login, {
      'email': email,
      'password': password,
    });
    
    // Save token if returned in response
    if (response.containsKey('token')) {
      await ApiClient.saveToken(response['token']);
    }
    
    return User.fromJson(response['user'] ?? response);
  }
}
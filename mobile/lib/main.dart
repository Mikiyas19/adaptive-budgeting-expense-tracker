import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'presentation/bloc/auth_provider.dart';
import 'presentation/bloc/budget_provider.dart';
import 'presentation/screens/login_screen.dart';
import 'presentation/screens/dashboard_screen.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => AuthProvider()),
        ChangeNotifierProvider(create: (_) => BudgetProvider()),
      ],
      child: MaterialApp(
        title: 'Adaptive Budgeting',
        theme: ThemeData(primarySwatch: Colors.blue, useMaterial3: true),
        home: Consumer<AuthProvider>(
          builder: (context, auth, _) {
            // Check if user is logged in (has user object)
            return auth.user != null ? const DashboardScreen() : const LoginScreen();
          },
        ),
      ),
    );
  }
}
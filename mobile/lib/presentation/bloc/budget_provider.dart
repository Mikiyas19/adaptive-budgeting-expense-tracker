import 'package:flutter/material.dart';
import '../../data/repositories/expense_repository.dart';
import '../../data/repositories/debt_repository.dart';
import '../../domain/entities/expense.dart';
import '../../domain/entities/debt.dart';

class BudgetProvider with ChangeNotifier {
  final ExpenseRepository _expenseRepo = ExpenseRepository();
  final DebtRepository _debtRepo = DebtRepository();

  List<Expense> _expenses = [];
  List<Debt> _debts = [];
  bool _isLoading = false;

  List<Expense> get expenses => _expenses;
  List<Debt> get debts => _debts;
  bool get isLoading => _isLoading;

  Future<void> loadData() async {
    _isLoading = true;
    notifyListeners();

    try {
      _expenses = await _expenseRepo.getLocalExpenses();
      _debts = await _debtRepo.getLocalDebts();
      
      // Try syncing in background
      await _expenseRepo.syncExpenses();
      _expenses = await _expenseRepo.getLocalExpenses(); // Refresh after sync
    } catch (e) {
      // Handle error silently or log
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<void> addExpense(String category, double amount, String? description) async {
    final newExpense = Expense(
      category: category,
      amount: amount,
      description: description,
      date: DateTime.now().toIso8601String(),
    );
    await _expenseRepo.insertExpense(newExpense);
    await loadData();
  }

  Future<void> addDebt(String creditorDebtor, double amount, String? dueDate, String status) async {
    final newDebt = Debt(
      creditorDebtor: creditorDebtor,
      amount: amount,
      dueDate: dueDate,
      status: status,
    );
    await _debtRepo.insertDebt(newDebt);
    await loadData();
  }
}
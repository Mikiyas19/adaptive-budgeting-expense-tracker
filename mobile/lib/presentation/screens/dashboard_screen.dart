import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../bloc/budget_provider.dart';
import '../bloc/auth_provider.dart';

class DashboardScreen extends StatefulWidget {
  const DashboardScreen({super.key});

  @override
  State<DashboardScreen> createState() => _DashboardScreenState();
}

class _DashboardScreenState extends State<DashboardScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      Provider.of<BudgetProvider>(context, listen: false).loadData();
    });
  }

  void _showAddExpenseDialog(BuildContext context) {
    final categoryController = TextEditingController();
    final amountController = TextEditingController();
    final descController = TextEditingController();

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Add Expense'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(controller: categoryController, decoration: const InputDecoration(labelText: 'Category')),
            TextField(controller: amountController, decoration: const InputDecoration(labelText: 'Amount'), keyboardType: TextInputType.number),
            TextField(controller: descController, decoration: const InputDecoration(labelText: 'Description (Optional)')),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(context), child: const Text('Cancel')),
          ElevatedButton(
            onPressed: () {
              final cat = categoryController.text.trim();
              final amt = double.tryParse(amountController.text.trim()) ?? 0.0;
              if (cat.isNotEmpty && amt > 0) {
                Provider.of<BudgetProvider>(context, listen: false)
                    .addExpense(cat, amt, descController.text.trim());
                Navigator.pop(context);
              }
            },
            child: const Text('Save'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(context) {
    final budgetProvider = Provider.of<BudgetProvider>(context);
    final authProvider = Provider.of<AuthProvider>(context, listen: false);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Adaptive Budgeting Dashboard'),
        actions: [
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: () => authProvider.logout(),
          )
        ],
      ),
      body: budgetProvider.isLoading
          ? const Center(child: CircularProgressIndicator())
          : ListView(
              padding: const EdgeInsets.all(16.0),
              children: [
                Text('Welcome, ${authProvider.user?.fullName ?? "User"}', style: const TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
                const SizedBox(height: 20),
                const Text('Recent Expenses', style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600)),
                const SizedBox(height: 10),
                ...budgetProvider.expenses.map((expense) => Card(
                      child: ListTile(
                        title: Text(expense.category),
                        subtitle: Text(expense.description ?? 'No description'),
                        trailing: Text('\$${expense.amount.toStringAsFixed(2)}', style: const TextStyle(fontWeight: FontWeight.bold)),
                      ),
                    )),
              ],
            ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showAddExpenseDialog(context),
        child: const Icon(Icons.add),
      ),
    );
  }
}
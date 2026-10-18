-- 1. Users Table (including password reset columns)
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100),
    monthly_budget DECIMAL(10, 2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_sync_timestamp TIMESTAMP,
    reset_token VARCHAR(255),
    reset_expires TIMESTAMP
);

-- 2. Categories Table (with is_fixed flag for fixed vs discretionary rules)
CREATE TABLE IF NOT EXISTS categories (
    category_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE, -- NULL for predefined categories
    name VARCHAR(30) NOT NULL,
    icon VARCHAR(50),
    is_predefined BOOLEAN DEFAULT FALSE,
    is_fixed BOOLEAN DEFAULT FALSE
);

-- 3. Expenses Table (with sync tracking for offline-first operation)
CREATE TABLE IF NOT EXISTS expenses (
    expense_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    category_id INT NOT NULL REFERENCES categories(category_id),
    amount DECIMAL(10, 2) NOT NULL CHECK (amount > 0),
    expense_date DATE NOT NULL DEFAULT CURRENT_DATE,
    description VARCHAR(200),
    is_synced BOOLEAN DEFAULT FALSE,
    sync_status VARCHAR(20) DEFAULT 'pending'
);

-- 4. Debts Table (for informal peer-to-peer lending)
CREATE TABLE IF NOT EXISTS debts (
    debt_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    counterparty_name VARCHAR(100) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL CHECK (amount > 0),
    is_lent BOOLEAN NOT NULL, -- TRUE: others owe me, FALSE: I owe others
    remaining_amount DECIMAL(10, 2) NOT NULL CHECK (remaining_amount >= 0),
    due_date DATE,
    status VARCHAR(20) DEFAULT 'pending',
    notes TEXT,
    is_synced BOOLEAN DEFAULT FALSE
);

-- 5. Reminders Table
CREATE TABLE IF NOT EXISTS reminders (
    reminder_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    expense_id INT REFERENCES expenses(expense_id) ON DELETE SET NULL,
    debt_id INT REFERENCES debts(debt_id) ON DELETE SET NULL,
    title VARCHAR(100) NOT NULL,
    reminder_datetime TIMESTAMP NOT NULL,
    recurrence_pattern VARCHAR(20) DEFAULT 'none',
    is_active BOOLEAN DEFAULT TRUE,
    last_triggered TIMESTAMP
);

-- 6. SyncLog Table (tracks synchronization audit history)
CREATE TABLE IF NOT EXISTS sync_logs (
    sync_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    last_sync_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    records_uploaded INT DEFAULT 0,
    records_downloaded INT DEFAULT 0,
    status VARCHAR(20) NOT NULL
);

-- Performance Indexes (to support rapid statistical aggregations and queries)
CREATE INDEX IF NOT EXISTS idx_expenses_user_date ON expenses(user_id, expense_date);
CREATE INDEX IF NOT EXISTS idx_expenses_category ON expenses(category_id);
CREATE INDEX IF NOT EXISTS idx_debts_user_status ON debts(user_id, status, due_date);
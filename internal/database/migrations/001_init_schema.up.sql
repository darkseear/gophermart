-- Создание таблицы пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    current_balance NUMERIC DEFAULT 0,
    withdrawn_balance NUMERIC DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы заказов
CREATE TABLE orders (
    number VARCHAR(255) PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('NEW', 'REGISTERED', 'PROCESSING', 'INVALID', 'PROCESSED')),
    accrual NUMERIC NOT NULL DEFAULT 0,
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы списаний
CREATE TABLE withdrawals (
    order_number VARCHAR(255) PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) NOT NULL,
    sum NUMERIC NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

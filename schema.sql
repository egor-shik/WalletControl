CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY, 
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY, 
    amount NUMERIC(10, 2) NOT NULL, 
    category_id INT REFERENCES categories(id), 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
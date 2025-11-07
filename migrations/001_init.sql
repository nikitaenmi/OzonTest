CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    provider_name VARCHAR(100) NOT NULL,
    amount DECIMAL(15,2) NOT NULL CHECK (amount >= 0),
    currency VARCHAR(3) NOT NULL,
    amount_rub DECIMAL(15,2) NOT NULL CHECK (amount_rub >= 0),
    payment_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
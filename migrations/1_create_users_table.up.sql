CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    passport_serie TEXT,
    passport_number TEXT NOT NULL,
    surname TEXT,
    name TEXT,
    patronymic TEXT,
    address TEXT,
);

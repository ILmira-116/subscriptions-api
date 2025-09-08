-- 001_create_subscriptions.sql
-- Создание расширения для генерации UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Создание таблицы подписок
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    service TEXT NOT NULL,
    price INT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индекс для ускорения поиска по пользователю
CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id); 
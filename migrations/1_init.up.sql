-- Таблица для хранения информации о пользователях
CREATE TABLE IF NOT EXISTS users (
    id              INTEGER PRIMARY KEY,
    first_name      TEXT NOT NULL,
    last_name       TEXT NOT NULL,
    patronymic      TEXT NOT NULL,
    birth_date      DATE NOT NULL,
    email           TEXT NOT NULL UNIQUE
);

-- Таблица для хранения информации об организациях
CREATE TABLE IF NOT EXISTS organizations (
    id          INTEGER PRIMARY KEY,
    name        TEXT NOT NULL,
    city        TEXT NOT NULL,
    office      TEXT NOT NULL,
    department  TEXT NOT NULL,
    UNIQUE (name, city, office, department)
);

-- Таблица для хранения связей пользователей с организациями
CREATE TABLE IF NOT EXISTS user_organizations (
    user_id         INTEGER NOT NULL,
    organization_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (organization_id) REFERENCES organizations(id),
    PRIMARY KEY (user_id, organization_id)
);

-- Таблица для хранения информации о мессенджерах и связанных с ними ID
CREATE TABLE IF NOT EXISTS user_messengers (
    user_id         INTEGER NOT NULL,
    messenger_type  TEXT NOT NULL, -- Например, 'telegram', 'whatsapp' и т.д.
    messenger_id    TEXT NOT NULL,
    chat_id         TEXT NOT NULL,
    is_activated    BOOLEAN NOT NULL DEFAULT 0,
    token           TEXT, 
    FOREIGN KEY (user_id) REFERENCES users(id),
    PRIMARY KEY (user_id, messenger_type)
);
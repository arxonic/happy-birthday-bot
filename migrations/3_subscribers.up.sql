-- Таблица для хранения информации о мессенджерах и связанных с ними  ID
CREATE TABLE IF NOT EXISTS subscribes (
    user_id         INTEGER NOT NULL, -- Именинник
    sub_id          INTEGER NOT NULL, -- Подписчик
    link            TEXT,
    expire          DATE,
    
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (sub_id) REFERENCES users(id),
    PRIMARY KEY (user_id, sub_id),
    UNIQUE (user_id, sub_id)
);
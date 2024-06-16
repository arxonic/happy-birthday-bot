-- Вставка пользователей
INSERT INTO users (first_name, last_name, patronymic, birth_date, email)
VALUES ('Ivan', 'Ivanov', 'Ivanovich', '1985-04-23', 'ivan.ivanov@example.com');

INSERT INTO users (first_name, last_name, patronymic, birth_date, email)
VALUES ('Petr', 'Petrov', 'Petrovich', '1990-06-15', 'petr.petrov@example.com');

INSERT INTO users (first_name, last_name, patronymic, birth_date, email)
VALUES ('Sidor', 'Sidorov', 'Sidorovich', '1988-11-30', 'sidor.sidorov@example.com');

-- Вставка организаций
INSERT INTO organizations (name, city, office, department)
VALUES ('Gazprom Media', 'Moscow', 'Main Office', 'IT Department');

INSERT INTO organizations (name, city, office, department)
VALUES ('Gazprom Media', 'Saint Petersburg', 'Main Office', 'Marketing Department');

-- Связь пользователей с организациями
INSERT INTO user_organizations (user_id, organization_id)
VALUES ((SELECT id FROM users WHERE email='ivan.ivanov@example.com'), (SELECT id FROM organizations WHERE name='Gazprom Media' AND city='Moscow'));

INSERT INTO user_organizations (user_id, organization_id)
VALUES ((SELECT id FROM users WHERE email='petr.petrov@example.com'), (SELECT id FROM organizations WHERE name='Gazprom Media' AND city='Saint Petersburg'));

INSERT INTO user_organizations (user_id, organization_id)
VALUES ((SELECT id FROM users WHERE email='sidor.sidorov@example.com'), (SELECT id FROM organizations WHERE name='Gazprom Media' AND city='Moscow'));

-- Вставка мессенджеров для пользователей
INSERT INTO user_messengers (user_id, messenger_type, messenger_id, chat_id, is_activated, token)
VALUES ((SELECT id FROM users WHERE email='ivan.ivanov@example.com'), 'telegram', '123456789', '987654321', 0, 'token123');

INSERT INTO user_messengers (user_id, messenger_type, messenger_id, chat_id, is_activated, token)
VALUES ((SELECT id FROM users WHERE email='petr.petrov@example.com'), 'telegram', '987654321', '123456789', 0, 'token456');

INSERT INTO user_messengers (user_id, messenger_type, messenger_id, chat_id, is_activated, token)
VALUES ((SELECT id FROM users WHERE email='sidor.sidorov@example.com'), 'telegram', '555555555', '111111111', 0, 'token789');

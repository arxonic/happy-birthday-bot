# Happy birtday Bot

**Запуск:** go run cmd/gmh/main.go --config="path" --tg-bot-key="key" --mail-passw="smtp_password"

Представим, что Вы работаете в отделе продаж и **каждый раз** в дни рождения коллег вам приходится как-то организовываться, спрашивать кто готов скинуться на подарок, а кто нет.

Данный телеграм бот интегрируется с системой работодателя (для получения сотрудников) и за неделю до ДР какого либо сотрудника **создает тг группу**, приглашая в нее пользователей, которые захотели получать уведомления (подписались) на своего коллегу.

В основе общения пользователей с Тг Ботом лежит **Детерминированный Конечный Автомат**, который перегоняет пользователей по состояниям.


## Приложение состоит из следующих сервисов:

**1) Сервис аутентификации** - проверяет аутентификационные данные пользователей и регистрирует их

**2) Email сервис** - позволяет отправлять сообщения пользователям (в т.ч при регистрации для отправки auth ссылки)

**3) Сервис работы с API организации,** которая предоставляет данные о сотрудниках. Ключем для получения данных является **рабочий Email пользователя** в организации

**4) Сервис подписок**

**5) Демон,** который за неделю до ДР создает тг группу, собрав всех подписавшихся


# Регистрация

- Пользователь впервые пишет боту
- Бот предлагает ввести рабочий email -> пользователь отправляет его (Email проходит валидацию)
- Employer сервис запращивает данные у организации, связанные с его рабочей почтой
- В случае успешного поиска, на почту пользователя отправляется письмо с аутентификационной ссылкой
- Сохраняем данные о сотруднике, а также его tgID, chatID
- Http Auth Endpoint обрабатывает переход пользователя по ссылке и активирует его аккаунт в системе

![image](https://github.com/arxonic/happy-birthday-bot/assets/115946622/3f4bab4a-867d-4398-aafa-07a7e9877772)


# Auth Middleware

По аналогии с REST API я хотел реализовать аутентифицирующий миддлвар. Он реализован, как начальное состояние ДКА, которое проверяет данные от аккаунта мессенджера пользователя (tgID) и статус регистрации (isActivated) в нашей системе


# Состояния пользователей (ДКА)

Состояния пользователей хранятся в памяти сервиса (что не является хорошей практикой). Я бы использовал для этого **Redis**
В случае падения приложения, состояния пользователей обнуляются, поэтому первое что встречает пользователей - Auth Middleware.


# Конфиг

Все секреты хранятся в переменных окружения и/или во флагах запуска. Остальные файлы конфигурации хранятся в конфиге.


# Хранилище

БД - Sqlite, в 3 н.ф. 

**Таблицы:**

1) **Users** - ФИО, дата рождения, email

2) **Organizations** - Название, город, офис, отдел 

3) **User_organizations** - связь User-Organization

4) **User_messengers** - uID, messenger_type, messenger_id, chat_id. Также в этой таблице есть 2 поля: is_activated и token, которые лучше вынести в отдельную таблицу.

5) **Subscribes** - подписка uID - subscriberID


# Не реализовал до конца


## Демон Cron

1) В начале каждого дня, демон собирает пользователей, у которых ДР < через неделю
2) Собирает подписчиков каждого пользователя
3) Создает группы. Заносит в БД ссылку на группу и дату потери ее актуальности (день ДР пользователя, чтобы через год создать новую беседу)
4) Пришлашает пользователей в беседы

## Система подписки

Зарегистрированный пользователь имеет в системе информацию о месте работы (название организации, город, офис, отдел). Основываяся на названии организации, в которой он работает, пользователь сможет найти своих коллег и подписаться на них.

Если при совершении подписки на коллегу, у коллеги ДР < чем через неделю, пользователю отправляется ссылка на уже существующую тг беседу. 

Если такой беседы нет, то она создается.

## Тесты

INSERT INTO clients (login, "password", first_name, last_name, middle_name, passport, birthday, status)
VALUES ('login1', 'password', 'Иван', 'Иванов', 'Иванович', '0001 00000001', '2000-02-20', 'ACTIVE'),
       ('login2', 'password', 'Петр', 'Петров', 'Петрович', '0002 00000002', '1999-10-01', 'ACTIVE'),
       ('login3', 'password', 'Василий', 'Васинов', 'Васильевич', '0003 00000003', '1960-06-06', 'INACTIVE');

INSERT INTO cards ("number", balance, issuer, holder, owner_id, status)
VALUES ('4929100232729184', 1000000, 'Visa', 'user1', 1, 'ACTIVE'),
       ('4539670260962412', 1000000, 'Visa', 'user2', 2, 'ACTIVE'),
       ('2720990990763192', 0, 'MasterCard', 'user2', 2, 'INACTIVE');

INSERT INTO icons (url)
VALUES ('http://i1.com'),
       ('http://i2.com'),
       ('http://i3.com'),
       ('http://i4.com');

INSERT INTO mcc (mcc, description)
VALUES ('6540', 'Пополнения'),
       ('5411', 'Супермаркеты'),
       ('4814', 'Мобильная связь'),
       ('4829', 'Переводы');

INSERT INTO transactions (card_id, amount, status, mcc_id, description, icon_id)
VALUES (1, 5000000, 'Исполнена', 1, 'Пополнение через Альфа-Банк', 1),
       (1, -100000, 'Исполнена', 2, 'Продукты', 2),
       (1, -100000, 'Исполнена', 3, 'Пополнение телефона', 3),
       (1, -100000, 'Исполнена', 4, 'Перевод', 4);
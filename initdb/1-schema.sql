CREATE TABLE clients
(
    id          BIGSERIAL PRIMARY KEY,
    login       TEXT      NOT NULL,
    password    TEXT      NOT NULL,
    first_name  TEXT      NOT NULL,
    last_name   TEXT      NOT NULL,
    middle_name TEXT      NOT NULL,
    passport    TEXT      NOT NULL,
    birthday    DATE      NOT NULL,
    status      TEXT      NOT NULL DEFAULT 'INACTIVE' CHECK (status IN ('INACTIVE', 'ACTIVE')),
    created     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cards
(
    id       BIGSERIAL PRIMARY KEY,
    number   BIGINT    NOT NULL,
    balance  BIGINT    NOT NULL DEFAULT 0,
    issuer   TEXT      NOT NULL CHECK (issuer IN ('Visa', 'MasterCard', 'MIR')),
    holder   TEXT      NOT NULL,
    owner_id BIGINT    NOT NULL REFERENCES clients,
    status   TEXT      NOT NULL DEFAULT 'INACTIVE' CHECK (status IN ('INACTIVE', 'ACTIVE')),
    created  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE mcc
(
    id          BIGSERIAL PRIMARY KEY,
    mcc         TEXT NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE icons
(
    id  BIGSERIAL PRIMARY KEY,
    url TEXT NOT NULL
);

CREATE TABLE transactions
(
    id          BIGSERIAL PRIMARY KEY,
    amount      INT       NOT NULL,
    card_id     BIGINT    NOT NULL REFERENCES cards,
    mcc_id      BIGINT    NOT NULL REFERENCES mcc,
    icon_id     BIGINT    NOT NULL REFERENCES icons,
    status      TEXT      NOT NULL DEFAULT 'Обрабатывается' CHECK (status IN ('Обрабатывается', 'Исполнена', 'Отклонена')),
    created     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    description TEXT      NOT NULL
);

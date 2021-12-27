CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "user"
(
    id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    password TEXT         NOT NULL
);

CREATE TABLE IF NOT EXISTS restore_data
(
    user_id          UUID         NOT NULL UNIQUE,
    CONSTRAINT fk__restore_data__user__one_to_one
        FOREIGN KEY (user_id)
            REFERENCES "user" (id)
            ON UPDATE CASCADE
            ON DELETE CASCADE,
    email            VARCHAR(255) NOT NULL UNIQUE,
    telephone_number VARCHAR(255) DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS authority
(
    id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS user_authority
(
    user_id      UUID NOT NULL UNIQUE,
    CONSTRAINT fk__user_authority__user__many_to_many
        FOREIGN KEY (user_id)
            REFERENCES "user" (id)
            ON UPDATE CASCADE
            ON DELETE CASCADE,
    authority_id UUID NOT NULL UNIQUE,
    CONSTRAINT fk__user_authority__authority__many_to_many
        FOREIGN KEY (authority_id)
            REFERENCES authority (id)
            ON UPDATE CASCADE
            ON DELETE CASCADE
);
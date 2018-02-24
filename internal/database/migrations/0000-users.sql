
-- +migrate Up
CREATE TABLE users (
    uuid  CHAR(36)    PRIMARY KEY,
    email VARCHAR(60) NOT NULL,
    hash  CHAR(60)    NOT NULL,
    CONSTRAINT uq_uuid  UNIQUE (uuid),
    CONSTRAINT uq_email UNIQUE (email)
);

-- +migrate Down
DROP TABLE users;
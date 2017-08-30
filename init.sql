DROP TABLE users;
DROP TABLE domains;

CREATE TABLE domains (
  id SERIAL PRIMARY KEY,
  name VARCHAR(32) NOT NULL,
  data JSON NOT NULL,
  CONSTRAINT uq_name UNIQUE (name)
);

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  domain_id INTEGER NOT NULL,
  email VARCHAR(32) NOT NULL,
  hash VARCHAR(60) NOT NULL,
  CONSTRAINT uq_domain_email UNIQUE (domain_id, email),
  CONSTRAINT fk_domain_id FOREIGN KEY (domain_id)
    REFERENCES domains (id)
    ON DELETE CASCADE
);

--

INSERT INTO domains (name, data) VALUES ('name1', '{}');
INSERT INTO domains (name, data) VALUES ('name2', '{}');

INSERT INTO users (domain_id, email, hash) VALUES (1, 'email1', '123456789abcdef');
INSERT INTO users (domain_id, email, hash) VALUES (1, 'email2', '123456789abcdef');
INSERT INTO users (domain_id, email, hash) VALUES (2, 'email1', '123456789abcdef');

CREATE TABLE IF NOT EXISTS trainers (
    id bigserial PRIMARY KEY,
    email citext UNIQUE NOT NULL,
    name varchar (50) NOT NULL,
    experience smallint NOT NULL,
    activities text[] NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    name VARCHAR (50) NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL
);

CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL
);

CREATE TABLE IF NOT EXISTS permissions (
    id bigserial PRIMARY KEY,
    code text NOT NULL
);

CREATE TABLE IF NOT EXISTS users_permissions (
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
    PRIMARY KEY (user_id, permission_id)
);

INSERT INTO permissions (code)
VALUES
    ('trainers:read'),
    ('trainers:write');

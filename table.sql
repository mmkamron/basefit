CREATE TABLE session (
                         value text,
                         expiry text,
                         username text PRIMARY KEY
);

CREATE TABLE users (
                       username text PRIMARY KEY,
                       password text
);

CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    name text NOT NULL,
    author text NOT NULL,
    username text REFERENCES session(username)
);

CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    name text NOT NULL,
    author text NOT NULL,

);

INSERT INTO books(name, author) VALUES
('Antifragile', 'Nassim Taleb'),
('Manhood in the making', 'David Gilmore'),
('The blank slate', 'Steven Pinker');

CREATE TABLE users (
    username text PRIMARY KEY,
    password text
);

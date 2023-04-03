CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE activities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL
);

CREATE TABLE exercise_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    activity_id INTEGER REFERENCES activities(id),
    date_time TIMESTAMP NOT NULL,
    distance FLOAT,
    time INTERVAL,
    weight_lifted FLOAT
);

CREATE TABLE nutrition_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    date_time TIMESTAMP NOT NULL,
    calories FLOAT,
    protein FLOAT,
);
CREATE TABLE IF NOT EXISTS trainer (
    id serial PRIMARY KEY,
    name VARCHAR (500) NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL
);

CREATE TABLE IF NOT EXISTS client (
    id SERIAL PRIMARY KEY,
    email citext UNIQUE NOT NULL,
    username VARCHAR (50) UNIQUE NOT NULL,
    first_name VARCHAR (50) NOT NULL,
    last_name VARCHAR (50) NOT NULL
);

CREATE TABLE IF NOT EXISTS schedule (
    id SERIAL PRIMARY KEY,
    trainer_id INT NOT NULL,
    date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    is_available BOOLEAN NOT NULL,
    FOREIGN KEY (trainer_id) REFERENCES trainer(id)
);

CREATE TABLE IF NOT EXISTS booking (
    id SERIAL PRIMARY KEY,
    trainer_id INT NOT NULL,
    client_id INT NOT NULL,
    schedule_id INT NOT NULL,
    booking_status VARCHAR(20) NOT NULL,
    FOREIGN KEY (trainer_id) REFERENCES trainer(id),
    FOREIGN KEY (client_id) REFERENCES client(id),
    FOREIGN KEY (client_id) REFERENCES schedule(id)
);

CREATE TABLE IF NOT EXISTS token (
    hash bytea PRIMARY KEY,
    trainer_id bigint NOT NULL REFERENCES trainer ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL
);

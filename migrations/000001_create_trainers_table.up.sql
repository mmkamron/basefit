CREATE TABLE IF NOT EXISTS trainers(
   trainer_id serial PRIMARY KEY,
   name VARCHAR (50) NOT NULL,
   specialization VARCHAR(100) NOT NULL,
   description TEXT,
   availability VARCHAR(255) NOT NULL,
   contact VARCHAR(100) UNIQUE
);

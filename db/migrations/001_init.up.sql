BEGIN;

CREATE TABLE genres
    (
id serial PRIMARY KEY,
name text
);

CREATE TABLE books
    (
id serial PRIMARY KEY,
name character(100) UNIQUE,
price numeric(5,2) CHECK (price>0),
genre integer references genres(id) on DELETE CASCADE,
amount integer CHECK (amount>=0)
);

COMMIT;
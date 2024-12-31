CREATE TABLE IF NOT EXISTS migrations (
    id integer CONSTRAINT migrations_constraint_pk PRIMARY KEY,
    name varchar(127) NOT NULL,
    applied_on timestamp DEFAULT current_timestamp
);

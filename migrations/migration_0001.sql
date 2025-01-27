CREATE TABLE IF NOT EXISTS users (
    -- indexes will be added in future as per the requirement
    id integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    name varchar(127) NOT NULL,
    email varchar(127) NOT NULL,
    password_hash varchar(255) NOT NULL,
    created_at timestamp DEFAULT current_timestamp,
    UNIQUE(email)
);


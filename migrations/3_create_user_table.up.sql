CREATE TABLE users(
    id serial primary key,
    email varchar(255) not null unique,
    name varchar(20),
    lastname varchar(50),
    created_at timestamptz not null,
    updated_at timestamptz
);


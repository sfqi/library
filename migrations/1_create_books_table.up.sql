CREATE TABLE books (
    id      serial primary key,
    title   varchar(255)    not null,
    author  varchar(255),
    isbn    varchar(10)     not null    unique,
    isbn_13 varchar(13)     not null,
    open_library_id         varchar(50),
    cover_id varchar(50),
    year    varchar(50),
    publisher varchar(255)
);
CREATE TABLE books (
    id integer not null AUTO_INCREMENT primary key,
    title varchar(255) not null,
    author varchar(255),
    isbn varchar(10) unique not null,
    isbn_13 varchar(13) not null,
    open_library_id varchar(50),
    cover_id varchar(50),
    year integer ,
    publisher varchar(255)
)
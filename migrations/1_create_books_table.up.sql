CREATE TABLE books (
    title varchar(50),
    author_id varchar(50),
    isbn varchar(50) unique,
    isbn_13 varchar(50),
    open_library_id varchar(50),
    cover_id varchar(50),
    year integer ,
    publsher varchar(50)
)
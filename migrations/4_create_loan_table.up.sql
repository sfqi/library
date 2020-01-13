CREATE TABLE loans(
    id serial primary key,
    transaction_id varchar(255) unique,
    user_id int not null,
    book_id int unique not null,
    type int,
    created_at timestamptz not null,
    FOREIGN KEY(book_id) references books(id)
);
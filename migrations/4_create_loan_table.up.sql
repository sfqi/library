CREATE TABLE loans(
    id serial primary key,
    transaction_id varchar(255),
    user_id int not null,
    book_id int not null,
    type int,
    created_at timestamptz not null,
    FOREIGN KEY(book_id) references books(id)
);
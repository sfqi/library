ALTER TABLE books ADD created_at timestamptz (timestamp) not null;
ALTER TABLE books ADD updated_at timestamptz (timestamp) not null;
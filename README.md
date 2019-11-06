# MigrateBooks
Use : $ migrate -source file://path/to/migrations -database postgres://user:password@localhost:5432/database up 2 to create table

For setting env var in .bashrc to be permanent , open .bashrc file using:
$sudo nano ~/.bashrc
at the and of file, enter this line at the end of file:
export LIBRARY='https://openlibrary.org/api/'
Logout, and login again 

We need to remove '/' at the end from https://openlibrary.org/api/ , and add it in bookPath
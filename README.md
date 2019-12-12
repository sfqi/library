# Migrate
Use : $ migrate -source file://path/to/migrations -database postgres://user:password@localhost:5432/database up run migrations

# Githooks
on unix systems, run:
```
cp .githooks/* .git/hooks/
chmod +x .git/hooks/*
```

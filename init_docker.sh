docker run --name=comics-db -e POSTGRES_PASSWORD='qwerty' -p 5436:5432 -d postgres;
migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up;
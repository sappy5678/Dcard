echo "Creating database dcard"
psql "postgres://postgres:password@sql:5432/postgres?sslmode=disable" -c "create database dcard;"

echo "Running migrations"
./migrate -source file://deploy/db/migrations -database "postgres://postgres:password@sql:5432/dcard?sslmode=disable" up

echo "Migrations completed"
# ZHurd
ZHurd - service for managing templates and task queue for ZPL-compatible printers.

## building

```sh
go build -o zhurd ./cmd/api/main.go
```

## run

```sh
./zhurd -config /usr/local/etc/zhurd/config.json
```
example of config look in `./share/config.json.example`

or via docker compose:
```sh
docker compose up
```

## migrations

You can build docker image in dbmigrations and run it to create a new migration file:
```sh
docker run --rm -v $PWD/dbmigrations:/dbmigrations zhurd-migrator create --dir /dbmigrations [migration_name] sql
```
migrations are applyed automatically when service run with docker compose
to apply them manually run goose, e.g.:
```sh
docker run --rm -e GOOSE_DRIVER=postgres -e GOOSE_DBSTRING='postgres://zhurd:passwordsecretdb@dbhost:5432/zhurd?sslmode=disable' -v $PWD/dbmigrations:/dbmigrations zhurd-migrator --dir /dbmigrations up
``` 

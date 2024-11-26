# ZHurd
ZHurd - service for managing templates and task queue for ZPL-compatible printers.

## building

```
go build -o zhurd ./cmd/api/main.go
```

## run

```
./zhurd -config /usr/local/etc/zhurd/config.json
```
example of config look in `./share/config.json.example`



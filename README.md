# Ready for DB

## .envrc

```sh
cp .envrc.example .envrc
direnv allow
```

## docker

```sh
docker compose up -d
```

## psql login

```sh
psql -h 127.0.0.1 -p 5432 -U $(POSTGRES_USER) $(POSTGRES_DB)
```

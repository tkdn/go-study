# Ready for DB

## .envrc

```sh
cp .envrc.example .envrc
direnv allow
```

## docker build

```sh
docker buildx build \
	--build-arg dbuser=$(POSTGRES_USER) \
	--build-arg dbpass=$(POSTGRES_PASSWORD) \
	--build-arg dbname=$(POSTGRES_DB) \
	-f ./docker/postgresql/Dockerfile \
	-t postgresql \
	--no-cache \
	./
```

## docker run

```sh
docker run -d -p 5433:5432 -it --name app-postgre postgresql
docker start app-postgre
```

## psql login

```sh
psql -h 127.0.0.1 -p 5433 -U $(POSTGRES_USER) $(POSTGRES_DB)
```

## Dependencies
- docker
- kind
- kubectl
- kustomize

## Setup
1. Pull docker images
```sh
make dev-docker-pull
```
2. run k8s cluster
```sh
make dev-up
```
3. 
```sh
make service
```
4. 
```sh
make dev-load
```
5. 
```sh
make dev-apply
```
## Database
1. Since we use extension postgis, we need to build our own image, or pull from vincent0426/tuber/postgres
```sh
make dev-db-build

docker start postgres
```

2. run migration
```sh
make db/migrations/up

make db/migrations/down
```

3. (optional) run seed
```sh
make db/seed/up

make db/seed/down
```
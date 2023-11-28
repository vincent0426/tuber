## Dependencies
- docker
- kind 
- kubectl
- kustomize
- psql

## Git Commit Message
https://www.conventionalcommits.org/en/v1.0.0/

## Database Setup

We need **postgis** for our postgres. 

If not using Apple M1 chip or you want to build your own pg, please run the following command, make sure to change the image name if needed and update the image name under container section in `zarf/k8s/dev/database/dev-database.yaml` to use your own image.

```sh
make service-db
```
else if you are using Apple M1 chip and does not want to change anything, get docker image for postgres
```
docker pull vincent0426/tuber-postgres
```

## Setup
1. create k8s cluster with Kind
```sh
make dev-up
```
2. create docker image
```sh
make all
```
3. load images built in step 2 to Kind container
```sh
make dev-load
```
4. apply all k8s config 
```sh
make dev-apply
```

## Test endpoint
### port-forward

for main service
```sh
make dev-port-forward

curl localhost:3000/v1/ping
```

for chat service
```sh
make dev-chat-port-forward

curl localhost:3002/v1/chat/ping
```

## Update
### code only
```sh
make dev-update
```
### update both config file and code
```sh
make dev-update-apply
```

## Database

### Migrations
```sh
make db-migrations-up

make db-migrations-down
```

### Seed
```sh
make db-seed-up

make db-seed-down
```

## Other commands
```sh
make dev-down  # delete k8s cluster

make dev-delete  # delete all k8s config, but not delete k8s cluster

make dev-status  # check k8s cluster status

make dev-logs  # check k8s cluster logs

make dev-chat-logs  # check chat service logs

make pgcli  # connect to postgres database if you have pgcli installed
```
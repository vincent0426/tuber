## Dependencies
- docker
- kind 
- kubectl
- kustomize

## Git Commit Message
https://www.conventionalcommits.org/en/v1.0.0/

# Not Complete
## Setup

We need **postgis** for our postgres. 

If not using Apple M1 chip, please run the following command
```sh
make service-db
```
else if using Apple M1 chip own get docker image for postgres
```
docker pull vincent0426/tuber-postgres
```
1. create k8s cluster with Kind
```sh
make dev-up
```
2. create docker image
```sh
make all
```
3. load images built in step 3 to kind
```sh
make dev-load
```
4. 
```sh
make dev-apply
```

5. port forward

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
if update code only
```sh
make dev-update
```
if update config and code
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
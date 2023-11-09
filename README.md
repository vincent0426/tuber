## Dependencies
- docker
- kind 
- kubectl
- kustomize

## Git Commit Message
https://www.conventionalcommits.org/en/v1.0.0/

# Not Complete
## Setup
1. run k8s cluster
```sh
make dev-up
```
2. 
```sh
make service
```
3. 
```sh
make dev-load
```
4. 
```sh
make dev-apply
```

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
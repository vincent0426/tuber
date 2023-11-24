DB_DSN=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable

# ==============================================================================
# Define dependencies
NAMESPACE			  := tuber-system
APP							:= tuber
BASE_IMAGE_NAME := tuber
KIND            := kindest/node:v1.27.3
KIND_CLUSTER    := tuber
# POSTGRES        := postgres:15.4
POSTGRES        := vincent0426/tuber-postgres
TEMPO           := grafana/tempo:2.2.0
LOKI            := grafana/loki:3.3.1

# VERSION         := dev
VERSION         := 0.0.1
SERVICE_NAME    := tuber-api
SERVICE_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)

dev-brew:
	brew update
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize
	brew list helmfile || brew install helmfile
	brew list pgcli || brew install pgcli

# ==============================================================================
# Building containers
all: service

service:
	docker build \
	-f zarf/docker/dockerfile.service \
	-t $(SERVICE_IMAGE) \
	--build-arg BUILD_REF=$(VERSION) \
	--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
	.

# ==============================================================================
# Running from within k8s/kind
dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml
	
	kubectl wait --timeout=300s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner
	
	kind load docker-image $(POSTGRES) --name $(KIND_CLUSTER)

dev-load:
	cd zarf/k8s/dev/tuber; kustomize edit set image service-image=$(SERVICE_IMAGE)
	kind load docker-image $(SERVICE_IMAGE) --name $(KIND_CLUSTER)

dev-apply:
	kustomize build zarf/k8s/dev/database | kubectl apply -f -
	kubectl rollout status --namespace=$(NAMESPACE) --watch --timeout=300s sts/database

	kustomize build zarf/k8s/dev/grafana | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=grafana --timeout=300s --for=condition=Ready

	kustomize build zarf/k8s/dev/tempo | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=tempo --timeout=300s --for=condition=Ready

# helmfile is used to deploy helm charts.
	helmfile -n $(NAMESPACE) -f zarf/k8s/dev/prometheus/prometheus-helmfile.yaml sync
	kubectl wait --for=condition=ready pod --selector=app.kubernetes.io/instance=kube-prometheus-stack --namespace $(NAMESPACE) --timeout=600s
	
	helmfile -n $(NAMESPACE) -f zarf/k8s/dev/loki/loki-helmfile.yaml sync
	kubectl wait --for=condition=ready pod --selector=app=loki --namespace $(NAMESPACE) --timeout=600s
	
# create redis secret
	kustomize build zarf/k8s/dev/redis | kubectl apply -f -
	helmfile -n $(NAMESPACE) -f zarf/k8s/dev/redis/redis-helmfile.yaml sync
	kubectl wait --for=condition=ready pod --selector=app.kubernetes.io/instance=redis --namespace $(NAMESPACE) --timeout=300s
	
	kustomize build zarf/k8s/dev/tuber | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(APP) --timeout=300s --for=condition=Ready

dev-delete:
# helmfile -n $(NAMESPACE) -f zarf/k8s/dev/prometheus/prometheus-helmfile.yaml destroy
	helmfile -n $(NAMESPACE) -f zarf/k8s/dev/redis/redis-helmfile.yaml destroy
	helmfile -n $(NAMESPACE) -f zarf/k8s/dev/loki/loki-helmfile.yaml destroy
	kustomize build zarf/k8s/dev/redis | kubectl delete -f -
	kustomize build zarf/k8s/dev/tempo | kubectl delete -f -
	kustomize build zarf/k8s/dev/grafana | kubectl delete -f -
	kustomize build zarf/k8s/dev/database | kubectl delete -f -
	kustomize build zarf/k8s/dev/tuber | kubectl delete -f -

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) --all-containers=true -f --tail=100

dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

dev-restart:
	kubectl rollout restart deployment $(APP) --namespace=$(NAMESPACE)

dev-update: all dev-load dev-restart

dev-update-apply: all dev-load dev-apply

dev-docker-pull:
	docker pull $(KIND)
	docker pull $(POSTGRES)

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

tidy:
	go mod tidy
	go mod vendor
	
pgcli:
	pgcli postgresql://postgres:postgres@localhost

# db-migrations-new name=$1: create a new database migration
db-migrations-new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./db/migrations ${name}
	
# db-migrations-up: apply all up database migrations
db-migrations-up:
	@echo 'Running up migrations...'
	migrate -path ./db/migrations -database ${DB_DSN} up
	
# db-migrations-down: apply all down database migrations
db-migrations-down:
	@echo 'Running down migrations...'
	migrate -path ./db/migrations -database ${DB_DSN} down

# receive version from command line
db-migrations-force:
	@echo 'Running force migrations...'
	migrate -path ./db/migrations -database ${DB_DSN} force $(version)

# db-seed-up: apply all up database seeds
db-seed-up:
	@echo 'Running up seeds...'
	psql -h localhost -p 5432 -U postgres -d postgres -a -f db/seed/up.sql

# db-seed-down: apply all down database seeds
db-seed-down:
	@echo 'Running down seeds...'
	psql -h localhost -p 5432 -U postgres -d postgres -a -f db/seed/down.sql
	
# db-drop-all: drop all tables
db-drop-all:
	@echo 'Dropping all tables...'
	psql -h localhost -p 5432 -U postgres -d postgres -a -f db/drop-all.sql
	
run-local:
	go run app/services/tuber-api/main.go
	
run-local-db:
	docker run -p 5432:5432 --name postgres -d vincent0426/tuber-postgres	

create-local:
	@NAME=$$(awk 'BEGIN { srand(); print "Name_" int(rand()*10000) }'); \
	EMAIL=$$(awk 'BEGIN { srand(); print "user" int(rand()*10000) "@example.com" }'); \
	curl -iX POST 'http://localhost:3000/v1/users' \
	-H 'Content-Type: application/json' \
	--data-raw '{ \
	  "name": "'$$NAME'", \
	  "email": "'$$EMAIL'", \
	  "bio": "Experienced software developer with a passion for AI.", \
	  "acceptNotification": true \
	}'

# how to use: make update-local id=1
# add -H "Origin: http://localhost:5173" to test CORS
update-local:
	@NAME=$$(awk 'BEGIN { srand(); print "Name_" int(rand()*10000) }'); \
	EMAIL=$$(awk 'BEGIN { srand(); print "user" int(rand()*10000) "@example.com" }'); \
	IMAGE=$$(awk 'BEGIN { srand(); print "https://example.com/" int(rand()*10000) ".jpg" }'); \
	curl -iX PUT 'http://localhost:3000/v1/users/$(id)' \
	-H 'Content-Type: application/json' \
	--data-raw '{ \
	  "name": "'$$NAME'", \
	  "email": "'$$EMAIL'", \
		"imageURL": "'$$IMAGE'", \
	  "bio": "Experienced software developer with a passion for AI.", \
	  "acceptNotification": true \
	}'

# how to use: make delete-local id=1
delete-local:
	@curl -iX DELETE 'http://localhost:3000/v1/users/$(id)'
# how to use: make query-local | jq
query-local:
	@curl 'http://localhost:3000/v1/users?page=1&rows=2&orderBy=name'
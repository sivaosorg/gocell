.PHONY: run build run-job build-job test tidy deps-upgrade deps-clean-cache swaggo

# ==============================================================================
# Start Rest
run:
	go run ./cmd/api/main.go start_server ./config/conf-local.yaml ./config/conf-params.yaml

build:
	go build ./cmd/api/main.go start_server ./config/conf-local.yaml ./config/conf-params.yaml

# ==============================================================================
# Start Job
run-job:
	go run ./cmd/job/main.go start_server_job ./config/conf-local.yaml

build-job:
	go build ./cmd/job/main.go start_server_job ./config/conf-local.yaml

# ==============================================================================
# Modules support
test:
	go test -cover ./...

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-clean-cache:
	go clean -modcache
 
# ==============================================================================
# Tools commands
swaggo:
	echo "Starting swagger generating"
	swag init -g **/**/*.go
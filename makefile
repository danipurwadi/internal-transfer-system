# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

GOPATH_DEFAULT := $(HOME)/go/pkg/mod

# Use the first argument passed to make as GOPATH_ARG, or default if not provided
GOPATH_ARG 		:= $(if $(GOPATH),$(GOPATH)/pkg/mod,$(GOPATH_DEFAULT))

# ==============================================================================
# Install dependencies
dev-gotooling:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest

# ==============================================================================
# Development Commands
tidy:
	go mod tidy

start:
	@GOPATH=$(GOPATH_ARG) docker compose -f zarf/docker/docker-compose.yml \
		-p internal-transfer-system up

stop:
	@GOPATH=$(GOPATH_ARG) docker compose -f zarf/docker/docker-compose.yml \
		-p internal-transfer-system stop

# Exit cleans up the docker-compose services using the down command and removes volumes.
exit:
	@GOPATH=$(GOPATH_ARG) docker compose -f zarf/docker/docker-compose.yml \
		-p internal-transfer-system down -v

# ==============================================================================
# Running tests within the local computer

test-r:
	CGO_ENABLED=1 go test -race -count=1 ./...

test-only:
	CGO_ENABLED=0 go test -count=1 ./...

lint:
	CGO_ENABLED=0 go vet ./...
	staticcheck -checks=all ./...

vuln-check:
	govulncheck ./...

test: test-only lint vuln-check

test-race: test-r lint vuln-check

# ==============================================================================
stats:
	open -a "Google Chrome" http://localhost:8090/debug/statsviz 
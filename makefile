# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

GOPATH_DEFAULT := $(HOME)/go/pkg/mod

# Use the first argument passed to make as GOPATH_ARG, or default if not provided
GOPATH_ARG 		:= $(if $(GOPATH),$(GOPATH)/pkg/mod,$(GOPATH_DEFAULT))

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

# Exit cleans up the docker-compose services using the down command.
exit:
	@GOPATH=$(GOPATH_ARG) docker compose -f zarf/docker/docker-compose.yml \
		-p internal-transfer-system down

# ======
stats:
	open -a "Google Chrome" http://localhost:8090/debug/statsviz 
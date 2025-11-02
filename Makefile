.PHONY: run build tidy

run:
	cd mcs-manager && go run ./cmd/mcs-manager

start:
	@echo "Starting backend and frontend (cross-platform)..."
	@if [ "$(OS)" = "Windows_NT" ]; then \
		pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/start-all.ps1 ; \
	else \
		sh scripts/start-all.sh ; \
	fi

build:
	cd mcs-manager && CGO_ENABLED=0 go build -o ../bin/mcs-manager ./cmd/mcs-manager

tidy:
	cd mcs-manager && go mod tidy

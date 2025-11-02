.PHONY: run build tidy

run:
	cd mcs-manager && go run ./cmd/mcs-manager

build:
	cd mcs-manager && CGO_ENABLED=0 go build -o ../bin/mcs-manager ./cmd/mcs-manager

tidy:
	cd mcs-manager && go mod tidy

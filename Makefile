include Makefile.ledger
all: install

install: go.sum
		GO111MODULE=on GOPROXY=direct go install -tags "$(build_tags)" ./cmd/hhd
		GO111MODULE=on GOPROXY=direct go install -tags "$(build_tags)" ./cmd/hhcli

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

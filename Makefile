.PHONY: install test server

mod-update:
	GO111MODULE=on go get -u -m
	GO111MODULE=on go mod tidy

mod-tidy:
	GO111MODULE=on go mod tidy

install:
	env GO111MODULE=on go build ./...
	env GO111MODULE=on go install github.com/mattn/go-sqlite3
	env GO111MODULE=on go get -u github.com/manifoldco/grafton
	env GO111MODULE=on go build -o grafton github.com/manifoldco/grafton/cmd
	./grafton generate

test:
	env GO111MODULE=on go run cmd/server/main.go --test=true

serve:
	env GO111MODULE=on go run cmd/server/main.go

.PHONY: install test server

install:
	env GO111MODULE=on go build ./...
	env GO111MODULE=on go install github.com/mattn/go-sqlite3
	go get github.com/manifoldco/grafton/cmd
	go build -o grafton github.com/manifoldco/grafton/cmd
	./grafton generate

test:
	go run cmd/server/main.go --test=true

serve:
	go run cmd/server/main.go

.PHONY: install test server

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

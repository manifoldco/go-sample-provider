install:
	go build ./...
	go install github.com/mattn/go-sqlite3
	go get github.com/manifoldco/grafton/cmd
	go build -o grafton github.com/manifoldco/grafton/cmd
	./grafton generate

test:
	./grafton test --product bear \
	  --plan ursa-minor --new-plan ursa-major \
	  --region all::global \
	  --client-id 21jtaatqj8y5t0kctb2ejr6jev5w8 \
	  --client-secret 3yTKSiJ6f5V5Bq-kWF0hmdrEUep3m3HKPTcPX7CdBZw \
	  --connector-port 3001 \
	  --exclude plan-change \
	  --exclude resource-measures \
	  http://localhost:8080/

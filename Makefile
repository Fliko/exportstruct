all: lint clean tidy test

lint:
	gofmt -s -w .
	golangci-lint run -E goconst,gocritic,gomnd,gosec,interfacer,maligned,misspell,prealloc,unconvert,unparam ./...

clean:
	rm -f exportstruct types.go

tidy:
	go mod tidy

docs:
	(sleep 1 && open http://localhost:6060/pkg/exportstruct/) & \
	godoc -http=:6060

test: clean lint
	echo "Make sure services db is started locally"
	go run main.go -u ebuser -p password --db ebdb -json -sql

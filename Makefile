HOST=localhost
PORT=9999

build:
	go build -o bin/protohackers cmd/protohackers/main.go

echo: build
	./bin/protohackers --handler=echo --host=$(HOST) --port=$(PORT)
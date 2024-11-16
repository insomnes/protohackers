HOST=localhost
PORT=9999

build:
	go build -o bin/protohackers cmd/protohackers/main.go
	go build -o bin/phchat cmd/phchat/main.go
	go build -o bin/phkv cmd/phkv/main.go
	go build -o bin/phmitm cmd/phmitm/main.go

echo: build
	./bin/protohackers --handler=echo --host=$(HOST) --port=$(PORT) --verbose

prime: build
	./bin/protohackers --handler=prime --host=$(HOST) --port=$(PORT) --verbose

means: build
	./bin/protohackers --handler=means --host=$(HOST) --port=$(PORT) --verbose

chat: build
	./bin/phchat --host=$(HOST) --port=$(PORT)

kv: build
	./bin/phkv --host=$(HOST) --port=$(PORT)

mitm: build
	./bin/phmitm --host=$(HOST) --port=$(PORT)

build:
	@go build -o ./bin/syncserver ./server
	@go build -o ./bin/syncclient ./client
	
run: build
	@./bin/syncserver
	@./bin/syncclient
	

test:
	@go test ./...
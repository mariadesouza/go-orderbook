redis: build server

build:
	@docker build -t md-redis-server .

server:
	@docker run --name orderbook-redis  -p 6379:6379 md-redis-server

client:
	@docker exec -it orderbook-redis redis-cli

stop:
	@docker stop $$(docker ps -a -q --filter ancestor=md-redis-server --format="{{.ID}}")

clean:
	@docker rm $$(docker stop $$(docker ps -a -q --filter ancestor=md-redis-server --format="{{.ID}}"))

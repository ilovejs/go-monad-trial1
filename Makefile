all: server client-cur client-monadic

vendor:
	glide install

server:
	go build srv/server.go

run-server:
	go run srv/server.go --source data/good.json

client-cur:
	go build current/client-cur.go

run-cur:
	go run current/client-cur.go http://localhost:8080 suzanne

client-monadic:
	go build monadic/client-monadic.go

run-monad:
	go run monadic/client-monadic.go http://localhost:8080 suzanne

clean:
	-rm -f server
	-rm -f client-cur
	-rm -f client-monadic

.PHONY: vendor all clean server client-cur client-monadic

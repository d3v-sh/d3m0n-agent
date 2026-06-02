build:
	go build -o d3m0n .

run:
	go run main.go

clean:
	rm -f sec-agent
	rm -rf data/ sessions/
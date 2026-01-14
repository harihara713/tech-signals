run: build
	./bin/goscrape

build: 
	go build -o ./bin/goscrape .
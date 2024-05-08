test:
	go test -v . -cover

build:
	go build -o golite-api main.go

docker:
	docker build . -t golite-api:latest

install:build
	mkdir -p /usr/local/golite/bin
	mv golite-api /usr/local/golite/bin/
	mkdir -p /usr/local/golite/conf
	cp -rf resources/conf/* /usr/local/golite/conf/
	mkdir -p /usr/local/golite/assets
	cp -rf assets/* /usr/local/golite/assets/

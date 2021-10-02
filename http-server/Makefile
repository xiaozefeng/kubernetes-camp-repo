export tag=v1.0

build:
	echo "building http server binary"
	mkdir bin/amd64
	CGO_EANBLED=0 GOOS=linux GOARCH=amd64  go build -v -o bin/amd64 .

release: build 
	echo "building http servre container"
	docker build -t cncamp/httpserver:${tag} .

push: release
	echo "pushing cncamp/httpserver"
	docker push cncamp/httpserver:${tag}


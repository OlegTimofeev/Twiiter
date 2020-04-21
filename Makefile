


build:
	go mod download
	go build

build-docker:
	docker build -t twitter -f Dockerfile .

run-docker:
	docker run  -p 8080:8000 twitter


SWAGGER_IMAGE=quay.io/goswagger/swagger:v0.23.0

gen-trusted-client-server:
	docker run --rm -v `pwd`:/go/ -w /go/ -t $(SWAGGER_IMAGE) \
	generate server \
	--target=twitter \
	-f twitter.swagger.yml

	docker run --rm -v `pwd`:/go/ -w /go/ -t $(SWAGGER_IMAGE) \
	generate client \
	--target=twitter \
	-f twitter.swagger.yml


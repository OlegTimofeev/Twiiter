


build:
	go mod download
	go build

build-docker:
	docker build -t twitter -f Dockerfile .

run-docker:
	docker run  -p 8080:8000 twitter



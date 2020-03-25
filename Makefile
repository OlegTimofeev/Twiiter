


build:
	go mod download
	go build

build-docker:
	docker build -t twitter -f Dockerfile .

run-docker:
	docker run  -p 8080:8000 twitter
# docker ps
# docker ps -a
# docker images
# docker rm c29730347625
# docker images prune
# docker rmi c29730347625
# docker container prune
# docker pull postgers:9.6



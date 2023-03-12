.PHONY: all docker-image
LDFLAGS=-static
all: pause docker-image
pause: pause.c
docker-image: pause Dockerfile
	docker image build . -t "minik8s/pause"

IMAGE=localhost/salsa-moves

build:
	go build -o _output/salsa-moves .

run:
	_output/salsa-moves

container-build:
	podman build -t $(IMAGE) -f Dockerfile

container-run:
	podman run --rm -it $(IMAGE) /usr/local/bin/salsa-moves

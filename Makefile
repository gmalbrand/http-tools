SUBDIRS = $(notdir $(wildcard ./cmd/*))

.PHONY: all
all: clean dep build

.PHONY: build
build: $(addprefix ./bin/, $(SUBDIRS))

$(SUBDIRS): %:./bin/%

./bin/%: ./cmd/%/main.go
	go build -o $@ $^

.PHONY: dep
dep:
	go mod download

.PHONY: clean 
clean:
	rm -rf ./bin

.PHONY: docker-image
docker-image:
	docker build -t gmalbrand/http-mirror:latest -f build/Dockerfile .

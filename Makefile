SRC=$(wildcard **/*.go)
LANGUAGE=$(wildcard language/*.yml)

all: homecontrol-tubestatus
run: all

homecontrol-tubestatus: $(SRC) main.go language/strings.go
	go build

language/strings.go: $(LANGUAGE)
	${GOPATH}/bin/go-bindata -o language/strings.go -pkg language -ignore language/*.go language/*.yml

run:
	./homecontrol-tubestatus

clean:
	rm -f language/strings.go
	rm -f homecontrol-tubestatus

SRC=$(wildcard **/*.go)
LANGUAGE=$(wildcard language/*.yml)

all: tfl-announcements

tfl-announcements: $(SRC) main.go language/strings.go
	go build

language/strings.go: $(LANGUAGE)
	${GOPATH}/bin/go-bindata -o language/strings.go -pkg language -ignore language/*.go language/*.yml

run: all
	./tfl-announcements --verbose

clean:
	rm -f language/strings.go
	rm -f tfl-announcements

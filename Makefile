SRC=$(wildcard **/*.go)
LANGUAGE=$(wildcard language/*.yml)
ASSETS=$(wildcard web/assets/**)

all: tfl-announcements

tfl-announcements: $(SRC) main.go language/strings.go web/assets/assets.go
	go build

language/strings.go: $(LANGUAGE)
	${GOPATH}/bin/go-bindata -o language/strings.go -pkg language -ignore language/*.go language/*.yml

web/assets/assets.go: $(ASSETS)
	${GOPATH}/bin/go-bindata -o web/assets.go -pkg web web/assets/...

run: all
	./tfl-announcements --verbose

clean:
	rm -f language/strings.go
	rm -f tfl-announcements

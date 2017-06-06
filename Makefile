all: dist/natural

dist/natural:
	go build -o dist/natural github.com/SimonRichardson/naturalsort/cmd/natural

clean: FORCE
	rm -rf dist

FORCE:

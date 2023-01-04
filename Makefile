
:ONESHELL:

dist:
	mkdir $@
	go install github.com/spf13/cobra-cli@v1.3.0
	- go mod init github.com/christianlc-highlights/stripseven

run:
	go run main.go

build: cleandist dist
	go build -o dist/stripseven main.go

lint:
	goimports -l -w .
	golint ./...
	go vet ./... ||:

cleandist:
	rm -rvf dist

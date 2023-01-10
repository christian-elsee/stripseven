
:ONESHELL:

dist:
	mkdir $@
	go install github.com/spf13/cobra-cli@v1.3.0
	- go mod init github.com/christianlc-highlights/stripseven

	# add bats struts
	# https://bats-core.readthedocs.io/en/stable/tutorial.html#quick-installation
	mkdir -p $@/test
	git submodule add -f https://github.com/bats-core/bats-core.git $@/test/bats
	git submodule add -f https://github.com/bats-core/bats-support.git $@/test/test_helper/bats-support
	git submodule add -f https://github.com/bats-core/bats-assert.git $@/test/test_helper/bats-assert

run:
	go run main.go

build: dist
	go build -o dist/stripseven main.go

check: dist build
	rsync -av bats/ dist/test
	dist/test/bats/bin/bats --tap dist/test

lint:
	goimports -l -w .
	golint ./...
	go vet ./... ||:

cleandist:
	rm -rvf dist

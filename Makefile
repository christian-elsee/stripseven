
export NAME := $(shell pwd | xargs basename)
export TS := $(shell date +%s)

.DEFAULT_GOAL := @goal
.ONESHELL:

## recipe
@goal: cleandist dist build install check

dist: export sha := $(shell git rev-parse --short HEAD)
dist: export portecho ?= 2000
dist:
	mkdir $@
	go mod init github.com/christianlc-highlights/stripseven ||:
	go mod tidy

	# add bats struts
	# https://bats-core.readthedocs.io/en/stable/tutorial.html#quick-installation
	mkdir -p $@/test/test_helper
	git clone https://github.com/bats-core/bats-core.git $@/test/bats
	git clone https://github.com/bats-core/bats-support.git $@/test/test_helper/bats-support
	git clone https://github.com/bats-core/bats-assert.git $@/test/test_helper/bats-assert

	# interpolate manifest
	cat manifest/* \
		| envsubst \
		| tee $@/manifest.yaml

run:
	go run main.go

build: dist
	go build -o dist/build main.go

install: dist build
	kubectl config set-context --current --namespace $(NAME)
	kubectl apply -f dist/manifest.yaml
	kubectl rollout status deployment

check: dist build install
	rsync -av bats/ dist/test
	dist/test/bats/bin/bats --tap dist/test ||:

lint:
	goimports -l -w .
	golint ./...
	go vet ./... ||:

cleandist:
	rm -rvf dist

clean:
	kubectl delete -f dist/manifest.yaml ||:


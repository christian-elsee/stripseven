
export NAME := $(shell pwd | xargs basename)
export TS := $(shell date +%s)

.DEFAULT_GOAL := @goal
.ONESHELL:

## recipe
@goal: distclean dist check-dist build check-build install check
@fubar: one two

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

	# add tests
	cp -rf test $@

	# interpolate manifest
	cat manifest/* \
		| envsubst \
		| tee $@/manifest.yaml

check-dist: distclean dist
	dist/test/bats/bin/bats --tap dist/test/dist.bats

build: dist
	go build -o dist/build main.go

check-build: build
	dist/test/bats/bin/bats --tap dist/test/build.bats

install: dist build
	kubectl config set-context --current --namespace $(NAME)
	kubectl apply -f dist/manifest.yaml --dry-run=server -oyaml \
		| kubectl apply -f-

	#kubectl rollout status deployment
	#kubectl get all -lapp.kubernetes.io/part-of=$(NAME)

check-install: install
	dist/test/bats/bin/bats --tap dist/test/install.bats

check: dist build install
	dist/test/bats/bin/bats --tap dist/test

distclean:
	rm -rvf dist

clean:
	kubectl delete -f dist/manifest.yaml ||:

lint:
	goimports -l -w .
	golint ./...
	go vet ./... ||:

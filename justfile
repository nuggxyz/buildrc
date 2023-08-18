all: binaries

build:
    ./hack/build

shell:
    ./hack/shell

binaries:
    docker buildx bake binaries

binaries-cross:
    docker buildx bake binaries-cross

install: binaries
    mkdir -p ~/bin
    install bin/build/your-app ~/bin/your-app

release VERSION=local:
    VERSION={{VERSION}} ./hack/release

validate-all: lint test validate-vendor validate-docs validate-gen

lint:
    docker buildx bake lint

validate-vendor:
    docker buildx bake validate-vendor

validate-docs:
    docker buildx bake validate-docs

validate-authors:
    docker buildx bake validate-authors

validate-gen:
    docker buildx bake validate-gen

update-all: vendor docs gen authors mod-outdated

vendor:
    ./hack/update-vendor

docs:
    ./hack/update-docs

authors:
    docker buildx bake update-authors

mod-outdated:
    docker buildx bake mod-outdated

gen:
    docker buildx bake update-gen --progress plain


test-driver:
    ./hack/test-driver
test:
    ./hack/test

test-unit:
    TESTPKGS=./... SKIP_INTEGRATION_TESTS=1 ./hack/test

test-integration:
    TESTPKGS=./tests ./hack/test

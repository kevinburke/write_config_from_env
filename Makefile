SHELL = /bin/bash -o pipefail

BUMP_VERSION := $(GOPATH)/bin/bump_version
MEGACHECK := $(GOPATH)/bin/megacheck
RELEASE := $(GOPATH)/bin/github-release

$(MEGACHECK):
	go get -u honnef.co/go/tools/cmd/megacheck

lint: $(MEGACHECK)
	go list ./... | grep -v vendor | xargs $(MEGACHECK) --ignore='github.com/kevinburke/write_config_from_env/*.go:S1002'
	go vet ./...

test: lint
	go list ./... | grep -v vendor | xargs go test

race-test: lint
	go list ./... | grep -v vendor | xargs go test -race

$(BUMP_VERSION):
	go get -u github.com/Shyp/bump_version

$(RELEASE):
	go get -u github.com/aktau/github-release

# Run "GITHUB_TOKEN=my-token make release version=0.x.y" to release a new version.
release: race-test | $(BUMP_VERSION) $(RELEASE)
ifndef version
	@echo "Please provide a version"
	exit 1
endif
ifndef GITHUB_TOKEN
	@echo "Please set GITHUB_TOKEN in the environment"
	exit 1
endif
	bump_version --version=$(version) main.go
	git push origin --tags
	mkdir -p releases/$(version)
	# Change the binary names below to match your tool name
	GOOS=linux GOARCH=amd64 go build -o releases/$(version)/write_config_from_env-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o releases/$(version)/write_config_from_env-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build -o releases/$(version)/write_config_from_env-windows-amd64 .
	# Change the Github username to match your username.
	# These commands are not idempotent, so ignore failures if an upload repeats
	github-release release --user kevinburke --repo write_config_from_env --tag $(version) || true
	github-release upload --user kevinburke --repo write_config_from_env --tag $(version) --name write_config_from_env-linux-amd64 --file releases/$(version)/write_config_from_env-linux-amd64 || true
	github-release upload --user kevinburke --repo write_config_from_env --tag $(version) --name write_config_from_env-darwin-amd64 --file releases/$(version)/write_config_from_env-darwin-amd64 || true
	github-release upload --user kevinburke --repo write_config_from_env --tag $(version) --name write_config_from_env-windows-amd64 --file releases/$(version)/write_config_from_env-windows-amd64 || true

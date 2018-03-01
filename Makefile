# Release variables (usually for binaries not libraries)
RELEASE		?= 0.0.1
COMMIT		?= $(shell git rev-parse --short HEAD)
BUILD_TIME	?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
PROJECT		?= $(shell go list)

# Tools
GOJUNITREPORT	:= $(GOPATH)/bin/go-junit-report
GOCOV		:= $(GOPATH)/bin/gocov
GOCOVXML	:= $(GOPATH)/bin/gocov-xml
STATICCHECK	:= $(GOPATH)/bin/staticcheck
GOLINT		:= $(GOPATH)/bin/golint
GOSIMPLE	:= $(GOPATH)/bin/gosimple

build:
	go build -ldflags "-s -w"

test:
	go test -v -race ./...

coverage:
	@echo 'mode: atomic' > cover.profile && go list ./... \
		| xargs -n1 -I{} sh -c 'go test -covermode=atomic \
			-coverprofile=cover.profile.tmp {} && \
			tail -n +2 cover.profile.tmp >> cover.profile' && \
			rm cover.profile.tmp
	@go tool cover -func=cover.profile | grep total | awk '{ printf("coverage: %s\n", $$3) }'

reports/tests.xml: $(GOJUNITREPORT)
	@mkdir -p reports
	go test -v -race ./... 2>&1 | $(GOJUNITREPORT) > reports/tests.xml

reports/coverage.xml: $(GOCOV) $(GOCOVXML) coverage
	@mkdir -p reports
	$(GOCOV) convert cover.profile | $(GOCOVXML) > reports/coverage.xml

check: vet staticcheck golint gosimple

$(GOJUNITREPORT):
	go get -u github.com/jstemmer/go-junit-report

$(GOCOV):
	go get -u github.com/axw/gocov/gocov

$(GOCOVXML): $(GOCOV)
	go get -u github.com/AlekSi/gocov-xml

$(STATICCHECK):
	go get -u honnef.co/go/staticcheck/cmd/staticcheck

$(GOLINT):
	go get -u github.com/golang/lint/golint

$(GOSIMPLE):
	go get -u honnef.co/go/simple/cmd/gosimple

vet:
	go vet ./...

staticcheck: $(STATICCHECK)
	$(STATICCHECK) ./...

golint: $(GOLINT)
	$(GOLINT) ./...

gosimple: $(GOSIMPLE)
	$(GOSIMPLE) ./...

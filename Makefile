COVERAGE_DIR  := ./coverage
COVERAGE_OUT  := ${COVERAGE_DIR}/coverage.out
COVERAGE_HTML := ${COVERAGE_DIR}/coverage.html

.PHONY: build
build:
	go build

.PHONY: lint
lint:
	golangci-lint run

.PHONY: coverage
coverage:
	@mkdir -p ${COVERAGE_DIR}
	go test -coverprofile=${COVERAGE_OUT} -count=1
	go tool cover -html $(COVERAGE_OUT) -o $(COVERAGE_HTML)

COVERAGE_DIR  := ./coverage
COVERAGE_OUT  := ${COVERAGE_DIR}/coverage.out
COVERAGE_HTML := ${COVERAGE_DIR}/coverage.html
COVERAGE_LCOV := ${COVERAGE_DIR}/coverage.lcov.info

.PHONY: build
build:
	go build

.PHONY: lint
lint:
	golangci-lint run

.PHONY: coverage
coverage: coverage-dir ${COVERAGE_HTML} ${COVERAGE_LCOV}

.PHONY: coverage-dir
coverage-dir:
	@mkdir -p ${COVERAGE_DIR}

${COVERAGE_OUT}:
	go test -coverprofile='${COVERAGE_OUT}' -count=1

${COVERAGE_HTML}: ${COVERAGE_OUT}
	go tool cover -html '${COVERAGE_OUT}' -o '${COVERAGE_HTML}'

${COVERAGE_LCOV}: ${COVERAGE_OUT}
	GOROOT=$$(go env GOROOT) gcov2lcov -infile '${COVERAGE_OUT}' -outfile '${COVERAGE_LCOV}'

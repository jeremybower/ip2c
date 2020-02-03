.PHONY: integration test unit

ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
COVERAGE_DIR=${ROOT_DIR}/coverage

clean:
	rm -rf ${COVERAGE_DIR}

integration:
	go test -cover -tags=integration

test: 
	mkdir -p ${COVERAGE_DIR}
	go test -tags="integration unit" -race -coverprofile=${COVERAGE_DIR}/coverage.out -covermode=atomic
	go tool cover -html=${COVERAGE_DIR}/coverage.out -o ${COVERAGE_DIR}/coverage.html

unit:
	go test -cover -tags=unit

NAME:=github-branch-bot
BUILD_DIR:=bin
LDFLAGS:=-ldflags "-s -w"

clean:
	@echo "===> Cleaning build directories"
	@-rm -rfv ${BUILD_DIR} vendor

deps:
	@echo "===> Running deps"
	@dep ensure -v

test:
	@echo "===> Running go test"
	@go test -cover $$(go list ./... | grep -v /vendor/)

build: deps test
	@echo "===> Running go build"
	@go build -v ${LDFLAGS} -o ${BUILD_DIR}/${NAME} ./cmd/${NAME}

package: build
	@echo "===> Building zip"
	@zip -J -r ${BUILD_DIR}/${NAME}.zip ${BUILD_DIR}

.PHONY: clean deps build package
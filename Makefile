NAME:=github-branch-bot
BUILD_DIR:=$$PWD/bin
LDFLAGS:=-ldflags "-s -w"
REPO:=/go/src/github.com/aaron-vaz/${NAME}

clean:
	@echo "===> Cleaning build directories"
	@-rm -rfv ${BUILD_DIR} vendor

setup:
	@mkdir -p ${BUILD_DIR}

deps:
	@echo "===> Running deps"
	@dep ensure -v

test: 
	@echo "===> Running go test"
	@go test -cover $$(go list ./... | grep -v /vendor/)

build: deps setup test
	@echo "===> Running go build"
	@go build -v ${LDFLAGS} -o ${BUILD_DIR}/${NAME} ./cmd/${NAME}

package: build
	@echo "===> Building zip"
	@zip -J -r ${BUILD_DIR}/${NAME}.zip ${BUILD_DIR}

docker: setup
	@echo "===> Building in docker container"
	@docker build --build-arg REPO=${REPO} -t ${NAME} .
	@docker run --rm --volume ${BUILD_DIR}:${REPO}/bin -t ${NAME}	

.PHONY: clean setup deps build test package docker
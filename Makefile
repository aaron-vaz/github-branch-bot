NAME:=github-branch-bot
BUILD_DIR:=bin
LDFLAGS:=-ldflags "-s -w"
REPO:=/go/src/github.com/aaron-vaz/${NAME}

clean:
	@echo "===> Cleaning build directories"
	@-rm -rfv ${BUILD_DIR} vendor .serverless

setup:
	@mkdir -p ${BUILD_DIR}

test: 
	@echo "===> Running go test"
	@go test -cover ./...

build: setup test
	@echo "===> Running go build"
	@go build -v ${LDFLAGS} -o ${BUILD_DIR}/${NAME} ./cmd/${NAME}

package: build
	@echo "===> Building zip"
	@zip -J -r ${BUILD_DIR}/${NAME}.zip ${BUILD_DIR}

docker: setup
	@echo "===> Building in docker container"
	@docker build --build-arg REPO=${REPO} -t ${NAME} .
	@docker run --rm --volume $$PWD/${BUILD_DIR}:${REPO}/bin -t ${NAME}

deploy: docker
	@echo "===> Deploying to AWS"
	@sls deploy	

.PHONY: clean setup build test package docker
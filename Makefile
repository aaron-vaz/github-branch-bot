NAME:=github-branch-bot
BUILD_DIR:=bin
LDFLAGS:=-ldflags "-s -w"
REPO:=/go/src/github.com/aaron-vaz/${NAME}

BRANCH_BOT:=${NAME}
BRANCH_CHECK=github-branch-check

clean:
	@echo "===> Cleaning build directories"
	@-rm -rfv ${BUILD_DIR} vendor .serverless

setup:
	@mkdir -p ${BUILD_DIR}

test: 
	@echo "===> Running go test"
	@go test -cover ./...

build: setup test
	@echo "===> Building ${BRANCH_BOT}"
	@go build -v ${LDFLAGS} -o ${BUILD_DIR}/${BRANCH_BOT} ./cmd/${BRANCH_BOT}

	@echo "===> Building ${BRANCH_CHECK}"
	@go build -v ${LDFLAGS} -o ${BUILD_DIR}/${BRANCH_CHECK} ./cmd/${BRANCH_CHECK}

package: build
	@echo "===> Building zips"
	@zip -J -r ${BUILD_DIR}/${BRANCH_BOT}.zip ${BUILD_DIR}/${BRANCH_BOT}
		@zip -J -r ${BUILD_DIR}/${BRANCH_CHECK}.zip ${BUILD_DIR}/${BRANCH_CHECK}

docker: setup
	@echo "===> Building in docker container"
	@docker build --build-arg REPO=${REPO} -t ${NAME} .
	@docker run --rm --volume $$PWD/${BUILD_DIR}:${REPO}/bin -t ${NAME}

deploy: docker
	@echo "===> Deploying to AWS"
	@sls deploy	

.PHONY: clean setup build test package docker
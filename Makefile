# ==================SETUP APPLICATION BEFORE DEVELOPMENT==================

# Setup Application name
APP_NAME=homestead
APP_PORT=8081
JWT_KEY=secret


# Setup MySQL Configuration
MYSQL_HOST=localhost
MYSQL_USER=root
MYSQL_PASS=root
MYSQL_DB=modulordb
MYSQL_PORT=3306

# # Setup Registry Deployment for deploy your server
# DOCKER_DEPLOY_NAME=xxxxx/xxxxx
# DOCKER_LOGIN_USERNAME=xxxxx
# DOCKER_LOGIN_PASSWORD=xxxxx
# DOCKER_DEPLOY_TAG=xxxxx


# ========================END OF LINE CONFIGURATION========================

service-up:
	SERVICE_NAME=${APP_NAME} MYSQL_PWD=${MYSQL_PASS} MYSQL_DB=${MYSQL_DB} \
	docker-compose -f services/docker-compose.yaml up -d

service-down:
	SERVICE_NAME=${APP_NAME} docker-compose -f services/docker-compose.yaml down

local-test:
	go clean -testcache
	SERVICE_NAME=${APP_NAME} docker-compose -f ./deployments/test/docker-compose.yaml up -d database
	BACKEND_PORT=8081 \
	JWT_KEY=${JWT_KEY} \
	DRIVER_TYPE=mysql \
	MYSQL_HOST=localhost \
	MYSQL_USER=root \
	MYSQL_PORT=3307 \
	MYSQL_PWD=251199 \
	MYSQL_DB=modulordb \
	go test -v ./app/test

local-run:
	BACKEND_PORT=${APP_PORT} \
	JWT_KEY=${JWT_KEY} \
	DRIVER_TYPE=mysql \
	MYSQL_HOST=${MYSQL_HOST} \
	MYSQL_USER=${MYSQL_USER} \
	MYSQL_PORT=${MYSQL_PORT} \
	MYSQL_PWD=${MYSQL_PASS} \
	MYSQL_DB=${MYSQL_DB} \
	go run cmd/main.go 

cloud-test :
	docker build  -t ${APP_NAME}:test -f ./deployments/test/Dockerfile .
	SERVICE_NAME=${APP_NAME} docker-compose -f deployments/test/docker-compose.yaml up -d database
	SERVICE_NAME=${APP_NAME} docker-compose -f deployments/test/docker-compose.yaml up app

cloud-deploy :
	docker login -u ${DOCKER_LOGIN_USERNAME} -p ${DOCKER_LOGIN_PASSWORD}
	docker tag ${APP_NAME}:latest ${DOCKER_DEPLOY_NAME}:${DOCKER_DEPLOY_TAG}
	docker push ${DOCKER_DEPLOY_NAME}:${DOCKER_DEPLOY_TAG}

image-build :
	docker build  -t ${APP_NAME}:latest -f ./deployments/build/Dockerfile .

image-remove :
	docker image rm ${APP_NAME}:latest

clean-test :
	docker-compose -f deployments/test/docker-compose.yaml down

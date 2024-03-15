BINARY_NAME := vaccine
MYSQL_CONTAINER :=vulnerables/web-dvwa
SQLITE_CONTAINER := gitlab.cylab.be:8081/cylab/play/sqlite-injection

# ================================ Sources =================================== #

SRC := main.go scraper.go http_client.go sql_injector.go basic_commands.go

# ================================ Rules ===================================== #
.SILENT:

all: run

run:
	go build -o ${BINARY_NAME} ${SRC}

docker-start:
		docker pull ${MYSQL_CONTAINER}
		docker run -it -d -p 8080:80 ${MYSQL_CONTAINER}
		docker run -it -d -p 8081:80 ${SQLITE_CONTAINER}

docker-stop:
	docker stop $(shell docker ps -q --filter ancestor=${MYSQL_CONTAINER})
	docker stop $(shell docker ps -q --filter ancestor=${SQLITE_CONTAINER})

clean:	docker-stop
	go clean
	rm -f ${BINARY_NAME}
	docker rm $(shell docker ps -a -q --filter ancestor=${MYSQL_CONTAINER})
	docker rm $(shell docker ps -a -q --filter ancestor=${SQLITE_CONTAINER})

re: clean run

.PHONY: all run clean re

sqlite:	./$(BINARY_NAME)
	./$(BINARY_NAME) http://127.0.0.1:8081/ -X POST -v

mysql:	./$(BINARY_NAME)
	./$(BINARY_NAME) http://127.0.0.1:8080/vulnerabilities/sqli/ -s c6derp6ibht228uh4nj5pjogq7 -X GET -v



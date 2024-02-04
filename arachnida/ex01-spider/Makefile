BINARY_NAME := spider

# ================================ Sources =================================== #

SRC := main.go parse.go extract.go set.go

# ================================ Rules ===================================== #
.SILENT:

all: run

run:
	go build -o ${BINARY_NAME} ${SRC}

clean:
	go clean
	rm -f ${BINARY_NAME}

re: clean run

.PHONY: all run clean re


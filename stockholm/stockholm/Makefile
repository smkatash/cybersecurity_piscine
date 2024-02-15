BINARY_NAME := stockholm

# ================================ Sources =================================== #

SRC := main.go parse.go infect.go crypto.go

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


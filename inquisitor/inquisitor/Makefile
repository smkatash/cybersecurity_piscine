BINARY_NAME := inquisitor
# ================================ Sources =================================== #

SRC := main.go mitm.go parse.go

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


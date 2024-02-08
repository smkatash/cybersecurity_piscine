BINARY_NAME := ft_otp

# ================================ Sources =================================== #

SRC := main.go hotp.go cypher.go

# ================================ Rules ===================================== #
.SILENT:

all: run

run:
	go build -o ${BINARY_NAME} ${SRC}

clean:
	go clean
	rm -f ${BINARY_NAME}

fclean: clean
	rm ft_otp.key qrcode.png secret.key

re: clean run

.PHONY: all run clean re


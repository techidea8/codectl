package parse

import "testing"

var dsn1 = "user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true"
var dsn2 = "root:pw@unix(/tmp/mysql.sock)/myDatabase?loc=Local"
var dsn3 = "user@unix(/path/to/socket)/dbname"
var dsn4 = "/dbname"

func TestDSN(t *testing.T) {
	//ParseMysql(dsn4)
	ParseMysql(dsn1)
	ParseMysql(dsn2)
	ParseMysql(dsn3)
}

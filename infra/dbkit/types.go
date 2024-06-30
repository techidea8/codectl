package dbkit

// MySQL, PostgreSQL, SQLite, SQL Server 和 TiDB
type DBTYPE string

const (
	MySQL      DBTYPE = "mysql"
	PostgreSQL DBTYPE = "postgres"
	SQLite     DBTYPE = "sqlite"
	SqlServer  DBTYPE = "sqlserver"
	Tidb       DBTYPE = "tidb"
	ClickHouse DBTYPE = "clickhouse"
)

package dbkit

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func GussDbtype(dsn string) DBTYPE {
	dbtype := SQLite
	if strings.Contains(dsn, ":") {
		dbtype = MySQL
	}
	if strings.Contains(dsn, "host=") {
		dbtype = PostgreSQL
	}
	return dbtype
}
func TOString(l logger.LogLevel) string {
	if l == logger.Error {
		return "error"
	} else if l == logger.Info {
		return "INFO"
	} else if l == logger.Silent {
		return "SILENY"
	} else if l == logger.Warn {
		return "warn"
	} else {
		return fmt.Sprintf("%d", l)
	}
}

// 切换loglevel
func TOGormLogLevel(level string) (r logger.LogLevel) {
	_level := strings.ToLower(level)
	switch _level {
	case "error":
		r = logger.Error
	case "warn":
		r = logger.Warn
	case "debug":
	case "info":
	case "trace":
		r = logger.Info
	case "fatal":
	case "panic":
		r = logger.Error
	default:
		r = logger.Info
	}
	return r
}

// gorm.Dialector
func OpenDb(dst DBTYPE, dsn string, opts ...Option) (db *gorm.DB, err error) {
	var dialector gorm.Dialector
	if dst == "" {
		dst = GussDbtype(dsn)
	}
	if dst == MySQL {
		dialector = mysql.Open(dsn)
	} else if dst == PostgreSQL {
		dialector = postgres.Open(dsn)
	} else if dst == SQLite {
		dialector = sqlite.Open(dsn)
	} else if dst == SqlServer {
		dialector = sqlserver.Open(dsn)
	} else {
		err = fmt.Errorf("not suport %s", dst)
		return
	}
	ctx := NewDbContext()
	for _, opt := range opts {
		opt(ctx)
	}
	newLogger := logger.New(
		log.New(ctx.Writer, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,                   // Slow SQL threshold
			LogLevel:                  ctx.LogLevel,                  // Log level
			IgnoreRecordNotFoundError: ctx.IgnoreRecordNotFoundError, // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      ctx.ParameterizedQueries,      // Don't include params in the SQL log
			Colorful:                  false,                         // Disable color
		},
	)
	//fmt.Println("level%d=%s", ctx.LogLevel, TOString(ctx.LogLevel))
	db, err = gorm.Open(dialector, &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: ctx.SingularTable,
			TablePrefix:   ctx.TablePrefix,
		},
	})

	if err != nil {
		return nil, err
	}
	if ctx.Debug {
		db = db.Debug()
	}
	if len(ctx.ModuleMigrates) > 0 {
		err = db.AutoMigrate(ctx.ModuleMigrates...)
	}
	if err != nil {
		return
	}
	return db, err
}

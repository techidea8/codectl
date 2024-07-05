package dbkit

import (
	"io"
	"os"

	ctllogger "github.com/techidea8/codectl/infra/logger"
	"gorm.io/gorm/logger"
)

type Option func(*DbContext)

func AutoMigrate(obj ...interface{}) Option {
	return func(ctx *DbContext) {
		if ctx.ModuleMigrates == nil {
			ctx.ModuleMigrates = make([]interface{}, 0)
		}
		ctx.ModuleMigrates = append(ctx.ModuleMigrates, obj...)
	}
}

// SetLogLevel
func SetLogLevel(level ctllogger.LogLevel) Option {
	return func(ctx *DbContext) {
		//fmt.Println("SetLogLevel(level int32),%s ", level, TOString(logger.LogLevel(level)))
		ctx.LogLevel = TOGormLogLevel(level)
	}
}

// SetLogLevel
func WithPrefix(prefix string) Option {
	return func(ctx *DbContext) {
		ctx.TablePrefix = prefix
	}
}

// SetLogLevel
func IgnoreRecordNotFoundError(flag bool) Option {
	return func(ctx *DbContext) {
		ctx.IgnoreRecordNotFoundError = flag
	}
}

// SetLogLevel
func SingularTable(flag bool) Option {
	return func(ctx *DbContext) {
		ctx.SingularTable = flag
	}
}

// SetLogLevel
func WithWriter(writers ...io.Writer) Option {
	return func(ctx *DbContext) {
		ctx.Writer = io.MultiWriter(writers...)
	}
}

// SetLogLevel
func ParameterizedQueries(f bool) Option {
	return func(ctx *DbContext) {
		ctx.ParameterizedQueries = f
	}
}

func Debug(f bool) Option {
	return func(ctx *DbContext) {
		ctx.Debug = f
	}
}

type DbContext struct {
	LogLevel                  logger.LogLevel
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
	Writer                    io.Writer
	TablePrefix               string
	SingularTable             bool
	ModuleMigrates            []interface{}
	Debug                     bool
}

func NewDbContext() *DbContext {
	return &DbContext{
		LogLevel:                  logger.Info,
		IgnoreRecordNotFoundError: true,
		ParameterizedQueries:      true,
		Writer:                    os.Stdout,
		TablePrefix:               "",
		SingularTable:             true,
		ModuleMigrates:            make([]interface{}, 0),
		Debug:                     false,
	}
}

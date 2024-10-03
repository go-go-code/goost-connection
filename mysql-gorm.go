package connection

import (
	"fmt"
	"os"

	logger "github.com/go-go-code/goost-logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _mysql_gorm *mysqlGormConnection

type mysqlGormConnection struct {
	Client *gorm.DB
}

func (m *mysqlGormConnection) Close() {
	if m.Client == nil {
		return
	}

	defer func() {
		m.Client = nil
	}()

	sqldb, err := m.Client.DB()
	if err != nil {
		logger.ErrorF("Closing MySQL Error: %v", err)
		return
	}

	if err := sqldb.Close(); err != nil {
		logger.ErrorF("Closing MySQL Error: %v", err)
	} else {
		logger.Info("MySQL Connection Is Close")
	}
}

func NewMySQLGormConnection() *gorm.DB {
	if _mysql_gorm == nil {
		initMySQLGormConnection()
	}

	if _mysql_gorm.Client == nil {
		return &gorm.DB{}
	}

	return _mysql_gorm.Client
}

func initMySQLGormConnection() {
	_mysql_gorm = &mysqlGormConnection{}

	if enabled, ok := cfg["mysql_enable"].(bool); !ok || !enabled {
		logger.Info("⚠️ MySQL is Disabled ⚠️")
		return
	}

	host, ok := cfg["mysql_host"].(string)
	if !ok || host == "" {
		host = "127.0.0.1"
	}

	port, ok := cfg["mysql_port"].(string)
	if !ok || port == "" {
		port = "3306"
	}

	username, _ := cfg["mysql_username"].(string)
	password, _ := cfg["mysql_password"].(string)
	database, _ := cfg["mysql_database"].(string)
	charset, _ := cfg["mysql_charset"].(string)
	loc, _ := cfg["mysql_loc"].(string)

	dsn := "%s:%s@tcp(%s:%s)/%s?charset=%s&multiStatements=true&parseTime=True&loc=%s"
	dsn = fmt.Sprintf(dsn, username, password, host, port, database, charset, loc)

	client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.ErrorF("Connecting MySQL Error (M001) : %v", err)
		return
	}

	sqldb, err := client.DB()
	if err != nil {
		logger.ErrorF("Connecting MySQL Error (M002) : %v", err)
		return
	}

	if err := sqldb.Ping(); err != nil {
		logger.ErrorF("Connecting MySQL Error (M003) : %v", err)
		return
	}

	maxopenconns := 100
	sqldb.SetMaxOpenConns(maxopenconns)

	maxidleconns := 50
	sqldb.SetMaxIdleConns(maxidleconns)

	if os.Getenv("APP_ENV") == "develop" {
		_mysql_gorm.Client = client.Debug()
	} else {
		_mysql_gorm.Client = client
	}

	add(_mysql_gorm)

	logger.Info("Connecting MySQL Success ")
}

package connection

import (
	"fmt"

	logger "github.com/go-go-code/goost-logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var _mysql_sqlx *msqlSqlxConnection

type msqlSqlxConnection struct {
	Client *sqlx.DB
}

func (m *msqlSqlxConnection) Close() {
	if m.Client == nil {
		return
	}

	if err := m.Client.Close(); err != nil {
		logger.ErrorF("Closing MySQL Error : %v", err)
	} else {
		logger.Info("MySQL Connection Is Close")
	}

	m.Client = nil
}

func NewMySQLSqlxConnection() *sqlx.DB {
	if _mysql_sqlx == nil {
		initMySQLConnection()
	}

	if _mysql_sqlx.Client == nil {
		return &sqlx.DB{}
	}

	return _mysql_sqlx.Client
}

func initMySQLConnection() {
	_mysql_sqlx = &msqlSqlxConnection{}

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

	client, err := sqlx.Open("mysql", dsn)
	if err != nil {
		logger.ErrorF("❌ Connecting MySQL Error (M001) ❌ : %v", err)
		return
	}

	if err := client.Ping(); err != nil {
		logger.ErrorF("❌ Connecting MySQL Error (M002) ❌ : %v", err)
		return
	}

	maxopenconns := 100

	client.SetMaxOpenConns(maxopenconns)

	maxidleconns := 50
	client.SetMaxIdleConns(maxidleconns)

	_mysql_sqlx.Client = client

	add(_mysql_sqlx)

	logger.Info("Connecting MySQL Success")
}

package main

/*
Плагин для работы с PostgreSQL.
Плагин сохраняет пакет в jsonb поле point у заданной в настройках таблице.

Раздел настроек, которые должны отвечають в конфиге для подключения плагина:

[store]
plugin = "postgresql.so"
host = "localhost"
port = "5432"
user = "postgres"
password = "postgres"
database = "receiver"
table = "points"
sslmode = "disable"
*/

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type PostgreSQLConnector struct {
	connection *sql.DB
	config     map[string]string
}

func (c *PostgreSQLConnector) Init(cfg map[string]string) error {
	var (
		err error
	)
	if cfg == nil {
		return fmt.Errorf("Не корректная ссылка на конфигурацию")
	}
	c.config = cfg
	//connStr := fmt.Sprintf("%s:%s@/%s",
	//	c.config["user"], c.config["password"], c.config["database"])
        connStr := "igorek:1@/glonas"
	if c.connection, err = sql.Open("mysql", connStr); err != nil {
		return fmt.Errorf("Ошибка подключения к postgresql: %v", err)
	}
	return err
}

func (c *PostgreSQLConnector) Save(msg interface{ ToBytes() ([]byte, error) }) error {
	if msg == nil {
		return fmt.Errorf("Не корректная ссылка на пакет")
	}

	innerPkg, err := msg.ToBytes()
	if err != nil {
		return fmt.Errorf("Ошибка сериализации  пакета: %v", err)
	}

	insertQuery := fmt.Sprintf("INSERT INTO %s (point) VALUES (?)", c.config["table"])
	if _, err = c.connection.Exec(insertQuery, innerPkg); err != nil {
		return fmt.Errorf("Не удалось вставить запись: %v", err)
	}
	return nil
}

func (c *PostgreSQLConnector) Close() error {
	return c.connection.Close()
}

var Connector PostgreSQLConnector

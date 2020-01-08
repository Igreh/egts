package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
)

type egtsParsePacket struct {
	Client              uint32    `json:"client"`
	PacketID            uint32    `json:"packet_id"`
	NavigationTimestamp int64     `json:"navigation_unix_time"`
	ReceivedTimestamp   int64     `json:"received_unix_time"`
	Latitude            float64   `json:"latitude"`
	Longitude           float64   `json:"longitude"`
	Speed               uint16    `json:"speed"`
	Pdop                uint16    `json:"pdop"`
	Hdop                uint16    `json:"hdop"`
	Vdop                uint16    `json:"vdop"`
	Nsat                uint8     `json:"nsat"`
	Ns                  uint16    `json:"ns"`
	Course              uint8     `json:"course"`
	GUID                uuid.UUID `json:"guid"`
	// AnSensors           []anSensor     `json:"an_sensors"`
	// LiquidSensors       []liquidSensor `json:"liquid_sensors"`
}

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

func (c *PostgreSQLConnector) InitSession(imei string) int64 {
	var (
		err error
		res sql.Result
		id  int64
	)
	insertQuery := fmt.Sprintf("insert into sessions (imei, created) values (?, now());")
	if res, err = c.connection.Exec(insertQuery, imei); err != nil {
		panic(fmt.Errorf("Не удалось инициировать сессию. IMEI: %s", imei))
		// return 0
	}
	id, err = res.LastInsertId()
	if err != nil {
		panic(fmt.Errorf("Не удалось получить ID сессии. IMEI: %s", imei))
	}
	return id
}

func (c *PostgreSQLConnector) Save(msg interface{ ToBytes() ([]byte, error) }, session_id int64) error {
	var (
		// latitude, longitude float64
		p    egtsParsePacket
		s    string
		s_id sql.NullInt64
	)
	if msg == nil {
		return fmt.Errorf("Не корректная ссылка на пакет")
	}

	if session_id == 0 {
		s_id = sql.NullInt64{}
	} else {
		s_id = sql.NullInt64{Int64: session_id, Valid: true}
	}

	// latitude = msg.Latitude
	innerPkg, err := msg.ToBytes()
	json.Unmarshal(innerPkg, &p)
	// latitude = innerPkg.Latitude
	if err != nil {
		return fmt.Errorf("Ошибка сериализации  пакета: %v", err)
	}

	s = fmt.Sprintf("POINT(%f %f)", p.Latitude, p.Longitude)

	insertQuery := fmt.Sprintf("INSERT INTO points (coords, session_id, json, created, navigation_ts, speed) VALUES (PointFromText(?), ?, ?, from_unixtime(?), from_unixtime(?), ?)")
	// insertQuery := fmt.Sprintf("INSERT INTO points (coords, session_id, created) VALUES (PointFromText('POINT(49.2343503980067 2.52738212494082)'), ?, now())")

	if _, err = c.connection.Exec(insertQuery, s, s_id, innerPkg, p.ReceivedTimestamp, p.NavigationTimestamp, p.Speed); err != nil {
		// if _, err = c.connection.Exec(insertQuery, session_id); err != nil {
		return fmt.Errorf("Не удалось вставить запись: %v", err)
	}
	return nil
}

func (c *PostgreSQLConnector) Close() error {
	return c.connection.Close()
}

var Connector PostgreSQLConnector

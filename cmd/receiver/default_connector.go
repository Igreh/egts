package main

import (
	"encoding/json"
	"fmt"
)

type defaultConnector struct{}

func (c defaultConnector) Init(cfg map[string]string) error {
	return nil
}

func (c defaultConnector) InitSession(imei string) int64 {
	return 0
}

func (c defaultConnector) Save(msg interface{ ToBytes() ([]byte, error) }, session_id int64) error {
	jsonPkg, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		return fmt.Errorf("Не сформировать отладочный json:\n %v", err)
	}

	fmt.Println("Export packet: ", string(jsonPkg))
	return nil
}

func (c defaultConnector) Save2(lat string, long string, s_id int64) error {
	return nil
}

func (c defaultConnector) Close() error {
	return nil
}

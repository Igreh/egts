package main

//Connector интерфейс для подключения внешних хранилищ
type Connector interface {
	// установка соединения с хранилищем
	Init(map[string]string) error

	InitSession(string) int64

	// сохранение в хранилище
	Save(interface{ ToBytes() ([]byte, error) }, int64) error

	// сохранение в хранилище
	// Save2(float64, float64, int64) error

	//закрытие соединения с хранилищем
	Close() error
}

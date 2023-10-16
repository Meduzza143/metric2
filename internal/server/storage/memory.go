package storage

import "fmt"

type MemStruct struct {
	MetricType string
	Value      string
}
type memStorage map[string]MemStruct

var storage memStorage = nil

//TODO: implement interface

// type Storage interface {
// 	SetValue(key, metricType, value string)
// 	GetAllValues() memStorage
// 	GetValue(key string) MemStruct
// }

func (memStorage) SetValue(key, metricType, value string) {
	storage[key] = MemStruct{metricType, value}

	fmt.Printf("value has ben set [%v]\n ", storage[key])
}

func (memStorage) GetAllValues() memStorage {
	fmt.Printf("getting all values [%v]\n ", storage)
	return storage
}

func (memStorage) GetValue(key string) MemStruct {
	fmt.Printf("getting value [%v]\n ", storage[key])
	return storage[key]
}

func GetInstance() memStorage {
	if storage == nil {
		storage = make(memStorage)
	}
	return storage
}

package storage

type MemStruct struct {
	MetricType string
	//Value      string
	GaugeValue   float64
	CounterValue int64
	MetricName   string
}
type memStorage map[string]MemStruct

var storage memStorage = nil

//TODO: implement interface

// type Storage interface {
// 	SetValue(key, metricType, value string)
// 	GetAllValues() memStorage
// 	GetValue(key string) MemStruct
// }

// func (memStorage) SetValue(MetricName, metricType string, gaugeValue float64, counterValue int64) {
func (memStorage) SetValue(MetricName string, metric MemStruct) {

	metric.CounterValue = storage[MetricName].CounterValue + metric.CounterValue
	storage[MetricName] = metric

}

func (memStorage) GetAllValues() memStorage {
	//fmt.Printf("getting all values [%v]\n ", storage)
	return storage
}

func (memStorage) GetValue(key string) MemStruct {
	//fmt.Printf("getting value [%v]\n ", storage[key])
	return storage[key]
}

func GetInstance() memStorage {
	if storage == nil {
		storage = make(memStorage)
	}
	return storage
}

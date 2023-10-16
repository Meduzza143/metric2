package agent

type MemStruct struct {
	metricType string
	value      string
}
type MemStorage map[string]MemStruct

func NewStorage() (storage MemStorage) {
	storage = make(map[string]MemStruct)
	return
}

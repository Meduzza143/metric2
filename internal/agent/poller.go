package agent

import (
	"math/rand"
	"runtime"
	"strconv"
)

func (storage MemStorage) Poll() {
	var mem runtime.MemStats

	runtime.ReadMemStats(&mem)
	storage["Alloc"] = MemStruct{"gauge", uintToStr(mem.Alloc)}
	storage["BuckHashSys"] = MemStruct{"gauge", uintToStr(mem.BuckHashSys)}
	storage["Frees"] = MemStruct{"gauge", uintToStr(mem.Frees)}
	storage["GCSys"] = MemStruct{"gauge", uintToStr(mem.GCSys)}
	storage["HeapAlloc"] = MemStruct{"gauge", uintToStr(mem.HeapAlloc)}
	storage["HeapIdle"] = MemStruct{"gauge", uintToStr(mem.HeapIdle)}
	storage["HeapInuse"] = MemStruct{"gauge", uintToStr(mem.HeapInuse)}
	storage["HeapObjects"] = MemStruct{"gauge", uintToStr(mem.HeapObjects)}
	storage["HeapReleased"] = MemStruct{"gauge", uintToStr(mem.HeapReleased)}
	storage["HeapSys"] = MemStruct{"gauge", uintToStr(mem.HeapSys)}
	storage["LastGC"] = MemStruct{"gauge", uintToStr(mem.LastGC)}
	storage["Lookups"] = MemStruct{"gauge", uintToStr(mem.Lookups)}
	storage["MCacheInuse"] = MemStruct{"gauge", uintToStr(mem.MCacheInuse)}
	storage["MCacheSys"] = MemStruct{"gauge", uintToStr(mem.MCacheSys)}
	storage["MSpanInuse"] = MemStruct{"gauge", uintToStr(mem.MSpanInuse)}

	storage["MSpanSys"] = MemStruct{"gauge", uintToStr(mem.MSpanSys)}
	storage["Mallocs"] = MemStruct{"gauge", uintToStr(mem.Mallocs)}
	storage["NextGC"] = MemStruct{"gauge", uintToStr(mem.NextGC)}
	storage["OtherSys"] = MemStruct{"gauge", uintToStr(mem.OtherSys)}
	storage["PauseTotalNs"] = MemStruct{"gauge", uintToStr(mem.PauseTotalNs)}
	storage["StackInuse"] = MemStruct{"gauge", uintToStr(mem.StackInuse)}
	storage["StackSys"] = MemStruct{"gauge", uintToStr(mem.StackSys)}
	storage["Sys"] = MemStruct{"gauge", uintToStr(mem.Sys)}
	storage["TotalAlloc"] = MemStruct{"gauge", uintToStr(mem.TotalAlloc)}

	storage["NumForcedGC"] = MemStruct{"gauge", uint32ToStr(mem.NumForcedGC)}
	storage["NumGC"] = MemStruct{"gauge", uint32ToStr(mem.NumGC)}

	storage["GCCPUFraction"] = MemStruct{"gauge", strconv.FormatFloat(mem.GCCPUFraction, 'f', -1, 64)}
	storage["RandomValue"] = MemStruct{"gauge", strconv.FormatFloat(rand.Float64(), 'f', -1, 64)}
	currCounter, _ := strconv.ParseUint(storage["PollCount"].value, 10, 64)

	storage["PollCount"] = MemStruct{"counter", uintToStr(currCounter)}
}

func uintToStr(value uint64) string {
	return strconv.FormatUint(value, 10)
}

func uint32ToStr(value uint32) string {
	return strconv.FormatUint(uint64(value), 10)
}

package poller

import (
	"math/rand"
	"runtime"

	"github.com/Meduzza143/metric/internal/agent/data"
)

func Poll() {
	storage := data.GetInstance()
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	storage.UpdateMetric("Alloc", "gauge", mem.Alloc)
	storage.UpdateMetric("BuckHashSys", "gauge", mem.BuckHashSys)
	storage.UpdateMetric("Frees", "gauge", mem.Frees)
	storage.UpdateMetric("GCSys", "gauge", mem.GCSys)
	storage.UpdateMetric("HeapAlloc", "gauge", mem.HeapAlloc)
	storage.UpdateMetric("HeapIdle", "gauge", mem.HeapIdle)
	storage.UpdateMetric("HeapInuse", "gauge", mem.HeapInuse)
	storage.UpdateMetric("HeapObjects", "gauge", mem.HeapObjects)
	storage.UpdateMetric("HeapReleased", "gauge", mem.HeapReleased)
	storage.UpdateMetric("HeapSys", "gauge", mem.HeapSys)
	storage.UpdateMetric("LastGC", "gauge", mem.LastGC)
	storage.UpdateMetric("Lookups", "gauge", mem.Lookups)
	storage.UpdateMetric("MCacheInuse", "gauge", mem.MCacheInuse)
	storage.UpdateMetric("MCacheSys", "gauge", mem.MCacheSys)
	storage.UpdateMetric("MSpanInuse", "gauge", mem.MSpanInuse)
	storage.UpdateMetric("MSpanSys", "gauge", mem.MSpanSys)
	storage.UpdateMetric("Mallocs", "gauge", mem.Mallocs)
	storage.UpdateMetric("NextGC", "gauge", mem.NextGC)
	storage.UpdateMetric("OtherSys", "gauge", mem.OtherSys)
	storage.UpdateMetric("PauseTotalNs", "gauge", mem.PauseTotalNs)
	storage.UpdateMetric("StackInuse", "gauge", mem.StackInuse)
	storage.UpdateMetric("StackSys", "gauge", mem.StackSys)
	storage.UpdateMetric("Sys", "gauge", mem.Sys)
	storage.UpdateMetric("TotalAlloc", "gauge", mem.TotalAlloc)
	storage.UpdateMetric("NumForcedGC", "gauge", mem.NumForcedGC)
	storage.UpdateMetric("NumGC", "gauge", mem.NumGC)
	storage.UpdateMetric("GCCPUFraction", "gauge", mem.GCCPUFraction)

	storage.UpdateMetric("RandomValue", "gauge", rand.Float64())
	currCounter := storage["PollCount"].CounterValue
	storage.UpdateMetric("PollCount", "counter", currCounter+1)
}

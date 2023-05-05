package appStat

import (
	"fmt"
	"github.com/hako/durafmt"
	"runtime"
	"time"
)

const (
	Version        = "1.0.0.412" //grepVersion
	DateTimeFormat = "2006-01-02 15:04:05 -0700"
)

var startTime = time.Now()
var startTimeString = time.Now().Format(DateTimeFormat)

type AppStatus struct {
	MemoryUsage  AppMemoryUsage
	NumGoroutine int
	NumCPU       int
	NumCgoCall   int64
	GoVersion    string
	Version      string
	Server       serverInfo
	Special      interface{}
}

type AppMemoryUsage struct {
	Alloc        string `json:"Alloc"`
	TotalAlloc   string `json:"TotalAlloc"`
	HeapAlloc    string `json:"HeapAlloc"`
	HeapReleased string `json:"HeapReleased"`
	Sys          string `json:"Sys"`
	NumGC        string `json:"NumGC"`
	LastGC       string `json:"LastGC"`
}

type serverInfo struct {
	Goos        string `json:"goos"`
	ServerTime  string `json:"serverTime"`
	ServerStart string `json:"serverStart"`
	Uptime      string `json:"serverUptime"`
	MainChat    string `json:"MainChat"`
	LogChat     string `json:"LogChat"`
	Pipeline    string `json:"pipeline"`
	LeadStatus  string `json:"leadStatus"`
	LeadStatusB string `json:"leadStatusBot"`
}

func Info() AppStatus {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	lastGC := time.Unix(0, int64(m.LastGC)).Format(DateTimeFormat)
	return AppStatus{
		MemoryUsage: AppMemoryUsage{
			Alloc:      fmt.Sprintf("%v MiB", bToMb(m.Alloc)),
			TotalAlloc: fmt.Sprintf("%v MiB", bToMb(m.TotalAlloc)),
			Sys:        fmt.Sprintf("%v MiB", bToMb(m.Sys)),
			HeapAlloc:  fmt.Sprintf("%v MiB", bToMb(m.HeapAlloc)),
			NumGC:      fmt.Sprintf("%v", m.NumGC),
			LastGC:     lastGC,
		},
		NumGoroutine: runtime.NumGoroutine(),
		NumCPU:       runtime.NumCPU(),
		NumCgoCall:   runtime.NumCgoCall(),
		GoVersion:    runtime.Version(),
		Version:      Version,
		Server: serverInfo{
			Goos:        runtime.GOOS,
			ServerTime:  time.Now().Format(DateTimeFormat),
			ServerStart: startTimeString,
			Uptime:      durafmt.Parse(time.Now().Sub(startTime)).LimitFirstN(2).String(),
		},
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

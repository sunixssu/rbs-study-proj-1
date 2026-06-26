package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"study/errs"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type JsonTelemetry struct {
	Cpu_Freq int64 `json:"cpu_freq"`
	Ram      int64 `json:"ram"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func HandleSendEcho(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(errs.ErrorWebSocketOpen + err.Error())
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	if r.Method == "OPTIONS" {
		//w.WriteHeader(http.StatusOK)
		return
	}

	for {
		messageType, body_byte, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(errs.ErrorReadingMessage, err)
			return
		}

		if err := conn.WriteMessage(messageType, body_byte); err != nil {
			fmt.Println(errs.ErrorWritingToWebSocketBody, err)
			return
		}
	}
}

func HandleSendTelemetry(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(errs.ErrorWebSocketOpen + err.Error())
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	for {
		v, err := mem.VirtualMemory()
		if err != nil {
			fmt.Println(errs.ErrorMemoryData, err)
			return
		}
		cpuInfo, err := cpu.Info()
		if err != nil {
			fmt.Println(errs.ErrorCPUData, err)
			return
		}
		cpu_mhz_int := int64(cpuInfo[0].Mhz)
		memtotal_int := int64(math.Ceil(float64(v.Total) / 1024 / 1024 / 1024))

		jsonTelemetry := JsonTelemetry{
			Cpu_Freq: cpu_mhz_int,
			Ram:      memtotal_int,
		}

		byte_arr, err := json.Marshal(jsonTelemetry)
		if err != nil {
			fmt.Println(errs.ErrorCreatingJsonFromStruct, err)
			return
		}

		if err := conn.WriteMessage(1, byte_arr); err != nil {
			fmt.Println(errs.ErrorWritingToWebSocketBody)
			return
		}
		time.Sleep(1 * time.Second)
	}
}

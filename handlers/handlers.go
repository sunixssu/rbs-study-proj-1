package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"study/errs"
	"time"

	"github.com/gorilla/websocket"
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
		//data, err := os.ReadFile("/proc/cpuinfo")
		data, err := os.ReadFile("/hostproc/cpuinfo")
		if err != nil {
			fmt.Println("Error while trying to read data about CPU:", err)
			return
		}

		data_string := string(data)
		// Делаем срез по строке, чтобы получить только число
		cache_size_idx := strings.Index(data_string, "cache size")
		cpu_mhz_idx := strings.Index(data_string, "cpu MHz")
		cpu_mhz_string := strings.Join(strings.Fields(data_string[cpu_mhz_idx+10:cache_size_idx-1]), " ")
		cpu_mhz_float, err := strconv.ParseFloat(cpu_mhz_string, 3)
		if err != nil {
			fmt.Println("Error while trying to convert string to float:", err)
			return
		}

		//data, err = os.ReadFile("/proc/meminfo")
		data, err = os.ReadFile("/hostproc/meminfo")
		if err != nil {
			fmt.Println("Error while trying to read data about RAM:", err)
			return
		}
		data_string = string(data)
		memfree_idx := strings.Index(data_string, "MemFree")
		memtotal_string := strings.Join(strings.Fields(data_string), " ")[10 : memfree_idx-11]
		memtotal_float, err := strconv.ParseFloat(memtotal_string, 3)
		if err != nil {
			fmt.Println("Error while trying to convert string to float:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cpu_mhz_int := int64(cpu_mhz_float)
		memtotal_int := int64(math.Ceil(memtotal_float / 1024 / 1024))

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

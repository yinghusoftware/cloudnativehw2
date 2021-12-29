package main

import (
	"encoding/json"
	"html"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		w.Header().Set("content-type", "application/json")
		for key, value := range r.Header {
			w.Header().Add(html.EscapeString(key), html.EscapeString(strings.Join(value, " ")))
		}
		w.Header().Add("VERSION", os.Getenv("VERSION"))
		w.WriteHeader(status)
		resp := make(map[string]string)
		resp["message"] = "success"
		jsonResp, _ := json.Marshal(resp)
		w.Write(jsonResp)
		logIPAndStatus(r, status)
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		w.WriteHeader(status)
		resp := make(map[string]string)
		resp["message"] = "success"
		jsonResp, _ := json.Marshal(resp)
		w.Write(jsonResp)
		logIPAndStatus(r, status)
	})

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	log.Fatal(http.ListenAndServe(":8081", nil))

}

func logIPAndStatus(r *http.Request, status int) {
	log.Printf("Client ip %s, return status %d\n", getIPAddress(r), status)
}

func getIPAddress(r *http.Request) string {
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() || realIP.IsPrivate() {
				continue
			}
			return ip
		}
	}
	return r.RemoteAddr
}

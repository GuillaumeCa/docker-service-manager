package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func newServer() {
	hub := newHub()

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		auth := cfg.Auth
		var postAuth struct {
			User     string
			Password string
		}
		err := json.NewDecoder(r.Body).Decode(&postAuth)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		if postAuth.User == auth.User && postAuth.Password == auth.Password {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
	})

	http.HandleFunc("/service", handleService)
	http.HandleFunc("/service/update", handleServiceUpdate)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("error upgrading connection: %v\n", err)
			return
		}
		client := newClient(hub, conn)
		go waitForMessage(client)
		hub.register <- client
	})

	go hub.run()

	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("Server started at localhost:8080...")
	log.Println(srv.ListenAndServe())
}

func handleService(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		m := newManager()
		srv, err := m.getServices()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resJSON(w, srv)
		return
	}
}

type messageWS struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func waitForMessage(client *Client) {
	var msg messageWS
	for {
		err := client.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("could not decode message: %v\n", err)
		}
		if msg.Type == "SERVICE_LOG" {
			handleLog(client, msg.Data)
		}
	}
}

type logLine struct {
	Content string
}

func handleLog(client *Client, id string) {
	m := newManager()

	rclogs, err := m.getLogs(id)
	if err != nil {
		log.Printf("error reading logs: %v\n", err)
		return
	}

	reader := bufio.NewReader(rclogs)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Read Error:", err)
			return
		}
		client.sendJSON(&logLine{
			Content: str[8:],
		})
	}
}

func handleServiceLogs(w http.ResponseWriter, r *http.Request, hub *Hub) {

	m := newManager()
	v, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Printf("error parsing query: %v\n", err)
		return
	}

	ID := v.Get("id")

	rclogs, err := m.getLogs(ID)
	if err != nil {
		log.Printf("error reading logs: %v\n", err)
		return
	}

	reader := bufio.NewReader(rclogs)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Read Error:", err)
			return
		}
		hub.broadcastJSON(logLine{
			Content: str,
		})
	}

}

func handleServiceUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		resStatus(w, http.StatusOK)
		return
	}
	if r.Method == http.MethodPut {
		m := newManager()
		m.updateServices(cfg)
		resStatus(w, http.StatusOK)
		return
	}
	resStatus(w, http.StatusBadRequest)
}

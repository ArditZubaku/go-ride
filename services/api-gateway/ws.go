package main

import (
	"log"
	"net/http"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/util"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleRidersWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WS upgrade failed: %v\n", err)
		return
	}
	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if !(len(userID) > 0) {
		log.Println("No user ID provided")
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from WS: %v\n", err)
			break
		}

		log.Printf("Received message: %s", msg)
	}
}

func handleDriversWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WS upgrade failed: %v\n", err)
		return
	}
	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if !(len(userID) > 0) {
		log.Println("No user ID provided")
		return
	}

	pkgSlug := r.URL.Query().Get("packageSlug")
	if !(len(pkgSlug) > 0) {
		log.Println("No package slug provided")
		return
	}

	type Driver struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		ProfilePicture string `json:"profilePicture"`
		CarPlate       string `json:"carPlate"`
		PackageSlug    string `json:"packageSlug"`
	}

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: Driver{
			ID:             userID,
			Name:           "TestUser",
			ProfilePicture: util.GetRandomAvatar(1),
			CarPlate:       "ABC123",
			PackageSlug:    pkgSlug,
		},
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("Error sending message: %v\n", err)
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from WS: %v\n", err)
			break
		}

		log.Printf("Received message: %s", msg)
	}
}

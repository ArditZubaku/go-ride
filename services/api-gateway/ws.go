package main

import (
	"log"
	"net/http"

	"ride-sharing/services/api-gateway/grpcclients"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/proto/driver"
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

	defer util.CloseAndLog(conn, "handleRidersWS")

	userID := r.URL.Query().Get("userID")
	if len(userID) <= 0 {
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

	defer util.CloseAndLog(conn, "handleDriversWS")

	userID := r.URL.Query().Get("userID")
	if len(userID) <= 0 {
		log.Println("No user ID provided")
		return
	}

	pkgSlug := r.URL.Query().Get("packageSlug")
	if len(pkgSlug) <= 0 {
		log.Println("No package slug provided")
		return
	}

	ctx := r.Context()

	driverService, err := grpcclients.NewDriverServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	// Closing connections
	defer func() {
		_, err := driverService.Client.UnregisterDriver(
			ctx,
			&driver.RegisterDriverRequest{
				DriverID:    userID,
				PackageSlug: pkgSlug,
			},
		)
		if err != nil {
			// TODO: Handle this better
			log.Println("Faied to unregister driver: ", userID)
		}

		driverService.Close()

		log.Println("Driver unregistered: ", userID)
	}()

	driverData, err := driverService.Client.RegisterDriver(
		ctx,
		&driver.RegisterDriverRequest{
			DriverID:    userID,
			PackageSlug: pkgSlug,
		},
	)
	if err != nil {
		log.Printf("Error registering driver: %v", err)
		return
	}

	msg := contracts.WSMessage[*driver.Driver]{
		Type: "driver.cmd.register",
		Data: driverData.Driver,
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

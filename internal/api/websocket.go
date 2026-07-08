package api

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"tailscope/internal/model"
)

var upgrader = websocket.Upgrader{

	CheckOrigin: func(
		r *http.Request,
	) bool {

		return true
	},
}

func (s *Server) streamLogs(
	w http.ResponseWriter,
	r *http.Request,
) {

	log.Println("client connected")

	conn, err := upgrader.Upgrade(
		w,
		r,
		nil,
	)

	if err != nil {
		log.Println("upgrade failed:", err)
		return
	}

	defer conn.Close()

	ch, unsubscribe := s.store.Subscribe()

	defer func() {

		unsubscribe()

		s.app.StopFollowing()

	}()

	log.Println("subscribed")

	for {

		entry, ok := <-ch

		if !ok {

			log.Println("channel closed")

			return
		}

		err := sendLog(
			conn,
			entry,
		)

		if err != nil {

			log.Println(
				"write failed:",
				err,
			)

			return
		}
	}
}

func sendLog(
	conn *websocket.Conn,
	entry model.LogEntry,
) error {

	return conn.WriteJSON(
		entry,
	)
}

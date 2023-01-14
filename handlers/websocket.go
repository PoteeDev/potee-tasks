package handlers

import (
	"console/shell"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (h *handler) WsEndpoint(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "Session")
	var name string
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		log.Println("not authenticated")
	} else {
		name = session.Values["name"].(string)
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Println(name, string(message))
		s := shell.InitShell(h.DB, name)
		answer := s.Execute(message)
		err = c.WriteMessage(mt, answer)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

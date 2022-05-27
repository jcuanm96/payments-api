package service

import (
	"context"
	"log"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/vamawebsocket"
)

func (svc *usecase) ReactWebsocket(ctx context.Context, conn *vamawebsocket.Conn, req request.ReactWebsocket) error {
	// conn.Locals is added to the *websocket.Conn
	log.Println(conn.Locals("allowed"))  // true
	log.Println(conn.Params("id"))       // 123
	log.Println(conn.Query("v"))         // 1.0
	log.Println(conn.Cookies("session")) // ""

	// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
	var (
		mt  int
		msg []byte
		err error
	)
	for {
		if mt, msg, err = conn.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", msg)

		if err = conn.WriteMessage(mt, msg); err != nil {
			log.Println("write:", err)
			break
		}
	}
	return nil
}

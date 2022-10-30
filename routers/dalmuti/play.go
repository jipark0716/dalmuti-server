package dalmuti

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jipark0716/dalmuti/service/dalmuti"
)

func Play(context *gin.Context) {
	upgrader := websocket.Upgrader{}
	ws, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}

	id := context.Param("id")
	game := dalmuti.GetGame(id)
	game.Join(ws)
}

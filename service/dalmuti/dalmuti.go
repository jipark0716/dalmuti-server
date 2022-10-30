package dalmuti

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

type Game struct {
	Id        string
	Players   []Player
	Status    GameStatus
	FirstGame bool
}

type Player struct {
	WebsocketConnection *websocket.Conn
	GameId              string
	Cards               []Card
}

type Card struct {
	Value    int
	WildCard bool
}

type Request struct {
	Type PlayerActionType `json:"type"`
}

type GameStatus int

const (
	GameStatusReady   GameStatus = 0
	GameStatusPlaying GameStatus = 1
)

type PlayerActionType string

const (
	PlayerActionStart PlayerActionType = "start"
)

var Games map[string]*Game

func init() {
	Games = make(map[string]*Game)
	rand.Seed(time.Now().UnixNano())
}

func GetGame(id string) *Game {
	if game, ok := Games[id]; ok {
		return game
	}

	Games[id] = NewGame(id)
	return Games[id]
}

func NewGame(id string) *Game {
	return &Game{
		FirstGame: true,
		Id:        id,
		Status:    GameStatusReady,
	}
}

func (g *Game) Join(ws *websocket.Conn) {
	player := NewPlayer(g.Id, ws)
	g.Players = append(g.Players, player)
	player.Open()
}

func (p *Player) Open() {
	for {
		//Read Message from client
		messageType, message, err := p.WebsocketConnection.ReadMessage()
		if err != nil {
			break
		}

		if messageType != websocket.TextMessage {
			continue
		}

		request := Request{}
		err = json.Unmarshal(message, &request)
		if err != nil {
			continue
		}

		p.Do(request)
	}
	p.Close()
}

func (p *Player) Do(request Request) {
	switch request.Type {
	case PlayerActionStart:
		err := GetGame(p.GameId).Start()
		if err != nil {
			p.WebsocketConnection.WriteMessage(
				websocket.TextMessage,
				[]byte(fmt.Sprintf("%+v", err)),
			)
		}
	}
}

func NewPlayer(gameId string, ws *websocket.Conn) Player {
	return Player{
		GameId:              gameId,
		WebsocketConnection: ws,
	}
}

func (g *Game) Close() {
	for _, player := range g.Players {
		player.Close()
	}

	delete(Games, g.Id)
}

func (p *Player) Close() {
	game := GetGame(p.GameId)
	for id, player := range game.Players {
		if player.WebsocketConnection == p.WebsocketConnection {
			copy(game.Players[id:], game.Players[id+1:])
			game.Players[len(game.Players)-1] = Player{}
			game.Players = game.Players[:len(game.Players)-1]
			break
		}
	}

	p.WebsocketConnection.Close()
}

func (p *Player) ClearCard() {
	p.Cards = []Card{}
}

func (g *Game) Start() error {
	if len(g.Players) < 3 {
		return fmt.Errorf("인원 부족")
	}

	if g.Status != GameStatusReady {
		return fmt.Errorf("시작할 수 있는 상태 아님")
	}

	g.Status = GameStatusPlaying

	for _, player := range g.Players {
		player.ClearCard()
	}

	// 첫게임이면 순서 랜덤
	if g.FirstGame {
		rand.Shuffle(len(g.Players), func(i, j int) {
			g.Players[i], g.Players[j] = g.Players[j], g.Players[i]
		})
	}

	for i, card := range g.GetCards() {
		player := g.Players[i%len(g.Players)]
		player.Cards = append(player.Cards, card)
	}

	// g.RespondGameStart()

	return nil
}

// func (g *Game) RespondGameStart() {
// 	for _, player := range g.Players {

// 	}
// }

func (g *Game) GetCards() []Card {
	joker := Card{
		Value:    13,
		WildCard: true,
	}
	cards := []Card{joker, joker}
	for i := 1; i <= 12; i++ {
		for j := 0; j < i; j++ {
			cards = append(cards, Card{
				Value:    i,
				WildCard: false,
			})
		}
	}

	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	return cards
}

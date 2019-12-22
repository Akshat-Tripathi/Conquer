package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

//Based off https://github.com/scotch-io/go-realtime-chat/blob/master/src/main.go

type hub struct {
	users map[string]*user
}

type user struct {
	name         string
	clients      map[*websocket.Conn]bool
	outputStream chan action
}

func (u *user) init(path string) {
	http.HandleFunc(path, u.handler)
	go u.readSend()
}

func (u *user) handler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	defer ws.Close()
	fmt.Println(time.Now(), ws.RemoteAddr().String()+" has connected as: "+u.name)
	u.clients[ws] = true
	for {
		var act action
		err := ws.ReadJSON(&act)
		if err != nil {
			delete(u.clients, ws)
			fmt.Println(time.Now(), ws.RemoteAddr().String()+" has left the game as: "+u.name)
			break
		}
		fmt.Println(act)
		if !validate(act) {
			ws.WriteJSON(action{
				Player:   "admin",
				Src:      "SS",
				Dest:     "SS",
				MoveType: -1,
				Numsrc:   -1,
				Numdest:  -1,
			})
		} else {
			//Insert save code here
		}
	}
}

func (u *user) readSend() {
	for {
		if len(u.clients) > 0 {
			select {
			case act, ready := <-u.outputStream:
				if ready {
					for client := range u.clients {
						err := client.WriteJSON(act)
						if err != nil {
							client.Close()
							delete(u.clients, client)
						}
					}
				}
			default:
			}
		}
		time.Sleep(time.Millisecond * 10) //Without this line the program eats the cpu
	}
}

func (h *hub) send(act action) {
	for _, u := range h.users {
		go func(u *user) {
			u.outputStream <- act
		}(u)
	}
}

func (h *hub) init() {
	for i := 0; i < players; i++ {
		h.users[accounts[i]] = &user{name: accounts[i], clients: make(map[*websocket.Conn]bool), outputStream: make(chan action)}
		go h.users[accounts[i]].init("/ws/" + accounts[i])
	}
}

func (h *hub) sendToPlayer(player string, act action) {
	go func() {
		h.users[player].outputStream <- act
	}()
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	timing   = 17
	duration = 10
)

var (
	players  = 8
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	countryMap             = loadCountries()
	accounts, hidden, pots = getAccounts()
	sender                 = hub{make(map[string]*user)}
	winner                 = ""
)

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func main() {
	fmt.Println(getOutboundIP().String())
	start := flag.Bool("start", false, "Is this the start of the game?")
	numPlayer := flag.Int("num", players, "How many players are playing?")
	numCountries := flag.Int("countries", 4, "How many countries does everyone get?")
	rapid := flag.Bool("rapid", false, "Enable rapid fire mode?")
	flag.Parse()
	fmt.Println(*start, *numPlayer, *numCountries, *rapid)
	if *start {
		allocateCountries(*numPlayer, *numCountries)
	}
	saveCountries(countryMap)
	sender.init()
	go housekeeping()
	go addPot(*start, *rapid)

	http.Handle("/static/css/", http.StripPrefix("/static/css/", http.FileServer(http.Dir("css/"))))
	http.Handle("/static/js/", http.StripPrefix("/static/js/", http.FileServer(http.Dir("js/"))))
	http.Handle("/static/initial/", http.StripPrefix("/static/initial/", http.FileServer(http.Dir("data/"))))

	http.HandleFunc("/", game)
	http.HandleFunc("/login/", authenticator)
	http.HandleFunc("/winner", win)
	http.HandleFunc("/lose", lose)
	http.HandleFunc("/wrong_login", badLogin)
	//log.Fatal(http.ListenAndServe(getOutboundIP().String()+":8080", nil))
	log.Fatal(http.ListenAndServe("192.168.1.16:80", nil))
}

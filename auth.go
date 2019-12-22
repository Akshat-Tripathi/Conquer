package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func setCookie(w http.ResponseWriter, name, value string) {
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Path:    "/",
		Expires: time.Now().Add(365 * 24 * time.Hour),
	}
	http.SetCookie(w, &cookie)
}

func authenticator(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:]
	code = strings.Replace(code, "login/", "", -1)
	for _, v := range accounts {
		if code == v {
			setCookie(w, "id", code)
			setCookie(w, "col", hidden[code])
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}
	http.Redirect(w, r, "/wrong_login", http.StatusFound)
}

//Gets the accounts, with their statuses as 0 and then converts them to their colours
//The colours are for the browser
func getAccounts() ([]string, map[string]string, map[string]int) {
	raw, err := ioutil.ReadFile("accounts.txt")
	if err != nil {
		log.Println(err)
	}
	raw2 := strings.Split(string(raw), "\r\n")
	state := make(map[string]int)
	colour := make(map[string]string)
	pots := make(map[string]int)
	accounts := make([]string, len(raw2))
	colour["nil"] = "#d2d7d3"
	for k, v := range raw2 {
		vals := strings.Split(v, ";")
		name, col, pot := vals[0], vals[1], vals[2]
		colour[name] = col
		state[name] = 0
		accounts[k] = name
		p, err := strconv.Atoi(pot)
		if err != nil {
			panic(err)
		}
		pots[name] = p
	}
	return accounts, colour, pots
}

func saveAccounts() {
	all := ""
	for {
		all = ""
		for k, v := range pots {
			all += k + ";" + hidden[k] + ";" + strconv.Itoa(v) + "\r\n"
		}
		if len(all) > 10 {
			break
		}
		fmt.Println("Sleeping", all)
		time.Sleep(time.Second)
	}
	err := ioutil.WriteFile("accounts.txt", []byte(all)[:len(all)-2], 0644)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("backup_accounts.txt", []byte(all)[:len(all)-2], 0644)
	if err != nil {
		panic(err)
	}
}

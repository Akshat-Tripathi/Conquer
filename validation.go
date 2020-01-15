package main

import (
	"fmt"
	"math/rand"
)

//Everything must be exported to send with json
type action struct {
	Player   string
	Src      string
	Dest     string
	MoveType int
	Numsrc   int
	Numdest  int
}

func (a *action) attack() {
	if countryMap[a.Dest].population == 0 {
		countryMap[a.Dest].owner = a.Player
		q := action{
			Player:   countryMap[a.Src].owner,
			Src:      a.Src,
			Dest:     a.Dest,
			MoveType: 3,
			Numsrc:   -1,
			Numdest:  1,
		}
		q.move()
		return
	}
	at := 3
	if countryMap[a.Src].population == 2 {
		at = 2
	}
	de := 2
	if countryMap[a.Dest].population == 1 {
		de = 1
	}

	attac := make([]int, at)
	defend := make([]int, de)
	for i := 0; i < at; i++ {
		attac[i] = rand.Intn(6)
	}
	for i := 0; i < de; i++ {
		defend[i] = rand.Intn(6)
	}

	attac = sort(attac)
	defend = sort(defend)
	fmt.Println(attac, defend)
	for k, v := range defend {
		if attac[k] <= v {
			a.Numsrc--
		} else {
			a.Numdest--
		}
	}

	countryMap[a.Src].population += a.Numsrc
	countryMap[a.Dest].population += a.Numdest
	fmt.Println(a.Numdest, a.Numsrc)
	//fmt.Println("attacking", a.Numdest, countryMap[a.Dest].population)
	if countryMap[a.Dest].population < 0 {
		a.Numdest -= countryMap[a.Dest].population
		countryMap[a.Dest].population = 0
		//fmt.Println("conquering", a.Numdest, countryMap[a.Dest].population)
		if countCountries()[a.Player] >= 35 {
			winner = a.Player
			fmt.Println(a.Player + " wins!")
		}
		sender.send(*a)
		a.attack()
		return
	}
	sender.send(*a)
	saveCountries(countryMap)
}

func (a *action) move() {
	if a.Src == "PO" {
		pots[a.Player] += a.Numsrc
		saveAccounts()
	} else {
		countryMap[a.Src].population += a.Numsrc
	}
	a.Player = hidden[a.Player]
	countryMap[a.Dest].population += a.Numdest
	sender.send(*a)
	saveCountries(countryMap)
}

func (a *action) donate() {
	src := a.Player
	dest := countryMap[a.Dest].owner
	pots[src] += a.Numsrc
	pots[dest] += a.Numdest
	sendPots(src)
	sendPots(dest)
	saveAccounts()
}

func validate(a action) bool {
	if a.MoveType < 0 || a.MoveType > 2 {
		fmt.Println("invalid movetype")
		return false
	}
	if a.Src != "PO" {
		if countryMap[a.Src].owner != a.Player {
			return false
		}
		if countryMap[a.Src].population < 1 {
			return false
		}
	}
	if a.MoveType == 0 {
		if canAttack(a) {
			a.attack()
			return true
		}
		return false
	} else if a.MoveType == 1 {
		if canDonate(a) {
			a.donate()
			return true
		}
		return false
	} else {
		if canMove(a) {
			a.move()
			return true
		}
		return false
	}
}

func canAttack(a action) bool {
	if a.Numsrc != 0 || a.Numdest != 0 {
		return false
	}
	if a.Src == a.Dest {
		return false
	}
	if countryMap[a.Src].owner == countryMap[a.Dest].owner {
		return false
	}
	if countryMap[a.Src].population <= 1 {
		return false
	}
	return in(countryMap[a.Src].neighbours, a.Dest)
}

func canDonate(a action) bool {
	if a.Numsrc >= 0 {
		fmt.Println("Numsrc", a.Numsrc)
		return false
	}
	if a.Numsrc*-1 != a.Numdest {
		fmt.Println("negative")
		return false
	}
	if a.Src != "PO" {
		if a.Player == countryMap[a.Dest].owner {
			fmt.Println("same dudes")
			return false
		}
	}
	if a.Src == a.Dest {
		fmt.Println("same dudes 2")
		return false
	}
	if a.Numdest > pots[a.Player] {
		fmt.Println("too much", a.Numdest, pots[a.Player])
		return false
	}
	if countryMap[a.Dest].owner == "nil" {
		fmt.Println("no owner")
		return false
	}
	return true
}

func canMove(a action) bool {
	if a.Numsrc*-1 != a.Numdest {
		fmt.Println("Invalid numsrc/numdest")
		return false
	}
	if a.Src == a.Dest {
		fmt.Println("Same src/dest")
		return false
	}
	if countryMap[a.Dest].owner != a.Player {
		if a.Src == "PO" {
			return !(a.Numdest > pots[a.Player])
		}
		fmt.Println("Player does not own territory")
		return false
	}
	if a.Src != "PO" {
		if a.Numdest >= countryMap[a.Src].population {
			fmt.Println("Moving too many troops")
			return false
		}
		return in(countryMap[a.Src].neighbours, a.Dest)
	}
	return !(a.Numdest > pots[a.Player])
}

func in(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func sort(arr []int) []int {
	if len(arr) == 1 {
		return arr
	}

	max := arr[0]
	if len(arr) == 2 {
		if max > arr[1] {
			return []int{max, arr[1]}
		}
		return []int{arr[1], max}
	}

	ret := make([]int, 1)
	dex := 0
	for i, v := range arr[1:] {
		if max < v {
			max = v
			dex = i + 1
		}
	}
	ret[0] = max
	ret = append(ret, sort(append(arr[:dex], arr[dex+1:]...))...)
	return ret
}

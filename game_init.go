package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	bonus = map[string]int{
		"Africa":        3,
		"Asia":          7,
		"North America": 5,
		"South America": 2,
		"Australia":     2,
		"Europe":        5,
	}
	continents = map[string][]string{
		"Africa":        []string{"NA", "WA", "CN", "MA", "EG", "SA"},
		"Asia":          []string{"IN", "MI", "UR", "SI", "NR", "PL", "CH", "JP", "AF", "SE", "PA"},
		"Europe":        []string{"UK", "RO", "SC", "PR", "RU"},
		"Australia":     []string{"WO", "EO", "CO", "NG"},
		"North America": []string{"AL", "CA", "SF", "NY", "ME", "GR"},
		"South America": []string{"VE", "PE", "AG", "BR"},
	}
)

func countCountries() map[string]int {
	countryCounts := make(map[string]int)
	cc := make([]int, 7)
	for i, v := range accounts {
		countryCounts[v] = i
	}
	for _, v := range countryMap {
		if v.owner != "nil" && v.population > 0 {
			cc[countryCounts[v.owner]]++
		}
	}
	//fmt.Println(cc, countryCounts)
	for k, v := range countryCounts {
		countryCounts[k] = cc[v]
	}
	return countryCounts
}

func countContinents() map[string]int {
	counts := make(map[string]int)
	cc := make([]int, 7)
	for i, v := range accounts {
		counts[v] = i
	}
	for k, v := range continents {
		u := countryMap[v[0]].owner
		if u != "nil" {
			has := true
			for _, i := range v[1:] {
				if countryMap[i].owner != u {
					has = false
					break
				}
			}
			if has {
				cc[counts[u]] += bonus[k]
			}
		}
	}
	for k, v := range counts {
		counts[k] = cc[v]
	}
	return counts
}

func calculate() map[string]int {
	cont := countContinents()
	count := countCountries()
	m := make(map[string]int)
	for _, v := range accounts {
		m[v] = 3 + cont[v] + count[v]/3
	}
	return m
}

func winCondition() {
	for _, v := range accounts {
		if countCountries()[v] >= 35 {
			winner = v
		}
	}
}

func loseCondition() {
	for _, v := range accounts {
		if countCountries()[v] < 0 {
			status[v] = -1
		}
	}
}

func sendPots(player string) {
	sender.sendToPlayer(player, action{
		Player:   player,
		Src:      "PO",
		Dest:     "PO",
		Numsrc:   1,
		Numdest:  pots[player] - 1,
		MoveType: 2,
	})
}

func housekeeping() {
	for {
		winCondition()
		loseCondition()
		time.Sleep(time.Minute)
	}
}

func addPot(start, rapid bool) {
	if start {
		for k := range pots {
			pots[k] = 10
		}
	}
	saveAccounts()
	fmt.Println("troops drop in: ", getTimeToNextPot(rapid))
	time.Sleep(getTimeToNextPot(rapid))
	for {
		additions := calculate()
		for k := range pots {
			pots[k] += additions[k]
			sendPots(k)
		}
		saveAccounts()
		time.Sleep(getTimeToNextPot(rapid))
	}
}

func getTimeToNextPot(rapid bool) time.Duration {
	if rapid {
		return delayBased()
	}
	return timeBased()
}

func timeBased() time.Duration {
	now := time.Now()
	current := now.Hour()
	if current < timing {
		return time.Until(time.Date(now.Year(), now.Month(), now.Day(), timing, 30, 0, 0, time.Local))
	}
	return time.Until(time.Date(now.Year(), now.Month(), now.Day()+1, timing, 30, 0, 0, time.Local))
}

func delayBased() time.Duration {
	return time.Second * 30
}

func allocateCountries(numPlayers, countries int) {
	allCountries := make([]string, 0)
	for k, v := range countryMap {
		if v.owner != "nil" {
			countryMap[k] = &country{
				population: 0,
				owner:      "nil",
				neighbours: v.neighbours,
			}
		}
		allCountries = append(allCountries, k)
	}
	players := accounts[:numPlayers]
	for _, player := range players {
		for i := 0; i < countries; i++ {
			r := rand.Intn(len(allCountries))
			countryMap[allCountries[r]].owner = player
			allCountries = append(allCountries[:r], allCountries[r+1:]...)
		}
		fmt.Println(player)
	}
}

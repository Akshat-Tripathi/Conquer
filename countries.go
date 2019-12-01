package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type country struct {
	population int
	owner      string
	neighbours []string
}

func loadCountries() map[string]*country {
	countries := make(map[string]*country)
	raw, err := ioutil.ReadFile("countries.txt")
	if err != nil {
		log.Println(err)
	}
	places := strings.Split(string(raw), "\n")
	for _, v := range places {
		temp := strings.Split(v, ";")
		population, err := strconv.Atoi(temp[1])
		if err != nil {
			panic(err)
		}
		for i := range temp[3:] {
			temp[i+3] = strings.Replace(temp[i+3], string([]byte{13}), "", -1)
		}
		countries[temp[0]] = &country{
			population: population,
			owner:      temp[2],
			neighbours: temp[3:],
		}
	}
	return countries
}

//Saves the contries to both the server and browser version
//The browser version contains colours instead of ids for security
func saveCountries(m map[string]*country) {
	var str string     //Where the stringified version will be saved
	var strSend string //Where the stringified version to be sent will be saved
	for k, v := range m {
		temp := []string{
			k,
			strconv.Itoa(v.population),
			v.owner,
		}
		temp1 := []string{
			k,
			strconv.Itoa(v.population),
			hidden[v.owner],
		}
		for _, neighbour := range v.neighbours {
			temp = append(temp, neighbour)
			temp1 = append(temp1, neighbour)
		}
		str1 := strings.Join(temp, ";")
		str = strings.Join([]string{str, str1}, "\n")
		strSend1 := strings.Join(temp1, ";")
		strSend = strings.Join([]string{strSend, strSend1}, "\n")
	}
	ioutil.WriteFile("countries.txt", []byte(str[1:]), 0644)
	ioutil.WriteFile("data/countries.txt", []byte(strSend[1:]), 0644)
}

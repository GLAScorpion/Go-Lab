package main

import (
	"fmt"
	"math/rand"
)

type Cliente struct {
	nome string
}

type Veicolo struct {
	tipo string
}

const (
	kBerlina      = "una Berlina"
	kSUV          = "un SUV"
	kStationWagon = "una Station Wagon"
)

var orders = make(map[Cliente]Veicolo)
var count = make(map[Veicolo]int)

func main() {
	num := 10
	for i := 0; i < num; i++ {
		noleggia(Cliente{clientGen()})
	}
	stampa()
}

func noleggia(cliente Cliente) {
	val := Veicolo{}
	switch rand.Int() % 3 {
	case 0:
		val.tipo = kBerlina
	case 1:
		val.tipo = kSUV
	case 2:
		val.tipo = kStationWagon
	}
	orders[cliente] = val
	count[val]++
	fmt.Printf("Il cliente %s ha noleggiato %s \n", cliente.nome, val.tipo)
}

func stampa() {
	for key, val := range count {
		switch key.tipo {
		case kBerlina:
			fmt.Printf("Il numero di Berline prenotate è %d\n", val)
		case kSUV:
			fmt.Printf("Il numero di SUV prenotati è %d\n", val)
		case kStationWagon:
			fmt.Printf("Il numero di Station Wagon prenotate è %d\n", val)

		}
	}
}

func clientGen() (res string) {
	nameL := (rand.Int() % 5) + 4
	for i := 0; i < nameL; i++ {
		res += string(rune((rand.Int() % 26) + 65))
	}
	return
}

package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type Cliente struct {
	nome string
}

type Veicolo struct {
	tipo string
}

type Ordine struct {
	cliente Cliente
	veicolo Veicolo
}

const (
	kBerlina      = "una Berlina"
	kSUV          = "un SUV"
	kStationWagon = "una Station Wagon"
)

var orders = make(map[Cliente]Veicolo) //mappa che lega un cliente alla sua scelta
var count = make(map[Veicolo]int)      // mappa che conta le prenotazioni dei veicoli

var ch = make(chan Ordine)
var wg = sync.WaitGroup{}

func main() {
	num := 10

	for i := 0; i < num; i++ {
		wg.Add(1)
		go noleggia(Cliente{clientGen()})
	}
	for i := 0; i < num; i++ {
		ordine := <-ch                          // riceve l'ordine
		orders[ordine.cliente] = ordine.veicolo // registra cliente e veicolo
		count[ordine.veicolo]++                 // aumenta il conteggio dei veicoli
	}
	wg.Wait()
	stampa()
}

func noleggia(cliente Cliente) {
	val := Veicolo{}
	switch rand.Int() % 3 { //sceglie un veicolo random tra tre
	case 0:
		val.tipo = kBerlina
	case 1:
		val.tipo = kSUV
	case 2:
		val.tipo = kStationWagon
	}
	ch <- Ordine{cliente, val} //invia l'ordine nel canale
	fmt.Printf("Il cliente %s ha noleggiato %s \n", cliente.nome, val.tipo)
	wg.Done()
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

func clientGen() (res string) { //genera randomicamente il nome di un cliente
	nameL := (rand.Int() % 5) + 4
	for i := 0; i < nameL; i++ {
		res += string(rune((rand.Int() % 26) + 65))
	}
	return
}

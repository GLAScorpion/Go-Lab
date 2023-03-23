package main

import (
	"fmt"
	"sync"
)

var wg = sync.WaitGroup{}
var count int = 0

func main() {
	string := "aaaaaaaaaaaaabbbbbbbbcccccddddccccccfff"
	toFind := 'c'
	ch := make(chan int, len(string)) //buffered channel per accogliere tutti gli esiti e processarli in seguito
	willClose := false
	for key, v := range string {
		wg.Add(1)
		if key == len(string)-1 { // la goroutine dell'ultimo carattere chiude il canale
			willClose = true
		}
		go checker(ch, toFind, v, willClose)
	}
	wg.Wait()
	for r := range ch { //somma gli esiti una volta terminate le routine
		count += r
	}
	fmt.Printf("Il carattere %c è stato trovato %d volte\n", toFind, count)
}

func checker(ch chan<- int, toFind rune, toCheck rune, closeCh bool) {
	if toFind == toCheck {
		ch <- 1
	} else {
		ch <- 0
	}
	wg.Done()
	if closeCh { //chiude il canale una volta che tutte le routine hanno terminato
		wg.Wait()
		close(ch)
	}
}

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
	ch := make(chan int, len(string))
	willClose := false
	for key, v := range string {
		wg.Add(1)
		if key == len(string)-1 {
			willClose = true
		}
		go checker(ch, toFind, v, willClose)
	}
	wg.Wait()
	for r := range ch {
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
	if closeCh {
		wg.Wait()
		close(ch)
	}
}

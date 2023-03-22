package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

var eur_usd = make(chan float32)
var gbp_usd = make(chan float32)
var jpy_usd = make(chan float32)
var quit = make(chan struct{})

func main() {
	go simulateMarketData()
	go quitter()
	for {
		select {
		case val := <-eur_usd:
			fmt.Println("EUR/USD:", val)
		case val := <-gbp_usd:
			fmt.Println("GBP/USD", val)
		case val := <-jpy_usd:
			fmt.Println("JPY/USD", val)
		case <-quit:
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func simulateMarketData() {
	for {
		eur_usd <- float32(int32(rand.Int31()))/float32(math.MaxInt32)*0.5 + 1.0
		gbp_usd <- float32(int32(rand.Int31()))/float32(math.MaxInt32)*0.5 + 1.0
		jpy_usd <- float32(int32(rand.Int31()))/float32(math.MaxInt32)*0.003 + 0.006
	}
}

func quitter() {
	reader := bufio.NewReader(os.Stdin)
	for {
		r, _, _ := reader.ReadRune()
		if r == 'q' {
			quit <- struct{}{}
			break
		}
	}
}

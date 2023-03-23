package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var eur_usd = make(chan float64, 5)
var gbp_usd = make(chan float64, 5)
var jpy_usd = make(chan float64, 5)

var quit = make(chan struct{})

var (
	holdEUR = 0.0
	holdGBP = 0.0
	holdJPY = 0.0
)

var (
	sellTimeEUR = time.Time{}
	buyTimeGBP  = time.Time{}
	buyTimeJPY  = time.Time{}
)

var (
	sellEUR = false
	buyGBP  = false
	buyJPY  = false
)

func main() {
	go simulateMarketData()
	go selectPair()
	time.Sleep(time.Minute)
}

func simulateMarketData() {
	for {
		if !sellEUR || (sellEUR && time.Since(sellTimeEUR).Seconds() >= 4.0) {
			eur_usd <- float64(int32(rand.Int31()))/float64(math.MaxInt32)*0.5 + 1.0
		}
		if !buyGBP || (buyGBP && time.Since(buyTimeGBP).Seconds() >= 3.0) {
			gbp_usd <- float64(int32(rand.Int31()))/float64(math.MaxInt32)*0.5 + 1.0
		}
		if !buyJPY || (buyJPY && time.Since(buyTimeJPY).Seconds() >= 3.0) {
			jpy_usd <- float64(int32(rand.Int31()))/float64(math.MaxInt32)*0.003 + 0.006
		}
		time.Sleep(time.Second)
	}
}

func selectPair() {

	for {
		select {
		case val := <-eur_usd:
			if val > 1.20 && !sellEUR {
				holdEUR = val
				sellTimeEUR = time.Now()
				sellEUR = true
			} else if sellEUR {
				sellEUR = false
				fmt.Println("Detected EUR/USD at", holdEUR, ", sold at", val, time.Since(sellTimeEUR))
			}
		case val := <-gbp_usd:
			if val < 1.35 && holdGBP == 0 {
				holdGBP = val
				buyTimeGBP = time.Now()
				buyGBP = true
			} else if buyGBP {
				buyGBP = false
				fmt.Println("Detected GBP/USD at", holdGBP, ", bought at", val, time.Since(buyTimeGBP))
			}
		case val := <-jpy_usd:
			if val < 0.0085 && holdJPY == 0 {
				holdJPY = val
				buyTimeJPY = time.Now()
				buyJPY = true
			} else if buyJPY {
				buyJPY = false
				fmt.Println("Detected JPY/USD at", holdJPY, ", bought at", val, time.Since(buyTimeJPY))
			}
		}
		//time.Sleep(500 * time.Millisecond)
	}
}

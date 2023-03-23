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

var (
	holdEUR = 0.0
	holdGBP = 0.0
	holdJPY = 0.0
)

var (
	timeEUR = time.Time{}
	timeGBP = time.Time{}
	timeJPY = time.Time{}
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
		eur_usd <- float64(int32(rand.Int31()))/float64(math.MaxInt32)*0.5 + 1.0
		gbp_usd <- float64(int32(rand.Int31()))/float64(math.MaxInt32)*0.5 + 1.0
		jpy_usd <- float64(int32(rand.Int31()))/float64(math.MaxInt32)*0.003 + 0.006
		time.Sleep(time.Second)
	}
}

func selectPair() {
	timeEUR = time.Now()
	timeGBP = time.Now()
	timeJPY = time.Now()
	for {
		select {
		case val := <-eur_usd:
			if val > 1.20 && !sellEUR {
				holdEUR = val
				timeEUR = time.Now()
				sellEUR = true
			} else if sellEUR && time.Since(timeEUR).Seconds() >= 4.0 {
				sellEUR = false
				fmt.Println("Detected EUR/USD at", holdEUR, ", sold at", val, ". It took", time.Since(timeEUR).Seconds(), "seconds")
			}
		case val := <-gbp_usd:
			if val < 1.35 && !buyGBP {
				holdGBP = val
				timeGBP = time.Now()
				buyGBP = true
			} else if buyGBP && time.Since(timeGBP).Seconds() >= 3.0 {
				buyGBP = false
				fmt.Println("Detected GBP/USD at", holdGBP, ", bought at", val, ". It took", time.Since(timeGBP).Seconds(), "seconds")
			}
		case val := <-jpy_usd:
			if val < 0.0085 && !buyJPY {
				holdJPY = val
				timeJPY = time.Now()
				buyJPY = true
			} else if buyJPY && time.Since(timeJPY).Seconds() >= 3.0 {
				buyJPY = false
				fmt.Println("Detected JPY/USD at", holdJPY, ", bought at", val, ". It took", time.Since(timeJPY).Seconds(), "seconds")
			}
		}
	}
}

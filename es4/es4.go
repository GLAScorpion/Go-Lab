package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
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
	timeEUR     = time.Time{}
	timeGBP     = time.Time{}
	timeJPY     = time.Time{}
)

var (
	sellEUR = false
	buyGBP  = false
	buyJPY  = false
)

var (
	mutexEUR = sync.RWMutex{}
	mutexGBP = sync.RWMutex{}
	mutexJPY = sync.RWMutex{}
)

func main() {
	go simulateMarketData()
	go selectPair()
	time.Sleep(time.Minute)
}

func simulateMarketData() {
	timeEUR = time.Now()
	timeGBP = time.Now()
	timeJPY = time.Now()
	for {
		mutexEUR.Lock()
		if (!sellEUR && time.Since(timeEUR).Seconds() >= 1.0) || (sellEUR && time.Since(sellTimeEUR).Seconds() >= 4.0) {
			timeEUR = time.Now()
			mutexEUR.Unlock()
			eur_usd <- float64(int32(rand.Int31()))/float64(math.MaxInt32)*0.5 + 1.0
		} else {
			mutexEUR.Unlock()
		}
		mutexGBP.Lock()
		if (!buyGBP && time.Since(timeGBP).Seconds() >= 1.0) || (buyGBP && time.Since(buyTimeGBP).Seconds() >= 3.0) {
			timeGBP = time.Now()
			mutexGBP.Unlock()
			gbp_usd <- float64(int32(rand.Int31()))/float64(math.MaxInt32)*0.5 + 1.0
		} else {
			mutexGBP.Unlock()
		}
		mutexJPY.Lock()
		if (!buyJPY && time.Since(timeJPY).Seconds() >= 1.0) || (buyJPY && time.Since(buyTimeJPY).Seconds() >= 3.0) {
			timeJPY = time.Now()
			mutexJPY.Unlock()
			jpy_usd <- float64(int32(rand.Int31()))/float64(math.MaxInt32)*0.003 + 0.006
		} else {
			mutexJPY.Unlock()
		}
		//time.Sleep(time.Second)
	}
}

func selectPair() {

	for {
		select {
		case val := <-eur_usd:
			mutexEUR.Lock()
			if val > 1.20 && !sellEUR {
				holdEUR = val
				sellTimeEUR = time.Now()
				sellEUR = true
			} else if sellEUR {
				sellEUR = false
				fmt.Println("Detected EUR/USD at", holdEUR, ", sold at", val, time.Since(sellTimeEUR))
			}
			mutexEUR.Unlock()
		case val := <-gbp_usd:
			mutexGBP.Lock()
			if val < 1.35 && !buyGBP {
				holdGBP = val
				buyTimeGBP = time.Now()
				buyGBP = true
			} else if buyGBP {
				buyGBP = false
				fmt.Println("Detected GBP/USD at", holdGBP, ", bought at", val, time.Since(buyTimeGBP))
			}
			mutexGBP.Unlock()
		case val := <-jpy_usd:
			mutexJPY.Lock()
			if val < 0.0085 && buyJPY {
				holdJPY = val
				buyTimeJPY = time.Now()
				buyJPY = true
			} else if buyJPY {
				buyJPY = false
				fmt.Println("Detected JPY/USD at", holdJPY, ", bought at", val, time.Since(buyTimeJPY))
			}
			mutexJPY.Unlock()
		}
		//time.Sleep(500 * time.Millisecond)
	}
}

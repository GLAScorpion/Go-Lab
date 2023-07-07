package main

/*
	Ho inteso le transazioni come possibili in simultanea per tipi differenti di valuta.
	Prima di poter effettuare una transazione con la stessa valuta bisognerà aspettare il
	tempo di necessario a completare l'operazione precedente, se in corso
*/
import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var ( // canali delle valute
	eurUSD = make(chan float64)
	gbpUSD = make(chan float64)
	jpyUSD = make(chan float64)
)
var ( // salvano il valore di rilevazione delle valute
	holdEUR = 0.0
	holdGBP = 0.0
	holdJPY = 0.0
)

var ( // tempo di acquisto delle valute
	timeEUR     = time.Time{}
	timeGBP     = time.Time{}
	timeJPY     = time.Time{}
	timeTracker = time.Time{}
)

var ( // stato delle transazioni (true = in corso, false = non in corso)
	sellEUR = false
	buyGBP  = false
	buyJPY  = false
)

func main() {
	go simulateMarketData()
	go selectPair()
	timeTracker = time.Now() // tracking del tempo di esecuzione
	time.Sleep(time.Minute)
}

func simulateMarketData() { // genera valori randomici ogni secondo
	for {
		eurUSD <- float64(int32(rand.Int31()))/float64(math.MaxInt32)*0.5 + 1.0
		gbpUSD <- float64(int32(rand.Int31()))/float64(math.MaxInt32)*0.5 + 1.0
		jpyUSD <- float64(int32(rand.Int31()))/float64(math.MaxInt32)*0.003 + 0.006
		time.Sleep(time.Second)
	}
}

func elapsedTimePrinter() {
	fmt.Print("[", int(time.Since(timeTracker).Seconds()), " seconds passed] ")
}

func selectPair() {
	// variabili di tempo azzerate
	timeEUR = time.Now()
	timeGBP = time.Now()
	timeJPY = time.Now()
	for {
		select {
		case val := <-eurUSD:
			if val > 1.20 && !sellEUR { // inizia operazione se il valore è accettabile e non ci sono operazioni simili già in corso
				holdEUR = val        // salvo valore a inizio operazione
				timeEUR = time.Now() // azzero il tempo
				sellEUR = true       // operazione in corso
			} else if sellEUR && time.Since(timeEUR).Seconds() >= 4.0 { // completa operazione solo se in corso e dopo un tempo accettabile
				sellEUR = false // termine operazione
				elapsedTimePrinter()
				fmt.Println("Detected EUR/USD at", holdEUR, ", sold at", val, ". It took", time.Since(timeEUR).Seconds(), "seconds")
			}
		case val := <-gbpUSD:
			if val < 1.35 && !buyGBP { // inizia operazione se il valore è accettabile e non ci sono operazioni simili già in corso
				holdGBP = val        // salvo valore a inizio operazione
				timeGBP = time.Now() // azzero il tempo
				buyGBP = true        // operazione in corso
			} else if buyGBP && time.Since(timeGBP).Seconds() >= 3.0 { // completa operazione solo se in corso e dopo un tempo accettabile
				buyGBP = false // termine operazione
				elapsedTimePrinter()
				fmt.Println("Detected GBP/USD at", holdGBP, ", bought at", val, ". It took", time.Since(timeGBP).Seconds(), "seconds")
			}
		case val := <-jpyUSD:
			if val < 0.0085 && !buyJPY { // inizia operazione se il valore è accettabile e non ci sono operazioni simili già in corso
				holdJPY = val        // salvo valore a inizio operazione
				timeJPY = time.Now() // azzero il tempo
				buyJPY = true        // operazione in corso
			} else if buyJPY && time.Since(timeJPY).Seconds() >= 3.0 { // completa operazione solo se in corso e dopo un tempo accettabile
				buyJPY = false // termine operazione
				elapsedTimePrinter()
				fmt.Println("Detected JPY/USD at", holdJPY, ", bought at", val, ". It took", time.Since(timeJPY).Seconds(), "seconds")
			}
		}
	}
}

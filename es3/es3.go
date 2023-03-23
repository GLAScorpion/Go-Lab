package main

import (
	"fmt"
	"sync"
	"time"
)

type cake struct {
	id int
}

type worker struct {
	willWork bool
}

const ( //costanti di tempo dei lavori in millisecondi
	kBakeTime     = 1000
	kGarnishTime  = 4000
	kDecorateTime = 8000
	kCakeNum      = 5 // numero di torte da processare
)

var uncookedCakes = make(map[int]cake)  // torte da preparare
var bakedCakes = make(map[int]cake)     // torte cotte
var garnishedCakes = make(map[int]cake) // torte guarnite
var doneCakes = make(map[int]cake)      // torte complete

var wg = sync.WaitGroup{}

var uncookMutex = sync.RWMutex{}
var bakeMutex = sync.RWMutex{}
var garnishMutex = sync.RWMutex{}
var doneMutex = sync.RWMutex{}

func main() {
	for i := 0; i < kCakeNum; i++ { //torte vengono inizializzate
		uncookedCakes[i] = cake{i}
	}
	startTime := time.Now()
	baker := worker{true} // variabili che permettono di controllare se il pasticcere precedente produrrà una torta
	garnisher := worker{true}
	wg.Add(3)
	go bake(&baker)
	go garnish(&garnisher, &baker)
	go decorate(&garnisher)
	go func() { // routine che stampa in tempo reale lo stato degli spazi a disposizione dei pasticceri
		for {
			printer()
			time.Sleep(500 * time.Millisecond)
		}
	}()
	wg.Wait()
	printer()
	elapsed := int(time.Since(startTime).Seconds()) // tempo trascorso per produrre le torte
	fmt.Println("It took", elapsed, "seconds")
}

func bake(baker *worker) {
	uncookMutex.Lock()
	for key, val := range uncookedCakes { // accede a tutte le torte iniziali senza vincoli
		uncookMutex.Unlock()
		seeker := 0 //indice per scorrere gli spazi dove mettere le torte complete
		for {
			bakeMutex.Lock()
			if _, ok := bakedCakes[seeker]; !ok { //lavora solo se c'è uno spazio libero tra i 2 slot a disposizione
				bakeMutex.Unlock()
				time.Sleep(kBakeTime * time.Millisecond) //lavoro
				bakeMutex.Lock()
				bakedCakes[seeker] = val //posiziona torta nello slot vuoto
				bakeMutex.Unlock()
				uncookMutex.Lock()
				delete(uncookedCakes, key)
				uncookMutex.Unlock()
				break
			} else {
				bakeMutex.Unlock()
			}
			seeker = (seeker + 1) % 2 // gli unici indici validi sono 2: 0 e 1
		}
		uncookMutex.Lock()
	}
	uncookMutex.Unlock()
	bakeMutex.Lock()
	baker.willWork = false //non produrrà altre torte
	bakeMutex.Unlock()
	wg.Done()
}

func garnish(garnisher *worker, baker *worker) {
	for {
		seeker := 0
	garnisherLoop:
		for {
			garnishMutex.Lock()
			if _, ok := garnishedCakes[seeker]; !ok {
				garnishMutex.Unlock()
				seeker2 := 0
				for {
					bakeMutex.Lock()
					if val, ok := bakedCakes[seeker2]; ok {
						delete(bakedCakes, seeker2)
						bakeMutex.Unlock()
						time.Sleep(kGarnishTime * time.Millisecond)
						garnishMutex.Lock()
						garnishedCakes[seeker] = val
						garnishMutex.Unlock()
						break garnisherLoop
					} else if !baker.willWork && seeker2 == 1 {
						bakeMutex.Unlock()
						garnisher.willWork = false
						wg.Done()
						return
					} else {
						bakeMutex.Unlock()
					}
					seeker2 = (seeker2 + 1) % 2
				}
			} else {
				garnishMutex.Unlock()
			}
			seeker = (seeker + 1) % 2
		}
	}
}

func decorate(garnisher *worker) {
	for i := 0; ; i++ {
		seeker := 0
		for {
			garnishMutex.Lock()
			if val, ok := garnishedCakes[seeker]; ok {
				delete(garnishedCakes, seeker)
				garnishMutex.Unlock()
				time.Sleep(kDecorateTime * time.Millisecond)
				doneMutex.Lock()
				doneCakes[i] = val
				doneMutex.Unlock()
				break
			} else if !garnisher.willWork && seeker == 1 {
				garnishMutex.Unlock()
				wg.Done()
				return
			} else {
				garnishMutex.Unlock()
			}

			seeker = (seeker + 1) % 2
		}
	}
}

func printer() {
	fmt.Print("\033[H\033[2J")
	uncookMutex.Lock()
	fmt.Print("Cakes to bake: ")
	for _, val := range uncookedCakes {
		fmt.Print("[Cake ", val.id+1, "]")
	}
	fmt.Println()
	uncookMutex.Unlock()
	bakeMutex.Lock()
	fmt.Print("Baked cakes: ")
	for _, val := range bakedCakes {
		fmt.Print("[Cake ", val.id+1, "]")
	}
	fmt.Println()
	bakeMutex.Unlock()
	garnishMutex.Lock()
	fmt.Print("Garnished cakes: ")
	for _, val := range garnishedCakes {
		fmt.Print("[Cake ", val.id+1, "]")
	}
	fmt.Println()
	garnishMutex.Unlock()
	doneMutex.Lock()
	fmt.Print("Completed cakes: ")
	for _, val := range doneCakes {
		fmt.Print("[Cake ", val.id+1, "]")
	}
	fmt.Println()
	doneMutex.Unlock()
}

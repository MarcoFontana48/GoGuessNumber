package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type Player struct {
	id          int
	maxNumBound int
	minNumBound int
}

func NewPlayer(id int) *Player {
	return &Player{id: id, maxNumBound: MAX, minNumBound: 0}
}

func (p *Player) start(startChan <-chan bool, guessChan chan<- int, resultChan <-chan string, adviceChan <-chan string, winChan <-chan string, wg *sync.WaitGroup, winWg *sync.WaitGroup) {
	wg.Done()

	winWg.Add(1)
	go func() {
		defer winWg.Done()
		playerResult := <-winChan
		fmt.Printf("PLAYER %d: received %s message\n", p.id, playerResult)
	}()

	for {
		// wait for the turn to start
		fmt.Printf("PLAYER %d: is waiting for the turn to start\n", p.id)
		<-startChan
		fmt.Printf("PLAYER %d: received start turn message\n", p.id)

		// guess a number
		guess := rand.Intn(p.maxNumBound-p.minNumBound) + p.minNumBound
		fmt.Printf("PLAYER %d: guessed %d\n", p.id, guess)
		guessChan <- guess

		// receive result
		result := <-resultChan
		fmt.Printf("PLAYER %d: received result %s\n", p.id, result)

		if result == "correct" {
			fmt.Printf("PLAYER %d: guessed the secret number\n", p.id)
		} else if result == "incorrect" {
			// receive advice
			advice := <-adviceChan
			fmt.Printf("PLAYER %d: received advice %s\n", p.id, advice)

			if advice == "higher" {
				p.minNumBound = guess + 1
				fmt.Printf("PLAYER %d: setting minNumBound to %d\n", p.id, p.minNumBound)
			} else {
				p.maxNumBound = guess
				fmt.Printf("PLAYER %d: setting maxNumBound to %d\n", p.id, p.maxNumBound)
			}
		}

		wg.Done()
	}
}

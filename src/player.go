package main

import (
	"fmt"
	//"github.com/rabbitmq/amqp091-go"
	"math/rand"
)

type Player struct {
	playerId        int
	minGuessableNum int
	maxGuessableNum int
}

func NewPlayer(id int) *Player {
	return &Player{playerId: id, minGuessableNum: 0, maxGuessableNum: MAX}
}

func (p *Player) start(guessChan chan<- *Guess, adviceChan <-chan string) {
	p.guess(guessChan, adviceChan)
}

func (p *Player) guess(guessChan chan<- *Guess, adviceChan <-chan string) {
	guessNum := rand.Intn(p.maxGuessableNum-p.minGuessableNum+1) + p.minGuessableNum
	fmt.Printf("PLAYER %d: Initial guess: '%d'\n", p.playerId, guessNum)
	guessChan <- &Guess{playerId: p.playerId, number: guessNum}

	for advice := range adviceChan {
		if advice == "CORRECT" {
			fmt.Printf("PLAYER %d: I guessed '%d' and it is correct!\n", p.playerId, guessNum)
			break
		} else if advice == "HIGHER" {
			p.SetMinGuessableNum(guessNum + 1)
		} else {
			p.SetMaxGuessableNum(guessNum - 1)
		}
	}
}

func (p *Player) SetMinGuessableNum(newMin int) {
	p.minGuessableNum = newMin
	fmt.Printf("PLAYER %d: set new minimum bound to %d\n", p.playerId, p.minGuessableNum)
}

func (p *Player) SetMaxGuessableNum(newMax int) {
	p.maxGuessableNum = newMax
	fmt.Printf("PLAYER %d: set new maximum bound to %d\n", p.playerId, p.maxGuessableNum)
}

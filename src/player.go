package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type Player struct {
	playerId        int
	guessedNum      int
	minGuessableNum int
	maxGuessableNum int
}

func NewPlayer(id int) *Player {
	return &Player{playerId: id, minGuessableNum: 0, maxGuessableNum: MAX}
}

func (p *Player) guess(guesses chan<- *Player, wg *sync.WaitGroup) {
	// decrease the WaitGroup counter when the goroutine completes.
	defer wg.Done()

	p.guessedNum = rand.Intn(p.maxGuessableNum-p.minGuessableNum+1) + p.minGuessableNum
	fmt.Printf("PLAYER %d: guessed %d\n", p.playerId, p.guessedNum)
	guesses <- p
}

func (p *Player) SetMinGuessableNum(newMin int) {
	p.minGuessableNum = newMin
	fmt.Printf("PLAYER %d: set new minimum bound to %d\n", p.playerId, p.minGuessableNum)
}

func (p *Player) SetMaxGuessableNum(newMax int) {
	p.maxGuessableNum = newMax
	fmt.Printf("PLAYER %d: set new maximum bound to %d\n", p.playerId, p.maxGuessableNum)
}

package main

import (
	"fmt"
	//"github.com/rabbitmq/amqp091-go"
	"math/rand"
)

type Oracle struct {
	secretNumber int
}

func NewOracle() *Oracle {
	return &Oracle{secretNumber: rand.Intn(MAX + 1)}
}

func (o *Oracle) start() {
	fmt.Printf("ORACLE: secret number is '%d'\n", o.secretNumber)

	// creates channels of size 'N' (N == number of players) to receive their guesses
	guessChan := make(chan *Guess, N)
	adviceChan := make(map[int]chan string)

	o.startGame(guessChan, adviceChan)
}

func (o *Oracle) startGame(guessChan chan *Guess, adviceChan map[int]chan string) {
	o.startPlayers(guessChan, adviceChan)

	for {
		fmt.Println("ORACLE: a new turn starts.")

		if o.checkGuessedNumbers(guessChan, adviceChan) {
			fmt.Println("ORACLE: someone guessed the right number, game over.")
			break
		}

		fmt.Println("ORACLE: no one guessed the right number, a new turn starts.")
	}
}

func (o *Oracle) startPlayers(guessChan chan *Guess, adviceChan map[int]chan string) {
	for i := 0; i < N; i++ {
		adviceChan[i] = make(chan string)
		go NewPlayer(i).start(guessChan, adviceChan[i])
	}
}

func (o *Oracle) checkGuessedNumbers(guessChan <-chan *Guess, adviceChan map[int]chan string) bool {
	for guess := range guessChan {
		if guess.number == o.secretNumber {
			adviceChan[guess.playerId] <- "CORRECT"
			return true
		} else if guess.number < o.secretNumber {
			adviceChan[guess.playerId] <- "HIGHER"
		} else {
			adviceChan[guess.playerId] <- "LOWER"
		}
	}
	return false
}

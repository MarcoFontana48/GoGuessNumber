package main

import (
	"fmt"
	//"github.com/rabbitmq/amqp091-go"
	"math/rand"
	"sync"
)

type Oracle struct {
	secretNumber int
}

func NewOracle() *Oracle {
	return &Oracle{secretNumber: rand.Intn(MAX + 1)}
}

func (o *Oracle) start() {
	fmt.Printf("ORACLE: secret number is '%d'\n", o.secretNumber)

	// creates a channel of 'Players' to receive their guesses
	guesses := make(chan *Player)

	// creates a WaitGroup to wait for all players to finish their guesses
	var wg sync.WaitGroup

	// starts the game
	for {
		fmt.Println("ORACLE: a new turn starts.")

		wg.Add(N)

		// creates N players (workers) and starts their guesses
		for i := 0; i < N; i++ {
			go NewPlayer(i+1).guess(guesses, &wg)
		}

		// closes the channel when all players have finished their guesses
		go func() {
			wg.Wait()
			close(guesses)
		}()

		winner := false
		for guess := range guesses {
			// checks if the player has guessed the secret number
			if o.guessNumber(guess) {
				fmt.Printf("ORACLE: player %d guessed the number correctly! The number was %d.\n", guess.playerId, o.secretNumber)
				winner = true
				break
			}
		}

		if winner {
			break
		} else {
			fmt.Println("ORACLE: no one guessed the number correctly.")
			guesses = make(chan *Player)
		}
	}
}

func (o *Oracle) guessNumber(p *Player) bool {
	if p.guessedNum < o.secretNumber {
		fmt.Printf("ORACLE: player %d the secret number is greater than %d\n", p.playerId, p.guessedNum)
		p.SetMinGuessableNum(p.guessedNum + 1)
	} else if p.guessedNum > o.secretNumber {
		fmt.Printf("ORACLE: player %d the secret number is less than %d\n", p.playerId, p.guessedNum)
		p.SetMaxGuessableNum(p.guessedNum - 1)
	}
	return p.guessedNum == o.secretNumber
}

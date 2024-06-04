package main

import (
	"fmt"
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
	fmt.Println("ORACLE: the secret number is ", o.secretNumber)

	// initialize channels
	startChan := make(map[int]chan bool)
	guessChan := make(map[int]chan int)
	adviceChan := make(map[int]chan string)
	resultChan := make(map[int]chan string)
	winChan := make(map[int]chan string)

	var oracleWaitGroup sync.WaitGroup
	oracleWaitGroup.Add(N)

	var winWaitGroup sync.WaitGroup

	for i := 0; i < N; i++ {
		// initialize the individual channels in the map
		startChan[i] = make(chan bool)
		guessChan[i] = make(chan int)
		adviceChan[i] = make(chan string)
		resultChan[i] = make(chan string)
		winChan[i] = make(chan string)

		// creates a new player
		player := NewPlayer(i)

		// starts the player
		go player.start(startChan[i], guessChan[i], resultChan[i], adviceChan[i], winChan[i], &oracleWaitGroup, &winWaitGroup)
	}

	oracleWaitGroup.Wait()
	winner := -1

	for {
		oracleWaitGroup.Add(N)

		// starts the game
		fmt.Printf("ORACLE: turn starts\n")
		for i := 0; i < N; i++ {
			// signals the players to start
			go func() {
				startChan[i] <- true
			}()
		}

		var playersWaitGroup sync.WaitGroup
		playersWaitGroup.Add(N)

		for i := 0; i < N; i++ {
			go func(i int) {
				defer playersWaitGroup.Done()

				// receives guesses from players
				guess := <-guessChan[i]
				fmt.Printf("ORACLE: received guess %d from %d\n", guess, i)

				// sends advice to the player
				if guess == o.secretNumber {
					fmt.Printf("ORACLE: player %d guessed the secret number (%d)\n", i, guess)
					resultChan[i] <- "correct"
					if winner == -1 {
						winner = i
					}
				} else {
					fmt.Printf("ORACLE: player %d guessed incorrectly (%d)\n", i, guess)
					resultChan[i] <- "incorrect"
					if guess < o.secretNumber {
						fmt.Printf("ORACLE: player %d should guess higher\n", i)
						adviceChan[i] <- "higher"
					} else {
						fmt.Printf("ORACLE: player %d should guess lower\n", i)
						adviceChan[i] <- "lower"
					}
				}
			}(i)
		}

		// wait for all players to finish
		playersWaitGroup.Wait()
		fmt.Printf("ORACLE: turn ends\n")

		if winner > -1 {
			break
		}
	}

	for i := 0; i < N; i++ {
		// signals the players who won and who lost
		go func() {
			if i == winner {
				winChan[i] <- "won"
			} else {
				winChan[i] <- "lost"
			}
		}()
	}

	winWaitGroup.Wait()

	fmt.Println("ORACLE: game over")
}

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
	startChan := make(map[int]chan bool)    //communicates when the game or turn starts
	guessChan := make(map[int]chan int)     //communicates guesses from players
	adviceChan := make(map[int]chan string) //oracle tells players suggestions based on guessed number
	resultChan := make(map[int]chan string) //oracle tells players if they have correctly guessed the number
	winChan := make(map[int]chan string)    //oracle tells to each player whether they have won or lost at the end of the game

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

	//oracle waits until all players have been initialized before proceeding
	oracleWaitGroup.Wait()
	winner := -1 //initialize winner player index

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
			//defines an anonimous function with 'chanIndex' as argument. Sends suggestions to players based on their guesses
			go func(playerIndex int) {
				defer playersWaitGroup.Done()

				// receives guesses from players
				guess := <-guessChan[playerIndex]
				fmt.Printf("ORACLE: received guess %d from %d\n", guess, playerIndex)

				// sends advice to the player
				if guess == o.secretNumber {
					fmt.Printf("ORACLE: player %d guessed the secret number (%d)\n", playerIndex, guess)
					resultChan[playerIndex] <- "correct"
					if winner == -1 {
						winner = playerIndex
					}
				} else {
					fmt.Printf("ORACLE: player %d guessed incorrectly (%d)\n", playerIndex, guess)
					resultChan[playerIndex] <- "incorrect"
					if guess < o.secretNumber {
						fmt.Printf("ORACLE: player %d should guess higher\n", playerIndex)
						adviceChan[playerIndex] <- "higher"
					} else {
						fmt.Printf("ORACLE: player %d should guess lower\n", playerIndex)
						adviceChan[playerIndex] <- "lower"
					}
				}
			}(i) //argument of the function
		}

		// wait for all players to receive the guesses
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

	//waits until game has ended
	winWaitGroup.Wait()

	fmt.Println("ORACLE: game over")
}

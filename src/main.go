package main

//import (
//	"github.com/rabbitmq/amqp091-go"
//)

const (
	//there has to be only a single winner (even if multiple guessed correctly) and it has to be the player whom the oracle received the correct secret number first
	//MAX = 0 			// Max secret number
	//N   = 3 			// Number of players

	//message order has to be evaluated correctly
	MAX = 25 // Max secret number
	N   = 5  // Number of players

	//massive
	//MAX = 1000 // Max secret number
	//N   = 1000 // Number of players
)

func main() {
	// creates a new oracle 'master'
	oracle := NewOracle()

	// starts the oracle
	oracle.start()
}

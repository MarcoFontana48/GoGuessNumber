package main

//import (
//	"github.com/rabbitmq/amqp091-go"
//)

const (
	MAX = 100 // Max secret number
	N   = 6   // Number of players
	//MAX = 100_000 	// Max secret number
	//N   = 10_000  	// Number of players
)

func main() {
	// creates a new oracle 'master'
	oracle := NewOracle()

	// starts the oracle
	oracle.start()
}

package main

const (
	//MAX = 0 			// Max secret number
	//N   = 3 			// Number of players

	MAX = 25 // Max secret number
	N   = 5  // Number of players

	//MAX = 1000 		// Max secret number
	//N   = 1000 		// Number of players
)

func main() {
	// creates a new oracle 'master'
	oracle := NewOracle()

	// starts the oracle
	oracle.start()
}

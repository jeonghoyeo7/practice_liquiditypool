package main

import (
	"time"
	"uniswap_test/test"
)

const simulationTime time.Duration = 10 // 10 seconds simulation

func main() {
	//liquidity.AvgArrival = 1000 // ms for new 1  arrival

	// Generate Transactions
	//tq := liquidity.CreateTransactionQueue()
	//go test.GenerateTransactions(tq)

	//time.Sleep(simulationTime * time.Second)

	test.Simpletest()
}

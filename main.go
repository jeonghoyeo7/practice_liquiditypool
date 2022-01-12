package main

import (
	"time"
	"uniswap_test/processings"
	transactions "uniswap_test/transactions"
)

const simulationTime time.Duration = 10 // 10 seconds simulation

func main() {
	//transactions.AvgArrival = 1000 // ms for new 1 transaction arrival
	tq := transactions.NewTransactionQueue()
	//testInt := sdk.NewInt(1).Add(sdk.NewInt(10))
	//testCoin := sdk.NewCoin("uatom", testInt)
	//testDec := sdk.NewDec(10).Quo(sdk.NewDec(3))

	//fmt.Println(testCoin, testInt, testDec)
	//fmt.Println(testCoin)

	// Generate Transactions
	go transactions.GenerateTransactions(tq)

	// blockchain goes.
	// go processings.ProcessBlocks()
	processings.ProcessBlocks()
	time.Sleep(simulationTime * time.Second)
}

package main

import (
	"time"
	"uniswap_test/liquidity"
)

const simulationTime time.Duration = 10 // 10 seconds simulation

func main() {
	//liquidity.AvgArrival = 1000 // ms for new 1  arrival
	tq := liquidity.NewTransactionQueue()
	//testInt := sdk.NewInt(1).Add(sdk.NewInt(10))
	//testCoin := sdk.NewCoin("uatom", testInt)
	//testDec := sdk.NewDec(10).Quo(sdk.NewDec(3))

	//testDec := sdk.NewDec(10).Quo(sdk.NewDec(4))
	//fmt.Println(testDec)
	//fmt.Println(testCoin, testInt, testDec)
	//fmt.Println(testCoin)

	// Generate Transactions
	go liquidity.GenerateTransactions(tq)

	// blockchain goes.
	// go processings.ProcessBlocks()
	//processings.ProcessBlocks()
	time.Sleep(simulationTime * time.Second)
}

package test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sync"
	"time"
	"uniswap_test/liquidity"
)

func RandomSimulate(simulateTime time.Duration, intervalTransactions time.Duration, intervalBlock time.Duration) {
	var wg sync.WaitGroup

	// create a liquidity pool --------------------------------------------
	lp := liquidity.CreatePool()
	// pool initial setup
	var CoinX, CoinY sdk.DecCoin
	CoinX.Amount = sdk.NewDec(int64(randInt(1000, 10000)))
	CoinX.Denom = "uatom"
	time.Sleep(2 * time.Nanosecond) // delay for 2 ns to make different random numbers
	CoinY.Amount = sdk.NewDec(int64(randInt(1000, 10000)))
	CoinY.Denom = "eth"
	Fee := sdk.NewDec(int64(0))
	lp.SetOnce(CoinX, CoinY, Fee)
	fmt.Println("The liquidity pool is randomly given.")
	lp.Print() // Print state of the pool

	// Create Transaction Queue -----------------------------------------
	tq := liquidity.CreateTransactionQueue()
	//tq := make(chan liquidity.TransactionQueue)    // transaction queue in real time
	tqResult := liquidity.CreateTransactionQueue() //result queue
	wg.Add(1)

	go processTransactions(&wg, tq, tqResult, &lp, simulateTime, intervalBlock)

	go generateTransactions(&wg, tq, simulateTime, intervalTransactions)

	wg.Wait()

}

func processTransactions(wg *sync.WaitGroup, tq *liquidity.TransactionQueue, tqResult *liquidity.TransactionQueue, lp *liquidity.Pool, simulateTime time.Duration, intervalBlock time.Duration) {
	tick := time.Tick(time.Millisecond * 100 * intervalBlock)
	terminate := time.After(time.Millisecond * 100 * simulateTime)

	for {
		select {
		case <-tick:
			// Processing the queue
			fmt.Println("Processing Queue...")
		case <-terminate:
			fmt.Println("Simulation is finished!")
			wg.Done()
			return
		}

	}
}

func generateTransactions(wg *sync.WaitGroup, tq *liquidity.TransactionQueue, simulateTime time.Duration, intervalTransactions time.Duration) {
	tick := time.Tick(time.Millisecond*100*intervalTransactions + time.Millisecond*50)
	terminate := time.After(time.Millisecond * 100 * simulateTime)

	for {
		select {
		case <-tick:
			// Generating transactions
			fmt.Println("Generating transactions.....")
		case <-terminate:
			fmt.Println("Simulation is finished!")
			wg.Done()
			return
		}

	}
}

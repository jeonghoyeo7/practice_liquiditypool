package test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sync"
	"time"
	"uniswap_test/liquidity"
)

var cntProcess int
var cntTransaction int

func RandomSimulation(simulateTime time.Duration, intervalTransactions time.Duration, intervalBlock time.Duration) {
	var wg sync.WaitGroup

	// counters
	cntProcess = 0
	cntTransaction = 0

	// Create Channel for transactions
	channelTransaction := make(chan liquidity.Transaction)

	wg.Add(2)

	// go routine to generate transactions
	go generateTransactions(&wg, channelTransaction, simulateTime, intervalTransactions)

	// go routine to process transactions
	go processTransactions(&wg, channelTransaction, simulateTime, intervalBlock)

	wg.Wait()
}

// go routine to process transactions
func processTransactions(wg *sync.WaitGroup, channelTransaction chan liquidity.Transaction, simulateTime time.Duration, intervalBlock time.Duration) {
	tick := time.Tick(time.Second*intervalBlock + time.Millisecond)
	terminate := time.After(time.Second*simulateTime + time.Second + time.Millisecond)

	// create a liquidity pool --------------------------------------------
	lp := liquidity.CreatePool()
	// pool initial setup
	var CoinX, CoinY sdk.DecCoin
	CoinX.Amount = sdk.NewDec(int64(randInt(1000, 10000)))
	CoinX.Denom = "atom"
	time.Sleep(2 * time.Nanosecond) // delay for 2 ns to make different random numbers
	CoinY.Amount = sdk.NewDec(int64(randInt(1000, 10000)))
	CoinY.Denom = "eth"
	Fee := sdk.NewDec(int64(0))
	lp.SetOnce(CoinX, CoinY, Fee)
	fmt.Println("The liquidity pool is randomly given.")
	lp.Print() // Print state of the pool
	initialPool := lp

	// Create Transaction Queue -----------------------------------------
	tq := liquidity.CreateTransactionQueue()

	// Creat Result Queue -----------------------------------------------
	tqResult := liquidity.CreateTransactionQueue() //result queue

	for {
		select {
		case <-tick:
			// Processing the queue

			fmt.Printf("We have now %d transactions to be committed. \n", tq.Len().(int))
			for i := 0; i < tq.Len().(int); i++ {
				tran := tq.Pop().(liquidity.Transaction)
				cntProcess += 1
				fmt.Println(cntProcess, " commit a transaction : ", tran.SprintLine())
				tranResult := processTransaction(tran, &lp)
				fmt.Println("        >> current Pool : ", lp.SprintLine())
				tqResult.Push(tranResult)
			}

		case <-terminate:
			fmt.Printf("Current Status of liquidity pool after %d processing\n", cntProcess)
			liquidity.ComparePoolsPrint(initialPool, lp)

			fmt.Println("All block commit is finished!")
			wg.Done()
			return

		case newTran := <-channelTransaction:
			tq.Push(newTran)
		}

	}
}

// go routine to generate transactions
func generateTransactions(wg *sync.WaitGroup, channelTransaction chan liquidity.Transaction, simulateTime time.Duration, intervalTransactions time.Duration) {
	tick := time.Tick(time.Second * intervalTransactions)
	terminate := time.After(time.Second * simulateTime)

	for {
		select {
		case <-tick:
			// Generation start -----------------------------------
			newTran := generateRandomTransaction()
			cntTransaction += 1
			fmt.Println(cntTransaction, " Broadcasting...", newTran.SprintLine())
			channelTransaction <- newTran
		//time.Sleep(time.Millisecond * 10 * intervalTransactions) // wait until the next transaction is coming
		// Generation end -----------------------------------
		case <-terminate:
			fmt.Println("Block broadcast ended.")
			wg.Done()
			return
		}
	}
}

// generate a transaction randomly one time
func generateRandomTransaction() liquidity.Transaction {
	var newTran liquidity.Transaction
	newTran.Init()

	trade := randInt(1, 4)
	time.Sleep(2 * time.Nanosecond)
	switch trade {
	case 1:
		newTran.Order = "swapXtoY"
		newTran.CoinX.Amount = sdk.NewDec(0).Add(sdk.NewDec(int64(randInt(1, 100))))
	case 2:
		newTran.Order = "swapYtoX"
		newTran.CoinY.Amount = sdk.NewDec(0).Add(sdk.NewDec(int64(randInt(1, 100))))
	case 3:
		newTran.Order = "deposit"
		newTran.CoinX.Amount = sdk.NewDec(0).Add(sdk.NewDec(int64(randInt(1, 100))))
		time.Sleep(2 * time.Nanosecond) // delay for 2 ns to make different random numbers
		newTran.CoinY.Amount = sdk.NewDec(0).Add(sdk.NewDec(int64(randInt(1, 100))))
	case 4:
		newTran.Order = "withdraw"
		newTran.LiquidityToken.Amount = sdk.NewDec(0).Add(sdk.NewDec(int64(randInt(1, 100))))
	}

	return newTran
}

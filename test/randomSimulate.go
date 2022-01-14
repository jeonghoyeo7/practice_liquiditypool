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

func RandomSimulate(simulateTime time.Duration, intervalTransactions time.Duration, intervalBlock time.Duration) {
	var wg sync.WaitGroup

	cntProcess = 0
	cntTransaction = 0

	fmt.Println(intervalTransactions, intervalBlock, simulateTime)

	// Create Channel for transactions
	channelTransaction := make(chan liquidity.Transaction)

	wg.Add(2)

	go processTransactions(&wg, channelTransaction, simulateTime, intervalBlock)

	go generateTransactions(&wg, channelTransaction, simulateTime, intervalTransactions)

	/*
		for {
			select {
			case <-tick:
				// Generation start -----------------------------------
				newTran := generateRandomTransaction()
				cntTransaction += 1
				fmt.Println(cntTransaction, " generating...", newTran)
				channelTransaction <- newTran
			//time.Sleep(time.Millisecond * 10 * intervalTransactions) // wait until the next transaction is coming
			// Generation end -----------------------------------
			case <-terminate:
				fmt.Println("Simulation ended.")
				wg.Done()
				break
			}
		}
	*/
	fmt.Println("waiting...")
	wg.Wait()

}

func processTransactions(wg *sync.WaitGroup, channelTransaction chan liquidity.Transaction, simulateTime time.Duration, intervalBlock time.Duration) {
	tick := time.Tick(time.Second * intervalBlock)
	terminate := time.After(time.Second*simulateTime + time.Second)

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
	//tqChannel := make(chan *liquidity.TransactionQueue) // transaction queue channel
	tq := liquidity.CreateTransactionQueue()

	// Creat Result Queue -----------------------------------------------
	//tqResultChannel := make(chan *liquidity.TransactionQueue) // transaction queue channel
	tqResult := liquidity.CreateTransactionQueue() //result queue

	for {
		select {
		case <-tick:
			// Processing the queue
			//fmt.Println("Processing Queue...")

			// Processing start -----------------------------------
			//fmt.Printf("Processing... %d transactions \n", tq.Len().(int))

			fmt.Printf("We have now %d transactions to be processed. \n", tq.Len())
			for i := 0; i < tq.Len().(int); i++ {
				tran := tq.Pop().(liquidity.Transaction)
				cntProcess += 1
				fmt.Println(cntProcess, " processing : ", tran)
				tranResult := processTransaction(tran, &lp)
				fmt.Println("current Pool : ", lp)
				tqResult.Push(tranResult)
			}

			//fmt.Println("Processing finished...")
			// Processing end -----------------------------------

		case <-terminate:
			fmt.Printf("Current Status of liquidity pool after %d processing\n", cntProcess)
			lp.Print() // Print state of the pool

			fmt.Println("Simulation is finished!")
			wg.Done()
			return

		case newTran := <-channelTransaction:
			tq.Push(newTran)
		}

	}
}

func generateTransactions(wg *sync.WaitGroup, channelTransaction chan liquidity.Transaction, simulateTime time.Duration, intervalTransactions time.Duration) {
	tick := time.Tick(time.Second*intervalTransactions + time.Millisecond)
	terminate := time.After(time.Second*simulateTime + time.Millisecond)

	for {
		select {
		case <-tick:
			// Generation start -----------------------------------
			newTran := generateRandomTransaction()
			cntTransaction += 1
			fmt.Println(cntTransaction, " generating...", newTran)
			channelTransaction <- newTran
		//time.Sleep(time.Millisecond * 10 * intervalTransactions) // wait until the next transaction is coming
		// Generation end -----------------------------------
		case <-terminate:
			fmt.Println("Generation ended.")
			wg.Done()
			return
		}
	}
}

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

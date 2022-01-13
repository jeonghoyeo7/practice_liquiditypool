package test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
	"uniswap_test/liquidity"
)

func SimpleTest() {
	// create a liquidity pool
	lp := liquidity.CreatePool()
	fmt.Println("A liquidity pool is created.")

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
	initialLP := lp
	// Create Transaction Queue
	//tq := liquidity.CreateTransactionQueue()
	//fmt.Println("Transaction Queue is generated.")

	// Create transaction
	newTran1 := createRandomTrans("swapXtoY", int64(randInt(1, 100)), int64(0), int64(0))
	//tq.Push(newTran1)
	fmt.Println("A transaction for swap(X->Y) is randomly given.")
	newTran1.Print() // Print

	newTran1 = processTransaction(newTran1, &lp)
	fmt.Println("The 1st transaction for swap(X->Y) is processed.")
	newTran1.PrintReceipt() // Print

	newTran2 := createRandomTrans("swapYtoX", int64(0), int64(0), int64(0))
	newTran2.CoinY = newTran1.RemainCoinY
	//tq.Push(newTran2)
	//fmt.Printf("Now we have %d tran in the queue.\n", tq.Len())
	fmt.Println("A transaction for swap(Y->X) is given.")
	newTran2.Print() // Print

	newTran2 = processTransaction(newTran2, &lp)
	fmt.Println("The 2nd transaction for swap(Y->X) is processed.")
	newTran2.PrintReceipt() // Print

	tempX := int64(randInt(1, 100))
	time.Sleep(2 * time.Nanosecond)
	tempY := int64(randInt(1, 100))
	newTran3 := createRandomTrans("deposit", tempX, tempY, int64(0))
	//tq.Push(newTran3)
	//fmt.Printf("Now we have %d tran in the queue.\n", tq.Len())
	fmt.Println("A transaction for deposit is given.")
	newTran3.Print() // Print

	newTran3 = processTransaction(newTran3, &lp)
	fmt.Println("The 3rd transaction for deposit is processed.")
	newTran3.PrintReceipt() // Print

	newTran4 := createRandomTrans("withdraw", int64(0), int64(0), int64(0))
	newTran4.LiquidityToken.Amount = newTran3.RemainLiquidityToken.Amount
	fmt.Println("A transaction for withdraw is given.")
	newTran4.Print() // Print
	//tq.Push(newTran4)
	//fmt.Printf("Now we have %d tran in the queue.\n", tq.Len())

	newTran4 = processTransaction(newTran4, &lp)
	fmt.Println("The 4th transaction for deposit is processed.")
	newTran4.PrintReceipt() // Print

	// Create Processed Queue
	//result := liquidity.CreateTransactionQueue()

	// print last lp status
	fmt.Println("Initial Liquidity Pool................")
	initialLP.Print()
	fmt.Println("Recent Liquidity Pool.................")
	lp.Print()
}

func createRandomTrans(trade string, x int64, y int64, lt int64) liquidity.Transaction {
	var newTran liquidity.Transaction
	newTran.Init()

	var newAmountX, newAmountY, newAmountLT sdk.Dec
	var newCoinX, newCoinY, newCoinLT sdk.DecCoin
	// 1: swap x to y
	newTran.Order = trade
	newAmountX = sdk.NewDec(x)
	newCoinX = sdk.NewDecCoinFromDec("uatom", newAmountX)

	newAmountY = sdk.NewDec(y)
	newCoinY = sdk.NewDecCoinFromDec("eth", newAmountY)

	newAmountLT = sdk.NewDec(lt)
	newCoinLT = sdk.NewDecCoinFromDec("ltoken", newAmountLT)

	newTran.CoinX = newCoinX
	newTran.CoinY = newCoinY
	newTran.LiquidityToken = newCoinLT

	return newTran
}

func processQueue(tq *liquidity.TransactionQueue, lp *liquidity.Pool) {
	var ProcessingTrans liquidity.Transaction

	ProcessingTrans = tq.Pop().(liquidity.Transaction)

	lp.Trade(&ProcessingTrans)
}

func processTransaction(ProcessingTrans liquidity.Transaction, lp *liquidity.Pool) liquidity.Transaction {
	lp.Trade(&ProcessingTrans)

	return ProcessingTrans
}

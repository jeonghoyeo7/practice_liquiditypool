package test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
	"uniswap_test/liquidity"
)

func Simpletest() {
	// create a liquidity pool
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

	// Original X,Y,K,Price,Liquidity Tokens
	var oriX, oriY, oriK, oriP, oriLT sdk.Dec
	oriX = lp.CoinX.Amount
	oriY = lp.CoinY.Amount
	oriK = lp.K
	oriP = lp.Price
	oriLT = lp.TotalLiquidityToken.Amount

	// Print state of the pool
	fmt.Println(lp)

	// Create Transaction Queue
	tq := liquidity.CreateTransactionQueue()

	// Create transaction
	newTran1 := createRandomTrans("swapXtoY", int64(randInt(1, 100)), int64(0), int64(0))
	tq.Push(newTran1)
	//fmt.Printf("Now we have %d tran in the queue.\n", tq.Len())

	fmt.Printf("\n\n The 1st transaction is starting.\n")
	newTran1 = processTransaction(newTran1, &lp)
	fmt.Printf("The 1st transaction is finished.\n")

	newTran2 := createRandomTrans("swapYtoX", int64(0), int64(0), int64(0))
	newTran2.CoinY = newTran1.RemainCoinY
	tq.Push(newTran2)
	//fmt.Printf("Now we have %d tran in the queue.\n", tq.Len())

	fmt.Printf("\n\n The 2nd transaction is starting.\n")
	newTran2 = processTransaction(newTran2, &lp)
	fmt.Printf("The 2nd transaction is finished.\n")

	fmt.Println(oriX, oriY, oriK, oriP, oriLT)
	fmt.Printf("Comparison... \n")
	fmt.Println(lp.CoinX.Amount, lp.CoinY.Amount, lp.K, lp.Price, lp.TotalLiquidityToken.Amount)
	if oriX.Equal(lp.CoinX.Amount) && oriY.Equal(lp.CoinY.Amount) && oriK.Equal(lp.K) && oriP.Equal(lp.Price) && oriLT.Equal(lp.TotalLiquidityToken.Amount) {
		fmt.Println("Same!!")
	} else {
		fmt.Println("Not Same.")
	}

	tempX := int64(randInt(1, 100))
	time.Sleep(2 * time.Nanosecond)
	tempY := int64(randInt(1, 100))
	newTran3 := createRandomTrans("deposit", tempX, tempY, int64(0))
	tq.Push(newTran3)
	//fmt.Printf("Now we have %d tran in the queue.\n", tq.Len())

	fmt.Printf("\n\n The 3rd transaction is starting.\n")
	newTran3 = processTransaction(newTran3, &lp)
	fmt.Printf("The 3rd transaction is finished.\n")

	newTran4 := createRandomTrans("withdraw", int64(0), int64(0), int64(0))
	newTran4.LiquidityToken.Amount = newTran3.RemainLiquidityToken.Amount
	tq.Push(newTran4)
	//fmt.Printf("Now we have %d tran in the queue.\n", tq.Len())

	fmt.Println(lp.TotalLiquidityToken.Amount)
	fmt.Printf("\n\n The 4th transaction is starting.\n")
	newTran4 = processTransaction(newTran4, &lp)
	fmt.Printf("The 4th transaction is finished.\n")

	fmt.Println(oriX, oriY, oriK, oriP, oriLT)
	fmt.Printf("Comparison... \n")
	fmt.Println(lp.CoinX.Amount, lp.CoinY.Amount, lp.K, lp.Price, lp.TotalLiquidityToken.Amount)
	if oriX.Equal(lp.CoinX.Amount) && oriY.Equal(lp.CoinY.Amount) && oriK.Equal(lp.K) && oriP.Equal(lp.Price) && oriLT.Equal(lp.TotalLiquidityToken.Amount) {
		fmt.Println("Same!!")
	} else {
		fmt.Println("Not Same.")
	}

	// Create Processed Queue
	//result := liquidity.CreateTransactionQueue()

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

	templp := lp
	ProcessingTrans = tq.Pop().(liquidity.Transaction)

	fmt.Println("Before Liquidity Pool.")
	fmt.Println(*templp)
	fmt.Println("Transaction to be processed")
	fmt.Println(ProcessingTrans)

	lp.Trade(&ProcessingTrans)

	fmt.Println("Updated Liquidity Pool.")
	fmt.Println(*lp)
	fmt.Println("Receipt for the transaction.")
	fmt.Println(ProcessingTrans)
}

func processTransaction(ProcessingTrans liquidity.Transaction, lp *liquidity.Pool) liquidity.Transaction {
	templp := lp

	fmt.Println("Before Liquidity Pool.")
	fmt.Println(*templp)
	fmt.Println("Transaction to be processed")
	fmt.Println(ProcessingTrans)

	lp.Trade(&ProcessingTrans)

	fmt.Println("Updated Liquidity Pool.")
	fmt.Println(*lp)
	fmt.Println("Receipt for the transaction.")
	fmt.Println(ProcessingTrans)

	return ProcessingTrans
}

package test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
	"uniswap_test/liquidity"
)

var AvgArrival time.Duration = 1000 // ms for new 1 transaction arrival

func generateTransaction() liquidity.Transaction {
	var newTran liquidity.Transaction
	newTran.Init()
	var newCoinX, newCoinY sdk.DecCoin
	var newAmountX, newAmountY sdk.Dec

	switch randInt(1, 4) {
	case 1:
		newTran.Order = "swapXtoY"
		newAmountX = sdk.NewDec(0).Add(sdk.NewDec(int64(randInt(1, 100))))
		newCoinX = sdk.NewDecCoinFromDec("uatom", newAmountX)

		newAmountY = sdk.NewDec(0).Add(sdk.NewDec(0))
		newCoinY = sdk.NewDecCoinFromDec("eth", newAmountY)
	case 2:
		newTran.Order = "swapYtoX"
		newAmountX = sdk.NewDec(0).Add(sdk.NewDec(0))
		newCoinX = sdk.NewDecCoinFromDec("uatom", newAmountX)

		newAmountY = sdk.NewDec(0).Add(sdk.NewDec(int64(randInt(1, 100))))
		newCoinY = sdk.NewDecCoinFromDec("eth", newAmountY)
	case 3:
		newTran.Order = "deposit"
		newAmountX = sdk.NewDec(0).Add(sdk.NewDec(int64(randInt(1, 100))))
		newCoinX = sdk.NewDecCoinFromDec("uatom", newAmountX)

		newAmountY = sdk.NewDec(0).Add(sdk.NewDec(0))
		newCoinY = sdk.NewDecCoinFromDec("eth", newAmountY)
	case 4:
		newTran.Order = "withdraw"
		newAmountX = sdk.NewDec(0).Add(sdk.NewDec(int64(randInt(1, 100))))
		newCoinX = sdk.NewDecCoinFromDec("uatom", newAmountX)

		newAmountY = sdk.NewDec(0).Add(sdk.NewDec(0))
		newCoinY = sdk.NewDecCoinFromDec("eth", newAmountY)

	}

	newTran.CoinX = newCoinX
	newTran.CoinY = newCoinY

	fmt.Println(newTran)
	return newTran
}

func GenerateTransactions(tq *liquidity.TransactionQueue) {

	for {
		tq.Push(generateTransaction())
		fmt.Printf("Now we have %d tran in the queue.\n", tq.Len())
		time.Sleep(time.Millisecond * AvgArrival) // wait until the next transaction is coming
	}
}

func ProcessBlocks(tq *liquidity.TransactionQueue) {

}

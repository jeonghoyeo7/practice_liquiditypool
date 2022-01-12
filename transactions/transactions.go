package transactions

import (
	"container/list"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/rand"
	"time"
)

var AvgArrival time.Duration = 1000 // ms for new 1 transaction arrival

type LiquidityPool struct {
	CoinX sdk.Coin // the first coin of the coin pair in this LP, called X
	CoinY sdk.Coin // the second coin of the coin pair in this LP, called Y
	Price float32  // pool price: X/Y
	K     float32  // Constant for constant product market maker
	Fee   float32  // swap fee. assumed as 0 for this test
}

type TransactionQueue struct {
	v *list.List
}

type Transaction struct {
	Order string   // one of "swapXtoY","swapYtoX", "deposit", "withdraw"
	CoinX sdk.Coin // the first coin of the coin pair in this LP, called X
	CoinY sdk.Coin // the second coin of the coin pair in this LP, called Y
}

func (tq *TransactionQueue) Push(val interface{}) {
	tq.v.PushBack(val)
}

func (tq *TransactionQueue) Pop() interface{} {
	front := tq.v.Front()
	if front != nil {
		return tq.v.Remove(front)
	}
	return nil
}

func (tq *TransactionQueue) Len() interface{} {
	return tq.v.Len()
}

func generateTransaction() Transaction {
	var newTran Transaction
	var newCoinX, newCoinY sdk.Coin
	var newAmountX, newAmountY sdk.Int

	switch randInt(1, 4) {
	case 1:
		newTran.Order = "swapXtoY"
		newAmountX = sdk.NewInt(0).Add(sdk.NewInt(int64(randInt(1, 100))))
		newCoinX = sdk.NewCoin("uatom", newAmountX)

		newAmountY = sdk.NewInt(0).Add(sdk.NewInt(0))
		newCoinY = sdk.NewCoin("eth", newAmountY)
	case 2:
		newTran.Order = "swapYtoX"
		newAmountX = sdk.NewInt(0).Add(sdk.NewInt(0))
		newCoinX = sdk.NewCoin("uatom", newAmountX)

		newAmountY = sdk.NewInt(0).Add(sdk.NewInt(int64(randInt(1, 100))))
		newCoinY = sdk.NewCoin("eth", newAmountY)
	case 3:
		newTran.Order = "deposit"
		newAmountX = sdk.NewInt(0).Add(sdk.NewInt(int64(randInt(1, 100))))
		newCoinX = sdk.NewCoin("uatom", newAmountX)

		newAmountY = sdk.NewInt(0).Add(sdk.NewInt(0))
		newCoinY = sdk.NewCoin("eth", newAmountY)
	case 4:
		newTran.Order = "withdraw"
		newAmountX = sdk.NewInt(0).Add(sdk.NewInt(int64(randInt(1, 100))))
		newCoinX = sdk.NewCoin("uatom", newAmountX)

		newAmountY = sdk.NewInt(0).Add(sdk.NewInt(0))
		newCoinY = sdk.NewCoin("eth", newAmountY)

	}

	newTran.CoinX = newCoinX
	newTran.CoinY = newCoinY

	fmt.Println(newTran)
	return newTran
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func NewTransactionQueue() *TransactionQueue {
	return &TransactionQueue{list.New()}
}

func GenerateTransactions(tq *TransactionQueue) {

	for {
		tq.Push(generateTransaction())
		fmt.Printf("Now we have %d transactions in the queue.\n", tq.Len())
		time.Sleep(time.Millisecond * AvgArrival) // wait until the next transaction is coming
	}
}

func ProcessBlocks(tq *TransactionQueue) {
	
}

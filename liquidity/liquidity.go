package liquidity

import (
	"container/list"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/rand"
	"time"
)

var AvgArrival time.Duration = 1000 // ms for new 1 transaction arrival

type Pool struct {
	CoinX          sdk.DecCoin // the first coin of the coin pair in this LP, called X
	CoinY          sdk.DecCoin // the second coin of the coin pair in this LP, called Y
	Price          sdk.Dec     // pool price: X/Y
	K              sdk.Dec     // Constant for constant product market maker
	LiquidityToken sdk.DecCoin // amount of Liquidity Token here
	Fee            sdk.Dec     // swap fee. assumed as 0 for this test
}

func (lp *Pool) Init() {
	lp.CoinX.Amount = sdk.NewDec(int64(randInt(1, 10000)))
	lp.CoinX.Denom = "uatom"

	lp.CoinY.Amount = sdk.NewDec(int64(randInt(1, 10000)))
	lp.CoinY.Denom = "eth"

	lp.Price = lp.CoinX.Amount.Quo(lp.CoinY.Amount) // X/Y
	lp.K = lp.CoinX.Amount.Mul(lp.CoinY.Amount)     // X*Y

	lp.LiquidityToken.Denom = "ltoken"
	lp.LiquidityToken.Amount, _ = lp.K.ApproxSqrt()

	lp.Fee = sdk.NewDec(int64(0))
}

func (lp *Pool) Swap(transaction Transaction) Transaction {
	var receipt Transaction
	receipt.Init()

	// update the Pool state
	// ...

	// receipt indicates the results of the transaction and the remaining amount of tokens
	return receipt
}

func (lp *Pool) Deposit(transaction Transaction) Transaction {
	var receipt Transaction
	receipt.Init()

	return receipt
}

func (lp *Pool) Withdraw(transaction Transaction) Transaction {
	var receipt Transaction
	receipt.Init()

	return receipt
}

type TransactionQueue struct {
	v *list.List
}

type Transaction struct {
	Order          string      // one of "swapXtoY","swapYtoX", "deposit", "withdraw"
	CoinX          sdk.DecCoin // the first coin of the coin pair in this LP, called X
	CoinY          sdk.DecCoin // the second coin of the coin pair in this LP, called Y
	LiquidityToken sdk.DecCoin // the liquidity token
	Result         string
}

func (tran *Transaction) Init() {
	tran.Order = "none"
	tran.CoinX.Amount = sdk.NewDec(int64(0))
	tran.CoinX.Denom = "uatom"
	tran.CoinY.Amount = sdk.NewDec(int64(0))
	tran.CoinY.Denom = "eth"
	tran.LiquidityToken.Amount = sdk.NewDec(int64(0))
	tran.LiquidityToken.Denom = "ltoken"
	tran.Result = "beforeTransaction"
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
		fmt.Printf("Now we have %d tran in the queue.\n", tq.Len())
		time.Sleep(time.Millisecond * AvgArrival) // wait until the next transaction is coming
	}
}

func ProcessBlocks(tq *TransactionQueue) {

}

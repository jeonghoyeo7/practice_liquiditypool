package liquidity

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Pool struct {
	CoinX                 sdk.DecCoin // the first coin of the coin pair in this LP, called X
	CoinY                 sdk.DecCoin // the second coin of the coin pair in this LP, called Y
	Price                 sdk.Dec     // pool price: X/Y
	K                     sdk.Dec     // Constant for constant product market maker
	InitialLiquidityToken sdk.DecCoin // amount of Liquidity Token in initial state
	TotalLiquidityToken   sdk.DecCoin // amount of Liquidity Token here
	Fee                   sdk.Dec     // swap fee. assumed as 0 for this test
}

func (lp Pool) Print() {
	fmt.Println("<Current Liquidity Pool Status>")
	fmt.Printf("Coin X = %f \n", lp.CoinX.Amount)
	fmt.Printf("Coin Y = %f \n", lp.CoinY.Amount)
	fmt.Printf("Price  = %f \n", lp.Price)
	fmt.Printf("K      = %f \n", lp.K)
	fmt.Printf("Fee    = %f \n", lp.Fee)
	fmt.Println("-------------------------------")
}

func CreatePool() Pool {
	var lp Pool
	lp.Init()
	return lp
}

func (lp *Pool) Init() {
	lp.CoinX.Amount = sdk.NewDec(int64(0))
	lp.CoinX.Denom = "none"

	lp.CoinY.Amount = sdk.NewDec(int64(0))
	lp.CoinY.Denom = "none"

	lp.Price = sdk.NewDec(int64(0)) // X/Y
	lp.K = sdk.NewDec(int64(0))     // X*Y

	lp.InitialLiquidityToken.Denom = "Initial ltoken"
	lp.InitialLiquidityToken.Amount = sdk.NewDec(int64(0))
	lp.TotalLiquidityToken.Denom = "Current Total ltoken"
	lp.TotalLiquidityToken.Amount = sdk.NewDec(int64(0))
	lp.Fee = sdk.NewDec(int64(0))
}

func (lp *Pool) SetOnce(CoinX sdk.DecCoin, CoinY sdk.DecCoin, Fee sdk.Dec) {
	lp.CoinX = CoinX
	lp.CoinY = CoinY

	lp.Fee = Fee

	lp.UpdateK() // update K and prices

	lp.InitialLiquidityToken.Denom = "ltoken"
	lp.InitialLiquidityToken.Amount, _ = lp.K.ApproxSqrt()

	lp.TotalLiquidityToken = lp.InitialLiquidityToken
}

func (lp *Pool) UpdateK() {
	lp.Price = lp.CoinX.Amount.Quo(lp.CoinY.Amount)      // X/Y
	lp.K = lp.CoinX.Amount.Mul(lp.CoinY.Amount)          // X*Y
	lp.TotalLiquidityToken.Amount, _ = lp.K.ApproxSqrt() // sqrt(K)
}

func (lp *Pool) Trade(transaction *Transaction) {
	switch transaction.Order {
	case "swapXtoY":
		lp.SwapXtoY(transaction)
	case "swapYtoX":
		lp.SwapYtoX(transaction)
	case "deposit":
		lp.Deposit(transaction)
	case "withdraw":
		lp.Withdraw(transaction)
	default:
		fmt.Println("Trading Error")
	}
}

func (lp *Pool) SwapXtoY(transaction *Transaction) {
	// The user gives dx to the pool and receives dy from the pool
	// X*Y=K --> (X+dx)(Y-dy)=K : dy
	// no change of LT but change to the prices

	lp.CoinX.Amount = lp.CoinX.Amount.Add(transaction.CoinX.Amount) // X' <-- X+dx

	newY := lp.K.Quo(lp.CoinX.Amount) // Y' <-- Y-dy = K/X'
	dy := lp.CoinY.Amount.Sub(newY)
	lp.CoinY.Amount = newY

	lp.UpdateK() // update K and prices

	// Update Transaction
	transaction.Result = "Swap X to Y Completed Successfully"
	transaction.RemainCoinX.Amount = sdk.NewDec(int64(0)) // The user gives dx to the pool
	transaction.RemainCoinY.Amount = dy                   // The user gets dy from the pool
}

func (lp *Pool) SwapYtoX(transaction *Transaction) {
	// The user gives dy to the pool and receives dx from the pool
	// X*Y=K --> (X-dx)(Y+dy)=K
	// no change of LT but change to the prices

	lp.CoinY.Amount = lp.CoinY.Amount.Add(transaction.CoinY.Amount) // Y' <-- Y+dy

	newX := lp.K.Quo(lp.CoinY.Amount) // X' <-- X-dx = K/Y'
	dx := lp.CoinX.Amount.Sub(newX)
	lp.CoinX.Amount = newX

	lp.UpdateK() // update K and prices

	// Update Transaction
	transaction.Result = "Swap Y to X Completed Successfully"
	transaction.RemainCoinY.Amount = sdk.NewDec(int64(0)) // The user gives dy to the pool
	transaction.RemainCoinX.Amount = dx                   // The user gets dx from the pool
}

func (lp *Pool) Deposit(transaction *Transaction) {
	// The user gives the pair of (dx, dy) to the pool and receives d(LT) from the pool
	// X*Y=K --> (X+dx)(Y+dy)=K'
	// no price change
	// new L tokens are created and given to the user.
	var dx, dy sdk.Dec

	tempPrice := transaction.CoinX.Amount.Quo(transaction.CoinY.Amount) // dx/dy

	if tempPrice.GT(lp.Price) { // if dx/dy > Price -> too much dx : dx'/dy = Price --> dx' = Price*dy
		dy = transaction.CoinY.Amount
		dx = lp.Price.Mul(dy)
	} else { // if dx/dy <= Price -> too much dy : dx/dy' = Price --> dy' = dx/Price
		dx = transaction.CoinX.Amount
		//dy = dx.Quo(lp.Price) //lp.CoinX.Amount.Quo(lp.CoinY.Amount)
		tempdy := dx.Mul(lp.CoinY.Amount)
		dy = tempdy.Quo(lp.CoinX.Amount)
	}

	beforeLiquidityToken := lp.TotalLiquidityToken.Amount
	lp.CoinX.Amount = lp.CoinX.Amount.Add(dx)
	lp.CoinY.Amount = lp.CoinY.Amount.Add(dy)
	lp.UpdateK() // update K and prices

	// update receipt in below
	transaction.RemainLiquidityToken.Amount = lp.TotalLiquidityToken.Amount.Sub(beforeLiquidityToken) // the

	transaction.Result = "Deposit Completed Successfully"
	transaction.RemainCoinX.Amount = transaction.CoinX.Amount.Sub(dx)
	transaction.RemainCoinY.Amount = transaction.CoinY.Amount.Sub(dy)

}

func (lp *Pool) Withdraw(transaction *Transaction) {
	// The user gives d(LT) to the pool and receives the pair of (dx, dy) from the pool.
	// X*Y=K --> (X-dx)(Y-dy)=K'
	// The user gives L tokens to the pool and the L tokens are burnt.
	var dx, dy sdk.Dec

	ratio := transaction.LiquidityToken.Amount.Quo(lp.TotalLiquidityToken.Amount)

	if ratio.GT(sdk.NewDec(int64(1))) { // not possible
		fmt.Println("Trading Error: Too much liquidity tokens entered.")
	} else {
		//dx = ratio.Mul(lp.CoinX.Amount)
		tempdx := transaction.LiquidityToken.Amount.Mul(lp.CoinX.Amount)
		dx = tempdx.Quo(lp.TotalLiquidityToken.Amount)
		//dy = ratio.Mul(lp.CoinY.Amount)
		tempdy := transaction.LiquidityToken.Amount.Mul(lp.CoinY.Amount)
		dy = tempdy.Quo(lp.TotalLiquidityToken.Amount)
	}

	beforeLiquidityToken := lp.TotalLiquidityToken.Amount
	lp.CoinX.Amount = lp.CoinX.Amount.Sub(dx)
	lp.CoinY.Amount = lp.CoinY.Amount.Sub(dy)
	lp.UpdateK() // update K and prices

	// update receipt in below

	transaction.RemainLiquidityToken.Amount = transaction.LiquidityToken.Amount.Sub(beforeLiquidityToken.Sub(lp.TotalLiquidityToken.Amount)) // the
	// This should be zero.

	transaction.Result = "Withdraw Completed Successfully "
	transaction.RemainCoinX.Amount = dx
	transaction.RemainCoinY.Amount = dy
}

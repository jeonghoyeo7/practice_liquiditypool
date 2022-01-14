package test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/rand"
	"time"
	"uniswap_test/liquidity"
)

func ArbitrageTest(simulateTime int) {
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

	// Reference market
	maxChange := 0.03 // pool price changes at most +- 3%
	referPrice := lp.Price
	fmt.Println("Current reference market price: ", referPrice)

	// Creat Wallet
	var wallet Wallet
	wallet.Init()
	wallet.Set(1000, 1000)
	wallet.ComparePrint()

	for i := 0; i < simulateTime; i++ {
		// ideal reference market price tracker
		rand.Seed(time.Now().UTC().UnixNano())
		changeRatio := (rand.Float64() - 0.5) * 2.0 * maxChange // -1.0 <= f < 1.0
		changeRatioDec := sdk.NewDec(int64(changeRatio * 1000.0)).Quo(sdk.NewDec(int64(1000)))
		referPrice = referPrice.Add(referPrice.Mul(changeRatioDec))

		fmt.Println("Trade Starting.")
		fmt.Println(" >> Current reference market price: ", referPrice)
		fmt.Println(" >> Current liquidity pool price: ", lp.Price)

		if lp.Price.GT(referPrice) {
			resultTrade := wallet.arbitrageTrade(&lp, referPrice)
			fmt.Printf("Swap Y(%14.6f) to X(%14.6f) in the pool. \n", liquidity.DecToFloat64(resultTrade[1]), liquidity.DecToFloat64(resultTrade[0]))
			fmt.Printf("Swap X(%14.6f) to Y(%14.6f) in the ideal market. \n", liquidity.DecToFloat64(resultTrade[0]), liquidity.DecToFloat64(resultTrade[3]))
		} else {
			resultTrade := wallet.arbitrageTrade(&lp, referPrice)
			fmt.Printf("Swap X(%14.6f) to Y(%14.6f) in the pool. \n", liquidity.DecToFloat64(resultTrade[0]), liquidity.DecToFloat64(resultTrade[1]))
			fmt.Printf("Swap Y(%14.6f) to X(%14.6f) in the ideal market. \n", liquidity.DecToFloat64(resultTrade[1]), liquidity.DecToFloat64(resultTrade[2]))
		}

		fmt.Println("  >>>>  Current reference market price: ", referPrice)
		fmt.Println("  >>>>  Current liquidity pool price: ", lp.Price)
		fmt.Println("Trade Completed.")
		fmt.Printf("Current wallet has %s.\n", wallet.SprintLine())
		time.Sleep(time.Second)

	}

}

func (wallet *Wallet) arbitrageTrade(lp *liquidity.Pool, referPrice sdk.Dec) [4]sdk.Dec {
	var transaction liquidity.Transaction
	var dx, dy, dxRE, dyRE sdk.Dec

	if lp.Price.GT(referPrice) { // LP price > reference market
		// buy dX (swap Y to X): dy = sqrt(K/refPrice) - Y
		temp := lp.K.Quo(referPrice)
		temp, _ = temp.ApproxSqrt()
		dy = temp.Sub(lp.CoinY.Amount)

		// if CoinY in the wallet is not enough
		if dy.GT(wallet.CoinY.Amount) {
			dy = wallet.CoinY.Amount
		}
		//fmt.Println("dy: ", dy)

		// Create Transaction
		transaction.Init()
		transaction.Order = "swapYtoX"
		transaction.CoinY.Amount = dy

		lp.SwapYtoX(&transaction) // Trade
		dx = transaction.RemainCoinX.Amount
		//fmt.Println("dx: ", dx)

		// sell dx and get dy' in ideal market
		dyRE = dx.Quo(referPrice)
		wallet.CoinY.Amount = wallet.CoinY.Amount.Sub(dy)   //withdraw dY from the wallet
		wallet.CoinY.Amount = wallet.CoinY.Amount.Add(dyRE) //withdraw dY from the wallet
		dxRE = sdk.NewDec(int64(0))
	} else {
		// sell dX: dx = X-sqrt(refPrice*K)
		dx = lp.CoinX.Amount
		temp := referPrice.Mul(lp.K)
		temp, _ = temp.ApproxSqrt()
		dx = temp.Sub(dx)

		// if CoinX in the wallet is not enough
		if dx.GT(wallet.CoinX.Amount) {
			dx = wallet.CoinX.Amount
		}

		// Create Transaction
		transaction.Init()
		transaction.Order = "swapXtoY"
		transaction.CoinX.Amount = dx

		lp.SwapXtoY(&transaction)
		dy = transaction.RemainCoinY.Amount

		// buy dx' and sell dy in ideal market
		dxRE = dy.Mul(referPrice)
		wallet.CoinX.Amount = wallet.CoinX.Amount.Sub(dx)   //withdraw X from the wallet
		wallet.CoinX.Amount = wallet.CoinX.Amount.Add(dxRE) //withdraw X from the wallet
		dyRE = sdk.NewDec(int64(0))
	}

	output := [4]sdk.Dec{dx, dy, dxRE, dyRE}
	return output
}

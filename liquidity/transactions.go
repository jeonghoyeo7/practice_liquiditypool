package liquidity

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Transaction struct {
	Order          string      // one of "swapXtoY","swapYtoX", "deposit", "withdraw"
	CoinX          sdk.DecCoin // the first coin of the coin pair in this LP, called X
	CoinY          sdk.DecCoin // the second coin of the coin pair in this LP, called Y
	LiquidityToken sdk.DecCoin // the liquidity token

	// save result of processing
	Result               string
	RemainCoinX          sdk.DecCoin
	RemainCoinY          sdk.DecCoin
	RemainLiquidityToken sdk.DecCoin
}

func (tran Transaction) Print() {
	fmt.Printf("<Current %s Transaction Status>\n", tran.Order)

	switch tran.Order {
	case "swapXtoY":
		fmt.Printf("Input Coin X = %f \n", tran.CoinX.Amount)
	case "swapYtoX":
		fmt.Printf("Input Coin Y = %f \n", tran.CoinY.Amount)
	case "deposit":
		fmt.Printf("Input Coin X = %f \n", tran.CoinX.Amount)
		fmt.Printf("Input Coin Y = %f \n", tran.CoinY.Amount)
	case "withdraw":
		fmt.Printf("Input Liquidity Token = %f \n", tran.LiquidityToken.Amount)
	default:

	}

	fmt.Println("-----------------------------------")
}

func (tran Transaction) PrintReceipt() {
	fmt.Printf("<Current %s Transaction Receipt>\n", tran.Order)
	fmt.Println(tran.Result)
	switch tran.Order {
	case "swapXtoY":
		fmt.Printf("Input Coin X  = %f \n", tran.CoinX.Amount)
		fmt.Printf("Output Coin Y = %f \n", tran.RemainCoinY.Amount)
	case "swapYtoX":
		fmt.Printf("Input Coin Y  = %f \n", tran.CoinY.Amount)
		fmt.Printf("Output Coin X = %f \n", tran.RemainCoinX.Amount)
	case "deposit":
		fmt.Printf("Input Coin X      = %f \n", tran.CoinX.Amount)
		fmt.Printf("Input Coin Y      = %f \n", tran.CoinY.Amount)
		fmt.Printf("Output Liq. Token = %f \n", tran.RemainLiquidityToken.Amount)
	case "withdraw":
		fmt.Printf("Input Liq. Token = %f \n", tran.LiquidityToken.Amount)
		fmt.Printf("Output Coin X    = %f \n", tran.RemainCoinX.Amount)
		fmt.Printf("Output Coin Y    = %f \n", tran.RemainCoinY.Amount)
	default:

	}

	fmt.Println("-----------------------------------")
}

func GenerateTransaction() *Transaction {
	var tran *Transaction
	tran.Init()
	return tran
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
	tran.RemainCoinX.Amount = sdk.NewDec(int64(0))
	tran.RemainCoinX.Denom = "uatom"
	tran.RemainCoinY.Amount = sdk.NewDec(int64(0))
	tran.RemainCoinY.Denom = "eth"
	tran.RemainLiquidityToken.Amount = sdk.NewDec(int64(0))
	tran.RemainLiquidityToken.Denom = "ltoken"
}

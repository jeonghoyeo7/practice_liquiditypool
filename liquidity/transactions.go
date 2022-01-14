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

func GenerateTransaction() *Transaction {
	var tran *Transaction
	tran.Init()
	return tran
}

func (tran *Transaction) Init() {
	tran.Order = "none"
	tran.CoinX.Amount = sdk.NewDec(int64(0))
	tran.CoinX.Denom = "atom"
	tran.CoinY.Amount = sdk.NewDec(int64(0))
	tran.CoinY.Denom = "eth"
	tran.LiquidityToken.Amount = sdk.NewDec(int64(0))
	tran.LiquidityToken.Denom = "ltoken"

	tran.Result = "beforeTransaction"
	tran.RemainCoinX.Amount = sdk.NewDec(int64(0))
	tran.RemainCoinX.Denom = "atom"
	tran.RemainCoinY.Amount = sdk.NewDec(int64(0))
	tran.RemainCoinY.Denom = "eth"
	tran.RemainLiquidityToken.Amount = sdk.NewDec(int64(0))
	tran.RemainLiquidityToken.Denom = "ltoken"
}

func (tran Transaction) Print() {
	fmt.Printf("   <Current %s Transaction Status>\n", tran.Order)

	switch tran.Order {
	case "swapXtoY":
		fmt.Printf("     Input Coin X  = %14.6f \n", DecToFloat64(tran.CoinX.Amount))
	case "swapYtoX":
		fmt.Printf("     Input Coin Y  = %14.6f \n", DecToFloat64(tran.CoinY.Amount))
	case "deposit":
		fmt.Printf("     Input Coin X  = %14.6f \n", DecToFloat64(tran.CoinX.Amount))
		fmt.Printf("     Input Coin Y  = %14.6f \n", DecToFloat64(tran.CoinY.Amount))
	case "withdraw":
		fmt.Printf("     Input Liq. Token = %14.6f \n", DecToFloat64(tran.LiquidityToken.Amount))
	default:

	}

	fmt.Println("-----------------------------------------------------------------")
}

func (tran Transaction) PrintReceipt() {
	fmt.Printf("   <Current %s Transaction Receipt>\n", tran.Order)
	fmt.Println("   >> ", tran.Result)
	switch tran.Order {
	case "swapXtoY":
		fmt.Printf("     Input Coin X  = %14.6f \n", DecToFloat64(tran.CoinX.Amount))
		fmt.Printf("     Output Coin Y = %14.6f \n", DecToFloat64(tran.RemainCoinY.Amount))
	case "swapYtoX":
		fmt.Printf("     Input Coin Y  = %14.6f \n", DecToFloat64(tran.CoinY.Amount))
		fmt.Printf("     Output Coin X = %14.6f \n", DecToFloat64(tran.RemainCoinX.Amount))
	case "deposit":
		fmt.Printf("     Input Coin X  = %14.6f \n", DecToFloat64(tran.CoinX.Amount))
		fmt.Printf("     Input Coin Y  = %14.6f \n", DecToFloat64(tran.CoinY.Amount))
		fmt.Printf("     Output Coin X = %14.6f \n", DecToFloat64(tran.RemainCoinX.Amount))
		fmt.Printf("     Output Coin Y = %14.6f \n", DecToFloat64(tran.RemainCoinY.Amount))
		fmt.Printf("     Output Liq. Token = %14.6f \n", DecToFloat64(tran.RemainLiquidityToken.Amount))
	case "withdraw":
		fmt.Printf("     Input Liq. Token = %14.6f \n", DecToFloat64(tran.LiquidityToken.Amount))
		fmt.Printf("     Output Coin X    = %14.6f \n", DecToFloat64(tran.RemainCoinX.Amount))
		fmt.Printf("     Output Coin Y    = %14.6f \n", DecToFloat64(tran.RemainCoinY.Amount))
	default:

	}

	fmt.Println("-----------------------------------------------------------------")
}

func (tran Transaction) SprintLine() string {
	var printOutput string

	switch tran.Order {
	case "swapXtoY":
		printOutput = fmt.Sprintf("%s, %14.6f%s", tran.Order, DecToFloat64(tran.CoinX.Amount), tran.CoinX.Denom)
	case "swapYtoX":
		printOutput = fmt.Sprintf("%s, %14.6f%s", tran.Order, DecToFloat64(tran.CoinY.Amount), tran.CoinY.Denom)
	case "deposit":
		printOutput = fmt.Sprintf("%s, %14.6f%s", tran.Order, DecToFloat64(tran.CoinX.Amount), tran.CoinX.Denom)
		printOutput = fmt.Sprintf("%s, %14.6f%s", tran.Order, DecToFloat64(tran.CoinY.Amount), tran.CoinY.Denom)
	case "withdraw":
		printOutput = fmt.Sprintf("%s, %14.6f%s", tran.Order, DecToFloat64(tran.LiquidityToken.Amount), tran.LiquidityToken.Denom)
	default:

	}

	return printOutput
}

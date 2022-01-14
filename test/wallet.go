package test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"uniswap_test/liquidity"
)

type Wallet struct {
	CoinX          sdk.DecCoin // the first coin of the coin pair in this LP, called X
	CoinY          sdk.DecCoin // the second coin of the coin pair in this LP, called Y
	LiquidityToken sdk.DecCoin // liquidity token

	InitialCoinX          sdk.DecCoin // Initial value: the first coin of the coin pair in this LP, called X
	InitialCoinY          sdk.DecCoin // Initial value: the second coin of the coin pair in this LP, called Y
	InitialLiquidityToken sdk.DecCoin // liquidity token
}

func (wallet Wallet) ComparePrint() {
	fmt.Println("<Original Wallet Status>")
	fmt.Printf("Initial Coin X = %14.6f \n", wallet.InitialCoinX.Amount)
	fmt.Printf("Initial Coin Y = %14.6f \n", wallet.InitialCoinY.Amount)

	fmt.Println("<Current Wallet Status>")
	fmt.Printf("Coin X = %14.6f \n", wallet.CoinX.Amount)
	fmt.Printf("Coin Y = %14.6f \n", wallet.CoinY.Amount)
	fmt.Println("-------------------------------")
}

func (wallet Wallet) SprintLine() string {
	return fmt.Sprintf("%14.6f%s, %14.6f%s", liquidity.DecToFloat64(wallet.CoinX.Amount), wallet.CoinX.Denom, liquidity.DecToFloat64(wallet.CoinY.Amount), wallet.CoinY.Denom)
}

func CreateWallet() Wallet {
	var wallet Wallet
	wallet.Init()
	return wallet
}

func (wallet *Wallet) Init() {
	wallet.CoinX.Amount = sdk.NewDec(int64(0))
	wallet.CoinX.Denom = "atom"

	wallet.CoinY.Amount = sdk.NewDec(int64(0))
	wallet.CoinY.Denom = "eth"

	wallet.LiquidityToken.Amount = sdk.NewDec(int64(0))
	wallet.LiquidityToken.Denom = "ltoken"

	wallet.InitialCoinX = wallet.CoinX
	wallet.InitialCoinY = wallet.CoinY
	wallet.InitialLiquidityToken = wallet.LiquidityToken
}

func (wallet *Wallet) Set(X int, Y int) {
	wallet.CoinX.Amount = sdk.NewDec(int64(X))
	wallet.CoinX.Denom = "atom"

	wallet.CoinY.Amount = sdk.NewDec(int64(Y))
	wallet.CoinY.Denom = "eth"

	wallet.InitialCoinX = wallet.CoinX
	wallet.InitialCoinY = wallet.CoinY
}

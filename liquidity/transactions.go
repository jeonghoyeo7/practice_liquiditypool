package liquidity

import (
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

package liquidity_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
	"uniswap_test/liquidity"
)

func TestPoolSwapXtoY(t *testing.T) {
	// create a liquidity pool
	lp := liquidity.CreatePool()

	// pool initial setup
	var CoinX, CoinY sdk.DecCoin
	CoinX.Amount = sdk.NewDec(20000)
	CoinX.Denom = "atom"
	CoinY.Amount = sdk.NewDec(10000)
	CoinY.Denom = "eth"
	Fee := sdk.NewDec(int64(0))
	lp.SetOnce(CoinX, CoinY, Fee)

	var tran liquidity.Transaction
	tran.Order = "swapXtoY"
	tran.CoinX = sdk.NewDecCoinFromDec("atom", sdk.NewDec(200))
	tran.CoinY = sdk.NewDecCoinFromDec("eth", sdk.NewDec(0))
	tran.LiquidityToken = sdk.NewDecCoinFromDec("ltoken", sdk.NewDec(0))

	lp.SwapXtoY(&tran)

	correctOutput := sdk.NewDec(10000).Sub(sdk.NewDec(200000000).Quo(sdk.NewDec(20000 + 200)))
	// dy = Y - Y', Y' = K/X' = (20000*10000)/(20000+200)

	if !correctOutput.Equal(tran.RemainCoinY.Amount) {
		t.Errorf("The output is %14.6f; want %14.6f.", liquidity.DecToFloat64(tran.RemainCoinY.Amount), liquidity.DecToFloat64(correctOutput))
	}
}

func TestPoolSwapYtoX(t *testing.T) {
	// create a liquidity pool
	lp := liquidity.CreatePool()

	// pool initial setup
	var CoinX, CoinY sdk.DecCoin
	CoinX.Amount = sdk.NewDec(20000)
	CoinX.Denom = "atom"
	CoinY.Amount = sdk.NewDec(10000)
	CoinY.Denom = "eth"
	Fee := sdk.NewDec(int64(0))
	lp.SetOnce(CoinX, CoinY, Fee)

	var tran liquidity.Transaction
	tran.Order = "swapYtoX"
	tran.CoinX = sdk.NewDecCoinFromDec("atom", sdk.NewDec(0))
	tran.CoinY = sdk.NewDecCoinFromDec("eth", sdk.NewDec(100))
	tran.LiquidityToken = sdk.NewDecCoinFromDec("ltoken", sdk.NewDec(0))

	lp.SwapYtoX(&tran)

	correctOutput := sdk.NewDec(20000).Sub(sdk.NewDec(200000000).Quo(sdk.NewDec(10000 + 100)))
	// dx = X - X', X' = K/Y' = (20000*10000)/(10000+100)

	if !correctOutput.Equal(tran.RemainCoinX.Amount) {
		t.Errorf("The output is %14.6f; want %14.6f.", liquidity.DecToFloat64(tran.RemainCoinX.Amount), liquidity.DecToFloat64(correctOutput))
	}
}

func TestPoolDeposit(t *testing.T) {
	// create a liquidity pool
	lp := liquidity.CreatePool()

	// pool initial setup
	var CoinX, CoinY sdk.DecCoin
	CoinX.Amount = sdk.NewDec(20000)
	CoinX.Denom = "atom"
	CoinY.Amount = sdk.NewDec(10000)
	CoinY.Denom = "eth"
	Fee := sdk.NewDec(int64(0))
	lp.SetOnce(CoinX, CoinY, Fee)

	var tran liquidity.Transaction
	tran.Order = "deposit"
	tran.CoinX = sdk.NewDecCoinFromDec("atom", sdk.NewDec(200))
	tran.CoinY = sdk.NewDecCoinFromDec("eth", sdk.NewDec(200))
	tran.LiquidityToken = sdk.NewDecCoinFromDec("ltoken", sdk.NewDec(0))

	lp.Deposit(&tran)

	correctOutputX := sdk.NewDec(0)
	correctOutputY := sdk.NewDec(100)

	temp1, _ := (sdk.NewDec(20000).Mul(sdk.NewDec(10000))).ApproxSqrt()
	temp2, _ := (sdk.NewDec(20200).Mul(sdk.NewDec(10100))).ApproxSqrt()
	correctOutputLT := temp2.Sub(temp1)

	if !(correctOutputX.Equal(tran.RemainCoinX.Amount) && correctOutputY.Equal(tran.RemainCoinY.Amount) && correctOutputLT.Equal(tran.RemainLiquidityToken.Amount)) {
		t.Errorf("The output X is %14.6f; want %14.6f.\n", liquidity.DecToFloat64(tran.RemainCoinX.Amount), liquidity.DecToFloat64(correctOutputX))
		t.Errorf("The output Y is %14.6f; want %14.6f.\n", liquidity.DecToFloat64(tran.RemainCoinY.Amount), liquidity.DecToFloat64(correctOutputY))
		t.Errorf("The output Liquidity token is %14.6f; want %14.6f.\n", liquidity.DecToFloat64(tran.RemainLiquidityToken.Amount), liquidity.DecToFloat64(correctOutputLT))
	}
}

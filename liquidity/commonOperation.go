package liquidity

import sdk "github.com/cosmos/cosmos-sdk/types"

func DecToFloat64(input sdk.Dec) float64 {
	var output float64
	output, _ = input.Float64()
	return output
}

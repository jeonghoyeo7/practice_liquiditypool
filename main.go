package main

import (
	"time"
	"uniswap_test/test"
)

const simulationTime time.Duration = 10 // 10 seconds simulation

func main() {
	// Test 1:
	//    - when price get back to original price → pool token balances exactly same as initial balances?
	//    - what is the relationship between future pool token balances and price change?
	//        - theoretical vs simulation result
	test.SimpleSimulation() // Please uncomment this to run Test 1

	// Observation: Almost the same. The reason why the pool token balances before/after price changes are different is because of precision in the calculations.

	// Test 2: Random generation of transactions, queueing, processing in periodic manner
	test.RandomSimulation(30, 1, 5) // Please uncomment this to run Test 2

	// Test 3: an arbitrage bot swapping on the AMM
	// Assumptions: one Liquidity Pool, one legacy market (like infinity pool),
	// Assumptions: no transactions from other users
	test.ArbitrageSimulation(100) // Please uncomment this to run Test 3

}

This is to implement liquidity pool and make transactions as a toy problem.

- batch style system (blockchain-like)
    - state change for each block
- structures and messages
    - structures to store current AMM states
    - messages for swap, deposit, withdraw
- build a simulator
    - random token pair price generator
    - an arbitrage bot swapping on the AMM
- testing
    - when price get back to original price → pool token balances exactly same as initial balances?
    - what is the relationship between future pool token balances and price change?
        - theoretical vs simulation result
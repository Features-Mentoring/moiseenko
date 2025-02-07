# slot-machine

The basic idea behind the program is that we make a distribution table based on the probability of getting one of the winning combinations in such a way that the overall RTP value is maintained.

This slot machine has three reels. Each spin is worth 1 coin. We use 4 types of symbols: A, B, C and Blank symbol, which obviously can represent any non-winning symbol.
As an initial game server setup we choose the winning combinations and the reward for them:
```
AAA - will win 5 coins (small win)
BBB - will win 10 coins (medium win)
CCC - will win 20 coins (large win)
```
We also set the probability of each type of reward. 
For example, we want small winnings to occur frequently, medium winnings to occur less frequently, and large winnings to occur infrequently.
In this way we can distribute the probability of each win within the desired RTP.

```
AAA - 50%
BBB - 30%
CCC - 15%
Total RTP 95%.
```

After that we calculate the distribution weight of each symbol on the machine wheel and calculate the probability of distribution of each symbol on the virtual wheel.

When the game starts, we randomly select a symbol on each wheel and save the result, including the combination and the resulting winnings, to the database.

Each goroutine on the client side represents a player making a certain number of spins. At the end of the game of all players, we count the resulting RTP to make sure that the value is kept within the desired range.

```
2025/02/07 11:14:43 Dist  A   0.464
2025/02/07 11:14:43 Dist  B   0.311
2025/02/07 11:14:43 Dist  C   0.196  
2025/02/07 11:14:43 Player 1 spins  200000
2025/02/07 11:14:43 Player 2 spins  200000
2025/02/07 11:14:43 Player 3 spins  200000
2025/02/07 11:14:43 Player 4 spins  200000
2025/02/07 11:14:43 Player 5 spins  200000
Total spent: 1000000 
Total spins: 1000000 
Total wins: 946880 
RTP: 0.95 
```
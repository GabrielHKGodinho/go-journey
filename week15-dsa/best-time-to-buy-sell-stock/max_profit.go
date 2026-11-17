package besttime

// maxProfit returns the best profit from one buy followed by one later sell.
// It returns 0 when no profitable trade exists.
//
// Invariant: when we evaluate day sell, buy holds the index of the
// cheapest price in prices[0:sell], so it is the best day to have bought.
func maxProfit(prices []int) int {
	buy := 0
	profit := 0

	for sell := 1; sell < len(prices); sell++ {
		profit = max(profit, prices[sell]-prices[buy])

		if prices[sell] < prices[buy] {
			buy = sell // a cheaper buy day beats every earlier one for all future sells
		}
	}

	return profit
}

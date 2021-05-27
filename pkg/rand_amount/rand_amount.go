/**
 *@Description
 *@ClassName rand_amount
 *@Date 2021/5/27 上午9:57
 *@Author ckhero
 */

package rand_amount

import (
	"math/rand"
	"time"
)

func GetRandAmount(data map[uint64]uint64) uint64 {
	var (
		weights []uint64
		amounts []uint64
		totalWeight uint64
		currTotalWeight uint64
	)
	for amount, weight := range data {
		totalWeight += weight
		weights = append(weights, weight)
		amounts = append(amounts, amount)
	}
	rand.Seed(time.Now().UnixNano())
	randVal := rand.Intn(int(totalWeight))
	for k, w := range weights {
		currTotalWeight += w
		if uint64(randVal) < currTotalWeight {
			return amounts[k]
		}
	}
	return 30
}

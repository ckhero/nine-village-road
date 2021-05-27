/**
 *@Description
 *@ClassName rand_amount_test
 *@Date 2021/5/27 上午10:50
 *@Author ckhero
 */

package rand_amount

import (
	"fmt"
	"testing"
)

func TestGetRandAmount(t *testing.T) {
	data := map[uint64]uint64{
		110 : 46,
		220 : 30,
		330 : 30,
		550 : 30,
		660 : 30,
		880 : 20,
		990 : 10,
		8800 : 1,
	}

	res := map[uint64]uint64{
		110 : 0,
		220 : 0,
		330 : 0,
		550 : 0,
		660 : 0,
		880 : 0,
		990 : 0,
		8800 : 0,
		30 : 0,
	}
	for i := 0; i < 60; i ++ {
		t := GetRandAmount(data)
		res[t] += 1
	}
	var total uint64
	for amount, nums := range res {
		total += amount * nums
		fmt.Printf("金额: %d 分， 数量: %d \n", amount, nums)
	}
	fmt.Println(total / 100)

}

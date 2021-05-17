/**
 *@Description
 *@ClassName user_scenic_test
 *@Date 2021/5/17 上午9:22
 *@Author ckhero
 */

package usecase

import (
	"fmt"
	"github.com/thoas/go-funk"
	"nine-village-road/pkg/constant"
	"testing"
)

func TestUserScenicUsecase_CheckAllScenicScaned(t *testing.T) {
	a :=[]string{constant.ScenicCWF, "test"}
	aa, bb := funk.DifferenceString(a, constant.AllScenic)
	fmt.Println(aa, bb)

}

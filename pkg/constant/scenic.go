/**
 *@Description
 *@ClassName scenic
 *@Date 2021/5/14 下午3:40
 *@Author ckhero
 */

package constant

// 景点
const (
	// 黄龙洞
	ScenicHLD = "HLD"
	// 楚王峰
	ScenicCWF = "CWF"
	// 两河口
	ScenicLHK = "LHK"
	// 悬崖索桥
	ScenicXYSQ = "XYSQ"
)

const (
	UserScenicStatusInvalid = "INVALID"
	UserScenicStatusValid = "VALID"
)

var AllScenic = []string{
	ScenicHLD,
	ScenicCWF,
	ScenicLHK,
	ScenicXYSQ,
}
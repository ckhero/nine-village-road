/**
 *@Description
 *@ClassName user_scenic
 *@Date 2021/5/14 下午4:26
 *@Author ckhero
 */

package domain

import (
	"github.com/ckhero/go-common/db/mysql"
)

// UserScenic [...]
type UserScenic struct {
	mysql.BaseEntity
	UserScenicID int64  `gorm:"primaryKey;column:user_scenic_id;type:bigint(11);not null"`
	Scenic       string `gorm:"index:idx_user_id_scenic;column:scenic;type:varchar(32);not null"` // 景点
	UserID       int64  `gorm:"index:idx_user_id_scenic;column:user_id;type:bigint(20);not null"`
	OpenID       string `gorm:"index:idx_openid;column:open_id;type:varchar(128);not null"`
	Status       string `gorm:"column:status;type:varchar(16);not null"` // [VALID;INVALID]
}

func (*UserScenic) TableName() string {
	return "user_scenic"
}

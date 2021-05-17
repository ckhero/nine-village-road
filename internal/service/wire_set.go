/**
 *@Description
 *@ClassName wire_set
 *@Date 2021/5/12 下午5:04
 *@Author ckhero
 */

package service

import (
	"github.com/ckhero/go-common/config"
	"github.com/ckhero/go-common/db"
	"github.com/ckhero/go-common/wx"
	"github.com/google/wire"
	"nine-village-road/internal/repo"
	"nine-village-road/internal/usecase"
)

var ProviderUserSet = wire.NewSet(NewUserService,
	usecase.NewUserUsecase,
	repo.NewUserRepo,
	db.NewDatabase,
	usecase.NewWeixinUsecase,
	repo.NewWeixinRepo,
	repo.NewUserRedPacketRepo,
	wx.NewMiniPayClient,
	wx.NewMiniClient,
	config.GetWeixinPayCfg,

	usecase.NewUserScenicUsecase,
	repo.NewUserScenicRepo,
	)
/**
 *@Description
 *@ClassName main
 *@Date 2021/4/29 下午3:47
 *@Author ckhero
 */

package main

import (
	"github.com/ckhero/go-common/config"
	"github.com/ckhero/go-common/constant/arg_const"
	"github.com/ckhero/go-common/plugin"
	"github.com/ckhero/go-common/util/arg"
	"github.com/ckhero/go-common/web"
	"log"
	"nine-village-road/pkg/router"
)

func init() {
	//path := util.GetArg(constant.ArgConfig, "./src/mini-program-gateway/config/dev.yaml")
	path := arg.GetArg(arg_const.ArgConfig, "./config/dev.yaml")
	config.InitConfig(path)
}

func main()  {
	Run(config.GetGlobalCfg())
	//a := []byte{1}
	//b := string(a)
	//fmt.Println(b)
}

func Run(cfg *config.Config) {

	// Mechanical domain.
	errc := make(chan error)
	// init plugins
	px := plugin.NewPluginCtx(errc,
		plugin.NewPluginMysql(),
	)
	px.InitPlugin()
	defer px.Release()
	go web.Run(errc, router.RegisterRouter(), nil)
	// Run!
	log.Println("exit", <-errc)
}
package main

import (
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"os"
	"sync"
	"tetris/config"
	"tetris/internal/game"
	"tetris/internal/web"
	"tetris/models"
	"tetris/pkg/log"
)

func main() {
	app := cli.NewApp()
	app.Name = "tetris"
	app.Author = "rezone games"
	app.Version = "0.0.1"
	app.Usage = "tetris"
	app.Action = serve
	app.Run(os.Args)
}

func serve(ctx *cli.Context) error {
	ctx1 := context.Background()
	log.InitLog()
	config.InitServerConfig()
	sc := config.ServerConfig
	models.StartMongo(ctx1, sc.Mongo.Uri, sc.Mongo.Db, false)
	models.StartRedis(ctx1, sc.Redis.Host, sc.Redis.Db, "")
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() { defer wg.Done(); game.StartUp() }() // 开启游戏服
	go func() { defer wg.Done(); web.StartUp() }()  // 开启web服务器
	wg.Wait()
	return nil
}

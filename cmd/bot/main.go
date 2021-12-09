package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/bot"
	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/cache"
	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/logger/logrus"
	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/repository/inmemory"
	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/service"
	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/tg"
	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/tg/updates/repost"
)

func main() {
	loggerFactory := logrus.NewFactory()
	mainLog := loggerFactory.New("main")
	check := func(err error) {
		if err != nil {
			mainLog.Error(err)
			os.Exit(1)
		}
	}

	if err := godotenv.Load(); err != nil {
		mainLog.Warnf("godotenv.Load: %v", err)
	}

	tgClientConf := tg.Config{
		BotToken: os.Getenv("TG_BOT_TOKEN"),
		DataPath: os.Getenv("DATA_PATH"),
	}
	tgClient, err := tg.NewClient(tgClientConf)

	check(err)
	tgBot := bot.NewBot(loggerFactory, tgClient)

	repostConfigRepository := inmemory.NewRepostConfigRepository()
	publicChatSearcher := cache.NewPublicChatSearcher(tgClient)
	configService := service.NewRepostConfigService(repostConfigRepository, publicChatSearcher)
	for _, source := range []string{"tmp_src", "karga4", "armeyskov"} {
		err = configService.Add(source, "tmp_dst")
		check(err)
	}

	repostHandler := repost.NewHandler(loggerFactory, tgClient, repostConfigRepository)
	tgBot.AddUpdateHandlers(repostHandler)

	ctx, stopNotify := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopNotify()

	check(tgBot.MainLoop(ctx))
}

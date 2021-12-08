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

	repostConfig := repost.Config{
		Sources:     []string{"tmp_src"},
		Destination: "tmp_dst",
	}
	repostHandler := repost.NewHandler(repostConfig, loggerFactory, tgClient, cache.NewPublicChatSearcher(tgClient))
	tgBot.AddUpdateHandlers(repostHandler)

	ctx, stopNotify := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopNotify()

	check(tgBot.MainLoop(ctx))
}

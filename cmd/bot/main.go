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

	repostConfigRepository := inmemory.NewRepostConfigRepository(loggerFactory)
	publicChatSearcher := cache.NewPublicChatSearcher(tgClient)
	configService := service.NewRepostConfigService(repostConfigRepository, publicChatSearcher)

	err = configService.Add("tmp_src", "tmp_dst")
	check(err)

	// TODO: move to db
	// TODO: rename to repost rule
	kargoChannels := []string{"armeyskov", "kargokult", "ikkinpi", "holarhia", "neoposta4", "neopostshit", "neopostart"}
	for _, source := range kargoChannels {
		err = configService.Add(source, "tmp_dst")
		check(err)

		err = configService.Add(source, "karga4")
		check(err)
	}

	for _, source := range []string{"armeyskov", "kargokult", "ikkinpi", "holarhia"} {
		err = configService.Add(source, "neopostshit")
		check(err)
	}

	repostHandler := repost.NewHandler(loggerFactory, tgClient, repostConfigRepository)
	//userRepository := inmemory.NewUserRepository()
	//
	//initialAdmin := model.User{
	//	Username: os.Getenv("INITIAL_ADMIN_USERNAME"),
	//	IsAdmin:  true,
	//}
	//initialAdmin.Id, err = strconv.ParseInt(os.Getenv("INITIAL_ADMIN_ID"), 10, 64)
	//check(err)
	//
	//err = userRepository.Add(initialAdmin)
	//check(err)
	//
	//commandHandler := command.NewHandler(loggerFactory, tgClient, repostConfigRepository, userRepository)
	tgBot := bot.NewBot(loggerFactory, tgClient, repostHandler)

	ctx, stopNotify := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopNotify()

	check(tgBot.MainLoop(ctx))
}

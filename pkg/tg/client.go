package tg

import (
	"fmt"
	"path/filepath"

	"github.com/Arman92/go-tdlib"
)

type Config struct {
	BotToken string
	DataPath string
}

type Client struct {
	*tdlib.Client
	updates chan tdlib.UpdateMsg
}

func NewClient(conf Config) (*Client, error) {
	config := tdlib.Config{
		APIID:               "187786",
		APIHash:             "e782045df67ba48e441ccb105da8fc85",
		SystemLanguageCode:  "en",
		DeviceModel:         "Server",
		SystemVersion:       "1.0.0",
		ApplicationVersion:  "1.0.0",
		UseMessageDatabase:  true,
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		DatabaseDirectory:   filepath.Join(conf.DataPath, "db"),
		FileDirectory:       filepath.Join(conf.DataPath, "files"),
		IgnoreFileNames:     true,
	}

	tdlib.SetLogVerbosityLevel(1)
	client := tdlib.NewClient(config)
	for {
		currentState, err := client.Authorize()
		if err != nil {
			return nil, fmt.Errorf("client authorize: %w", err)
		}
		if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPhoneNumberType {
			_, err := client.CheckAuthenticationBotToken(conf.BotToken)
			if err != nil {
				return nil, fmt.Errorf("check bot token: %w", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateReadyType {
			break
		}
	}

	return &Client{Client: client, updates: client.GetRawUpdatesChannel(100)}, nil
}

func (c *Client) Updates() <-chan tdlib.UpdateMsg {
	return c.updates
}

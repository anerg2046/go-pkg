package notification

import (
	resty "github.com/anerg2046/go-pkg/httpclient"
)

const apiURI = "https://tgapi.anerg.com/"

type Telegram struct {
	BotToken string
	ChatID   string
}

func (t *Telegram) SendMessage(text string) {
	uri := apiURI + t.BotToken + "/sendMessage"
	body := map[string]string{
		"chat_id": t.ChatID,
		"text":    text,
	}
	resty.HttpClient.R().SetBody(body).Post(uri)
}

package context

import (
	"tgbotBARAHOLKA/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserState struct {
	Name      string
	Level     int
	Data      map[string]interface{}
	MessageID int
}

type Context struct {
	BotAPI     *tgbotapi.BotAPI
	UserStates map[int64]*UserState
	Config     *config.Config
}

func NewContext(botAPI *tgbotapi.BotAPI, cfg *config.Config) *Context {
	return &Context{
		BotAPI:     botAPI,
		UserStates: make(map[int64]*UserState),
		Config:     cfg,
	}
}

func GetUserState(userID int64, ctx *Context) *UserState {
	if state, exists := ctx.UserStates[userID]; exists {
		return state
	}
	ctx.UserStates[userID] = &UserState{
		Name:      "start",
		Level:     0,
		Data:      make(map[string]interface{}),
		MessageID: 0,
	}
	return ctx.UserStates[userID]
}

func UpdateUserLevel(userID int64, ctx *Context, newLevel int) {
	state := GetUserState(userID, ctx)
	state.Level = newLevel
}

func UpdateUserName(userID int64, ctx *Context, newName string) {
	state := GetUserState(userID, ctx)
	state.Name = newName
	state.Level = 0
}

func ClearAllUserData(userID int64, ctx *Context) {
	state := GetUserState(userID, ctx)
	state.Data = make(map[string]interface{})
}

func SaveMessageID(userID int64, ctx *Context, messageID int) {
	state := GetUserState(userID, ctx)
	state.MessageID = messageID
}

func (ctx *Context) SendMessage(msg tgbotapi.MessageConfig) (int, error) {
	sentMessage, err := ctx.BotAPI.Send(msg)
	if err != nil {
		return 0, err
	}

	SaveMessageID(msg.ChatID, ctx, sentMessage.MessageID)
	return sentMessage.MessageID, nil
}

package slackbot

import (
	"fmt"
	"sync"

	"github.com/slack-go/slack"

	"github.com/spy16/pkg/lua"
)

// LuaHandler implements the Bot handler using a lua scripting layer.
type LuaHandler struct {
	InitFile    string
	HandlerFunc string

	once sync.Once
	lua  *lua.Lua
}

func (h *LuaHandler) Handle(bot *Bot, msg slack.MessageEvent, from slack.User) error {
	if err := h.initOnce(); err != nil {
		return err
	}

	lMsg, err := h.lua.Call(h.HandlerFunc, bot, from, msg)
	if err != nil {
		sendText := fmt.Sprintf(":pensive: I experienced an internal issue: %v", err)
		return bot.SendMessage(sendText, msg.Msg)
	}
	return bot.SendMessage(lMsg.String(), msg.Msg)
}

func (h *LuaHandler) initOnce() error {
	var err error
	h.once.Do(func() {
		h.lua, err = lua.New(
			lua.Path("."),
		)
		if err != nil {
			return
		}

		err = h.lua.ExecuteFile(h.InitFile)
		if err != nil {
			return
		}
	})
	return err
}

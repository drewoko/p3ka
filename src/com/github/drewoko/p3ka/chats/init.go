package chats

import (
	Core "../core"
)

func InitChats(messages_channel chan Core.Msg, messages_delete_channel chan Core.Msg, db *Core.DataBase, config *Core.Config) {
	go initPeka2Tv(messages_channel, messages_delete_channel, config)
	go initGoodGame(messages_channel, messages_delete_channel, db, config)
}

package frame

import "github.com/OmineDev/neomega-plugin-sdk/qq"

// DefaultCQMessageCb receives the message source, sender name and message text.
type DefaultCQMessageCb func(source, name, message string)

// Messaging sends messages and queries account, group and guild information through the configured
// messaging service.
type Messaging interface {
	// SendGroupMessage sends a message to a group and reports success and message ID in onCb.
	SendGroupMessage(groupID int64, message string, onCb func(ok bool, msgID int64))
	// SendGuildMessage sends a message to a guild channel and reports success and message ID in onCb.
	SendGuildMessage(guildID, channelID string, message string, onCb func(ok bool, msgID int64))
	// SendPrivateMessage sends a message to a user and reports success and message ID in onCb.
	SendPrivateMessage(userID int64, message string, onCb func(ok bool, msgID int64))
	// SendToDefault sends a message to the configured default destination.
	SendToDefault(message string)
	// OnDefaultMessage registers a callback for messages from the configured default destination.
	OnDefaultMessage(cb DefaultCQMessageCb)
	// SendTo sends a message to the destination string understood by the configured messaging service.
	SendTo(target, message string)
	// RegisterPacketNoBlockCB registers a callback for incoming service packets. The callback must return
	// promptly.
	RegisterPacketNoBlockCB(cb func(pk qq.CQPacket, data []byte))
	// GetGroupMember fetches the member cards for a group and returns them through onCb.
	GetGroupMember(groupID int64, onCb func(qq.GroupMemberCards))
	// GetGuildList fetches the guilds visible to the account and returns them through onCb.
	GetGuildList(onCb func(guilds qq.GuildList))
	// GetGuildChannels fetches the channels in guildID and returns them through onCb.
	GetGuildChannels(guildID string, onCb func(channels qq.GuildChannels))
	// GetGuildMemberProfile fetches the profile for userID in guildID and returns it through onCB.
	GetGuildMemberProfile(guildID, userID string, onCB func(qq.GuildMemberProfile))
}

// CQHTTPAccess is an alias for Messaging.
type CQHTTPAccess = Messaging

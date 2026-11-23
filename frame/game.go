package frame

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// InfoSender sends chat and on-screen text to players.
type InfoSender interface {
	// BotSay sends a chat message as the bot.
	BotSay(msg string)
	// SayTo sends plain text to a player or target selector.
	SayTo(target string, msg string)
	// RawSayTo sends a Bedrock rawtext JSON message to a player or target selector.
	RawSayTo(target string, msg string)
	// ActionBarTo displays msg above the target's hotbar.
	ActionBarTo(target string, msg string)
	// TitleTo displays msg as the target's title.
	TitleTo(target string, msg string)
	// SubTitleTo displays subTitle beneath title for the target. The subtitle argument comes first.
	SubTitleTo(target string, subTitle string, title string)
}

// CmdSenderFor sends game commands. R is the output type for commands that request a reply.
type CmdSenderFor[R any] interface {
	// SendWOCmd sends a settings command without requesting output.
	SendWOCmd(cmd string)
	// SendWebSocketCmdOmitResponse sends a command with a WebSocket origin without requesting its result.
	SendWebSocketCmdOmitResponse(cmd string)
	// SendPlayerCmdOmitResponse sends a command with a player origin without requesting its result.
	SendPlayerCmdOmitResponse(cmd string)
	// SendAICommandOmitResponse sends a command through the AI command interface without requesting its
	// result.
	SendAICommandOmitResponse(cmd string)
	// SendWebSocketCmdNeedResponse returns an operation that sends a WebSocket-origin command and receives its
	// output.
	SendWebSocketCmdNeedResponse(cmd string) AsyncResult[R]
	// SendPlayerCmdNeedResponse returns an operation that sends a player-origin command and receives its
	// output.
	SendPlayerCmdNeedResponse(cmd string) AsyncResult[R]
	// SendAICommandNeedResponse returns an operation that sends an AI command and receives its output.
	SendAICommandNeedResponse(cmd string) AsyncResult[R]
}

// GameCtrl sends packets, commands and messages from the bot's connection.
type GameCtrl interface {
	InfoSender
	CmdSenderFor[*packet.CommandOutput]
	// SendPacket sends a gophertunnel client packet. The host logs conversion errors; use TrySendPacket to
	// handle them yourself.
	SendPacket(packet.Packet)
	// TrySendPacket converts a gophertunnel client packet and passes it to the host for sending. An error
	// means the packet was not sent; nil is not a server acknowledgement.
	TrySendPacket(packet.Packet) error
}

// CmdSender sends commands whose replies are Bedrock CommandOutput packets.
type CmdSender = CmdSenderFor[*packet.CommandOutput]

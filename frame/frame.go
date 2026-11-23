// Package frame defines the Go interfaces shared by OMG and its in-process plugins.
// Embed BasicComponent for the default lifecycle and implement Inject to register handlers.
package frame

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// NoBlockAndDetachablePacketCallbackFor handles a packet without blocking the receive loop. Returning an
// error unsubscribes the callback.
type NoBlockAndDetachablePacketCallbackFor[P any] func(P) error

// PacketListenerFor dispatches incoming packets of type P; plugins use PacketDispatcher for Bedrock
// packets.
type PacketListenerFor[P any] interface {
	// GetPacketFilterSettings returns the switches that control delivery of high-volume packet types.
	GetPacketFilterSettings() PacketFilterSettings
	// SetAnyPacketCallback registers a handler for all dispatched packets. Pass true to run the handler in a
	// separate goroutine.
	SetAnyPacketCallback(func(P), bool)
	// SetTypedPacketCallback registers a handler for one packet ID. Pass false to run it on the receive loop;
	// in that case it must not block.
	SetTypedPacketCallback(uint32, func(P), bool)
	// GetMCPacketNameIDMapping returns the packet names and IDs understood by this listener.
	GetMCPacketNameIDMapping() map[string]uint32
	// TranslateStringWantsToIDSet builds a set of packet IDs from names in order. For example, "all",
	// "!MovePlayer" selects all but MovePlayer.
	TranslateStringWantsToIDSet([]string) map[uint32]bool
	// AddNewNoBlockAndDetachablePacketCallback subscribes to the selected packet IDs without blocking the
	// receive loop. An error from the handler removes the subscription.
	AddNewNoBlockAndDetachablePacketCallback(map[uint32]bool, NoBlockAndDetachablePacketCallbackFor[P])
}

// PacketDispatcher delivers incoming packets as gophertunnel packet.Packet values.
type PacketDispatcher = PacketListenerFor[packet.Packet]

// ExtendOmega is the host passed to a plugin's Inject method.
type ExtendOmega interface {
	PreInitOmega
	GameMenuSetter
	Messaging
	CanCloseWithError
	// GetReactCore returns the packet listener together with the host's close notifications.
	GetReactCore() ReactCore
	// GetPlayerInteract returns the player directory and chat listeners.
	GetPlayerInteract() PlayerInteract
	// GetGameControl returns the host for sending packets, commands and in-game messages.
	GetGameControl() GameCtrl
	// GetGameListener returns the listener for incoming Bedrock packets.
	GetGameListener() PacketDispatcher
	// GetMicroUQHolder returns the host's current view of the bot, players and world.
	GetMicroUQHolder() WorldInfo
	// GetLowLevelAreaRequester returns the interface for requesting structures and sub-chunks from the server.
	GetLowLevelAreaRequester() LowLevelAreaRequester
	// GetBlobHashHolder returns the cache for chunk blobs addressed by hash.
	GetBlobHashHolder() BlobCache
	// GetBotAction returns the bot's inventory, block interaction and command helpers.
	GetBotAction() BotActionComplex
}

// PacketFilterSettings controls delivery of frequently updated entity data.
type PacketFilterSettings interface {
	// MovePlayerEnabled reports whether MovePlayer packets are enabled.
	MovePlayerEnabled() (enabled bool)
	// SwitchMovePlayer enables or disables MovePlayer packets.
	SwitchMovePlayer(enabled bool)
	// MovePlayerDeltaEnabled reports whether MoveActorDelta packets are enabled.
	MovePlayerDeltaEnabled() (enabled bool)
	// SwitchMovePlayerDelta enables or disables MoveActorDelta packets.
	SwitchMovePlayerDelta(enable bool)
	// SetActorMotionEnabled reports whether SetActorMotion packets are enabled.
	SetActorMotionEnabled() (enabled bool)
	// SwitchSetActorData enables or disables SetActorData packets.
	SwitchSetActorData(enabled bool)
	// SetActorDataEnabled reports whether SetActorData packets are enabled.
	SetActorDataEnabled() (enabled bool)
	// SwitchSetActorMotion enables or disables SetActorMotion packets.
	SwitchSetActorMotion(enable bool)
	// SetEnabledMode sets the filter for a packet ID. It returns false if that ID has no switch.
	SetEnabledMode(pkID uint32, enabled bool) (success bool)
}

// ConfigRead locates enabled component configurations by source.
type ConfigRead interface {
	// GetEnabledConfigBySource returns the paths of enabled configuration files matching source.
	GetEnabledConfigBySource(source string) []string
}

// CanClose exposes the host's shutdown state and the reason it closed.
type CanClose interface {
	// Close requests shutdown without supplying an error.
	Close()
	// WaitClosed returns a channel that delivers the shutdown error when the host closes. The error may be
	// nil.
	WaitClosed() <-chan error
	// Closed reports whether the host has closed.
	Closed() bool
	// CloseError returns the error that caused shutdown, or nil if there is none.
	CloseError() error
}

// CanCloseWithError can shut down the host with a specific cause.
type CanCloseWithError interface {
	CanClose
	// CloseWithError closes the host and records err as the cause.
	CloseWithError(err error)
}

// ReactCore combines packet subscriptions with host shutdown notifications.
type ReactCore interface {
	CanCloseWithError
	PacketDispatcher
}

// PreInitOmega provides node services, storage and terminal access without requiring a game connection.
type PreInitOmega interface {
	ConfigRead
	StorageAndLogAccess
	FlexAccess
}

package frame

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// GameRule holds a game rule's value as text and whether players may change it.
type GameRule struct {
	// CanBeModifiedByPlayer reports whether a player may change this rule.
	CanBeModifiedByPlayer bool
	// Value is the rule's value as text.
	Value string
}

// WorldInfo is the host's read-only view of the world, the bot and known players.
type WorldInfo interface {
	WorldReader
	PlayersReader
	// GetBotBasicInfo returns the bot's login and entity identifiers.
	GetBotBasicInfo() BotReader
	// GetPlayersInfo returns the directory of known players.
	GetPlayersInfo() PlayersReader
	// GetExtendInfo returns the world settings, clock and bot position.
	GetExtendInfo() WorldReader
	BotReader
}

// WorldReaderFor reads world state using C for open-container data and D for dimensions. Methods returning
// found distinguish missing data from a zero value.
type WorldReaderFor[C, D any] interface {
	// GetCurrentOpenedContainer returns the most recent open-container packet if the bot has a container open.
	GetCurrentOpenedContainer() (C, bool)
	// GetBotDimension returns the dimension the bot is in, if known.
	GetBotDimension() (D, bool)
	// GetWorldName returns the world's display name, if known.
	GetWorldName() (worldName string, found bool)
	// GetWorldSeed returns the world generation seed, if known.
	GetWorldSeed() (worldSeed int64, found bool)
	// GetWorldGenerator returns the server's world generator ID, if known.
	GetWorldGenerator() (worldGenerator int32, found bool)
	// GetLevelID returns the level identifier reported by the server, if known.
	GetLevelID() (levelID string, found bool)
	// GetCompressThreshold returns the negotiated packet compression threshold in bytes, if known.
	GetCompressThreshold() (compressThreshold uint16, found bool)
	// GetWorldGameMode returns the world's game mode, not a particular player's mode.
	GetWorldGameMode() (worldGameMode int32, found bool)
	// GetWorldDifficulty returns the world's difficulty setting, if known.
	GetWorldDifficulty() (worldDifficulty uint32, found bool)
	// GetTime returns the last world time reported by the server, in ticks.
	GetTime() (time int32, found bool)
	// GetDayTime returns the tick within the 24,000-tick day.
	GetDayTime() (dayTime int32, found bool)
	// GetDayTimePercent returns the fraction of the day elapsed, not a percentage from 0 to 100.
	GetDayTimePercent() (dayTimePercent float32, found bool)
	// GetGameRules returns the known game rules keyed by name.
	GetGameRules() (gameRules map[string]*GameRule, found bool)
	// GetCurrentTick returns the host's estimate of the current server tick, if available.
	GetCurrentTick() (currentTick int64, found bool)
	// GetSyncRatio returns the observed tick rate relative to 20 ticks per second. The ratio is unknown until
	// the host has taken enough samples.
	GetSyncRatio() (ratio float32, known bool)
	// GetBotPosition returns the bot's last known position and its age in ticks.
	GetBotPosition() (pos mgl32.Vec3, outOfSyncTick int64)
}

// WorldReader reads world state using Bedrock container packets and SDK dimensions.
type WorldReader = WorldReaderFor[*packet.ContainerOpen, Dimension]

// BotReader identifies the bot's account and its in-game entity.
type BotReader interface {
	// GetBotName returns the bot's in-game name.
	GetBotName() string
	// GetBotRuntimeID returns the bot entity's runtime ID for this session.
	GetBotRuntimeID() uint64
	// GetBotUniqueID returns the bot entity's unique ID.
	GetBotUniqueID() int64
	// GetBotIdentity returns the identity string from the bot's login.
	GetBotIdentity() string
	// GetBotUUIDStr returns the bot's identity UUID as a string.
	GetBotUUIDStr() string
	// GetBotUID returns the bot's numeric account ID, distinct from its entity IDs.
	GetBotUID() int64
}

// PlayersReader looks up players for read-only access to their state.
type PlayersReader interface {
	// GetAllOnlinePlayers returns readers for the players currently online.
	GetAllOnlinePlayers() []PlayerUQReader
	// GetPlayerByUUID looks up a player's state by UUID.
	GetPlayerByUUID(uuid.UUID) (player PlayerUQReader, found bool)
	// GetPlayerByUUIDString looks up a player's state by UUID string.
	GetPlayerByUUIDString(uuidStr string) (player PlayerUQReader, found bool)
	// GetPlayerByUniqueID looks up a player's state by entity unique ID.
	GetPlayerByUniqueID(uniqueID int64) (player PlayerUQReader, found bool)
	// GetPlayerByName looks up a player's state by name.
	GetPlayerByName(name string) (player PlayerUQReader, found bool)
}

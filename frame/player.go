package frame

import (
	"time"

	"github.com/google/uuid"
)

// GameChat holds the original text packet fields alongside the name and message used by chat handlers.
type GameChat struct {
	// Name is the sender's name without display prefixes or formatting.
	Name string
	// Msg contains the message split into words. In a menu callback it contains the command arguments.
	Msg []string
	// Type is the Bedrock text type, such as packet.TextTypeChat.
	Type byte
	// RawMsg is the message before parsing or translation.
	RawMsg string
	// RawName is the source name as received, including formatting.
	RawName string
	// RawParameters are the text packet's translation parameters.
	RawParameters []string
	// Aux holds optional host-specific information. Most messages leave it nil.
	Aux any
	// ParsedMsg is the text after the host has interpreted rawtext or translation parameters.
	ParsedMsg string
}

// PlayerMsgListener registers chat handlers and waits for individual replies.
type PlayerMsgListener interface {
	// GetInput waits for the next input from playerName. With breakOnLeave set, a disconnect ends the wait
	// with an error.
	GetInput(playerName string, breakOnLeave bool) AsyncResult[*GameChat]
	// SetOnChatCallback adds a handler for general chat. A message consumed by a pending input request is not
	// sent to these handlers.
	SetOnChatCallback(func(chat *GameChat))
	// SetOnSpecificCommandBlockTellCallback adds a handler for whispers from the named command block.
	SetOnSpecificCommandBlockTellCallback(commandBlockName string, cb func(chat *GameChat))
	// SetOnSpecificItemMsgCallback adds a handler for messages whose raw source name matches itemName.
	SetOnSpecificItemMsgCallback(itemName string, cb func(chat *GameChat))
}

// PlayerInteract looks up online players and reports changes to the player list.
type PlayerInteract interface {
	PlayerMsgListener
	// ListAllPlayers returns handles for the players currently known to be online.
	ListAllPlayers() []PlayerKit
	// ListenPlayerChange registers a callback. Players already present are reported as "exist"; later joins
	// and departures are reported as "online" and "offline".
	ListenPlayerChange(cb func(player PlayerKit, action string))
	// GetPlayerKit looks up an online player by name. Check found before using the handle.
	GetPlayerKit(name string) (playerKit PlayerKit, found bool)
	// GetPlayerKitByUUID looks up an online player by UUID.
	GetPlayerKitByUUID(ud uuid.UUID) (playerKit PlayerKit, found bool)
	// GetPlayerKitByUUIDString looks up an online player by the string form of their UUID.
	GetPlayerKitByUUIDString(ud string) (playerKit PlayerKit, found bool)
	// GetPlayerKitByUniqueID looks up an online player by entity unique ID, not entity runtime ID.
	GetPlayerKitByUniqueID(uniqueID int64) (playerKit PlayerKit, found bool)
}

// QueryResult describes an entity matched by a querytarget command. Position may be nil.
type QueryResult struct {
	// Dimension is the entity's dimension ID.
	Dimension int `json:"dimension"`
	// Position is the entity's position, or nil if it was not returned.
	Position *struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	} `json:"position"`
	// UUID is the entity's UUID string.
	UUID string `json:"uniqueId"`
	// YRot is the entity's yaw in degrees.
	YRot float64 `json:"yRot"`
}

// PlayerKit is a handle for reading a player's state, sending messages and changing abilities.
type PlayerKit interface {
	PlayerUQReader
	// Say sends plain text to this player.
	Say(msg string)
	// RawSay sends Bedrock rawtext JSON to this player.
	RawSay(msg string)
	// ActionBar displays a message above the player's hotbar.
	ActionBar(msg string)
	// Title displays a title on the player's screen.
	Title(msg string)
	// SubTitle displays subTitle beneath title on the player's screen.
	SubTitle(subTitle string, title string)
	// GetInput waits for the player's next input. If breakOnLeave is true, leaving the game returns an error
	// instead of waiting for a rejoin.
	GetInput(breakOnLeave bool) AsyncResult[*GameChat]
	// CheckCondition calls onResult with whether the player matches the given selector arguments, such as
	// "m=c" or "tag=staff".
	CheckCondition(onResult func(bool), conditions ...string)
	// Query asks the server for the player's position and dimension, subject to the given selector arguments.
	Query(onResult func([]QueryResult), conditions ...string)
	// SetAbility(adventureFlagsUpdateMap, actionPermissionUpdateMap map[uint32]bool) (sent bool)
	// SetAbilityString(adventureFlagsUpdateMap, actionPermissionUpdateMap map[string]bool) (sent bool)
	// SetBuildAbility allows or disallows placing blocks.
	SetBuildAbility(allow bool)
	// SetMineAbility allows or disallows breaking blocks.
	SetMineAbility(allow bool)
	// SetDoorsAndSwitchesAbility allows or disallows using doors and switches.
	SetDoorsAndSwitchesAbility(allow bool)
	// SetOpenContainersAbility allows or disallows opening containers.
	SetOpenContainersAbility(allow bool)
	// SetAttackPlayersAbility allows or disallows attacking other players.
	SetAttackPlayersAbility(allow bool)
	// SetAttackMobsAbility allows or disallows attacking mobs.
	SetAttackMobsAbility(allow bool)
	// SetOperatorCommandsAbility allows or disallows using operator commands.
	SetOperatorCommandsAbility(allow bool)
	// SetTeleportAbility allows or disallows teleporting.
	SetTeleportAbility(allow bool)
}

// PlayerUQReader reads the host's last known state for one player. A false found result means the value is
// unknown, not that the ability or status is off.
type PlayerUQReader interface {
	// GetUUID returns the player's UUID, if known.
	GetUUID() (id uuid.UUID, found bool)
	// GetUUIDString returns the player's UUID as a string, if known.
	GetUUIDString() (id string, found bool)
	// GetEntityUniqueID returns the player's entity unique ID, if known.
	GetEntityUniqueID() (id int64, found bool)
	// GetLoginTime returns the time the host recorded the player joining.
	GetLoginTime() (time time.Time, found bool)
	// GetUsername returns the player's name as recorded by the host.
	GetUsername() (name string, found bool)
	// GetPlatformChatID returns the platform-specific chat identifier, if known.
	GetPlatformChatID() (id string, found bool)
	// GetBuildPlatform returns the Bedrock platform code reported for the player's client.
	GetBuildPlatform() (platform int32, found bool)
	// GetSkinID returns the identifier of the player's skin, if known.
	GetSkinID() (id string, found bool)
	// GetDeviceID returns the player's device identifier, if known.
	GetDeviceID() (id string, found bool)
	// GetEntityRuntimeID returns the ID used for the player's entity in packets for this session.
	GetEntityRuntimeID() (id uint64, found bool)
	// GetEntityMetadata returns the last known entity metadata keyed by Bedrock metadata index.
	GetEntityMetadata() (entityMetadata map[uint32]any, found bool)
	// IsOP reports whether the player is an operator and whether this is known.
	IsOP() (op bool, found bool)
	// CanBuild reports whether the player may place blocks.
	CanBuild() (hasAbility bool, found bool)
	// CanMine reports whether the player may break blocks.
	CanMine() (hasAbility bool, found bool)
	// CanDoorsAndSwitches reports whether the player may use doors and switches.
	CanDoorsAndSwitches() (hasAbility bool, found bool)
	// CanOpenContainers reports whether the player may open containers.
	CanOpenContainers() (hasAbility bool, found bool)
	// CanAttackPlayers reports whether the player may attack other players.
	CanAttackPlayers() (hasAbility bool, found bool)
	// CanAttackMobs reports whether the player may attack mobs.
	CanAttackMobs() (hasAbility bool, found bool)
	// CanOperatorCommands reports whether the player may use operator commands.
	CanOperatorCommands() (hasAbility bool, found bool)
	// CanTeleport reports whether the player may teleport.
	CanTeleport() (hasAbility bool, found bool)
	// StatusInvulnerable reports whether the player is invulnerable, if known.
	StatusInvulnerable() (hasStatus bool, found bool)
	// StatusFlying reports whether the player is currently flying, if known.
	StatusFlying() (hasStatus bool, found bool)
	// StatusMayFly reports whether the player is allowed to fly, regardless of whether they are flying now.
	StatusMayFly() (hasStatus bool, found bool)
	// StillOnline reports whether the host still considers this player online.
	StillOnline() bool
}

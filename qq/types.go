// Package qq defines the data passed to plugin messaging callbacks.
package qq

// User describes a QQ or guild-chat account returned by the messaging service.
type User struct {
	// UserID is the account's numeric user ID.
	UserID int64 `json:"user_id"`
	// Nickname is the account's display name.
	Nickname string `json:"nickname"`
	// Sex is the sex string returned by the messaging service.
	Sex string `json:"sex"`
	// Age is the age value returned by the messaging service.
	Age int `json:"age"`
}

// GroupUser adds a QQ group card and role to a user.
type GroupUser struct {
	User
	// Card is the member's group card, if set.
	Card string `json:"card"`
	// Area is the location string returned for the member.
	Area string `json:"area"`
	// Level is the group level string.
	Level string `json:"level"`
	// Role is the member's group role.
	Role string `json:"role"`
	// Title is the member's group title.
	Title string `json:"title"`
	// GroupID is the group ID containing this member.
	GroupID int64 `json:"group_id"`
}

// GroupMemberCard adds group membership and moderation details to a group user.
type GroupMemberCard struct {
	GroupUser
	// JoinTime is the member's join time as a Unix timestamp.
	JoinTime int32 `json:"join_time"`
	// LastSentTime is the member's last message time as a Unix timestamp.
	LastSentTime int32 `json:"last_sent_time"`
	// Unfriendly reports the service's unfriendly flag.
	Unfriendly bool `json:"unfriendly"`
	// TitleExpireTime is the member title's expiration time as a Unix timestamp.
	TitleExpireTime int64 `json:"title_expire_time"`
	// CardChangeable reports whether the member may edit their card.
	CardChangeable bool `json:"card_changeable"`
	// ShutUpTimestamp is the member's mute expiration as a Unix timestamp.
	ShutUpTimestamp int64 `json:"shut_up_timestamp"`
}

// GroupMemberCards is a list of group member cards.
type GroupMemberCards []GroupMemberCard

// SlowModeInfo describes a guild channel's slow-mode setting.
type SlowModeInfo struct {
	// SlowModeKey is the service's slow-mode key.
	SlowModeKey int32 `json:"slow_mode_key"`
	// SlowModeText is the text shown for the slow-mode setting.
	SlowModeText string `json:"slow_mode_text"`
	// SpeakFrequency is the allowed speaking frequency.
	SpeakFrequency int32 `json:"speak_frequency"`
	// SlowModeCircle is the slow-mode interval value.
	SlowModeCircle int32 `json:"slow_mode_circle"`
}

// ChannelInfo describes a channel belonging to a guild.
type ChannelInfo struct {
	// OwnerGuildID is the ID of the guild that owns the channel.
	OwnerGuildID string `json:"owner_guild_id"`
	// ChannelID is the channel ID.
	ChannelID string `json:"channel_id"`
	// ChannelType is the channel type code.
	ChannelType int32 `json:"channel_type"`
	// ChannelName is the channel's display name.
	ChannelName string `json:"channel_name"`
	// CreateTime is the channel creation time as a Unix timestamp.
	CreateTime int64 `json:"create_time"`
	// CreatorTinyID is the tiny ID of the channel creator.
	CreatorTinyID string `json:"creator_tiny_id"`
	// TalkPermission is the channel's speaking permission code.
	TalkPermission int32 `json:"talk_permission"`
	// VisibleType is the channel visibility code.
	VisibleType int32 `json:"visible_type"`
	// CurrentSlowMode is the current slow-mode key.
	CurrentSlowMode int32 `json:"current_slow_mode"`
	// SlowModes lists the channel's available slow-mode settings.
	SlowModes []SlowModeInfo `json:"slow_modes"`
}

// GuildChannels is a list of channels in a guild.
type GuildChannels []ChannelInfo

// GuildInfo describes a guild visible to the current account.
type GuildInfo struct {
	// GuildID is the guild ID.
	GuildID string `json:"guild_id"`
	// GuildName is the guild's display name.
	GuildName string `json:"guild_name"`
	// GuildDisplayID is the numeric display ID.
	GuildDisplayID int64 `json:"guild_display_id"`
}

// GuildList is a list of guilds visible to the current account.
type GuildList []GuildInfo

// RoleInfo describes a role assigned to a guild member.
type RoleInfo struct {
	// RoleID is the role ID.
	RoleID string `json:"role_id"`
	// RoleName is the role's display name.
	RoleName string `json:"role_name"`
}

// GuildMemberProfile describes a member's public profile and guild roles.
type GuildMemberProfile struct {
	// TinyID is the member's tiny ID.
	TinyID string `json:"tiny_id"`
	// Nickname is the member's display name.
	Nickname string `json:"nickname"`
	// AvatarURL is the member's avatar URL.
	AvatarURL string `json:"avatar_url"`
	// JoinTime is the member's join time as a Unix timestamp.
	JoinTime int64 `json:"join_time"`
	// Roles lists the roles assigned to the member.
	Roles []RoleInfo `json:"roles"`
}

// CQPacketType names a layer or packet type in the messaging protocol.
type CQPacketType string

// CQPacket exposes the identifier and nested layers of a messaging packet.
type CQPacket interface {
	// ID returns the packet's type identifier.
	ID() CQPacketType
	// Layers returns the packet's nested packet layers keyed by type.
	Layers() map[CQPacketType]CQPacket
}

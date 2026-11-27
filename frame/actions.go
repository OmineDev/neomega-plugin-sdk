package frame

import (
	"time"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// CmdCannotGetResponse is a command that can be sent without collecting its output.
type CmdCannotGetResponse interface {
	// Send sends the command and discards its output.
	Send()
}

// CmdResponseFor sends a command and returns its response as R.
type CmdResponseFor[R any] interface {
	CmdCannotGetResponse
	// SendAndGetResponse returns a single-use operation for sending and receiving the command output.
	SendAndGetResponse() AsyncResult[R]
}

// CmdCanGetResponse is a command whose reply is a Bedrock CommandOutput packet.
type CmdCanGetResponse = CmdResponseFor[*packet.CommandOutput]

// GeneralCommandFor can be sent through one of the supported command origins.
type GeneralCommandFor[C any] interface {
	// Send sends through the AI command interface.
	Send()
	// AsPlayer selects the player command origin.
	AsPlayer() C
	// AsWebSocket selects the WebSocket command origin.
	AsWebSocket() C
	// AsAI selects the AI command origin.
	AsAI() C
}

// GeneralCommand is a general command specialized for Bedrock CommandOutput replies.
type GeneralCommand = GeneralCommandFor[CmdCanGetResponse]

// CommandHelperFor constructs commands and exposes common block, item and structure operations. Generic
// arguments let the host supply its own position and response types.
type CommandHelperFor[P, Props, Reply, General any] interface {
	// ConstructDimensionLimitedWOCommand prepares a command that uses the settings origin in the
	// Overworld and the WebSocket origin elsewhere. Call Send to execute it.
	ConstructDimensionLimitedWOCommand(cmd string) CmdCannotGetResponse
	// ConstructDimensionLimitedGeneralCommand prepares cmd for the bot's current dimension.
	ConstructDimensionLimitedGeneralCommand(cmd string) General
	// ConstructGeneralCommand constructs a command without adding a dimension prefix.
	ConstructGeneralCommand(cmd string) General
	// GenAutoUnfilteredUUID creates a UUID-shaped name that can pass through the game's text filter.
	GenAutoUnfilteredUUID() string
	// ReplaceHotBarItemCmd constructs a command that replaces an item in a hotbar slot.
	ReplaceHotBarItemCmd(slotID int32, item string) Reply
	// ReplaceBotHotBarItemFullCmd constructs a hotbar replacement command with count, damage and item
	// properties.
	ReplaceBotHotBarItemFullCmd(slotID int32, itemName string, count uint8, value int32, components Props) Reply
	// ReplaceContainerBlockItemCmd constructs a command that replaces an item in a container block.
	ReplaceContainerBlockItemCmd(pos P, slotID int32, item string) Reply
	// ReplaceContainerItemFullCmd constructs a container replacement command with item properties.
	ReplaceContainerItemFullCmd(pos P, slotID int32, itemName string, count uint8, value int32, components Props) Reply
	// BackupStructureWithGivenNameCmd constructs a command that saves a structure region under name.
	BackupStructureWithGivenNameCmd(start P, size P, name string) Reply
	// BackupStructureWithAutoNameCmd creates a name and constructs the matching structure-save command.
	BackupStructureWithAutoNameCmd(start P, size P) (name string, cmd Reply)
	// RevertStructureWithGivenNameCmd constructs a command that loads a saved structure at start.
	RevertStructureWithGivenNameCmd(start P, name string) Reply
	// SetBlockCmd constructs a dimension-aware command that sets one block.
	SetBlockCmd(pos P, blockString string) General
	// SetBlockRelativeCmd constructs a dimension-aware command that sets a block relative to the bot.
	SetBlockRelativeCmd(pos P, blockString string) General
	// FillBlocksWithRangeCmd constructs a dimension-aware fill command between inclusive start and end
	// positions.
	FillBlocksWithRangeCmd(startPos P, endPos P, blockString string) General
	// FillBlocksWithSizeCmd constructs a dimension-aware fill command. The size is converted to an inclusive
	// end position.
	FillBlocksWithSizeCmd(startPos P, size P, blockString string) General
}

// CommandHelper is the command helper specialized for SDK positions and item properties.
type CommandHelper = CommandHelperFor[CubePos, *ItemPropsInGiveOrReplace, CmdCanGetResponse, GeneralCommand]

// BotActionFor provides low-level inventory and block operations. Generic arguments keep host-specific
// value types out of the SDK.
type BotActionFor[P, Settings, Inventory any] interface {
	// SelectHotBar selects one of the bot's hotbar slots.
	SelectHotBar(slotID uint8) error
	// SleepTick waits for the requested number of game ticks (nominally 50 ms each).
	SleepTick(ticks int)
	// UseHotBarItem uses the item in slot.
	UseHotBarItem(slot uint8) (err error)
	// DropItemFromHotBar drops the item in slot.
	DropItemFromHotBar(slot uint8) error
	// MoveItemInsideHotBarOrInventory moves count items between inventory slots.
	MoveItemInsideHotBarOrInventory(sourceSlot, targetSlot, count uint8) error
	// UseHotBarItemOnBlock uses the item in slot on the block at blockPos.
	UseHotBarItemOnBlock(blockPos P, blockNEMCRuntimeID uint32, face int32, slot uint8) (err error)
	// UseHotBarItemOnBlockWithBotOffset uses the item on the given face of blockPos, reporting
	// blockPos + botOffset as the bot's position for the interaction.
	UseHotBarItemOnBlockWithBotOffset(blockPos P, botOffset P, blockNEMCRuntimeID uint32, face int32, slot uint8) (err error)
	// MoveItemFromInventoryToEmptyContainerSlots moves items using the supplied source-to-destination slot
	// map.
	MoveItemFromInventoryToEmptyContainerSlots(pos P, blockNemcRtid uint32, blockName string, moveOperations map[uint8]uint8) error
	// UseAnvil uses the anvil at pos to rename the item in slot.
	UseAnvil(pos P, blockNemcRtid uint32, slot uint8, newName string) error
	// SetStructureBlockData applies settings to a structure block.
	SetStructureBlockData(pos P, settings Settings)
	// GetInventoryContent returns an item instance for windowID and slotID, if the host has that slot.
	GetInventoryContent(windowID uint32, slotID uint8) (instance Inventory, found bool)
	// WeaveBannerLayer adds a layer using the currently open loom. The banner and dye
	// occupy single-item stacks. A negative templateSlot selects a pattern without a template.
	WeaveBannerLayer(bannerSlot, dyeSlot byte, templateSlot int32, pattern string, color int32) error
	// DecorateShield applies the banner in bannerSlot to the shield in shieldSlot.
	DecorateShield(bannerSlot, shieldSlot byte) error
	// TrimArmour applies a trim using the currently open smithing table.
	TrimArmour(armourSlot, materialSlot, templateSlot byte, pattern, material string) error
	// TakeOminousBanner puts an ominous banner from the creative inventory in an empty slot.
	TakeOminousBanner(slot byte) error
	// TakeCreativeItem takes count items from a server-advertised creative entry into an empty slot.
	TakeCreativeItem(creativeID uint32, slot, count byte) error
	// CraftGrid crafts a server-advertised recipe in the open crafting grid.
	CraftGrid(recipe string, inputs []GridInput, resultSlot byte) error
	// TransferStack moves count items, swaps two stacks, or drops items according to mode.
	TransferStack(from, to StackSlot, count byte, mode StackTransfer) error
	// MoveStacks applies the ordered moves in one item-stack request.
	MoveStacks(moves []StackMove) error
}

// BotAction is the low-level action set specialized for SDK positions, structure settings and Bedrock item
// instances.
type BotAction = BotActionFor[CubePos, *StructureBlockSupportedData, *protocol.ItemInstance]

// BotActionHighLevelFor provides multi-step bot operations. These methods may issue commands, move the bot
// or change world state.
type BotActionHighLevelFor[P, Sign, Block, Enchantments, Action, Slot, ItemT, Provider any] interface {
	// HighLevelWriteBook writes pages into the book in slotID.
	HighLevelWriteBook(slotID uint8, pages []string) (err error)
	// HighLevelWriteBookAndClose writes pages, sets title and author, and closes the book.
	HighLevelWriteBookAndClose(slotID uint8, pages []string, bookTitle string, bookAuthor string) (err error)
	// HighLevelChangeFlyingStates requests the given flying state for the bot.
	HighLevelChangeFlyingStates(isFlying bool) error
	// HighLevelPlaceSign places a sign block and writes its front and back text.
	HighLevelPlaceSign(targetPos P, signBlock string, data Sign) (err error)
	// HighLevelPlaceCommandBlock places and configures a command block, retrying placement up to maxRetry
	// times.
	HighLevelPlaceCommandBlock(targetPos P, option Block, maxRetry int) error
	// HighLevelMoveItemToContainer moves items into the container at pos using the supplied slot map.
	HighLevelMoveItemToContainer(pos P, moveOperations map[uint8]uint8) error
	// HighLevelEnsureBotNearby moves the bot near pos when its distance exceeds threshold blocks.
	HighLevelEnsureBotNearby(pos P, threshold float32) error
	// HighLevelRemoveSpecificBlockSideEffect backs up the block at pos and temporarily replaces it. Call the
	// returned function to restore it.
	HighLevelRemoveSpecificBlockSideEffect(pos P, backupName string) (deferFunc func(), err error)
	// HighLevelRenameItemWithAnvil renames the item in slot using an anvil, creating an anvil when
	// autoGenAnvil is true.
	HighLevelRenameItemWithAnvil(pos P, slot uint8, newName string, autoGenAnvil bool) (err error)
	// HighLevelEnchantItem applies enchants to the item in slot.
	HighLevelEnchantItem(slot uint8, enchants Enchantments) (err error)
	// HighLevelListenItemPicked waits up to timeout for an inventory pickup action. The returned cancel
	// function stops the wait and releases the bot for other work.
	HighLevelListenItemPicked(timeout time.Duration) (actionChan chan Action, cancel func(), err error)
	// HighLevelPickBlock is equivalent to PickBlock with withNBT set to true.
	HighLevelPickBlock(pos P, targetHotBar uint8, retryTimes int) error
	// PickBlock picks the block at pos into targetHotBar, requesting block NBT only if withNBT is true.
	PickBlock(pos P, targetHotBar uint8, withNBT bool, retryTimes int) error
	// HighLevelBlockBreakAndPickInHotBar breaks the block at pos and optionally restores it, then places the
	// drop in targetSlot.
	HighLevelBlockBreakAndPickInHotBar(pos P, recoverBlock bool, targetSlot uint8, maxRetriesTotal int) (err error)
	// HighLevelSetContainerContent replaces the contents of the container at pos.
	HighLevelSetContainerContent(pos P, containerInfo map[uint8]Slot) (err error)
	// HighLevelGenContainer places block at pos and fills it with containerInfo.
	HighLevelGenContainer(pos P, containerInfo map[uint8]Slot, block string) (err error)
	// HighLevelPlaceItemFrameItem places the item in slotID into the frame at pos, with the requested
	// rotation.
	HighLevelPlaceItemFrameItem(pos P, slotID uint8, rotation int) error
	// HighLevelMakeItem crafts or prepares item using the anvil and next-container positions as needed.
	HighLevelMakeItem(item ItemT, slotID uint8, anvilPos, nextContainerPos P) error
	// HighLevelRequestLargeArea requests a region of blocks and writes the result to dst. It may move the bot
	// when withMove is true.
	HighLevelRequestLargeArea(startPos P, size P, dst Provider, withMove bool) error
}

// BotActionHighLevel is the high-level action set specialized for SDK data types.
type BotActionHighLevel = BotActionHighLevelFor[CubePos, *SignBlockSupportedData, *CommandBlockSupportedData, Enchants, protocol.InventoryAction, *ContainerSlotItemStack, *Item, ChunkProvider]

// BotActionComplex combines command helpers with low- and high-level bot actions.
type BotActionComplex interface {
	BotAction
	BotActionHighLevel
	CommandHelper
}

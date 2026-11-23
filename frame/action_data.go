package frame

// ItemLockPlace is the protection applied when giving or replacing an item.
type ItemLockPlace string

const (
	// ItemLockPlaceInventory locks the item against moving it out of the inventory.
	ItemLockPlaceInventory ItemLockPlace = "lock_in_inventory"
	// ItemLockPlaceSlot locks the item to its inventory slot.
	ItemLockPlaceSlot ItemLockPlace = "lock_in_slot"
)

// Enchant identifies an enchantment by its game ID.
type Enchant int32

// Enchants maps enchantment IDs to levels.
type Enchants map[Enchant]int32

// CommandBlockSupportedData contains the command and execution settings written to a command block.
type CommandBlockSupportedData struct {
	// BlockName is the command block's block identifier.
	BlockName string
	// BlockState contains the block-state string required to place the block.
	BlockState string
	// NeedsRedstone reports whether the block needs a redstone signal to execute.
	NeedsRedstone bool `json:"NeedRedStone" mapstructure:"NeedRedStone"`
	// Conditional reports whether execution depends on the previous block succeeding.
	Conditional bool
	// Command is the command stored in the block.
	Command string
	// Name is the command block's custom name.
	Name string
	// TickDelay is the delay, in ticks, before execution.
	TickDelay int
	// ShouldTrackOutput controls whether the block stores its last output.
	ShouldTrackOutput bool
	// ExecuteOnFirstTick controls whether it executes on its first tick after activation.
	ExecuteOnFirstTick bool
}

// StructureBlockSupportedData contains the settings written to a structure block.
type StructureBlockSupportedData struct {
	// DataField selects the structure block's data mode.
	DataField string `mapstructure:"dataField"`
	// IgnoreEntities controls whether entities are omitted when saving or loading.
	IgnoreEntities uint8 `mapstructure:"ignoreEntities"`
	// IncludePlayers controls whether players are included in the operation.
	IncludePlayers uint8 `mapstructure:"includePlayers"`
	// Integrity is the fraction of blocks to keep when applying the structure.
	Integrity float32 `mapstructure:"integrity"`
	// Mirror selects the structure's mirror operation.
	Mirror uint8 `mapstructure:"mirror"`
	// RedstoneSaveMode selects how redstone affects saving.
	RedstoneSaveMode int32 `mapstructure:"redstoneSaveMode"`
	// IgnoreBlocks controls whether blocks are removed during the operation.
	IgnoreBlocks int32 `mapstructure:"removeBlocks"`
	// Rotation selects the structure's rotation.
	Rotation uint8 `mapstructure:"rotation"`
	// ShowBoundingBox controls the visible bounding box.
	ShowBoundingBox uint8 `mapstructure:"showBoundingBox"`
	// AnimationMode selects the structure animation mode.
	AnimationMode uint8 `mapstructure:"animationMode"`
	// AnimationDuration is the animation duration in seconds.
	AnimationDuration float32 `mapstructure:"animationSeconds"`
	// Seed is the seed used when applying integrity.
	Seed int64 `mapstructure:"seed"`
	// XStructureSize is the structure size along X.
	XStructureSize int32 `mapstructure:"xStructureSize"`
	// YStructureSize is the structure size along Y.
	YStructureSize int32 `mapstructure:"yStructureSize"`
	// ZStructureSize is the structure size along Z.
	ZStructureSize int32 `mapstructure:"zStructureSize"`
	// XStructureOffset is the structure offset along X.
	XStructureOffset int32 `mapstructure:"xStructureOffset"`
	// YStructureOffset is the structure offset along Y.
	YStructureOffset int32 `mapstructure:"yStructureOffset"`
	// ZStructureOffset is the structure offset along Z.
	ZStructureOffset int32 `mapstructure:"zStructureOffset"`
	// StructureName is the name of the saved structure.
	StructureName string `mapstructure:"structureName"`
	// StructureBlockType identifies the structure block mode.
	StructureBlockType int32
}

// SignBlockText contains one side's text and formatting settings.
type SignBlockText struct {
	// HideGlowOutline controls the glowing text outline.
	HideGlowOutline uint8
	// IgnoreLighting controls whether text ignores scene lighting.
	IgnoreLighting uint8
	// PersistFormatting controls whether formatting is preserved.
	PersistFormatting uint8
	// SignTextColor is the sign text color encoded as an integer.
	SignTextColor int32
	// Text is the text displayed on this side of the sign.
	Text string
	// TextOwner         string
}

// SignBlockSupportedDataFor combines front text, back text and the waxed state.
type SignBlockSupportedDataFor[T any] struct {
	// FrontText and BackText contain the text and formatting for the sign faces.
	FrontText, BackText T
	// IsWaxed is nonzero when the sign is waxed.
	IsWaxed uint8
}

// SignBlockSupportedData contains front and back sign settings using SDK text types.
type SignBlockSupportedData SignBlockSupportedDataFor[SignBlockText]

// ItemPropsInGiveOrReplace contains the placement, break, lock and death-retention properties for an item.
type ItemPropsInGiveOrReplace struct {
	// CanPlaceOn limits which blocks can be placed against.
	CanPlaceOn []string `json:"can_place_on,omitempty"`
	// CanDestroy lists blocks the item can break.
	CanDestroy []string `json:"can_destroy,omitempty"`
	// ItemLock controls whether the item is locked in the inventory or slot.
	ItemLock ItemLockPlace `json:"item_lock,omitempty"`
	// KeepOnDeath keeps the item when its holder dies.
	KeepOnDeath bool `json:"keep_on_death,omitempty"`
}

// ItemFor describes an item and its optional game-specific data. Generic parameters let hosts use their
// native data types.
type ItemFor[Props, Specific, Enchantments, Complex any] struct {
	// Name is the item's namespaced identifier.
	Name string `json:"name"`
	// Value is the item's legacy damage or variant value.
	Value int32 `json:"val"`
	// BaseProps contains optional placement and inventory properties.
	BaseProps Props `json:"base_props,omitempty"`
	// IsBlock reports whether the item can be placed as a block.
	IsBlock bool `json:"is_block"`
	// RelatedBlockBedrockStateString contains the Bedrock state for the placeable block.
	RelatedBlockBedrockStateString string `json:"block_bedrock_state_string,omitempty"`
	// END Basic Item
	// SpecificKnownNonBlockItemData contains details for supported items such as books.
	SpecificKnownNonBlockItemData Specific `json:"specific_known_item_data,omitempty"`
	// End Known Item
	// Enchants contains the item's enchantments.
	Enchants Enchantments `json:"enchants,omitempty"`
	// End NeedHotBarItem
	// DisplayName is the item's custom name.
	DisplayName string `json:"display_name,omitempty"`
	// End NeedAuxItem
	// RelateComplexBlockData contains block data stored with the item, such as container contents.
	RelateComplexBlockData Complex `json:"complex_block_data,omitempty"`
	// End Complex Block
}

// Item is the SDK item type, specialized for SDK item properties and container data.
type Item ItemFor[*ItemPropsInGiveOrReplace, *SpecificKnownNonBlockItemData, Enchants, *ComplexBlockData]

// ContainerSlotItemStackFor pairs an item with its count in one container slot.
type ContainerSlotItemStackFor[I any] struct {
	// Item is the item stored in the slot.
	Item I `json:"item"`
	// Count is the number of items in the slot.
	Count uint8 `json:"count"`
}

// ContainerSlotItemStack is a container slot using the SDK Item type.
type ContainerSlotItemStack ContainerSlotItemStackFor[*Item]

// ComplexBlockDataFor describes optional structured data embedded in a block item.
type ComplexBlockDataFor[Slot any] struct {
	// Container maps container slot IDs to their item stacks.
	Container map[uint8]Slot `json:"container,omitempty"`
	// Unknown preserves additional named data not represented by known fields.
	Unknown map[string]any `json:"unknown_nbt,omitempty"`
}

// ComplexBlockData is complex block data using SDK container item types.
type ComplexBlockData ComplexBlockDataFor[*ContainerSlotItemStack]

// SpecificKnownNonBlockItemData holds additional data for supported non-block items.
type SpecificKnownNonBlockItemData struct {
	// Pages contains the pages of a book.
	Pages []string `json:"pages,omitempty"`
	// BookAuthor is the author stored in a written book.
	BookAuthor string `json:"book_author,omitempty"`
	// BookName is the title stored in a written book.
	BookName string `json:"book_name,omitempty"`
	// Unknown preserves additional NBT fields for this item.
	Unknown map[string]any `json:"known_nbt,omitempty"`
}

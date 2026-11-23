package frame

import "time"

// Chunk exposes decoded block storage in one chunk column.
type Chunk interface {
	// Range returns the inclusive world-height range represented by this chunk.
	Range() [2]int
	// Layers returns the number of block layers at block height y.
	Layers(y int16) int
	// Block returns the runtime ID at local chunk coordinates x,z, world height y and layer.
	Block(x uint8, y int16, z, layer uint8) uint32
	// SetBlock writes a runtime ID at local chunk coordinates x,z, world height y and layer.
	SetBlock(x uint8, y int16, z, layer uint8, id uint32)
}

// ChunkWithAuxInfo combines block storage with block actor data, its coordinate and last sync time.
type ChunkWithAuxInfo struct {
	// AirRuntimeID is the block runtime ID used when no block storage is present.
	AirRuntimeID uint32
	// Chunk is the decoded block storage. It is nil when the chunk contains no decoded storage.
	Chunk Chunk
	// BlockNbts maps absolute block positions to block actor data.
	BlockNbts map[CubePos]map[string]any
	// SyncTime is the Unix time when the chunk was last synchronized.
	SyncTime int64
	// ChunkPos is the chunk's X,Z coordinate.
	ChunkPos ChunkPos
}

// ChunkProviderFor reads and writes chunks using host-selected position and data types.
type ChunkProviderFor[P, C any] interface {
	// Get returns the chunk at pos, or nil when it is not available.
	Get(P) C
	// Save writes data to the provider.
	Save(C) error
}

// ChunkProvider is the chunk provider specialized for SDK chunk coordinates and chunk data.
type ChunkProvider = ChunkProviderFor[ChunkPos, *ChunkWithAuxInfo]

// GetBlock returns the foreground block at the absolute position, or AirRuntimeID if storage is absent.
func (c *ChunkWithAuxInfo) GetBlock(p CubePos) uint32 {
	if c.Chunk == nil {
		return c.AirRuntimeID
	}
	return c.Chunk.Block(uint8(p.X()&15), int16(p.Y()), uint8(p.Z()&15), 0)
}

// GetBlockWithNbt returns the foreground block and its block actor data at pos.
func (c *ChunkWithAuxInfo) GetBlockWithNbt(p CubePos) (uint32, map[string]any) {
	return c.GetBlock(p), c.BlockNbts[p]
}

// GetSyncTime converts SyncTime to a time.Time.
func (c *ChunkWithAuxInfo) GetSyncTime() time.Time { return time.Unix(c.SyncTime, 0) }

// SetSyncTime stores t as a Unix timestamp.
func (c *ChunkWithAuxInfo) SetSyncTime(t time.Time) { c.SyncTime = t.Unix() }

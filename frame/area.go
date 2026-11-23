package frame

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// SubChunk exposes decoded block layers for an area-result callback.
type SubChunk interface {
	// Layers returns the number of block layers stored in this subchunk.
	Layers() int
	// Block returns the runtime ID at local X,Y,Z and layer.
	Block(x, y, z, layer uint8) uint32
	// SetBlock writes a runtime ID at local X,Y,Z and layer.
	SetBlock(x, y, z, layer uint8, block uint32)
}

// DecodedStructureFor reads a decoded structure using generic position and chunk-provider types.
type DecodedStructureFor[P comparable, Provider any] interface {
	// IndexOfRelativePos returns the flat block index for a position relative to the structure origin.
	IndexOfRelativePos(pos P) int
	// BlockOfRelativePos returns the foreground and background runtime IDs at a relative position.
	BlockOfRelativePos(pos P) (foreGround, backGround uint32)
	// NBTsInAbsolutePos returns block actor NBT keyed by absolute block position.
	NBTsInAbsolutePos() map[P]map[string]any
	// ForeGroundRtidNested returns foreground runtime IDs in the structure's storage order.
	ForeGroundRtidNested() []uint32
	// BackGroundRtidNested returns background runtime IDs in the structure's storage order.
	BackGroundRtidNested() []uint32
	// DumpToChunkProviderAbsolutePos writes the decoded structure at its absolute origin to chunkProvider.
	DumpToChunkProviderAbsolutePos(chunkProvider Provider) (err error)
}

// DecodedStructure specializes the structure API for CubePos and ChunkProvider.
type DecodedStructure = DecodedStructureFor[CubePos, ChunkProvider]

// StructureResponseFor carries a raw structure packet and can decode it into Decoded.
type StructureResponseFor[RawPacket, Decoded any] interface {
	// Raw returns the underlying structure response packet.
	Raw() RawPacket
	// Decode decodes the response into the requested structure representation.
	Decode() (Decoded, error)
}

// StructureResponse returns a Bedrock StructureTemplateDataResponse and a decoded SDK structure.
type StructureResponse = StructureResponseFor[*packet.StructureTemplateDataResponse, DecodedStructure]

// SubChunkResultFor describes one response for a requested subchunk. Its generic types describe the
// position, subchunk coordinate, chunk coordinate and decoded subchunk.
type SubChunkResultFor[SubPos any, CP comparable, P comparable, Sub any] interface {
	// SubChunkPos returns the absolute subchunk position.
	SubChunkPos() SubPos
	// ChunkPos returns the chunk coordinate containing this subchunk.
	ChunkPos() CP
	// Error returns the request or decode error, if one occurred.
	Error() error
	// AttachDecodeError attaches an error found while decoding the subchunk data.
	AttachDecodeError(error)
	// ResultCode returns the server's subchunk result code.
	ResultCode() byte
	// NBTsInAbsolutePos returns block actor NBT keyed by absolute block position.
	NBTsInAbsolutePos() map[P]map[string]any
	// SubChunk returns the decoded subchunk block storage.
	SubChunk() Sub
}

// SubChunkResult specializes the result for SDK positions and SubChunk.
type SubChunkResult = SubChunkResultFor[protocol.SubChunkPos, ChunkPos, CubePos, SubChunk]

// SubChunkBatchResultFor collects the results for a group of subchunk requests.
type SubChunkBatchResultFor[SP comparable, CP comparable, P comparable, Result any, Sub any, ChunkData any] interface {
	// MapResults returns results indexed by subchunk position.
	MapResults() map[SP]Result
	// Results returns all results as a slice.
	Results() []Result
	// AllOk reports whether all results completed without an error.
	AllOk() bool
	// AllErrors includes every requested position; its value is nil for a successful result.
	AllErrors() map[SP]error
	// ToChunks assembles the results into chunks. If optionalAlterFn is non-nil, it can alter each decoded
	// subchunk and its NBT before assembly.
	ToChunks(
		optionalAlterFn func(r Result) (Sub, map[P]map[string]any, error),
	) map[CP]ChunkData
}

// SubChunkBatchResult is the batch result for SDK positions, decoded chunks and block NBT.
type SubChunkBatchResult = SubChunkBatchResultFor[protocol.SubChunkPos, ChunkPos, CubePos, SubChunkResult, SubChunk, *ChunkWithAuxInfo]

// SubChunkBatchReqHandlerFor builds a subchunk request by dimension and coordinate ranges.
type SubChunkBatchReqHandlerFor[D, Result any] interface {
	// OmitResult starts the request without retaining a result handle.
	OmitResult()
	// GetResult returns a single-use operation that sends the request and collects the batch.
	GetResult() AsyncResult[Result]
	// AutoDimension uses the bot's current dimension when resolving the request.
	AutoDimension() SubChunkBatchReqHandlerFor[D, Result]
	// InDimension selects the dimension for the request.
	InDimension(dim D) SubChunkBatchReqHandlerFor[D, Result]
	// X requests one X offset from the base chunk.
	X(xOffset int8) SubChunkBatchReqHandlerFor[D, Result]
	// Z requests one Z offset from the base chunk.
	Z(zOffset int8) SubChunkBatchReqHandlerFor[D, Result]
	// Y selects an absolute subchunk Y index, not a block height.
	Y(zOffset int8) SubChunkBatchReqHandlerFor[D, Result]
	// XRange requests subchunks from startOffset up to, but not including, endNotIncludedOffset on X.
	XRange(startOffset int8, endNotIncludedOffset int8) SubChunkBatchReqHandlerFor[D, Result]
	// YRange requests subchunks from startOffset up to, but not including, endNotIncludedOffset on Y.
	YRange(startOffset int8, endNotIncludedOffset int8) SubChunkBatchReqHandlerFor[D, Result]
	// ZRange requests subchunks from startOffset up to, but not including, endNotIncludedOffset on Z.
	ZRange(startOffset int8, endNotIncludedOffset int8) SubChunkBatchReqHandlerFor[D, Result]
	// FullY requests every supported subchunk height for the selected dimension.
	FullY() SubChunkBatchReqHandlerFor[D, Result]
}

// SubChunkBatchReqHandler specializes the batch request for SDK dimensions and result types.
type SubChunkBatchReqHandler = SubChunkBatchReqHandlerFor[Dimension, SubChunkBatchResult]

// LowLevelAreaRequesterFor makes structure and subchunk requests without frequency, size or movement
// controls.
type LowLevelAreaRequesterFor[P, CP, Structure, Batch, Result any] interface {
	// LowLevelRequestStructure requests a structure at pos with the given size and name.
	LowLevelRequestStructure(pos P, size P, structureName string) AsyncResult[Structure]
	// LowLevelRequestStructureWithAutoName requests a structure at pos and size, choosing its name
	// automatically.
	LowLevelRequestStructureWithAutoName(pos P, size P) AsyncResult[Structure]
	// LowLevelRequestChunk prepares a batch at baseChunkPos. Select the dimension and offsets before
	// calling OmitResult or consuming GetResult.
	LowLevelRequestChunk(baseChunkPos CP) Batch
	// SetOnSubChunkResult sets a callback for subchunk results as they arrive.
	SetOnSubChunkResult(nonBlockingCallback func(Result))
	// AttachSubChunkResultListener subscribes to subchunk results and returns a function that detaches the
	// listener.
	AttachSubChunkResultListener(nonBlockingCallback func(Result)) (detachFn func())
}

// LowLevelAreaRequester specializes the requester for SDK positions, response packets and chunks.
type LowLevelAreaRequester = LowLevelAreaRequesterFor[CubePos, ChunkPos, StructureResponse, SubChunkBatchReqHandler, SubChunkResult]

package frame

import "github.com/sandertv/gophertunnel/minecraft/protocol"

// BlobHashHolderInfo reports the roles this node currently has in the blob cache.
type BlobHashHolderInfo interface {
	// IsServer 返回自身是否是 blob cache 缓存数据集的权威持有人(服务者)
	IsServer() bool
	// IsDiskHolder 返回当前是否是镜像存档资源的持有人
	IsDiskHolder() bool
}

// BlobStore reads and updates payloads in the local blob cache.
type BlobStore interface {
	BlobHashHolderInfo
	// LoadBlobCache 从底层缓存数据集检索 hash 所指示的数据负载。
	// 如果不存在，则返回 nil
	LoadBlobCache(hash uint64) []byte
	// UpdateBlobCache 将 hash 与 payLoad 关联起来，
	// 并将它们放入底层的缓存数据集合。
	//
	// 保证 UpdateBlobCache 在返回值后，UpdateBlobCache
	// 的调用者可以安全的修改 payload 所指示的切片。
	//
	// UpdateBlobCache 的调用者有责任校验来自服务器的 hash
	// 是否与 xxhash.Sum64(payload) 等价。
	//
	// 如果 UpdateBlobCache 返回假，则说明使用者发送了不正确的
	// (hash, payload) 对，此时底层 blob hash cache 缓存数据
	// 集将拒绝接受并不发生任何变化。
	//
	// 作为一种特殊情况，长度为 0 的数据荷载不会被接受，
	// UpdateBlobCache 会永远为它返回假
	UpdateBlobCache(hash uint64, payload []byte) bool
}

// BlobHash is the 64-bit xxHash of a subchunk payload.
type BlobHash uint64

// HashWithPosition identifies a cached payload by hash, absolute subchunk position and dimension.
type HashWithPosition struct {
	// Hash is the payload's 64-bit xxHash.
	Hash BlobHash
	// SubChunkPos is the absolute subchunk coordinate. X and Z are chunk coordinates; Y is a subchunk index.
	SubChunkPos protocol.SubChunkPos
	// Dimension is the dimension code used by the blob cache protocol.
	Dimension uint8
}

// PayloadByHash pairs a subchunk identity with its encoded payload bytes.
type PayloadByHash struct {
	// Hash identifies the payload and its world position.
	Hash HashWithPosition
	// Payload is the encoded subchunk data.
	Payload []byte
}

// BlobQuery looks up payloads in the mirrored world cache.
type BlobQuery interface {
	// QueryDiskHashExist 向镜像存档持有人发起检索请求，
	// 目的仅在于检索 hashes 是否命中镜像存档中的存储
	QueryDiskHashExist([]HashWithPosition) (hit, miss []HashWithPosition)
	// GetDiskHashPayload 从镜像存档持有人请求
	// hashes 的二进制数据荷载。
	// hit 指示命中的部分，miss 指示未命中的部分
	GetDiskHashPayload([]HashWithPosition) (hit []PayloadByHash, miss []HashWithPosition)
}

// BlobServer exposes the query and synchronization operations available on a cache server.
type BlobServer interface {
	BlobQuery
	// RequireSyncHashToDisk asks the mirror to persist the supplied payloads.
	RequireSyncHashToDisk([]PayloadByHash)
}

// HandleQueryDiskHashExist handles a request to check which hashes are present in the mirror.
type HandleQueryDiskHashExist func([]HashWithPosition) []bool

// HandleGetDiskHashPayload handles a request for payloads stored in the mirror.
type HandleGetDiskHashPayload func([]HashWithPosition) []PayloadByHash

// HandleRequireSyncHashToDisk handles a request to persist payloads in the mirror.
type HandleRequireSyncHashToDisk func([]PayloadByHash)

// HandleCleanBlobHashAndApplyToWorld handles updates that remove cached hashes and change the mirrored
// world.
type HandleCleanBlobHashAndApplyToWorld func([]HashWithPosition)

// HandleServerDisconnect handles loss of the cache server connection.
type HandleServerDisconnect func()

// BlobMirrorHandlers installs the callbacks used to serve requests for mirrored world data.
type BlobMirrorHandlers interface {
	// SetHandler installs callbacks for existence checks, payload reads, disk writes, world updates and server
	// disconnects.
	SetHandler(HandleQueryDiskHashExist, HandleGetDiskHashPayload, HandleRequireSyncHashToDisk, HandleCleanBlobHashAndApplyToWorld, HandleServerDisconnect)
}

// BlobMirror is the mirror-side view of blob cache operations.
type BlobMirror interface {
	BlobMirrorHandlers
	// RegisterListener registers the packet listeners required to keep the mirror synchronized. Call it at
	// most once.
	RegisterListener()
}

// BlobCacheFor combines local storage, cache queries and optional server- and mirror-side interfaces.
type BlobCacheFor[M any] interface {
	BlobStore
	BlobQuery
	// GetHashPayload fetches payloads for the requested hashes and adds valid results to the local cache.
	GetHashPayload([]HashWithPosition) map[HashWithPosition][]byte
	// AsServerSide returns the server operations, or nil when this node is not the server.
	AsServerSide() BlobServer
	// AsMirrorWorldSide returns the mirror-side operations provided by the host.
	AsMirrorWorldSide() M
	// SetHolderRequest asks the server to make this node the mirror holder. A successful request cannot be
	// revoked while the server is running.
	SetHolderRequest() bool
}

// BlobCache is the host blob cache specialized for the SDK's BlobMirror interface.
type BlobCache = BlobCacheFor[BlobMirror]

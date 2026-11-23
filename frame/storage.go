package frame

// KVDBLike stores string values under string keys.
type KVDBLike interface {
	// Get returns the value for key, or an empty string if it is not present.
	Get(key string) (value string)
	// Delete removes key from the store.
	Delete(key string)
	// Set stores value under key, replacing any previous value.
	Set(key string, value string)
	// Iter calls the function for each entry. Return false to stop iterating.
	Iter(func(key, value string) bool)
}

// StorageAndPathAccess manages component files under the host's storage directories.
type StorageAndPathAccess interface {
	// GetLoggerPath returns the log path for topic.
	GetLoggerPath(topic string) string
	// GetFileData reads the named file from the data directory.
	GetFileData(topic string) ([]byte, error)
	// GetJsonData decodes the named JSON file in the data directory into data. Pass a pointer to the
	// destination.
	GetJsonData(topic string, data any) error
	// WriteFileData writes data to the named file in the data directory.
	WriteFileData(topic string, data []byte) error
	// WriteJsonData encodes data as JSON and writes it to the named data file.
	WriteJsonData(topic string, data any) error
	// WriteJsonDataWithTMP writes JSON to topic + tmpSuffix, then renames the temporary file to topic.
	WriteJsonDataWithTMP(topic string, tmpSuffix string, data any) error
	// GetFilePath joins elem under the data directory.
	GetFilePath(elem ...string) string
	// GetOmegaCachePath joins elem under the host's cache directory.
	GetOmegaCachePath(elem ...string) string
	// GetArchivePath joins elem under the archive directory.
	GetArchivePath(elem ...string) string
	// GetConfigPath joins elem under the configuration directory.
	GetConfigPath(elem ...string) string
	// NewTempDir returns a fresh directory path under the host's temporary storage.
	NewTempDir() string
	// GetLangSpecificPath joins elem under the directory for language-specific plugin files.
	GetLangSpecificPath(elem ...string) string
	// GetKVDBLike opens a key-value store under the data directory using the requested backend.
	GetKVDBLike(saveDir string, dbType string) (KVDBLike, error)
}

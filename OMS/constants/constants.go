package constants

const (
	ErrInvalidID             = "Invalid ID"
	ErrHubNotFound           = "Hub Not Found"
	ErrSKUNotFound           = "SKU Not Found"
	ErrServerError           = "Server Error"
	ErrParsingJSON           = "Issue While Parsing JSON"
	ErrRecordExists          = "Record Already Exists"
	ErrHubCreate             = "Failed to create hub"
	ErrHubUpdate             = "Failed to update hub"
	ErrSKUCreate             = "Failed to create SKU"
	ErrSKUUpdate             = "Failed to update SKU"
	ErrGetAllHubs            = "Failed to get all hubs"
	ErrGetAllSKUs            = "Failed to get all SKUs"
	ErrInventoryNotFound     = "Inventory Not Found"
	ErrInvalidRequest        = "Invalid Request"
	ErrInsufficientInventory = "Insufficient Inventory"
	ErrBatchOperation        = "Batch Operation Failed"

	RedisHost        = "localhost:6379"
	RedisPoolSize    = 50
	RedisMinIdleConn = 10

	ServerPort         = ":8000"
	ServerReadTimeout  = 10
	ServerWriteTimeout = 10
	ServerIdleTimeout  = 70

	CacheTTLHubs = 5
	CacheTTLSKUs = 5

	CacheKeyHubID   = "hub:id:"
	CacheKeyHubName = "hub:name:"
	CacheKeySKUID   = "sku:id:"
	CacheKeySKUCode = "sku:code:"

	DefaultQuantity = 0

	InventoryOpUpsert = "upsert"
	InventoryOpReduce = "reduce"
	InventoryOpView   = "view"

	CacheKeySKUValidation = "sku:validation:"
	CacheKeyHubValidation = "hub:validation:"
	CacheKeyInventoryView = "inventory:view:"

	MaxBatchSize     = 1000
	DefaultBatchSize = 100

	ErrInventoryUpsert = "Failed to upsert inventory"
	ErrInventoryReduce = "Failed to reduce inventory"
	ErrInventoryView   = "Failed to view inventory"
	ErrAtomicOperation = "Atomic operation failed"
)

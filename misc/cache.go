package misc

import (
	"context"
	"log"
	"time"

	"github.com/allegro/bigcache/v3"
)

var Cache *bigcache.BigCache

func init() {
	config := bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: 1024,
		// time after which entry can be evicted
		LifeWindow: 1 * time.Minute,
		// Interval between removing expired entries (clean up).
		// If set to <= 0 then no action is performed.
		// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
		CleanWindow: 1 * time.Second,
		// rps * lifeWindow, used only in initial memory allocation
		MaxEntriesInWindow: 1000 * 10 * 60,
		// max entry size in bytes, used only in initial memory allocation
		MaxEntrySize: 500,
		// prints information about additional memory allocation
		Verbose: true,
		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		HardMaxCacheSize: 8192,
		// callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		OnRemove: nil,
		// OnRemoveWithReason is a callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A constant representing the reason will be passed through.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		// Ignored if OnRemove is specified.
		OnRemoveWithReason: func(key string, entry []byte, reason bigcache.RemoveReason) {
			var rea string
			switch reason {
			case bigcache.RemoveReason(1):
				rea = "Expired"
			case bigcache.RemoveReason(2):
				rea = "NoSpace"
			case bigcache.RemoveReason(3):
				rea = "Deleted"
			}
			Logger.Info().Str("from", "cache").Str("key", key).Str("entry", string(entry)).Str("reason", rea).Msg("")
		},
	}
	cache, initErr := bigcache.New(context.Background(), config)
	if initErr != nil {
		log.Fatal(initErr)
	}
	Cache = cache
}

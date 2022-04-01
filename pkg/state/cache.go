package state

import (
	"fmt"
	"os"
	"path"

	"github.com/rs/zerolog/log"
)

var (
	// Cache is the dir where stacks are stored locally
	Cache string
)

func init() {
	if err := initCache(); err != nil {
		log.Fatal().Err(err).Msg("failed to initialize cache home")
	}
}

func initCache() error {
	Cache = os.Getenv("HLN_CACHE_HOME")
	if Cache == "" {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return fmt.Errorf("failed to get user cache dir: %w", err)
		}
		Cache = path.Join(cacheDir, "heighliner")
	}
	err := os.MkdirAll(Cache, 0755)
	if err != nil {
		return fmt.Errorf("failed to create dir %s: %w", Cache, err)
	}
	return nil
}

// CleanCache cleans all cached cuemods and stacks
func CleanCache() {
	if err := os.RemoveAll(Cache); err != nil {
		panic(err)
	}
}
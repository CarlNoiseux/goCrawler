package context

import (
	"goCrawler/frontierExplorer"
	"goCrawler/storage/storageTypes"
)

type Context struct {
	Storage              *storageTypes.StorageInterface
	FrontierStateManager chan frontierExplorer.State
}

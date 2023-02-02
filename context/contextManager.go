package context

import (
	"goCrawler/frontierExplorer"
	"goCrawler/storage/storageInterfaces"
)

type Context struct {
	Storage              *storageInterfaces.StorageInterface
	FrontierStateManager *chan frontierExplorer.State
}

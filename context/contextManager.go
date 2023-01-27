package context

import (
	"goCrawler/storage/storageTypes"
)

type Context struct {
	Storage *storageTypes.StorageInterface
}

package types

type FolderLink struct {
	ParentFolderID uint `json:"parentFolderId"`
	ChildFolderID  uint `json:"childFolderId"`
}

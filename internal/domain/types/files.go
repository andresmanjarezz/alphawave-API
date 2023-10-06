package types

type CreateFileDTO struct {
	TeamID    string
	FileName  string
	Location  string
	Extension string
	Size      int
}

type CreateFolderDTO struct {
	TeamID     string
	FolderName string
	Location   string
}

type GetFileDTO struct {
	ID        string
	Name      string
	Size      int
	Extension string
	File      *[]byte
}

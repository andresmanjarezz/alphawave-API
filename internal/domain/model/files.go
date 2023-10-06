package model

const (
	BUCKET_STORAGE = "storage"
)

type File struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	TeamId    string `json:"teamId" bson:"teamId"`
	Name      string `json:"name" bson:"name"`
	FilePath  string `json:"filePath" bson:"filePath"`
	Key       string `json:"key" bson:"key"`
	Type      string `json:"type" bson:"type"`
	Size      int    `json:"size" bson:"size"`
	Extension string `json:"extension" bson:"extension"`
	CreatedBy string `json:"createdBy" bson:"createdBy"`
}

package upload_entity

const (
	MinUploadSize = 10 * 1024       // 10KB
	MaxUploadSize = 2 * 1024 * 1024 // 2MB
	UploadPath    = "./uploads"
)

type UploadedImage struct {
	ImageURL string `json:"imageUrl"`
}

package response

type UploadImageResponse struct {
	Name string `json:"name"` // 文件名
	Url  string `json:"url"`  // 访问图片的url
}

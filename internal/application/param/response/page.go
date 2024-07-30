package response

type PageResponse struct {
	Content     interface{} `json:"content"`
	HasPrevious bool        `json:"hasPrevious"`
	HasNext     bool        `json:"hasNext"`
	Page        int         `json:"page"`  // 当前页数，从0开始
	Pages       int         `json:"pages"` // 总页数
	Total       int64       `json:"total"` // 总消息数
}

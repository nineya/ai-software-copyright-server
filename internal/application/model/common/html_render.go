package common

type HtmlRender interface {
	// 渲染数据
	Render(templateName string, param any) (string, error)
	// 刷新数据
	Refresh()
}

package middleware

import (
	_func "ai-software-copyright-server/internal/application/router/html/func"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
)

type SiteHtmlRender struct {
	template   *template.Template
	delimLeft  string
	delimRight string
}

func (s *SiteHtmlRender) Instance(name string, data any) render.Render {
	return render.HTML{
		Template: s.template,
		Name:     name,
		Data:     data,
	}
}

func (s *SiteHtmlRender) Render(templateName string, param any) (string, error) {
	buf := new(bytes.Buffer)
	err := s.template.ExecuteTemplate(buf, templateName, param)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *SiteHtmlRender) Refresh() {
	s.template = loadHTMLGlob(s.delimLeft, s.delimRight)
}

type DevSiteHtmlRender struct {
	delimLeft  string
	delimRight string
}

func (s *DevSiteHtmlRender) Instance(name string, data any) render.Render {
	return render.HTML{
		Template: loadHTMLGlob(s.delimLeft, s.delimRight),
		Name:     name,
		Data:     data,
	}
}

func (s *DevSiteHtmlRender) Render(templateName string, param any) (string, error) {
	templ := loadHTMLGlob(s.delimLeft, s.delimRight)
	buf := new(bytes.Buffer)
	err := templ.ExecuteTemplate(buf, templateName, param)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *DevSiteHtmlRender) Refresh() {
	// 开发模式下每次访问都会从新加载模板，无须刷新
}

func HtmlRender(engine *gin.Engine) {
	HtmlRenderAndDelims(engine, "{{", "}}")
}

func HtmlRenderAndDelims(engine *gin.Engine, delimLeft, delimRight string) {
	if global.CONFIG.Server.Mode == "dev" {
		htmlRender := &DevSiteHtmlRender{
			delimLeft:  delimLeft,
			delimRight: delimRight,
		}
		global.HTML_RENDER = htmlRender
		engine.HTMLRender = htmlRender
		return
	}

	// 生成模板模板
	htmlRender := &SiteHtmlRender{
		template:   loadHTMLGlob(delimLeft, delimRight),
		delimLeft:  delimLeft,
		delimRight: delimRight,
	}
	global.HTML_RENDER = htmlRender
	engine.HTMLRender = htmlRender
}

func loadHTMLGlob(delimLeft, delimRight string) *template.Template {
	Func := _func.BaseFunc{}
	funcMap := template.FuncMap{
		"TimeFormat":   Func.TimeFormat,
		"TimeAgo":      Func.TimeAgo,
		"Pagination":   Func.Pagination,
		"Add":          Func.Add,
		"Sub":          Func.Sub,
		"Mod":          Func.Mod,
		"Html":         Func.Html,
		"Css":          Func.Css,
		"Js":           Func.Js,
		"Base64Encode": Func.Base64Encode,
		"Br":           Func.Br,
		"Default":      Func.Default,
		"Switch":       Func.Switch,
		"Params":       Func.Params,
		"NetdiskName":  Func.NetdiskName,
		"Split":        strings.Split,
		"Contains":     strings.Contains,
		"ReplaceAll":   strings.ReplaceAll,
		"HasPrefix":    strings.HasPrefix,
	}
	templ := template.New("").Delims(delimLeft, delimRight).Funcs(funcMap)
	// 加载内置模板
	loadInternalTemplate(templ, "resource/templates/internal", "internal")
	loadInternalTemplate(templ, "resource/templates/feature", "feature")
	return templ
}

func loadInternalTemplate(templ *template.Template, sourcePath, prefixName string) {
	internalFiles, err := global.FS.ReadDir(sourcePath)
	utils.PanicErr(err)
	for _, file := range internalFiles {
		filePath := sourcePath + "/" + file.Name()
		name := prefixName + "/" + file.Name()
		if file.IsDir() {
			loadInternalTemplate(templ, filePath, name)
			continue
		}
		fileContent, err := utils.ReadEmbedFileContent(filePath)
		utils.PanicErr(err)
		template.Must(templ.New(name).Parse(fileContent))
	}
}

func FilePathList(filePath string) map[string]string {
	list := make(map[string]string, 0)
	filepath.Walk(filePath, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			list[path] = strings.Replace(path[len(filePath)+1:], "\\", "/", -1)
		}
		return nil
	})
	return list
}

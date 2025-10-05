package software_copyright

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	difyPlugin "ai-software-copyright-server/internal/application/plugin/dify"
	"ai-software-copyright-server/internal/application/service"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ZeroHawkeye/wordZero/pkg/document"
	"github.com/ZeroHawkeye/wordZero/pkg/markdown"
	"github.com/ZeroHawkeye/wordZero/pkg/style"
	"github.com/chromedp/chromedp"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"
	"xorm.io/xorm"
)

type SoftwareCopyrightService struct {
	service.UserCrudService[table.SoftwareCopyright]
	ApiKey string
}

var onceSoftwareCopyright = sync.Once{}
var softwareCopyrightService *SoftwareCopyrightService

// 获取单例
func GetSoftwareCopyrightService() *SoftwareCopyrightService {
	onceSoftwareCopyright.Do(func() {
		softwareCopyrightService = new(SoftwareCopyrightService)
		softwareCopyrightService.Db = global.DB
		softwareCopyrightService.ApiKey = "app-kPGnBkdf9bSG850c5kgCS3SC"
	})
	return softwareCopyrightService
}

func (s *SoftwareCopyrightService) Create(userId int64, param table.SoftwareCopyright) (*response.UserBuyResponse, error) {
	expenseCredits := 50
	// 预检余额
	_, err := userSev.GetUserService().GetAndCheckBalance(userId, expenseCredits)
	if err != nil {
		return nil, err
	}
	param.UserId = userId
	param.Status = enum.SoftwareCopyrightStatus(1)

	result := &response.UserBuyResponse{}
	err = s.DbTransaction(func(session *xorm.Session) error {
		_, err = session.Insert(&param)
		if err != nil {
			return err
		}
		remark := fmt.Sprintf("购买软著申请创建服务，花费%d积分", expenseCredits)
		user, err := userSev.GetUserService().PaymentCreditsRunning(userId, session, enum.BuyType(1), expenseCredits, remark)
		if err != nil {
			return err
		}
		result.BuyCredits = expenseCredits
		result.BalanceCredits = user.Credits
		result.BuyMessage = remark
		return nil
	})
	if err != nil {
		return nil, err
	}
	go s.GenerateFileTask(userId, param)
	return result, err
}

// 创建文档任务
func (s *SoftwareCopyrightService) GenerateFileTask(userId int64, sc table.SoftwareCopyright) {
	var err error
	defer func() {
		sc.Progress = 100
		if err != nil {
			sc.Status = enum.SoftwareCopyrightStatus(3)
		} else if err := recover(); err != nil {
			sc.Status = enum.SoftwareCopyrightStatus(3)
		} else {
			sc.Status = enum.SoftwareCopyrightStatus(2)
		}
		_, err := s.WhereUserSession(userId).ID(sc.Id).Update(&sc)
		if err != nil {
			global.LOG.Error(fmt.Sprintf("更新软著申请状态失败：%+v", err))
		}
	}()
	// 创建目录
	storePath := utils.GetSoftwareCopyrightPath(sc.Id)
	demoPath := storePath + "/demo"
	if err = os.MkdirAll(demoPath, 0755); err != nil {
		global.LOG.Error(fmt.Sprintf("创建软著会话目录失败：%+v", err))
		return
	}

	// 基础软著信息
	param := difyPlugin.DifyChatMessageParam{
		Query: "请帮我编写内容",
		Inputs: map[string]any{
			"name":        sc.Name,
			"short_name":  sc.ShortName,
			"version":     sc.Version,
			"category":    sc.Category,
			"code_lang":   sc.CodeLang,
			"description": sc.Description,
			"owner":       sc.Owner,
			"mode":        "需求分析",
		},
		User: fmt.Sprintf("用户%d", userId),
	}

	// 分析用户需求
	result, err := difyPlugin.GetDifyPlugin().SendChat(s.ApiKey, param)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("分析用户需求失败：%+v", err))
		return
	}
	requirements := make([]response.SCRequirementItemResponse, 0)
	err = json.Unmarshal([]byte(result.Answer), &requirements)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("用户需求结果解析失败：%+v", err))
		return
	}
	sc.ConversationId = result.ConversationId
	param.ConversationId = result.ConversationId
	// 软著进度
	progressCount := 7 + (len(requirements) * 4)
	progressCurrent := 1 + len(requirements)
	sc.Progress = 100 * progressCount / progressCurrent
	_, err = s.WhereUserSession(userId).ID(sc.Id).Update(&sc)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("更新软著会话ID失败：%+v", err))
		return
	}

	// 生成源代码
	param.Inputs["mode"] = "源代码"
	codeStr, err := difyPlugin.GetDifyPlugin().SendSSEChat(s.ApiKey, param)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("生成软件源代码失败：%+v", err))
		return
	}
	codeLines := strings.Split(codeStr, "\n")
	// 创建新文档
	codeDoc := document.New()
	codeDoc.SetPageMargins(25, 25, 20, 25)
	codeDoc.SetDocGrid(document.DocGridDefault, 5, 40)
	// 添加页眉
	codeDoc.AddHeader(document.HeaderFooterTypeDefault, sc.Name+" "+sc.Version)
	// 添加带页码的页脚
	codeDoc.AddFooterWithPageNumber(document.HeaderFooterTypeDefault, "", true)
	// 添加正文段落
	codeContentStyle := &document.TextFormat{FontFamily: "宋体", FontSize: 11}
	for _, line := range codeLines {
		if line == "" || strings.HasPrefix(line, "```") || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "# ") {
			continue
		}
		codeDoc.AddFormattedParagraph(line, codeContentStyle)
	}
	// 保存文档
	if err = codeDoc.Save(storePath + "/程序鉴别材料.docx"); err != nil {
		global.LOG.Error(fmt.Sprintf("生成软件源代码失败：%+v", err))
	}
	// 更新进度
	progressCurrent += 1 + len(requirements)
	sc.Progress = 100 * progressCount / progressCurrent
	_, _ = s.WhereUserSession(userId).ID(sc.Id).Update(&sc)

	// 生成用户手册
	param.Inputs["mode"] = "用户手册"
	// 创建新文档
	bookDoc := document.New()
	// 添加页眉
	bookDoc.AddHeader(document.HeaderFooterTypeDefault, sc.Name+" "+sc.Version)
	// 添加带页码的页脚
	bookDoc.AddFooterWithPageNumber(document.HeaderFooterTypeDefault, "", true)
	bookDoc.SetDifferentFirstPage(true)
	bookDoc.GetStyleManager().AddStyle(&style.Style{
		Type:    "paragraph",
		StyleID: "Break",
		Name:    &style.StyleName{Val: "Break"},
		ParagraphPr: &style.ParagraphProperties{
			PageBreak: &style.PageBreak{},
		},
	})
	bookDoc.GetStyleManager().AddStyle(&style.Style{
		Type:    "paragraph",
		StyleID: "Normal",
		Default: true,
		Name: &style.StyleName{
			Val: "Normal",
		},
		ParagraphPr: &style.ParagraphProperties{
			Spacing: &style.Spacing{Line: "360"},
			Indentation: &style.Indentation{
				FirstLine: "480", // 左缩进240 TWIPs (12磅)
			},
		},
		RunPr: &style.RunProperties{
			FontSize: &style.FontSize{
				Val: "24", // 五号字体（10.5磅，Word中以半磅为单位）
			},
			FontFamily: &style.FontFamily{ASCII: "宋体", EastAsia: "宋体", HAnsi: "宋体", CS: "宋体"},
		},
	})
	bookDoc.GetStyleManager().AddStyle(&style.Style{
		Type:    "paragraph",
		StyleID: "Heading1",
		Name:    &style.StyleName{Val: "heading 1"},
		BasedOn: &style.BasedOn{Val: "Normal"},
		Next:    &style.Next{Val: "Normal"},
		ParagraphPr: &style.ParagraphProperties{
			KeepNext:  &style.KeepNext{},
			KeepLines: &style.KeepLines{},
			PageBreak: &style.PageBreak{},
			Spacing: &style.Spacing{
				Before: "340", // 17磅段前间距
				After:  "330", // 16.5磅段段后间距
			},
			Justification: &style.Justification{Val: "center"},
			OutlineLevel:  &style.OutlineLevel{Val: "0"},
		},
		RunPr: &style.RunProperties{
			Bold: &style.Bold{},
			FontSize: &style.FontSize{
				Val: "32", // 16磅
			},
			FontFamily: &style.FontFamily{ASCII: "宋体"},
			Color:      &style.Color{Val: "000000"},
		},
	})
	bookDoc.GetStyleManager().AddStyle(&style.Style{
		Type:    "paragraph",
		StyleID: "Heading2",
		Name:    &style.StyleName{Val: "heading 2"},
		BasedOn: &style.BasedOn{Val: "Normal"},
		Next:    &style.Next{Val: "Normal"},
		ParagraphPr: &style.ParagraphProperties{
			KeepNext:  &style.KeepNext{},
			KeepLines: &style.KeepLines{},
			Spacing: &style.Spacing{
				Before: "260", // 13磅段前间距
				After:  "260", // 13磅段段后间距
			},
			Indentation: &style.Indentation{
				FirstLine: "0", // 左缩进240 TWIPs (12磅)
			},
			OutlineLevel: &style.OutlineLevel{Val: "1"},
		},
		RunPr: &style.RunProperties{
			Bold: &style.Bold{},
			FontSize: &style.FontSize{
				Val: "28", // 14磅
			},
			FontFamily: &style.FontFamily{ASCII: "黑体", EastAsia: "黑体", HAnsi: "黑体", CS: "黑体"},
			Color:      &style.Color{Val: "000000"},
		},
	})
	bookDoc.GetStyleManager().AddStyle(&style.Style{
		Type:    "paragraph",
		StyleID: "Heading3",
		Name:    &style.StyleName{Val: "heading 3"},
		BasedOn: &style.BasedOn{Val: "Normal"},
		Next:    &style.Next{Val: "Normal"},
		ParagraphPr: &style.ParagraphProperties{
			KeepNext:  &style.KeepNext{},
			KeepLines: &style.KeepLines{},
			Spacing: &style.Spacing{
				Before: "220", // 13磅段前间距
				After:  "220", // 13磅段段后间距
			},
			Indentation: &style.Indentation{
				FirstLine: "0", // 左缩进240 TWIPs (12磅)
			},
			OutlineLevel: &style.OutlineLevel{Val: "2"},
		},
		RunPr: &style.RunProperties{
			Bold: &style.Bold{},
			FontSize: &style.FontSize{
				Val: "24", // 14磅
			},
			FontFamily: &style.FontFamily{ASCII: "黑体", EastAsia: "黑体", HAnsi: "黑体", CS: "黑体"},
			Color:      &style.Color{Val: "000000"},
		},
	})
	// 封面页
	coverStyle := &document.TextFormat{Bold: true, FontFamily: "宋体", FontSize: 26}
	defaultSpacing := &document.SpacingConfig{LineSpacing: 1.5}
	bookDoc.AddFormattedParagraph("\n", coverStyle)
	paragraph := bookDoc.AddFormattedParagraph(sc.Name, coverStyle)
	paragraph.SetAlignment(document.AlignCenter)
	paragraph.SetSpacing(defaultSpacing)
	paragraph = bookDoc.AddFormattedParagraph(sc.Version, coverStyle)
	paragraph.SetAlignment(document.AlignCenter)
	paragraph.SetSpacing(defaultSpacing)
	bookDoc.AddFormattedParagraph("\n", coverStyle)
	paragraph = bookDoc.AddFormattedParagraph("用", coverStyle)
	paragraph.SetAlignment(document.AlignCenter)
	paragraph.SetSpacing(defaultSpacing)
	paragraph = bookDoc.AddFormattedParagraph("户", coverStyle)
	paragraph.SetAlignment(document.AlignCenter)
	paragraph.SetSpacing(defaultSpacing)
	paragraph = bookDoc.AddFormattedParagraph("手", coverStyle)
	paragraph.SetAlignment(document.AlignCenter)
	paragraph.SetSpacing(defaultSpacing)
	paragraph = bookDoc.AddFormattedParagraph("册", coverStyle)
	paragraph.SetAlignment(document.AlignCenter)
	paragraph.SetSpacing(defaultSpacing)
	bookDoc.AddFormattedParagraph("\n", coverStyle)
	bookDoc.AddFormattedParagraph("\n", coverStyle)
	bookDoc.AddFormattedParagraph("\n", coverStyle)
	bookDoc.AddFormattedParagraph(time.Now().Format("2006年01月02日"), coverStyle).SetAlignment(document.AlignCenter)
	// 目录页
	paragraph = bookDoc.AddFormattedParagraph("目录", &document.TextFormat{Bold: true, FontFamily: "宋体", FontSize: 16})
	paragraph.SetStyle("Break")
	paragraph.SetAlignment(document.AlignCenter)
	paragraph.SetSpacing(defaultSpacing)
	config := &document.TOCConfig{
		MaxLevel:    3,
		ShowPageNum: true,
	}
	bookDoc.GenerateTOC(config)
	converter := markdown.NewConverter(markdown.DefaultOptions())
	// 更新进度
	progressCurrent += 1
	sc.Progress = 100 * progressCount / progressCurrent
	_, _ = s.WhereUserSession(userId).ID(sc.Id).Update(&sc)
	// 添加引言
	bookDoc.AddHeadingParagraph("第一章 引言", 1)
	param.Query = `请结合软件的功能和信息，帮我完成第一章引言的编写，要求包含编写目的、背景、目标用户等内容。
编写目的：介绍编写手册的目的，手册的作用
背景：描述开发这款软件的原因，有哪些市场和商业背景
目标用户：这款软件主要用户群体有哪些，分群体回答

内容要丰富、有深度，可以分多段回答。

严格按以下格式回复，不用包含大标题，不要有其他任何的解释说明：
## 1.1 编写目的
...
## 1.2 背景
...
## 1.3 目标用户
...`
	bookStr, err := difyPlugin.GetDifyPlugin().SendSSEChat(s.ApiKey, param)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("生成用户手册的引言失败：%+v", err))
		return
	}
	handleMarkdownToWord(bookStr, converter, bookDoc)
	// 更新进度
	progressCurrent += 1
	sc.Progress = 100 * progressCount / progressCurrent
	_, _ = s.WhereUserSession(userId).ID(sc.Id).Update(&sc)
	bookDoc.AddHeadingParagraph("第二章 软件概述", 1)
	param.Query = `请结合软件的功能和信息，帮我完成软件概述章节的编写。容要丰富、有深度，可以分多段回答。
直接回复章节内容，不用包含大标题，不要有其他任何的解释说明。`
	bookStr, err = difyPlugin.GetDifyPlugin().SendSSEChat(s.ApiKey, param)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("生成用户手册的软件概述失败：%+v", err))
		return
	}
	handleMarkdownToWord(bookStr, converter, bookDoc)
	// 更新进度
	progressCurrent += 1
	sc.Progress = 100 * progressCount / progressCurrent
	_, _ = s.WhereUserSession(userId).ID(sc.Id).Update(&sc)
	bookDoc.AddHeadingParagraph("第三章 软件运行的软硬件环境", 1)
	param.Query = `请结合软件的功能和信息，帮我完成第三章软件运行的软硬件环境的编写，要求包含运行硬件环境、软件环境等内容，请从各个角度给出软硬件具体的版本、参数要求。

严格按以下格式回复，不用包含大标题，不要有其他任何的解释说明：
## 3.1 运行硬件环境
...
## 3.2 运行软件环境
...`
	bookStr, err = difyPlugin.GetDifyPlugin().SendSSEChat(s.ApiKey, param)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("生成用户手册的软件概述失败：%+v", err))
		return
	}
	handleMarkdownToWord(bookStr, converter, bookDoc)
	// 更新进度
	progressCurrent += 1
	sc.Progress = 100 * progressCount / progressCurrent
	_, _ = s.WhereUserSession(userId).ID(sc.Id).Update(&sc)
	bookDoc.AddHeadingParagraph("第四章 主要功能与特点", 1)
	htmlContent := ""
	var width int64 = 430
	var height int64 = 930
	if sc.Category == "WEB应用" || sc.Category == "桌面软件" {
		width = 1920
		height = 1080
	}

	// 创建chrome上下文
	chromeCtx, chromeCancel := chromedp.NewContext(context.Background())
	defer chromeCancel()
	for i, item := range requirements {
		bookDoc.AddHeadingParagraph(fmt.Sprintf("4.%d %s", i+1, item.Name), 2)
		// 生成demo代码
		param.Query = fmt.Sprintf(`
请帮我完成%s功能的html前端代码编写。

## 功能介绍
%s

## 操作流程
%s

## UI风格参考
%s
`, item.Name, item.Desc, item.Operation, htmlContent)
		param.Inputs["mode"] = "demo"
		htmlResult, err := difyPlugin.GetDifyPlugin().SendSSEChat(s.ApiKey, param)
		if err != nil {
			global.LOG.Error(fmt.Sprintf("生成用户手册%s的demo失败：%+v", item.Name, err))
		} else {
			htmlContent = htmlResult
			// 将字符串写入文件，如果文件不存在会创建，存在则覆盖
			err = os.WriteFile(demoPath+"/"+item.Name+".html", []byte(htmlContent), 0644)
			if err != nil {
				log.Fatal("写入文件失败:", err)
			}
			var imageBytes []byte
			err = chromedp.Run(chromeCtx,
				// 设置视口
				chromedp.EmulateViewport(width, height),
				// 设置内容
				chromedp.Navigate("data:text/html;charset=utf-8;base64,"+
					base64.StdEncoding.EncodeToString([]byte(htmlContent))),
				// 等待页面加载完成
				chromedp.WaitReady("body"),
				// 等待 JavaScript 执行
				chromedp.Sleep(2*time.Second),
				// 截图
				chromedp.FullScreenshot(&imageBytes, 100),
			)
			if err != nil {
				global.LOG.Error(fmt.Sprintf("生成用户手册%s的demo运行截图失败：%+v", item.Name, err))
			} else {
				addImageIntoWord(imageBytes, bookDoc)
			}
		}
		// 更新进度
		progressCurrent += 1
		sc.Progress = 100 * progressCount / progressCurrent
		_, _ = s.WhereUserSession(userId).ID(sc.Id).Update(&sc)
		bookDoc.AddHeadingParagraph(fmt.Sprintf("4.%d.1 功能介绍", i+1), 3)
		handleMarkdownToWord(item.Desc, converter, bookDoc)
		bookDoc.AddHeadingParagraph(fmt.Sprintf("4.%d.2 操作说明", i+1), 3)
		// 生成操作流程图
		param.Query = item.Operation
		param.Inputs["mode"] = "流程图"
		base64Text, err := difyPlugin.GetDifyPlugin().SendSSEChat(s.ApiKey, param)
		if err != nil {
			global.LOG.Error(fmt.Sprintf("生成用户手册%s的流程图失败：%+v", item.Name, err))
		} else {
			imageBytes, err := base64.StdEncoding.DecodeString(base64Text)
			if err != nil {
				global.LOG.Error(fmt.Sprintf("图片Base64解析失败：%+v", err))
			} else {
				addImageIntoWord(imageBytes, bookDoc)
			}
		}
		// 更新进度
		progressCurrent += 1
		sc.Progress = 100 * progressCount / progressCurrent
		_, _ = s.WhereUserSession(userId).ID(sc.Id).Update(&sc)
		//bookDoc.AddFormattedParagraph("操作说明：", &document.TextFormat{Bold: true})
		handleMarkdownToWord(item.Operation, converter, bookDoc)
	}
	bookDoc.UpdateTOC(config)

	// 保存文档
	if err := bookDoc.Save(storePath + "/文档鉴别材料.docx"); err != nil {
		log.Fatal(err)
	}
}

// 后台分页查询列表
func (s *SoftwareCopyrightService) GetByPage(userId int64, param request.QueryPageParam) (*response.PageResponse, error) {
	// 同一个时间可能包含多条数据，必须加上id做分页
	session := s.WhereUserSession(userId).Desc("create_time").Asc("id")
	if param.Keyword != "" {
		session.And("name like concat('%',?,'%')", param.Keyword)
	}
	list := make([]table.SoftwareCopyright, 0)
	return s.HandlePageable(param.PageableParam, &list, session)
}

// 将md文档转为word内容
func handleMarkdownToWord(bookStr string, converter *markdown.Converter, bookDoc *document.Document) {
	dd, _ := converter.ConvertString(bookStr, nil)
	for _, element := range dd.Body.Elements {
		if p, ok := element.(*document.Paragraph); ok {
			if p.Properties != nil && p.Properties.ParagraphStyle != nil && p.Properties.ParagraphStyle.Val == style.StyleHeading2 {
				bookDoc.AddHeadingParagraph(p.Runs[len(p.Runs)-1].Text.Content, 2)
				continue
			}
		}
		bookDoc.Body.Elements = append(bookDoc.Body.Elements, element)
	}
}

func addImageIntoWord(imageBytes []byte, bookDoc *document.Document) {
	width, height, err := utils.ImagePngWithDecode(imageBytes)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("获取图片宽高失败：%+v", err))
		return
	}
	// 假设 imageData 是图片的字节数据
	_, err = bookDoc.AddImageFromData(
		imageBytes,
		"image.png",             // 文件名
		document.ImageFormatPNG, // 图片格式
		width, height,           // 原始宽度和高度（像素）
		&document.ImageConfig{
			Size: &document.ImageSize{
				Width:           130.0 * math.Min(float64(width)/float64(height), 1), // 显示宽度（毫米）
				KeepAspectRatio: true,
			},
			Position:  document.ImagePositionInline,
			Alignment: document.AlignCenter,
		},
	)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("插入图片到文档失败：%+v", err))
	}
}

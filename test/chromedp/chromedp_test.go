package test

import (
	"context"
	"encoding/base64"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"io/ioutil"
	"log"
	"math"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

func TestContainer(t *testing.T) {
	t.Run("HTMLToScreenshot", HTMLToScreenshot)
	//t.Run("BookFile", BookFile)
	//t.Run("CodeFile", CodeFile)
	//t.Run("CodeTemplate", CodeTemplate)
}

func HTMLToScreenshot(t *testing.T) {
	htmlContent := `<!DOCTYPE html>
	<html>
	<head>
	   <title>测试页面</title>
	   <style>
	       body { font-family: Arial, sans-serif; padding: 20px; }
	       .container { max-width: 600px; margin: 0 auto; }
	       .header { background: #f0f0f0; padding: 20px; text-align: center; }
	       .content { padding: 20px; border: 1px solid #ddd; margin-top: 20px; }
	   </style>
	</head>
	<body>
	   <div class="container">
	       <div class="header">
	           <h1>Hello, World!</h1>
	           <p>这是一个测试页面</p>
	       </div>
	       <div class="content">
	           <p>当前时间: <span id="time"></span></p>
	           <ul>
	               <li>项目 1</li>
	               <li>项目 2</li>
	               <li>项目 3</li>
	           </ul>
	       </div>
	   </div>
	   <script>
	       document.getElementById('time').textContent = new Date().toLocaleString();
	   </script>
	</body>
	</html>`
	// 创建上下文
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	err := chromedp.Run(ctx,
		// 设置视口
		chromedp.EmulateViewport(1920, 1080),
		// 设置内容
		chromedp.Navigate("data:text/html;charset=utf-8;base64,"+
			base64.StdEncoding.EncodeToString([]byte(htmlContent))),
		// 等待页面加载完成
		chromedp.WaitReady("body"),
		// 等待 JavaScript 执行
		chromedp.Sleep(2*time.Second),
		// 截图
		chromedp.FullScreenshot(&buf, 100),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 保存截图到文件
	err = ioutil.WriteFile("screenshot.png", buf, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("截图已保存为 screenshot.png")
}

func HTMLToScreenshot2(t *testing.T) {
	htmlContent := `<!DOCTYPE html>
	<html>
	<head>
	   <title>测试页面</title>
	   <style>
	       body { font-family: Arial, sans-serif; padding: 20px; }
	       .container { max-width: 600px; margin: 0 auto; }
	       .header { background: #f0f0f0; padding: 20px; text-align: center; }
	       .content { padding: 20px; border: 1px solid #ddd; margin-top: 20px; }
	   </style>
	</head>
	<body>
	   <div class="container">
	       <div class="header">
	           <h1>Hello, World!</h1>
	           <p>这是一个测试页面</p>
	       </div>
	       <div class="content">
	           <p>当前时间: <span id="time"></span></p>
	           <ul>
	               <li>项目 1</li>
	               <li>项目 2</li>
	               <li>项目 3</li>
	           </ul>
	       </div>
	   </div>
	   <script>
	       document.getElementById('time').textContent = new Date().toLocaleString();
	   </script>
	</body>
	</html>`
	outputPath := "screenshot.png"

	// 创建上下文，取消headless模式进行调试
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoDefaultBrowserCheck, //不检查默认浏览器
		chromedp.Flag("headless", true),
		chromedp.Flag("blink-settings", "imagesEnabled=true"), //开启图像界面,重点是开启这个
		chromedp.Flag("ignore-certificate-errors", true),      //忽略错误
		chromedp.Flag("disable-web-security", true),           //禁用网络安全标志
		chromedp.Flag("disable-extensions", true),             //开启插件支持
		chromedp.Flag("disable-default-apps", true),
		chromedp.WindowSize(1920, 1080),    // 设置浏览器分辨率（窗口大小）
		chromedp.Flag("disable-gpu", true), //开启gpu渲染
		chromedp.Flag("hide-scrollbars", true),
		chromedp.Flag("mute-audio", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.NoFirstRun, //设置网站不是首次运行
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36"), //设置UserAgent
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var buf []byte
	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(htmlContent, 100, &buf)); err != nil {
		log.Fatal(err)
	}

	//// 执行任务
	//err := chromedp.Run(ctx,
	//	// 设置视口大小
	//	chromedp.EmulateViewport(1920, 1080),
	//
	//	// 导航到HTML内容
	//	chromedp.Navigate("data:text/html;charset=utf-8,"+htmlContent),
	//
	//	// 等待页面完全加载
	//	chromedp.WaitReady("body", chromedp.ByQuery),
	//
	//	// 等待JavaScript执行完成
	//	chromedp.Sleep(3*time.Second),
	//
	//	// 检查页面内容
	//	chromedp.ActionFunc(func(ctx context.Context) error {
	//		var bodyHTML string
	//		if err := chromedp.Evaluate(`document.body.innerHTML`, &bodyHTML).Do(ctx); err != nil {
	//			return err
	//		}
	//		log.Printf("Body内容长度: %d", len(bodyHTML))
	//
	//		var hasVisibleContent bool
	//		if err := chromedp.Evaluate(`
	//            document.body.innerText.length > 0 ||
	//            document.querySelector('svg') !== null ||
	//            document.querySelector('canvas') !== null ||
	//            document.querySelector('img') !== null
	//        `, &hasVisibleContent).Do(ctx); err != nil {
	//			return err
	//		}
	//
	//		if !hasVisibleContent {
	//			return fmt.Errorf("页面没有可见内容")
	//		}
	//
	//		log.Println("页面有可见内容，继续截图...")
	//		return nil
	//	}),
	//
	//	chromedp.FullScreenshot(&buf, 90),
	//)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//if len(buf) == 0 {
	//	log.Fatal("截图数据为空")
	//}
	//
	// 保存截图
	if err := ioutil.WriteFile(outputPath, buf, 0644); err != nil {
		log.Fatal(err)
	}
}

// 获取整个浏览器窗口的截图（全屏）
// 这将模拟浏览器操作设置。
func fullScreenshot(htmlContent string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("data:text/html;charset=utf-8;base64," +
			base64.StdEncoding.EncodeToString([]byte(htmlContent))),
		//chromedp.WaitVisible("style"),
		chromedp.Sleep(10 * time.Second),
		//chromedp.OuterHTML(`document.querySelector("body")`, &htmlContent, chromedp.ByJSPath),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// 得到布局页面
			_, _, _, _, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// 浏览器视窗设置模拟
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// 捕捉屏幕截图
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}

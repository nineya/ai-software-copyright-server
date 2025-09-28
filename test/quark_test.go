package test

import (
	netdSev "ai-software-copyright-server/internal/application/plugin/quark"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/initialize"
	"fmt"
	"testing"
)

func TestContainer(t *testing.T) {
	initialize.InitSystemConfig()
	initialize.InitLogger()
	//DumpAndShare(nil)
	t.Run("DumpAndShare", DumpAndShare)
}

func DumpAndShare(t *testing.T) {
	shareLink, err := netdSev.GetQuarkService().DumpAndShare("https://pan.quark.cn/s/bbbb03941696")
	global.LOG.Info(shareLink)
	if err != nil {
		global.LOG.Sugar().Error("%v", err)
		return
	}
}

func GetStokenApi(t *testing.T) {
	stoken, err := netdSev.GetQuarkService().GetStokenApi("0f94e8c83cb0", "")
	fmt.Println(err)
	fmt.Println(stoken)
}

func DetailApi(t *testing.T) {
	stoken, err := netdSev.GetQuarkService().GetStokenApi("0f94e8c83cb0", "")
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	detail, next, err := netdSev.GetQuarkService().DetailApi("0f94e8c83cb0", stoken, 1)
	fmt.Printf("%v", detail)
	fmt.Println(next)
}

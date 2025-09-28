package _func

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/utils"
)

func (f BaseFunc) NetdiskName(typ enum.NetdiskType) string {
	return utils.TransformNetdiskName(typ)
}

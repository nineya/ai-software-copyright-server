package initialize

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/global"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"xorm.io/xorm"
	"xorm.io/xorm/caches"
)

func InitDatabase() {
	datasource := global.CONFIG.Datasource
	global.LOG.Info(fmt.Sprintf("The database type is %s", datasource.Type))
	db, err := xorm.NewEngine(datasource.Type, datasource.Url)
	if err != nil {
		panic(errors.Wrap(err, "Database startup exception"))
	}
	if err = db.Ping(); err != nil {
		panic(errors.Wrap(err, "Database connection failure"))
	}
	db.SetMaxIdleConns(datasource.MaxIdleConns)
	db.SetMaxOpenConns(datasource.MaxOpenConns)
	if datasource.UseCache && global.CONFIG.Server.Mode != "dev" {
		global.LOG.Info(fmt.Sprintf("Use the xorm cache. The cache size is %d", datasource.CacheSize))
		db.SetDefaultCacher(caches.NewLRUCacher(caches.NewMemoryStore(), datasource.CacheSize))
	}
	// 设置日志级别
	//db.Logger().SetLevel(log.LOG_DEBUG)
	db.ShowSQL(datasource.ShowSql)
	global.DB = db
	syncTable(db)
	loadData()
}

func syncTable(db *xorm.Engine) {
	_, err := db.SyncWithOptions(xorm.SyncOptions{IgnoreDropIndices: true}, new(table.NetdiskResource))
	if err != nil {
		global.LOG.Error("Failed to synchronize the database table structure.", zap.Any("err", err))
		return
	}
	err = db.Sync(
		//new(table.Admin),
		//new(table.Buy),
		//new(table.Cdkey),
		//new(table.CdkeyRecord),
		//new(table.ClientInfo),
		new(table.CreditsChange),
		new(table.CreditsOrder),
		new(table.CreditsPrice),
		//new(table.InviteRecord),
		//new(table.NetdiskHelperConfigure),
		//new(table.NetdiskResourceSearch),
		//new(table.NetdiskSearchAppConfigure),
		//new(table.NetdiskSearchSiteConfigure),
		//new(table.NetdiskSearchWxampConfigure),
		//new(table.NetdiskShortLinkConfigure),
		//new(table.AdminLog),
		//new(table.Notice),
		//new(table.Qrcode),
		//new(table.RedbookCookie),
		//new(table.RedbookProhibited),
		//new(table.RedbookVisitsTask),
		//new(table.ShortLink),
		new(table.Statistic),
		//new(table.StudyResource),
		//new(table.TimeClock),
		//new(table.TimeClockMember),
		//new(table.TimeClockRecord),
		new(table.User),
		new(table.UserLog),
	)
	if err != nil {
		global.LOG.Error("Failed to synchronize the database table structure.", zap.Any("err", err))
	} else {
		global.LOG.Info("Synchronizing the database table structure is complete.")
	}
}

func loadData() {
}

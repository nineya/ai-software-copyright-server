package netdisk

import (
	"ai-software-copyright-server/internal/application/model/common"
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/service"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"mime/multipart"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
	"xorm.io/xorm"
)

type ResourceService struct {
	service.UserCrudService[table.NetdiskResource]
	awaitCheckQuark          []table.NetdiskResource
	awaitCheckQuarkMutex     sync.Mutex
	awaitCheckQuarkEmptyTime time.Time
	awaitCheckBaidu          []table.NetdiskResource
	awaitCheckBaiduMutex     sync.Mutex
	awaitCheckBaiduEmptyTime time.Time
	importStatus             sync.Map
}

var onceResource = sync.Once{}
var resourceService *ResourceService

// 获取单例
func GetResourceService() *ResourceService {
	onceResource.Do(func() {
		resourceService = new(ResourceService)
		resourceService.Db = global.DB
	})
	return resourceService
}

func (s *ResourceService) Create(userId int64, param table.NetdiskResource) (*response.NetdiskResourceCreateResponse, error) {
	// 免费额度
	freeCredit := 10000

	if param.Name == "" {
		return nil, errors.New("缺少资源名称")
	}
	mod := table.NetdiskResource{TargetUrl: strings.TrimSpace(param.TargetUrl)}
	exist, err := s.WhereUserSession(userId).Get(&mod)
	if err != nil {
		return nil, err
	}
	result := &response.NetdiskResourceCreateResponse{}
	// 更新资源名字
	if exist {
		result.BuyMessage = "资源已存在，更新资源名称"
		mod.Name = strings.TrimSpace(param.Name)
		_, err = s.Db.ID(mod.Id).Update(&mod)
		result.IsUpdate = true
		return result, err
	}
	err = s.DbTransaction(func(session *xorm.Session) error {
		// 新增
		mod.UserId = userId
		mod.UserName = param.UserName
		mod.Name = strings.TrimSpace(param.Name)
		mod.Type = utils.TransformNetdiskType(mod.TargetUrl)
		if param.Origin != 0 {
			mod.Origin = param.Origin
		} else {
			mod.Origin = enum.NetdiskOrigin(1)
		}
		mod.Status = enum.NetdiskStatus(1)
		_, err = session.Insert(&mod)
		if err != nil {
			return err
		}
		count, err := s.AddWhereUser(userId, session).Count(table.NetdiskResource{})
		if err != nil {
			return err
		}
		// 小于10000是免费的
		if int(count) <= freeCredit {
			result.BuyMessage = "新增资源成功"
			return nil
		}

		// 扣款
		user, err := userSev.GetUserService().PaymentCredits(userId, enum.BuyType(14), 1, "新增1条网盘资源，花费1币")
		if err != nil {
			return err
		}
		result.BuyCredits = 1
		result.BalanceCredits = user.Credits
		result.BuyMessage = "新增资源成功，花费1积分"
		return nil
	})
	return result, err
}

func (s *ResourceService) Import(userId int64, file *multipart.FileHeader) (*response.UserBuyResponse, error) {
	_, exist := s.importStatus.LoadOrStore(userId, true)
	if exist {
		return nil, errors.New("已有导入任务，请等待已有任务完成")
	}
	defer s.importStatus.Delete(userId)

	// 免费额度
	freeCredit := 10000
	result := &response.UserBuyResponse{}

	// * 处理xlxs，预检并获取资源列表
	rows, err := handleImportFile(file)
	if err != nil {
		return nil, err
	}

	// ** 获取用户余额，取得最大可插入资源数量
	user, err := userSev.GetUserService().GetById(userId)
	if err != nil {
		return nil, errors.Wrap(err, "获取用户信息失败")
	}
	count, err := s.WhereUserSession(userId).Count(table.NetdiskResource{})
	if err != nil {
		return nil, errors.Wrap(err, "获取资源信息失败")
	}
	maxInsertCount := user.Credits
	if freeCredit > int(count) {
		maxInsertCount += freeCredit - int(count)
	}

	// 新增资源列表
	insertResources := make([]table.NetdiskResource, 0, len(rows))
	updateCount := 0
	// 因为积分不够，减少资源数量
	cutCount := 0
	// 初始批次
	batch := 0
	// 数据处理批大小
	batchSize := 2000

	for true {
		start := batch * batchSize
		end := start + batchSize
		if end > len(rows) {
			end = len(rows)
		}
		// 查询用的目标地址数组
		selectTargetUrls := make([]string, 0, batchSize)
		// 查询资源
		for i := start; i < end; i++ {
			selectTargetUrls = append(selectTargetUrls, rows[i][1])
		}
		existTargetUrls := make([]string, 0)
		err = s.WhereUserSession(userId).Table("netdisk_resource").Select("target_url").In("target_url", selectTargetUrls).Find(&existTargetUrls)
		if err != nil {
			return nil, err
		}
		// 更新资源
		for i := start; i < end; i++ {
			exist := false
			for j := 0; j < len(existTargetUrls); j++ {
				if existTargetUrls[j] == rows[i][1] {
					exist = true
					break
				}
			}
			// 数据存在，更新
			if exist {
				_, err = s.WhereUserSession(userId).And("target_url = ?", rows[i][1]).Update(table.NetdiskResource{Name: rows[i][0]})
				if err != nil {
					return nil, err
				}
				updateCount++
			} else {
				// 资源不够，退出循环
				if maxInsertCount <= 0 {
					cutCount = len(rows) - i
					break
				}
				maxInsertCount--
				insertResources = append(insertResources, table.NetdiskResource{
					UserId:    userId,
					Name:      rows[i][0],
					TargetUrl: rows[i][1],
					Type:      utils.TransformNetdiskType(rows[i][1]),
					Origin:    enum.NetdiskOrigin(1),
					Status:    enum.NetdiskStatus(1),
				})

			}
		}
		if end == len(rows) {
			break
		}
		batch++
	}
	// 资源数量
	insertCount := len(insertResources)

	if cutCount > 0 {
		result.BuyMessage = fmt.Sprintf("新增%d条，更新%d条，有%d条因为积分不足未处理", insertCount, updateCount, cutCount)
	} else {
		result.BuyMessage = fmt.Sprintf("新增%d条，更新%d条资源", insertCount, updateCount)
	}
	// 不需要新增资源，直接退出
	if insertCount == 0 {
		return result, nil
	}

	err = s.DbTransaction(func(session *xorm.Session) error {
		// 重新取得当前资源总数量
		count, err = s.WhereUserSession(userId).Count(table.NetdiskResource{})
		// 计算费用
		expenseCredits := insertCount
		// 如果还剩余免费额度，则扣除免费额度
		if freeCredit > int(count) {
			expenseCredits -= freeCredit - int(count)
		}
		// 需要付费，执行付费逻辑
		if expenseCredits > 0 {
			user, err := userSev.GetUserService().GetById(userId)
			if err != nil {
				return err
			}
			// 余额不够，裁剪数组扣除金额对应的资源数量
			if user.Credits < expenseCredits {
				// 再减一次资源数量
				cutCount += expenseCredits - user.Credits
				// 新的新增资源数量
				insertCount -= expenseCredits - user.Credits
				insertResources = insertResources[:insertCount]
				expenseCredits = user.Credits
			}
			if cutCount > 0 {
				result.BuyMessage = fmt.Sprintf("新增%d条，更新%d条，花费%d币，有%d条因为积分不足未处理", insertCount, updateCount, expenseCredits, cutCount)
			} else {
				result.BuyMessage = fmt.Sprintf("新增%d条，更新%d条，花费%d币", insertCount, updateCount, expenseCredits)
			}
			// 扣款
			user, err = userSev.GetUserService().PaymentCreditsRunning(userId, session, enum.BuyType(14), expenseCredits, fmt.Sprintf("新增%d条网盘资源，花费%d币", insertCount, expenseCredits))
			if err != nil {
				return err
			}
			result.BuyCredits = expenseCredits
			result.BalanceCredits = user.Credits
		}
		batch = 0
		batchSize = 2000
		// 新增资源
		for true {
			start := batch * batchSize
			end := start + batchSize
			if end > len(insertResources) {
				end = len(insertResources)
			}
			inserts := insertResources[start:end]
			_, err = session.InsertMulti(inserts)
			if err != nil || end == len(insertResources) {
				return err
			}
			batch++
		}
		return nil
	})

	return result, err
}

// 如果目标链接为空，或等于外部分享链接，则判断分享链接是否保存，已保存就更新名字，未保存就转存并保存资源
// 只有在资源没有保存时才设置资源来源，不主动更新资源来源，避免已经修改的资源来源被更新为搜索采集来源
func (s *ResourceService) Save(userId int64, param request.NetdiskResourceSaveParam) (*table.NetdiskResource, error) {
	// 免费额度
	freeCredit := 10000

	if param.TargetUrl == "" && param.ShareTargetUrl == "" {
		return nil, errors.New("缺少资源地址")
	}
	// 更新基础信息
	param.UserId = userId
	param.Name = strings.TrimSpace(param.Name)
	param.TargetUrl = strings.TrimSpace(param.TargetUrl)
	param.ShareTargetUrl = strings.TrimSpace(param.ShareTargetUrl)
	if param.ShareTargetUrl != "" {
		param.Type = utils.TransformNetdiskType(param.ShareTargetUrl)
	} else {
		param.Type = utils.TransformNetdiskType(param.TargetUrl)
	}
	param.Status = enum.NetdiskStatus(1)
	// 目标链接为空，或等于外部分享链接，需要转存
	if param.TargetUrl == "" || param.TargetUrl == param.ShareTargetUrl {
		// 存在客户端
		if global.SOCKET.GetClient(userId) != nil {
			// 通过分享链接判断是否已经转存过了
			mod := table.NetdiskResource{ShareTargetUrl: param.ShareTargetUrl}
			exist, err := s.WhereUserSession(param.UserId).Get(&mod)
			if err != nil {
				return nil, err
			}
			if exist {
				// 更新链接名字
				if param.Name != "" {
					mod.Name = param.Name
				}
				//mod.Origin = param.Origin
				_, err = s.Db.ID(mod.Id).Update(&mod)
				return &mod, err
			}
			serializedData, _ := json.Marshal(param)
			result := ""
			message := common.SocketMessage{
				NeedResult: true,
				Timeout:    60 * time.Second,
				Type:       enum.SocketMessageType(7),
				Data:       string(serializedData),
			}
			existResult, cliErr := global.SOCKET.SendMessage(userId, message, &result)
			// 查询错误或者没有等到结果
			if cliErr != nil || !existResult || result == "" {
				global.LOG.Warn(fmt.Sprintf("客户端查询未能获取到结果: %+v", cliErr))
				return nil, errors.New("该资源已失效！")
			}
			param.TargetUrl = result
		}
		// 到这里已经成功转存了，如果资源来源为空，更新来源为搜索转存
		if param.Origin == 0 {
			param.Origin = enum.NetdiskOrigin(2)
		}
	} else if param.TargetUrl != "" {
		//目标链接不为空，判断目标链接是否已经存储，是的话就更新链接名字
		mod := table.NetdiskResource{TargetUrl: param.TargetUrl}
		exist, err := s.WhereUserSession(userId).Get(&mod)
		if err != nil {
			return nil, err
		}
		if exist {
			// 更新链接名字
			if param.Name != "" {
				mod.Name = param.Name
			}
			//mod.Origin = param.Origin
			_, err = s.Db.ID(mod.Id).Update(&mod)
			return &mod, err
		}
	}

	count, err := s.WhereUserSession(param.UserId).Count(table.NetdiskResource{})
	if err != nil {
		return &param.NetdiskResource, nil
	}
	// 大于10000要收费
	if int(count) > freeCredit {
		// 扣款
		_, err = userSev.GetUserService().PaymentCredits(param.UserId, enum.BuyType(14), 1, "新增1条网盘资源，花费1币")
		if err != nil {
			return &param.NetdiskResource, nil
		}
	}
	// 到这里表示资源存过了，但是在数据库没有，要新增，将资源来源更新为用户上传
	if param.Origin == 0 {
		param.Origin = enum.NetdiskOrigin(1)
	}
	// 存储到数据库，不记录结果
	_, _ = s.Db.Insert(&param.NetdiskResource)
	return &param.NetdiskResource, nil
}

func (s *ResourceService) Clear(userId int64, origin enum.NetdiskOrigin, status enum.NetdiskStatus) (string, error) {
	count, err := s.WhereUserSession(userId).Delete(&table.NetdiskResource{Origin: origin, Status: status})
	return fmt.Sprintf("成功清除%d条资源", count), err
}

func (s *ResourceService) UpdateStatusInBatch(userId int64, ids []int64, status enum.NetdiskStatus) error {
	_, err := s.WhereUserSession(userId).In("id", ids).Update(table.NetdiskResource{Status: status})
	return err
}

func (s *ResourceService) UpdateCheckResult(param table.NetdiskResource) error {
	_, err := s.Db.ID(param.Id).NoAutoTime().Update(table.NetdiskResource{
		UserName:  param.UserName,
		Name:      param.Name,
		Status:    param.Status,
		CheckTime: param.CheckTime,
	})
	return err
}

// 取得用户指定夸克资源的数量
func (s *ResourceService) GetQuarkResourceByUserName(userId int64, userName string, size int) ([]table.NetdiskResource, error) {
	list := make([]table.NetdiskResource, 0)
	err := s.Db.Where("user_id = ? and user_name = ?", userId, userName).Limit(size).Find(&list)
	return list, err
}

// 取得指定数量待检测的夸克网盘资源
func (s *ResourceService) GetCheckQuarkResource(size int) []table.NetdiskResource {
	s.awaitCheckQuarkMutex.Lock()
	defer s.awaitCheckQuarkMutex.Unlock()
	length := len(s.awaitCheckQuark)
	batch := 300
	// 如果总长度不超过单次长度，且距离上次取得数据不全超过10分钟，追加获取数据
	if length < size && time.Now().Sub(s.awaitCheckQuarkEmptyTime).Minutes() > 10 {
		list := make([]table.NetdiskResource, 0)
		err := s.Db.SQL(`
select n.* from netdisk_resource n left join short_link sl on n.short_link = sl.alias
where type = 2 and status in (1,2) and TO_DAYS(NOW())-TO_DAYS(n.update_time) < 300 and (sl.update_time is null or TO_DAYS(NOW())-TO_DAYS(sl.update_time) < 600) and (check_time is null or TO_DAYS(NOW())-TO_DAYS(check_time) > 30 or
(TO_DAYS(NOW())-TO_DAYS(n.update_time) < 15 and TO_DAYS(NOW())-TO_DAYS(check_time) > 5) or
(TO_DAYS(NOW())-TO_DAYS(n.update_time) < 50 and TO_DAYS(NOW())-TO_DAYS(check_time) > 10) or
(TO_DAYS(NOW())-TO_DAYS(n.update_time) < 80 and TO_DAYS(NOW())-TO_DAYS(check_time) > 20) or
((sl.visits > 3000 or TO_DAYS(NOW())-TO_DAYS(sl.update_time) < 5) and TO_DAYS(NOW())-TO_DAYS(check_time) > 1) or
((sl.visits > 2500 or TO_DAYS(NOW())-TO_DAYS(sl.update_time) < 15) and TO_DAYS(NOW())-TO_DAYS(check_time) > 3) or
((sl.visits > 2000 or TO_DAYS(NOW())-TO_DAYS(sl.update_time) < 30) and TO_DAYS(NOW())-TO_DAYS(check_time) > 5) or
((sl.visits > 1000 or TO_DAYS(NOW())-TO_DAYS(sl.update_time) < 80) and TO_DAYS(NOW())-TO_DAYS(check_time) > 10) or
((sl.visits > 500 or TO_DAYS(NOW())-TO_DAYS(sl.update_time) < 100) and TO_DAYS(NOW())-TO_DAYS(check_time) > 20)) order by sl.update_time desc,n.update_time desc limit ?, ?
`, length, batch).Find(&list)
		if err != nil {
			global.LOG.Warn(fmt.Sprintf("从数据库获取待检测夸克数据失败: %+v", err))
		} else {
			s.awaitCheckQuark = append(s.awaitCheckQuark, utils.ListFilter(list, func(item table.NetdiskResource) bool {
				for _, resource := range s.awaitCheckQuark {
					if resource.Id == item.Id {
						return false
					}
				}
				return true
			})...)
			length = len(s.awaitCheckQuark)
		}
		if len(list) < batch {
			s.awaitCheckQuarkEmptyTime = time.Now()
		}
	}
	// 数据不够，全部返回
	if length < size {
		list := s.awaitCheckQuark
		s.awaitCheckQuark = make([]table.NetdiskResource, 0)
		return list
	}
	// 数据充足，拆分
	list := s.awaitCheckQuark[:size]
	s.awaitCheckQuark = s.awaitCheckQuark[size:]
	return list
}

// 取得指定数量待检测的百度网盘资源
func (s *ResourceService) GetCheckBaiduResource(size int) []table.NetdiskResource {
	s.awaitCheckBaiduMutex.Lock()
	defer s.awaitCheckBaiduMutex.Unlock()
	length := len(s.awaitCheckBaidu)
	batch := 300
	// 如果总长度不超过单次长度，且距离上次取得数据不全超过10分钟，追加获取数据
	if length < size && time.Now().Sub(s.awaitCheckBaiduEmptyTime).Minutes() > 10 {
		list := make([]table.NetdiskResource, 0)
		err := s.Db.SQL(`
select n.* from netdisk_resource n left join short_link sl on n.short_link = sl.alias
where type = 4 and status in (1,2) and TO_DAYS(NOW())-TO_DAYS(n.update_time) < 300 and (sl.update_time is null or TO_DAYS(NOW())-TO_DAYS(sl.update_time) < 600) and (check_time is null or TO_DAYS(NOW())-TO_DAYS(check_time) > 30 or
(TO_DAYS(NOW())-TO_DAYS(n.update_time) < 15 and TO_DAYS(NOW())-TO_DAYS(check_time) > 5) or
(TO_DAYS(NOW())-TO_DAYS(n.update_time) < 50 and TO_DAYS(NOW())-TO_DAYS(check_time) > 10) or
(TO_DAYS(NOW())-TO_DAYS(n.update_time) < 80 and TO_DAYS(NOW())-TO_DAYS(check_time) > 20) or
((sl.visits > 3000 or TO_DAYS(NOW())-TO_DAYS(sl.update_time) < 5) and TO_DAYS(NOW())-TO_DAYS(check_time) > 1) or
((sl.visits > 2500 or TO_DAYS(NOW())-TO_DAYS(sl.update_time) < 15) and TO_DAYS(NOW())-TO_DAYS(check_time) > 3) or
((sl.visits > 2000 or TO_DAYS(NOW())-TO_DAYS(sl.update_time) < 30) and TO_DAYS(NOW())-TO_DAYS(check_time) > 5) or
((sl.visits > 1000 or TO_DAYS(NOW())-TO_DAYS(sl.update_time) < 80) and TO_DAYS(NOW())-TO_DAYS(check_time) > 10) or
((sl.visits > 500 or TO_DAYS(NOW())-TO_DAYS(sl.update_time) < 100) and TO_DAYS(NOW())-TO_DAYS(check_time) > 20)) order by sl.update_time desc,n.update_time desc limit ?, ?
`, length, batch).Find(&list)
		if err != nil {
			global.LOG.Warn(fmt.Sprintf("从数据库获取待检测百度数据失败: %+v", err))
		} else {
			s.awaitCheckBaidu = append(s.awaitCheckBaidu, utils.ListFilter(list, func(item table.NetdiskResource) bool {
				for _, resource := range s.awaitCheckBaidu {
					if resource.Id == item.Id {
						return false
					}
				}
				return true
			})...)
			length = len(s.awaitCheckBaidu)
		}
		if len(list) < batch {
			s.awaitCheckBaiduEmptyTime = time.Now()
		}
	}
	// 数据不够，全部返回
	if length < size {
		list := s.awaitCheckBaidu
		s.awaitCheckBaidu = make([]table.NetdiskResource, 0)
		return list
	}
	// 数据充足，拆分
	list := s.awaitCheckBaidu[:size]
	s.awaitCheckBaidu = s.awaitCheckBaidu[size:]
	return list
}

// 取得需要删除的夸克资源
func (s *ResourceService) GetDeleteSearchResource(userId int64, typ enum.NetdiskType, minute int) []table.NetdiskResource {
	list := make([]table.NetdiskResource, 0)
	// TODO 1.0.0客户端有bug，不会限制时间，所以服务端要判断一下
	if minute <= 0 {
		return list
	}
	_ = s.WhereUserSession(userId).And("type = ? and origin = ? and short_link = '' and create_time < NOW() - INTERVAL ? MINUTE", typ, enum.NetdiskOrigin(2), minute).Find(&list)
	return list
}

// 通过目标路径取得网盘资源
func (s *ResourceService) GetByTargetUrl(userId int64, targetUrl string) (*table.NetdiskResource, error) {
	mod := &table.NetdiskResource{UserId: userId, TargetUrl: targetUrl}
	_, err := s.Db.Get(mod)
	return mod, err
}

// 通过id取得网盘资源
func (s *ResourceService) GetByOnlyId(id int64) (*table.NetdiskResource, error) {
	mod := &table.NetdiskResource{}
	_, err := s.Db.ID(id).Get(mod)
	return mod, err
}

// 管理员后台分页查询列表
func (s *ResourceService) GetByPage(userId int64, param request.NetdiskResourceQueryPageParam) (*response.PageResponse, error) {
	// 同一个时间可能包含多条数据，必须加上id做分页
	session := s.WhereUserSession(userId).Desc("create_time").Asc("id")
	if param.Type != "" {
		netdiskType, err := enum.NetdiskTypeValue(param.Type)
		if err != nil {
			return nil, err
		}
		session = session.And("type = ?", netdiskType)
	}
	if param.Status != "" {
		netdiskStatus, err := enum.NetdiskStatusValue(param.Status)
		if err != nil {
			return nil, err
		}
		session = session.And("status = ?", netdiskStatus)
	}
	if param.Origin != "" {
		netdiskOrigin, err := enum.NetdiskOriginValue(param.Origin)
		if err != nil {
			return nil, err
		}
		session = session.And("origin = ?", netdiskOrigin)
	}
	if param.UserName != "" {
		session = session.And("user_name = ?", param.UserName)
	}
	if param.Keyword != "" {
		session.And("name like concat('%',?,'%')", param.Keyword)
	}
	list := make([]table.NetdiskResource, 0)
	return s.HandlePageable(param.PageableParam, &list, session)
}

// 给客户端用的数据搜索接口
func (s *ResourceService) Search(userId int64, origin enum.ClientType, param request.NetdiskResourceSearchParam) (*response.PageResponse, error) {
	secureMode, _ := enum.NetdiskSecureModeValue(param.SecureMode)
	// 如果是审核模式，使用内部资源，关闭数据共享
	if secureMode == enum.NetdiskSecureMode(3) {
		userId = -1
		param.CollectTypes = nil
	}
	// 搜索数据库
	session := s.WhereUserSession(userId).And("status = ? and name != ''", enum.NetdiskStatus(1))
	if utf8.RuneCountInString(param.Keyword) > 3 && len(param.CollectTypes) == 0 { // 仅在没有使用共享资源库时采用
		session.And("MATCH(name) against(?)", param.Keyword)
	} else if param.Keyword != "" {
		session.And("name like concat('%',?,'%')", param.Keyword).Desc("create_time").Asc("id")
	} else {
		session.Desc("create_time").Asc("id")
	}
	list := make([]table.NetdiskResource, 0)
	result, err := s.HandlePageable(param.PageableParam, &list, session)
	if err != nil {
		return nil, err
	}
	// 如果搜索词是空的，不查询附加的内容，也不记录搜索结果
	if param.Keyword == "" {
		return result, nil
	}
	// 如果不是最后一页，就退出，不需要联网查询
	if result.HasNext {
		return result, nil
	}
	// 不使用共享采集，直接退出
	if len(param.CollectTypes) == 0 {
		// 第一页，记录搜索记录
		if param.Page == 0 {
			_, _ = s.Db.Insert(&table.NetdiskResourceSearch{UserId: userId, Keyword: param.Keyword, ResourceCount: int(result.Total), Origin: origin})
		}
		return result, nil
	}
	// 生成采集参数
	collectParam := request.NetdiskCollectParam{Keyword: param.Keyword, Types: utils.ListTransform(param.CollectTypes, func(item string) enum.NetdiskType {
		typ, _ := enum.NetdiskTypeValue(item)
		return typ
	})}
	resources := make([]table.NetdiskResource, 0)
	if userId == global.CONFIG.Plugin.Quark.UserId { // 是内部用户
		// 尝试采集
		resources = GetCollectService().Collect(collectParam)
	} else if global.SOCKET.GetClient(userId) != nil { // 如果是有客户端的用户，联网客户端进行查询
		// 旧版本依旧使用旧的采集方式
		serializedData, _ := json.Marshal(param.Keyword)
		if utils.VersionCode(global.SOCKET.GetClient(userId).Version) >= utils.VersionCode("1.1.4") {
			serializedData, _ = json.Marshal(collectParam)
		}
		message := common.SocketMessage{
			NeedResult: true,
			Timeout:    60 * time.Second,
			Type:       enum.SocketMessageType(6),
			Data:       string(serializedData),
		}
		existResult, cliErr := global.SOCKET.SendMessage(userId, message, &resources)
		// 查询错误或者没有等到结果
		if cliErr != nil || !existResult {
			global.LOG.Warn(fmt.Sprintf("客户端查询未能获取到结果: %+v", cliErr))
			return result, nil
		}
	} else if result.Total == 0 { // 是没有资源使用共享资源的用户
		// TODO 使用共享资源时暂时只提供夸克资源
		param.CollectTypes = []string{"QUARK"}
		// 尝试采集
		resources = GetCollectService().Collect(collectParam)
	}
	list = append(list, resources...)
	if len(resources) > 0 {
		// 第一页，记录搜索记录
		if param.Page == 0 {
			_, _ = s.Db.Insert(&table.NetdiskResourceSearch{UserId: userId, Keyword: param.Keyword, ResourceCount: len(list), Origin: origin})
		}
		return s.HandleContentPageable(&list, int64(len(list)), request.PageableParam{Page: 0, Size: len(list)}), nil
	}
	return result, nil
}

func handleImportFile(file *multipart.FileHeader) ([][]string, error) {
	f, err := file.Open() // 读取文件
	if err != nil {
		return nil, errors.Wrap(err, "打开文件失败")
	}
	defer f.Close()
	xlsx, err := excelize.OpenReader(f)
	if err != nil {
		return nil, errors.Wrap(err, "读取 Excel 文件失败")
	}
	defer xlsx.Close()
	// 获取活动工作表的名称
	sheetName := xlsx.GetSheetName(xlsx.GetActiveSheetIndex())
	// 获取工作表中的所有行
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		return nil, errors.Wrap(err, "读取 Excel 内容失败")
	}
	if len(rows) <= 1 || strings.TrimSpace(rows[0][0]) != "资源名称" || strings.TrimSpace(rows[0][1]) != "资源地址" {
		return nil, errors.New("资源导入文件格式错误，请重新下载文件")
	}
	if len(rows) > 20001 {
		return nil, errors.New("资源导入失败，一次性最多导入20000条")
	}
	// 处理数据前后空格并去重
	encountered := make(map[string]bool)
	writeIndex := 0
	for i := 1; i < len(rows); i++ {
		if len(rows[i]) < 2 {
			return nil, errors.New(fmt.Sprintf("资源导入失败，第%d条资源不完整"))
		}
		rows[i][0] = strings.TrimSpace(rows[i][0])
		rows[i][1] = strings.TrimSpace(rows[i][1])
		// 去重
		if rows[i][0] != "" && rows[i][1] != "" && !encountered[rows[i][1]] {
			encountered[rows[i][1]] = true
			rows[writeIndex] = rows[i]
			writeIndex++
		}
	}
	rows = rows[:writeIndex]
	return rows, nil
}

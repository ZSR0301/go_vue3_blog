package service

import (
	"errors"
	"fmt"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/model/response"
	"server/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type UserService struct {
}

// 使用预编译参数防止SQL注入
func (userService *UserService) Register(u database.User) (database.User, error) {
	if !errors.Is(global.DB.Where("email = ?", u.Email).First(&database.User{}).Error, gorm.ErrRecordNotFound) {
		return database.User{}, errors.New("this email address is already registered, please check the information you filled in, or retrieve your password")
	} //邮箱唯一性检查

	u.Password = utils.BcryptHash(u.Password) //使用 bcrypt 算法对明文密码进行哈希
	u.UUID = uuid.Must(uuid.NewV4())
	u.Avatar = "/image/avatar.jpg"
	u.RoleID = appTypes.User
	u.Register = appTypes.Email

	if err := global.DB.Create(&u).Error; err != nil {
		return database.User{}, err
	} //使用GORM的Create方法插入记录

	return u, nil
}

func (userService *UserService) EmailLogin(u database.User) (database.User, error) {
	var user database.User
	err := global.DB.Where("email = ?", u.Email).First(&user).Error
	//用户存在，验证密码
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok { //使用 BcryptCheck 比较输入密码和存储的哈希值
			return database.User{}, errors.New("incorrect email or password")
		}
		return user, nil
	}
	return database.User{}, err //用户不存在，返回错误信息
}

func (userService *UserService) QQLogin(accessTokenResponse other.AccessTokenResponse) (database.User, error) {
	var user database.User

	// 尝试查找用户
	err := global.DB.Where("openid = ?", accessTokenResponse.Openid).First(&user).Error
	//如果错误不是记录不存在，且不是其他错误，则返回错误信息
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return database.User{}, err
	}

	// 如果用户不存在，则创建新用户
	if errors.Is(err, gorm.ErrRecordNotFound) {
		userInfoResponse, err := ServiceGroupApp.QQService.GetUserInfoByAccessTokenAndOpenid(accessTokenResponse.AccessToken, accessTokenResponse.Openid)
		if err != nil {
			return database.User{}, err
		}
		user.UUID = uuid.Must(uuid.NewV4())
		user.Username = userInfoResponse.Nickname
		user.Openid = accessTokenResponse.Openid
		user.Avatar = userInfoResponse.FigureurlQQ2
		user.RoleID = appTypes.User
		user.Register = appTypes.QQ

		if err := global.DB.Create(&user).Error; err != nil {
			return database.User{}, err
		}
	}

	return user, nil
}

func (userService *UserService) ForgotPassword(req request.ForgotPassword) error {
	var user database.User
	if err := global.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return err
	}
	user.Password = utils.BcryptHash(req.NewPassword)
	return global.DB.Save(&user).Error
}

func (userService *UserService) UserCard(req request.UserCard) (response.UserCard, error) {
	var user database.User
	if err := global.DB.Where("uuid = ?", req.UUID).Select("uuid", "username", "avatar", "address", "signature").First(&user).Error; err != nil {
		return response.UserCard{}, err
	}
	return response.UserCard{
		UUID:      user.UUID,
		Username:  user.Username,
		Avatar:    user.Avatar,
		Address:   user.Address,
		Signature: user.Signature,
	}, nil
}

func (userService *UserService) Logout(c *gin.Context) {
	uuid := utils.GetUUID(c)                                                           //获取用户UUID
	jwtStr := utils.GetRefreshToken(c)                                                 //获取刷新令牌
	utils.ClearRefreshToken(c)                                                         //清除刷新令牌
	global.Redis.Del(uuid.String())                                                    //删除Redis中的用户信息
	_ = ServiceGroupApp.JwtService.JoinInBlacklist(database.JwtBlacklist{Jwt: jwtStr}) //使用 _ 忽略错误
}

func (userService *UserService) UserResetPassword(req request.UserResetPassword) error {
	var user database.User
	//根据用户ID查询单条记录，等效于First但不需要指定主键字段
	if err := global.DB.Take(&user, req.UserID).Error; err != nil {
		return err
	}
	if ok := utils.BcryptCheck(req.Password, user.Password); !ok {
		return errors.New("original password does not match the current account")
	}
	user.Password = utils.BcryptHash(req.NewPassword)
	return global.DB.Save(&user).Error
}

// 获取用户信息
func (userService *UserService) UserInfo(userID uint) (database.User, error) {
	var user database.User
	if err := global.DB.Take(&user, userID).Error; err != nil { //根据用户ID获取用户信息，如果用户不存在，则返回空值，并返回错误。
		return database.User{}, err
	}
	return user, nil //返回用户信息
}

// 更新用户信息
func (userService *UserService) UserChangeInfo(req request.UserChangeInfo) error {
	var user database.User
	if err := global.DB.Take(&user, req.UserID).Error; err != nil {
		return err
	}
	return global.DB.Model(&user).Updates(req).Error //Updates方法用于更新指定字段的值，更改User结构体中的username、address、signature字段的值。
}

func (userService *UserService) UserWeather(ip string) (string, error) {
	// 从redis中获取天气数据，如果没有数据，则调用高德api进行查询，天气不会突变，缓存一小时，避免频繁调用高德api。
	result, err := global.Redis.Get("weather-" + ip).Result()
	if err != nil { // // 缓存不存在的处理逻辑
		ipResponse, err := ServiceGroupApp.GaodeService.GetLocationByIP(ip)
		if err != nil {
			return "", err //如果出错，返回空字符串和错误信息
		}
		live, err := ServiceGroupApp.GaodeService.GetWeatherByAdcode(ipResponse.Adcode) //获取的地区编码（Adcode）查询天气信息
		if err != nil {
			return "", err //如果出错，返回空字符串和错误信息
		}

		weather := "地区：" + live.Province + "-" + live.City + " 天气：" + live.Weather + " 温度：" + live.Temperature + "°C" + " 风向：" + live.WindDirection + " 风级：" + live.WindPower + " 湿度：" + live.Humidity + "%"
		//将获取的天气数据拼接成易读的字符串格式
		// 将天气数据存入redis

		//此时返回的是一个 *StatusCmd 对象（表示一个 Redis 命令状态），但命令还未真正执行，所以需要调用它的 Err() 方法获取命令执行结果。
		if err := global.Redis.Set("weather-"+ip, weather, time.Hour*1).Err(); err != nil {
			return "", err
		}
		return weather, nil
	}
	return result, nil // 直接返回缓存结果
}

func (userService *UserService) UserChart(req request.UserChart) (response.UserChart, error) {
	// 构建查询条件
	where := global.DB.Where(fmt.Sprintf("date_sub(curdate(), interval %d day) <= created_at", req.Date))
	//使用date_sub函数计算日期范围，动态插入req.Date参数
	var res response.UserChart

	// 生成日期列表
	startDate := time.Now().AddDate(0, 0, -req.Date)
	for i := 1; i <= req.Date; i++ {
		res.DateList = append(res.DateList, startDate.AddDate(0, 0, i).Format("2006-01-02"))
	}
	// 获取登录数据
	loginCounts := utils.FetchDateCounts(global.DB.Model(&database.Login{}), where)
	// 获取注册数据
	registerCounts := utils.FetchDateCounts(global.DB.Model(&database.User{}), where)

	for _, date := range res.DateList {
		loginCount := loginCounts[date]
		registerCount := registerCounts[date]
		res.LoginData = append(res.LoginData, loginCount)
		res.RegisterData = append(res.RegisterData, registerCount)
	}

	return res, nil
}

func (userService *UserService) UserList(info request.UserList) (interface{}, int64, error) {
	db := global.DB

	if info.UUID != nil {
		db = db.Where("uuid = ?", info.UUID)
	}

	if info.Register != nil {
		db = db.Where("register = ?", info.Register)
	}

	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}

	return utils.MySQLPagination(&database.User{}, option)
}

func (userService *UserService) UserFreeze(req request.UserOperation) error {
	var user database.User
	if err := global.DB.Take(&user, req.ID).Update("freeze", true).Error; err != nil {
		return err
	}
	//数据库标记与令牌失效同步进行，即使令牌未过期也无法使用
	jwtStr, _ := ServiceGroupApp.JwtService.GetRedisJWT(user.UUID)
	if jwtStr != "" {
		_ = ServiceGroupApp.JwtService.JoinInBlacklist(database.JwtBlacklist{Jwt: jwtStr})
	}

	return nil
}

func (userService *UserService) UserUnfreeze(req request.UserOperation) error {
	return global.DB.Take(&database.User{}, req.ID).Update("freeze", false).Error
}

func (userService *UserService) UserLoginList(info request.UserLoginList) (interface{}, int64, error) {

	db := global.DB

	if info.UUID != nil {
		var userID uint
		//使用Pluck直接提取单个字段值
		if err := global.DB.Model(database.User{}).Where("uuid = ?", *info.UUID).Pluck("id", &userID); err != nil {
			return nil, 0, nil
		}
		db = db.Where("user_id = ?", userID)
	}

	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
		Preload:  []string{"User"},
	}

	return utils.MySQLPagination(&database.Login{}, option) //分页查询执行，指定查询的登录记录模型，包含分页和查询条件的配置；返回登录记录列表、总记录数、可能的错误信息
}

package main

import (
	"ShopBackground/config"
	"ShopBackground/controller"
	"ShopBackground/datasource"
	"ShopBackground/model"
	"ShopBackground/service"
	"ShopBackground/utils"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"strconv"
	"time"
)

func main() {
	app := newApp()

	//应用App设置
	configation(app)

	//路由设置
	mvcHandle(app)

	config := config.InitConfig()
	addr := ":" + config.Port
	app.Run(
		iris.Addr(addr),                               //在端口8080进行监听
		iris.WithoutServerError(iris.ErrServerClosed), //无服务错误提示
		iris.WithOptimizations,                        //对json数据序列化更快的配置
	)
}

//构建App
func newApp() *iris.Application {
	app := iris.New()

	//设置日志级别  开发阶段为debug
	app.Logger().SetLevel("debug")

	//注册静态资源
	app.StaticWeb("/static", "./static")
	app.StaticWeb("/manage/static", "./static")
	app.StaticWeb("/img", "./static/img")

	//注册视图文件
	app.RegisterView(iris.HTML("./static", ".html"))
	app.Get("/", func(context context.Context) {
		context.View("index.html")
	})

	return app
}

/**
 * MVC 架构模式处理
 */
func mvcHandle(app *iris.Application) {

	//启用session
	sessManager := sessions.New(sessions.Config{
		Cookie:  "sessioncookie",
		Expires: 24 * time.Hour,
	})

	//获取redis实例
	redis := datasource.NewRedis()
	//设置session的同步设置为redis
	sessManager.UseDatabase(redis)

	//实例化mysql数据库引擎
	engine := datasource.NewMysqlEngine()

	//管理员模块功能
	adminService := service.NewAdminService(engine)

	admin := mvc.New(app.Party("/admin"))
	admin.Register(
		adminService,
		sessManager.Start,
	)
	admin.Handle(new(controller.AdminController))

	//用户功能模块
	userService := service.NewUserService(engine)

	user := mvc.New(app.Party("/v1/users"))
	user.Register(
		userService,
		sessManager.Start,
	)
	user.Handle(new(controller.UserController))

	//获取用户详细信息
	app.Get("/v1/user/{user_name}", func(context context.Context) {
		userName := context.Params().Get("user_name")
		var user model.User
		_, err := engine.Where(" user_name = ? ", userName).Get(&user)
		if err != nil {
			context.JSON(iris.Map{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_USERINFO,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_USERINFO),
			})
		} else {
			context.JSON(user)
		}
	})
	//获取地址信息
	app.Get("/v1/addresse/{address_id}", func(context context.Context) {
		address_id := context.Params().Get("address_id")

		addressID, err := strconv.Atoi(address_id)
		if err != nil {
			context.JSON(iris.Map{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_ORDERINFO,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_ORDERINFO),
			})
		}

		var address model.Address
		_, err = engine.Id(addressID).Get(&address)
		if err != nil {
			context.JSON(iris.Map{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_ORDERINFO,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_ORDERINFO),
			})
		}
		//查询数据成功
		context.JSON(address)
	})

	//统计功能模块
	statisService := service.NewStatisService(engine)
	statis := mvc.New(app.Party("/statis/{model}/{date}/"))
	statis.Register(
		statisService,
		sessManager.Start,
	)
	statis.Handle(new(controller.StatisController))

	//	订单模块
	orderService := service.NewOrderService(engine)
	order := mvc.New(app.Party("/bos/orders/"))
	order.Register(
		orderService,
		sessManager.Start,
	)
	order.Handle(new(controller.OrderController)) //控制器

	//商铺模块
	shopService := service.NewShopService(engine)
	shop := mvc.New(app.Party("/shopping/restaurants"))
	shop.Register(
		shopService,
		sessManager.Start,
	)
	shop.Handle(new(controller.ShopController)) //控制器

}

/**
 * 项目设置
 */
func configation(app *iris.Application) {

	//配置 字符编码
	app.Configure(iris.WithConfiguration(iris.Configuration{
		Charset: "UTF-8",
	}))

	//错误配置
	//未发现错误
	app.OnErrorCode(iris.StatusNotFound, func(context context.Context) {
		context.JSON(iris.Map{
			"errmsg": iris.StatusNotFound,
			"msg":    " not found ",
			"data":   iris.Map{},
		})
	})

	app.OnErrorCode(iris.StatusInternalServerError, func(context context.Context) {
		context.JSON(iris.Map{
			"errmsg": iris.StatusInternalServerError,
			"msg":    " interal error ",
			"data":   iris.Map{},
		})
	})
}

package main

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12/middleware/logger"

	//"github.com/kataras/iris/v12/middleware/logger"
	recover "github.com/kataras/iris/v12/middleware/recover"

	//"github.com/kataras/golog@v0.0.10"
	"net/http"
	"os"
	"time"
	. "zyg/datamodels"
	"zyg/repositories"
	"zyg/service"
	"zyg/web/controllers"

	"github.com/kataras/iris/v12"

	//	"github.com/kataras/iris/cache"
	"github.com/BurntSushi/toml"
	"github.com/dchest/captcha"

	//	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris/v12/core/host"
	"gopkg.in/olivere/elastic.v5"

	"github.com/BelieveR44/goseaweedfs"

	"github.com/kataras/iris/v12/mvc"
	"github.com/urfave/cli"
)

type configs struct {
	Name          string   `toml:"name"`
	Url           string   `toml:"url"`
	UrlM          string   `toml:"urlm"`
	Addr          string   `toml:"addr"`
	Class         []string `toml:"class"`
	MarkMax       int      `toml:"markMax"`
	OnSon         string   `toml:"onSon"`
	UrlImg        string   `toml:"imgUrl"`
	MainPort      string   `toml:"mainport"`
	SonPort       string   `toml:"sonport"`
	ImgPort       string   `toml:"imgport"`
	Rec           []int    `toml:"rec"`
	Dir           dirConfig
	Weed          weedConfig
	Mysql         mysqlConfig
	Redis         redisConfig
	ElasticSearch elasticSearchConfig
}
type dirConfig struct {
	Template string `toml:"template"`
	Main     string `toml:"main"`
	Log      string `toml:"log"`
	Mainm    string `toml:"mainm"`
	Logm     string `toml:"logm"`
	Img      string `toml:"img"`
}
type weedConfig struct {
	Master  string   `toml:"master"`
	Filer   []string `toml:"filer"`
	Volume  []string `toml:"volume"`
	Timeout duration `toml:"timeout"`
}
type mysqlConfig struct {
	Addr string `toml:"addr"`
}
type redisConfig struct {
	Addr       string   `toml:"addr"`
	Pass       string   `toml:"pass"`
	SessDb     int      `toml:"sessDb"`
	SessExpiry duration `toml:"sessExpiry"`
	BookDb     int      `toml:"bookDb"`
	BookExpiry duration `toml:"bookExpiry"`
	Tag        string   `toml:"tag"`
}
type elasticSearchConfig struct {
	Host  string `toml:"host"`
	Index string `toml:"index"`
}
type duration struct {
	time.Duration
}

var (
	Config configs
	DB     sqlx.DB
	///-------------------------------
//	templateDir = "views"
//	serverDir   = ""
//	logDir      = ""
//	addr        = "yxj.com:8080"
)

func main() {
	//-------------------初始化变量-------------------------------------
	var configDir string
	appcli := cli.NewApp()
	appcli.Name = "RatelWebServer"
	appcli.Version = "1.0.0"
	appcli.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config,c",
			Value:       "config.toml",
			Usage:       "is config dir",
			Destination: &configDir,
		},
	}
	appcli.Action = func(c *cli.Context) error {
		fmt.Printf("start RatelWebServer")
		return nil
	}
	// 接受os.Args启动程序
	appcli.Run(os.Args)
	////////////////解析toml////////////
	if _, err := toml.DecodeFile(configDir, &Config); err != nil {
		fmt.Println(err)
		return
	}
	//--------------------------------------------------------
	repositories.RedisConfig.Addr = Config.Redis.Addr
	repositories.RedisConfig.Pass = Config.Redis.Pass
	repositories.RedisConfig.Tag = Config.Redis.Tag
	repositories.RedisConfig.BookDb = Config.Redis.BookDb
	repositories.RedisConfig.SessDb = Config.Redis.SessDb
	repositories.Sw = goseaweedfs.NewSeaweed("http", Config.Weed.Master, Config.Weed.Filer, 2*1024*1024, Config.Weed.Timeout.Duration)
	repositories.SwConfig.Master = Config.Weed.Master
	repositories.SwConfig.Filer = Config.Weed.Filer
	repositories.SwConfig.Volume = Config.Weed.Volume
	repositories.SwConfig.Timeout = Config.Weed.Timeout.Duration
	repositories.DB = mysqlLink(Config.Mysql.Addr)
	service.CacheTime.RedisBookExpiry = Config.Redis.BookExpiry.Duration
	service.CacheTime.RedisSessExpiry = Config.Redis.SessExpiry.Duration
	service.RecIdList = Config.Rec
	repositories.EsClient = EsLink(Config.ElasticSearch.Host)
	Scon.EsIndex = Config.ElasticSearch.Index
	Web.Name = Config.Name
	Web.Class = Config.Class
	Web.Url = Config.Url
	Web.UrlM = Config.UrlM
	Web.UrlImg = Config.UrlImg
	//fmt.Printf("%v", Config)
	//f := newLogFile()
	app := iris.New()
	sonapp := iris.New()
	imgapp := iris.New()
	//（可选）添加两个内置处理程序
	// 可以从任何http相关的恐慌中恢复
	// 并将请求记录到终端。
	app.Use(recover.New())
	app.Use(logger.New())
	app.Logger().SetLevel("error")
	sonapp.Use(recover.New())
	sonapp.Use(logger.New())
	sonapp.Logger().SetLevel("error")
	imgapp.Use(recover.New())
	imgapp.Use(logger.New())
	imgapp.Logger().SetLevel("error")
	//app.Logger().SetLevel("debug")
	//golog.SetLevel("error")
	//app.Logger().SetLevel("error")
	//golog.SetOutput(f)
	//app.Logger().SetOutput(f)
	////注册域名头
	//www := app.Subdomain("www")
	//www := app.WWW()
	/////错误控制
	app.OnErrorCode(iris.StatusNotFound, NotFoundHandler)
	app.OnErrorCode(iris.StatusInternalServerError, InternalServerError)
	app.OnErrorCode(iris.StatusServiceUnavailable, ServiceUnavailable)
	app.OnErrorCode(iris.StatusForbidden, StatusForbiddenHandler)
	sonapp.OnErrorCode(iris.StatusNotFound, NotFoundHandler)
	sonapp.OnErrorCode(iris.StatusInternalServerError, InternalServerError)
	sonapp.OnErrorCode(iris.StatusServiceUnavailable, ServiceUnavailable)
	sonapp.OnErrorCode(iris.StatusForbidden, StatusForbiddenHandler)
	imgapp.OnErrorCode(iris.StatusNotFound, NotFoundHandler)
	imgapp.OnErrorCode(iris.StatusInternalServerError, InternalServerError)
	imgapp.OnErrorCode(iris.StatusServiceUnavailable, ServiceUnavailable)
	imgapp.OnErrorCode(iris.StatusForbidden, StatusForbiddenHandler)
	/////载入模板
	//app.RegisterView(iris.HTML(Config.Dir.Template, ".html").Reload(true))
	app.RegisterView(iris.HTML(Config.Dir.Template, ".html"))
	sonapp.RegisterView(iris.HTML(Config.Dir.Template, ".html"))
	imgapp.RegisterView(iris.HTML(Config.Dir.Template, ".html"))
	app.Get("/captcha/{f:path}", iris.FromStd(captcha.Server(126, 39)))

	fileserverjs := iris.FileServer(Config.Dir.Main+"js", iris.DirOptions{ShowList: false, Gzip: true})
	js := iris.StripPrefix("/js", fileserverjs)
	app.Get("/js/{f:path}", js)

	fileservercss := iris.FileServer(Config.Dir.Main+"css", iris.DirOptions{ShowList: false, Gzip: true})
	css := iris.StripPrefix("/css", fileservercss)
	app.Get("/css/{f:path}", css)
	fileserveri := iris.FileServer(Config.Dir.Main+"i", iris.DirOptions{ShowList: false, Gzip: true})
	i := iris.StripPrefix("/i", fileserveri)
	app.Get("/i/{f:path}", i)
	fileserverimg := iris.FileServer(Config.Dir.Main+"img", iris.DirOptions{ShowList: false, Gzip: true})
	img := iris.StripPrefix("/img", fileserverimg)
	app.Get("/img/{f:path}", img)

	fileserverfonts := iris.FileServer(Config.Dir.Main+"fonts", iris.DirOptions{ShowList: false, Gzip: true})
	fonts := iris.StripPrefix("/fonts", fileserverfonts)
	app.Get("/fonts/{f:path}", fonts)
	// app.Get("/static/{f:path}", h)
	//app.SPA(fileserver)
	//app.RegisterView(iris.HTML("./views", ".html").Reload(true))
	// 基于根路由器服务控制器， "/".
	app.Get("/sitemap/index/{index:string}", controllers.GetSitemapIndexMain)
	//app.Get("/sitemap/index/{index:string}", controllers.GetSitemapIndexAll)
	//注册控制器
	mvc.New(app).Handle(new(controllers.MainController))
	mvc.New(app.Party("/read")).Handle(new(controllers.ReadController))
	//	mvc.New(www.Party("/list", cache.Handler(time.Hour*6))).Handle(new(controllers.ListController))
	//	mvc.New(www.Party("/full", cache.Handler(time.Hour*6))).Handle(new(controllers.FullController))
	//mvc.New(www.Party("/list")).Handle(new(controllers.ListController))
	//mvc.New(www.Party("/full")).Handle(new(controllers.FullController))
	mvc.New(app.Party("/bookcase")).Handle(new(controllers.BookCaseController))

	//app.SubdomainRedirect(app, www)
	if Config.OnSon != "" {
		//M := app.Subdomain(Config.OnSon)
		/////错误控制
		//		app.OnErrorCode(iris.StatusNotFound, NotFoundHandler)
		//		app.OnErrorCode(iris.StatusInternalServerError, InternalServerError)
		//		app.OnErrorCode(iris.StatusServiceUnavailable, ServiceUnavailable)
		//		app.OnErrorCode(iris.StatusForbidden, StatusForbiddenHandler)
		/////载入模板
		//app.RegisterView(iris.HTML(Config.Dir.Template, ".html"))
		sonapp.Get("/captcha/{f:path}", iris.FromStd(captcha.Server(126, 39)))
		fileserverjs := iris.FileServer(Config.Dir.Mainm+"js", iris.DirOptions{ShowList: false, Gzip: true})
		js := iris.StripPrefix("/js", fileserverjs)
		sonapp.Get("/js/{f:path}", js)

		fileservercss := iris.FileServer(Config.Dir.Mainm+"css", iris.DirOptions{ShowList: false, Gzip: true})
		css := iris.StripPrefix("/css", fileservercss)
		sonapp.Get("/css/{f:path}", css)

		fileserveri := iris.FileServer(Config.Dir.Mainm+"i", iris.DirOptions{ShowList: false, Gzip: true})
		i := iris.StripPrefix("/i", fileserveri)
		sonapp.Get("/i/{f:path}", i)

		fileserverimg := iris.FileServer(Config.Dir.Mainm+"img", iris.DirOptions{ShowList: false, Gzip: true})
		img := iris.StripPrefix("/img", fileserverimg)
		sonapp.Get("/img/{f:path}", img)

		fileserverfonts := iris.FileServer(Config.Dir.Mainm+"fonts", iris.DirOptions{ShowList: false, Gzip: true})
		fonts := iris.StripPrefix("/fonts", fileserverfonts)
		sonapp.Get("/fonts/{f:path}", fonts)
		// app.Get("/static/{f:path}", h)
		//app.SPA(fileserver)
		//app.RegisterView(iris.HTML("./views", ".html").Reload(true))
		// 基于根路由器服务控制器， "/".
		//注册控制器
		sonapp.Get("/sitemap/index/{index:string}", controllers.GetSitemapIndexSon)
		mvc.New(sonapp).Handle(new(controllers.MainMController))
		mvc.New(sonapp.Party("/read")).Handle(new(controllers.ReadMController))
		//mvc.New(M.Party("/list")).Handle(new(controllers.ListMController))
		mvc.New(sonapp.Party("/bookcase")).Handle(new(controllers.BookCaseMController))
		//app.SubdomainRedirect(app, M)
	}
	//Img := app.Subdomain("img")
	fileserverimg = iris.FileServer(Config.Dir.Img, iris.DirOptions{ShowList: false, Gzip: true})
	ImgDir := iris.StripPrefix("/", fileserverimg)
	imgapp.Get("/{f:path}", ImgDir)
	//app.SubdomainRedirect(app, Img)
	//----------build
	if err := app.Build(); err != nil {
		fmt.Println("app.Build() err")
		panic(err)
	}
	if err := sonapp.Build(); err != nil {
		fmt.Println("sonapp.Build() err")
		panic(err)
	}
	if err := imgapp.Build(); err != nil {
		fmt.Println("imgapp.Build() err")
		panic(err)
	}
	//---------
	// 80启动方式普通启动
	server := &http.Server{
		Addr: Config.MainPort,
		//Addr: "yxj.com:8080",
		Handler:     app,
		ReadHeaderTimeout:  30 * time.Second,
		//WriteTimeout: 45 * time.Second,
		IdleTimeout: 60 * time.Second,
	}
	sonserver := &http.Server{
		Addr: Config.SonPort,
		//Addr: "yxj.com:8080",
		Handler:     sonapp,
		ReadHeaderTimeout:  30 * time.Second,
		//WriteTimeout: 45 * time.Second,
		IdleTimeout: 60 * time.Second,
	}
	imgserver := &http.Server{
		Addr: Config.ImgPort,
		//Addr: "yxj.com:8080",
		Handler:     imgapp,
		ReadHeaderTimeout:  30 * time.Second,
		//WriteTimeout: 45 * time.Second,
		IdleTimeout: 60 * time.Second,
	}
	//sitemap

	//sonserver.ListenAndServe()
	//testserver.ListenAndServe()
	go sonserver.ListenAndServe()
	go imgserver.ListenAndServe()
	//go sonapp.Run(iris.Server(sonserver), iris.WithoutServerError(iris.ErrServerClosed))
	//go imgapp.Run(iris.Server(imgserver), iris.WithoutServerError(iris.ErrServerClosed))
	app.Run(iris.Server(server), iris.WithoutServerError(iris.ErrServerClosed))
	//server.ListenAndServe()

	//app.Listen(server, iris.WithoutServerError(iris.ErrServerClosed))

	//app.Run(localAutoTLS("ziyouge.com:443", "m.ziyouge.com www.ziyouge.com img.ziyouge.com ziyouge.com", "pptomenz@protonmail.com"), iris.WithoutServerError(iris.ErrServerClosed))
	//app.Run(iris.Addr("yxj.com:8080"))
	//---------服务器
	//app.Run(iris.Addr("yixuanju.com:80"))

}
func localAutoTLS( //构造带有timeout的 AutoTls
	addr string,
	domain string, email string,
	hostConfigs ...host.Configurator) iris.Runner {
	return func(app *iris.Application) error {
		return app.NewHost(&http.Server{Addr: addr, ReadTimeout: 20 * time.Second, WriteTimeout: 20 * time.Second}).
			Configure(hostConfigs...).
			ListenAndServeAutoTLS(domain, email, "letscache")
	}
}
func StatusForbiddenHandler(ctx iris.Context) {
	errinfo := PagesState{Err: "403",
		ErrInfo:     "没有找到您所访问的页面",
		ErrSynopsis: "没有找到您所访问的页面,他可能已经迁移到其他位置!",
		Web:         Web}
	ctx.ViewData("errinfo", errinfo)
	ctx.View("pubstate.html")
}
func NotFoundHandler(ctx iris.Context) {
	errinfo := PagesState{Err: "404",
		ErrInfo:     "没有找到您所访问的页面",
		ErrSynopsis: "没有找到您所访问的页面,他可能已经迁移到其他位置!",
		Web:         Web}
	ctx.ViewData("errinfo", errinfo)
	ctx.View("pubstate.html")
}
func InternalServerError(ctx iris.Context) {
	//ctx.HTML("500气死我也")
	errinfo := PagesState{Err: "500",
		ErrInfo:     "服务器出现问题",
		ErrSynopsis: "服务器出现问题,请稍后重试!",
		Web:         Web}
	ctx.ViewData("errinfo", errinfo)
	ctx.View("pubstate.html")
}
func ServiceUnavailable(ctx iris.Context) {
	errinfo := PagesState{Err: "503",
		ErrInfo:     "服务器繁忙请稍后重试",
		ErrSynopsis: "服务器繁忙请稍后重试",
		Web:         Web}
	ctx.ViewData("errinfo", errinfo)
	ctx.View("pubstate.html")
}
func todayFilename() string {
	today := time.Now().Format("Jan 02 2006")
	return today + ".txt"
}

func newLogFile() *os.File {
	filename := todayFilename()
	// Open the file, this will append to the today's file if server restarted.
	f, err := os.OpenFile(Config.Dir.Log+filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return f
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
func mysqlLink(Dsn string) *sqlx.DB {
	//var err error
	db, err := sqlx.Open("mysql", Config.Mysql.Addr)
	if err != nil {
		panic(err)
		return db
	}
	db.SetConnMaxLifetime(time.Second * 599)
	db.SetMaxIdleConns(100)
	return db
}
func EsLink(host string) *elastic.Client {
	var err error
	Client, err := elastic.NewClient(elastic.SetURL(host))
	if err != nil {
		panic(err)
	}
	info, code, err := Client.Ping(host).Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	return Client
}

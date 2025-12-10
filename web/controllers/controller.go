package controllers

import (
	//	"fmt"
	//	"fmt"
	//	"fmt"
	"time"
	. "zyg/datamodels"
	"zyg/repositories"
	"zyg/service"

	//	"github.com/go-redis/redis"
	"html/template"

	//	"github.com/dchest/captcha"
	valid "github.com/asaskevich/govalidator"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/cache"
	"github.com/kataras/iris/v12/mvc"
)

type MainController struct {
	Ctx iris.Context
	//Cache iris.Cache
}
type ReadController struct {
	Ctx iris.Context
}

//type ListController struct {
//	Ctx iris.Context
//}
//type FullController struct {
//	Ctx iris.Context
//}

// Get serves
// Method:   GET
// Resource: http://localhost:8080
func (c *MainController) Get() mvc.Result {
	c.Ctx.AddHandler(cache.Handler(time.Hour * 6))
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	res8img, err := service.GetRecId(8)
	if err != nil {
		return mvc.Response{Code: 503}
	}
	res5, err := service.GetRecId(5)
	if err != nil {
		return mvc.Response{Code: 503}
	}
	res9img, _ := service.GetRecId(9)
	if err != nil {
		return mvc.Response{Code: 503}
	}
	res16, _ := service.GetRecId(16)
	if err != nil {
		return mvc.Response{Code: 503}
	}
	res15, _ := service.GetLastUpdate(15)
	if err != nil {
		return mvc.Response{Code: 503}
	}
	c.Ctx.ViewData("res8img", res8img)
	c.Ctx.ViewData("res9img", res9img)
	c.Ctx.ViewData("res5", res5)
	c.Ctx.ViewData("res16", res16)
	c.Ctx.ViewData("res15", res15)
	//x := ServerConfig()
	c.Ctx.ViewData("WebInfo", Web)
	c.Ctx.ViewData("Date", time.Now().Unix())
	return mvc.View{Name: "ziyouge/Index/index.html"}
}

// 章节列表页阅读页
// Method:   GET
// resource:http://localhost:8080/book/1
func (c *MainController) GetBookBy(Id int) mvc.Result {
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	book, err := service.GetBook(Id)
	if err == BState.NotFound {
		return mvc.Response{Code: 404}
	} else if err != nil {
		//fmt.Print("\nbookerr:", err)
		return mvc.Response{Code: 503}
	}
	res6, err := service.GetRecId(6)
	if err != nil {
		return mvc.Response{Code: 503}
	}
	c.Ctx.ViewData("Intro", template.HTML(book.Article.Intro))
	c.Ctx.ViewData("Book", book)
	c.Ctx.ViewData("res6", res6)
	return mvc.View{
		Name: "ziyouge/Book/index.html",
	}

}

// 章节阅读页
// Method:   GET
// resource:http://localhost:8080/book/1/2
func (c *ReadController) GetBy(Id int, Cid int) mvc.Result {
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	book, err := service.GetBookChapter(Id, Cid)
	//fmt.Printf("%v",book)
	if err == BState.NotFound {
		return mvc.Response{Code: 404}
	} else if err != nil {
		//ctx.Application().Logger().Infof("Begin request for path: %s", ctx.Path())
		c.Ctx.Application().Logger().Infof("chapter503,err:%v", err)
		return mvc.Response{Code: 503}
	} else {
		c.Ctx.ViewData("Content", template.HTML(book.CurrentChapter.CurrentChapterContent))
		c.Ctx.ViewData("Book", book)
		return mvc.View{Name: "ziyouge/Book/read.html"}
	}
}
func (c *MainController) GetSo() mvc.Result {
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	res10, err := service.GetRecId(10)
	if err != nil {
		return mvc.Response{Code: 503}
	}
	c.Ctx.ViewData("res10", res10)
	c.Ctx.ViewData("WebInfo", Web)
	return mvc.View{Name: "ziyouge/Index/search.html"}
}
func (c *MainController) PostSoapi() mvc.Result {
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	var Search JsonSoapi
	err := c.Ctx.ReadJSON(&Search)
	if err != nil {
		return mvc.Response{Object: service.ReAjax("出错啦", 0, "")}
	}
	result, err := valid.ValidateStruct(Search)
	if err != nil {
		return mvc.Response{Object: service.ReAjax("请正确提交", 0, "")}
	}
	if !result {
		return mvc.Response{Object: service.ReAjax("查询格式不正确", 0, "")}
	}

	res, err := service.SearchByES("zyg", "article", Search.Q, 5, Search.P)
	if err != nil {
		return mvc.Response{Code: 503}
	}
	if res.Pages != 0 {
		res.Status = 1
	}
	return mvc.Response{Object: res}
}

// GetHello serves
// Method:   GET
// Resource: http://localhost:8080/hello
//func (c *MainController) GetHello() interface{} {
//	return map[string]string{"message": "Hello Iris!"}

//}

//在控制器适应主应用程序之前调用一次BeforeActivation
//当然在服务器运行之前。
//在版本9之后，您还可以为特定控制器的方法添加自定义路由。
//在这里您可以注册自定义方法的处理程序
//使用带有`ca.Router`的标准路由器做一些你可以做的事情，没有mvc，
//并添加将绑定到控制器的字段或方法函数的输入参数的依赖项。
func (c *MainController) BeforeActivation(b mvc.BeforeActivation) {
	//b.Handle("GET", "/", "CustomGet", cache.Handler(time.Second*3000))
	//b.Handle("GET", "/list", "CustomList", cache.Handler(time.Hour*3))
	b.Handle("GET", "/list", "CustomNewList1", cache.Handler(time.Hour*6))
	b.Handle("GET", "/list/{class:int}", "CustomNewList2", cache.Handler(time.Hour*6))
	b.Handle("GET", "/list/{class:int}/{page:int}", "CustomNewList3", cache.Handler(time.Hour*6))
	b.Handle("GET", "/full", "CustomFull1", cache.Handler(time.Hour*6))
	b.Handle("GET", "/full/{class:int}", "CustomFull2", cache.Handler(time.Hour*6))
	b.Handle("GET", "/full/{class:int}/{page:int}", "CustomFull3", cache.Handler(time.Hour*6))
	b.Handle("GET", "/bookcase", "Bookcase1")
	b.Handle("GET", "/bookcase/{sort:int}/{class:int}", "Bookcase2")
	b.Handle("GET", "/bookcase/{sort:int}/{class:int}/{page:int}", "Bookcase3")
	//甚至添加基于此控制器路由器的全局中间件，
	//在这个例子中是根“/”：
	// b.Router（）。使用（myMiddleware）
	//b.Router().Get("/book{id:int}",customFunc)
	//这个功能简直无敌
}

///////////////
//用户书架
//url:/bookcase
func (c *MainController) Bookcase1() mvc.Result {
	return c.customBookcase(0, 0, 0)
}
func (c *MainController) Bookcase2(sort int, class int) mvc.Result {
	return c.customBookcase(sort, class, 0)
}
func (c *MainController) Bookcase3(sort int, class int, page int) mvc.Result {
	return c.customBookcase(sort, class, page)
}

/////////////////
// 书籍列表页
// Method:   GET
// resource:http://localhost:8080/List/1/0/0/0
func (c *MainController) CustomNewList1() mvc.Result {
	return c.CustomList(0, 0, 0)
}
func (c *MainController) CustomNewList2(class int) mvc.Result {
	return c.CustomList(class, 0, 0)
}
func (c *MainController) CustomNewList3(class int, page int) mvc.Result {
	return c.CustomList(class, 0, page)
}
func (c *MainController) CustomFull1() mvc.Result {
	return c.CustomList(0, 2, 0)
}
func (c *MainController) CustomFull2(class int) mvc.Result {
	return c.CustomList(class, 2, 0)
}
func (c *MainController) CustomFull3(class int, page int) mvc.Result {
	return c.CustomList(class, 2, page)
}
func (c *MainController) CustomList(class int, full int, page int) mvc.Result {
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	var SortName string
	c.Ctx.ViewData("Date", time.Now().Unix())
	//articlelist, pages, err := repositories.QueryList(repositories.DB, class, full, label, size, 0)
	articlelist, pages, err := repositories.QueryList(repositories.DB, class, full, 0, 0, page)
	if err != nil {
		return mvc.Response{Code: 404}
	}

	//	for x, y := range Web.Class {
	//		if x == class {
	//			SortName = y
	//			break
	//		}
	//	}
	SortName = Web.Class[class]
	if SortName == "" {
		return mvc.Response{Code: 503}
	}
	pages.ClassName = Web.Class
	pages.FullName = []string{"全部", "连载", "全本"}
	pages.LabelName = []string{"全部", "收藏", "点击", "推荐", "入库"}
	pages.SizeName = []string{"全部", "10万以下", "10万-30万", "30万-50万", "50万-100万", "100万-200万", "200万以上"}
	pages.Class = class
	//	pages.Full = full
	//	pages.Label = label
	//	pages.Size = size
	for x, y := range articlelist {
		if y.Siteid >= 0 && y.Siteid <= 10 {
			articlelist[x].Class = Web.Class[y.Sortid]
		} else {
			articlelist[x].Class = Web.Class[0]
		}
		articlelist[x].LastTime = time.Unix(int64(y.Lastupdate), 0)
		articlelist[x].ImgUrl = service.GetImgUrl(y.Articleid, y.Imgflag, Web.UrlImg)
	}
	res8img, _ := service.GetRecId(8)
	if err != nil {
		return mvc.Response{Code: 503}
	}
	listdir := []string{"list", "dnf", "full"}
	listlable := []string{"最新", "连载", "全本"}
	c.Ctx.ViewData("date", time.Now().Unix())
	c.Ctx.ViewData("WebInfo", Web)
	c.Ctx.ViewData("ArticleList", articlelist)
	c.Ctx.ViewData("Pages", pages)
	c.Ctx.ViewData("SortName", SortName)
	c.Ctx.ViewData("Full", full)
	c.Ctx.ViewData("res8img", res8img)
	c.Ctx.ViewData("listdir", listdir)
	c.Ctx.ViewData("listlable", listlable)
	//fmt.Printf("\n%v", pages)
	return mvc.View{Name: "ziyouge/Index/list.html"}
	//return mvc.Response{Text: fmt.Sprintf("%v", err)}
}

// 书籍列表页子页page
// Method:   GET
// resource:http://localhost:8080/List/1/0/0/0/0
func (c *MainController) GetBy(class int, full int, label int, size int, page int) mvc.Result {
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	var SortName string
	articlelist, pages, err := repositories.QueryList(repositories.DB, class, full, label, size, page)
	if err != nil {
		return mvc.Response{Code: 404}
	}
	config := Web
	for x, y := range config.Class {
		if x == class {
			SortName = y
			break
		}
	}
	if SortName == "" {
		return mvc.Response{Code: 503}
	}
	pages.ClassName = config.Class
	pages.FullName = []string{"全部", "连载", "全本"}
	pages.LabelName = []string{"全部", "收藏", "点击", "推荐", "入库"}
	pages.SizeName = []string{"全部", "10万以下", "10万-30万", "30万-50万", "50万-100万", "100万-200万", "200万以上"}
	pages.Class = class
	pages.Full = full
	pages.Label = label
	pages.Size = size
	for x, y := range articlelist {
		if y.Siteid >= 0 && y.Siteid <= 10 {
			articlelist[x].Class = config.Class[y.Sortid]
		} else {
			articlelist[x].Class = config.Class[0]
		}
		articlelist[x].LastTime = time.Unix(int64(y.Lastupdate), 0)
	}
	c.Ctx.ViewData("date", time.Now())
	c.Ctx.ViewData("WebInfo", config)
	c.Ctx.ViewData("ArticleList", articlelist)
	c.Ctx.ViewData("Pages", pages)
	c.Ctx.ViewData("SortName", SortName)
	return mvc.View{Name: "ziyouge/Index/list.html"}
	//return mvc.Response{Text: fmt.Sprintf("%v", err)}
}

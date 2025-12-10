package controllers

import (
	"fmt"
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

type MainMController struct {
	Ctx iris.Context
	//Cache iris.Cache
}
type ReadMController struct {
	Ctx iris.Context
}

//type ListMController struct {
//	Ctx iris.Context
//}

// Get serves
// Method:   GET
// Resource: http://localhost:8080
func (c *MainMController) Get() mvc.Result {
	//c.Ctx.AddHandler(cache.Handler(time.Second * 300))
	c.Ctx.AddHandler(cache.Handler(time.Hour * 6))
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	res5img, err := service.GetRecId(5)
	if err != nil {
		return mvc.Response{Code: 503}
	}
	res4img, err := service.GetRecId(4)
	if err != nil {
		return mvc.Response{Code: 503}
	}
	res15, err := service.GetLastUpdate(15)
	if err != nil {
		return mvc.Response{Code: 503}
	}
	c.Ctx.ViewData("res5img", res5img)
	c.Ctx.ViewData("res4img", res4img)
	c.Ctx.ViewData("res15", res15)
	c.Ctx.ViewData("WebInfo", Web)
	c.Ctx.ViewData("Date", time.Now().Unix())
	return mvc.View{Name: "ziyougem/Index/index.html"}
}

// 章节列表页阅读页
// Method:   GET
// resource:http://localhost:8080/book/1
func (c *MainMController) GetNovelBy(Id int) mvc.Result {
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
	chapter10, ListPages, err := service.GetChapterListPage(book.Chapter, 10, 1)
	if err != nil {
		return mvc.Response{Code: 404}
	}

	c.Ctx.ViewData("Intro", template.HTML(book.Article.Intro))
	c.Ctx.ViewData("Book", book)
	c.Ctx.ViewData("Chapter10", chapter10)
	c.Ctx.ViewData("Pages", ListPages)
	c.Ctx.ViewData("QQreminder", service.QQreminder(book.Article.Articleid, book.Article.Articlename))

	return mvc.View{Name: "ziyougem/Book/index.html"}

}

// 章节阅读页
// Method:   GET
// resource:http://localhost:8080/book/1/2
func (c *ReadMController) GetBy(Id int, Cid int) mvc.Result {
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	book, err := service.GetBookChapter(Id, Cid)
	if err == BState.NotFound {
		return mvc.Response{Code: 404}
	} else if err != nil {
		//ctx.Application().Logger().Infof("Begin request for path: %s", ctx.Path())
		c.Ctx.Application().Logger().Infof("chapter503,err:%v", err)
		return mvc.Response{Code: 503}
	} else {
		c.Ctx.ViewData("Content", template.HTML(book.CurrentChapter.CurrentChapterContent))
		c.Ctx.ViewData("Book", book)
		return mvc.View{Name: "ziyougem/Book/read.html"}
	}
}
func (c *MainMController) GetSearch() mvc.Result {
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	c.Ctx.ViewData("WebInfo", Web)
	return mvc.View{Name: "ziyougem/Index/search.html"}
}
func (c *MainMController) PostSoapi() mvc.Result {
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
//func (c *MainMController) GetHello() interface{} {
//	return map[string]string{"message": "Hello Iris!"}

//}

//在控制器适应主应用程序之前调用一次BeforeActivation
//当然在服务器运行之前。
//在版本9之后，您还可以为特定控制器的方法添加自定义路由。
//在这里您可以注册自定义方法的处理程序
//使用带有`ca.Router`的标准路由器做一些你可以做的事情，没有mvc，
//并添加将绑定到控制器的字段或方法函数的输入参数的依赖项。
func (c *MainMController) BeforeActivation(b mvc.BeforeActivation) {
	//b.Handle("GET", "/", "CustomGet", cache.Handler(time.Second*3000))
	b.Handle("GET", "/list/{class:int}", "CustomList1", cache.Handler(time.Hour*6))
	b.Handle("GET", "/list/{class:int}/{page:int}", "CustomList2", cache.Handler(time.Hour*6))
	b.Handle("GET", "/novel/{id:int}/{page:int}", "NovelSon")
	b.Handle("GET", "/bookcase", "Bookcase1")
	b.Handle("GET", "/bookcase/{page:int}", "Bookcase2")
	//甚至添加基于此控制器路由器的全局中间件，
	//在这个例子中是根“/”：
	// b.Router（）。使用（myMiddleware）
	//b.Router().Get("/book{id:int}",customFunc)
	//这个功能简直无敌
}
func (c *MainMController) Bookcase1() mvc.Result {
	return c.customBookcase(0, 100, 0)
}
func (c *MainMController) Bookcase2(page int) mvc.Result {
	return c.customBookcase(0, 100, page)
}
func (c *MainMController) NovelSon(id, page int) mvc.Result {
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	book, err := service.GetBook(id)
	if err == BState.NotFound {
		return mvc.Response{Code: 404}
	} else if err != nil {
		return mvc.Response{Code: 503}
	}
	var (
		pages int = 10
		json  BookListPageJson
	)
	if page <= 0 {
		page = 1
	}
	chapterJson, ListPages, err := service.GetChapterListPage(book.Chapter, pages, page)
	if err != nil {
		return mvc.Response{Code: 404}
	}
	json.List = chapterJson
	json.Page = fmt.Sprintf("%d/%d", page, ListPages.TotalPages)
	return mvc.Response{Object: json}
}

func (c *MainMController) CustomList1(class int) mvc.Result {
	return c.CustomList(class, 0)
}
func (c *MainMController) CustomList2(class int, page int) mvc.Result {
	return c.CustomList(class, page)
}

/////////////////
// 书籍列表页
// Method:   GET
// resource:http://localhost:8080/List/1/0/0/0
func (c *MainMController) CustomList(class int, page int) mvc.Result {
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	var SortName string
	c.Ctx.ViewData("Date", time.Now().Unix())
	articlelist, pages, err := repositories.QueryList(repositories.DB, class, 0, 0, 0, page)
	if err != nil {
		return mvc.Response{Code: 404}
	}

	//	Web := ServerConfig()

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
	pageSelect := []string{}
	for i := 0; i <= int(pages.TotalPages); i++ {
		if i == 0 {
			pageSelect = append(pageSelect, "")
		} else {
			sel := fmt.Sprintf("%d/%d", i, pages.TotalPages)
			pageSelect = append(pageSelect, sel)
		}

	}
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
	if page-1 <= 0 {
		pages.PrevPage = 0
	} else {
		pages.PrevPage = int64(page - 1)
	}
	if page+1 > int(pages.TotalPages) {
		pages.NextPage = 0
	} else {
		pages.NextPage = int64(page + 1)
	}
	c.Ctx.ViewData("date", time.Now().Unix())
	c.Ctx.ViewData("WebInfo", Web)
	c.Ctx.ViewData("ArticleList", articlelist)
	c.Ctx.ViewData("Pages", pages)
	c.Ctx.ViewData("PageSelect", pageSelect)
	c.Ctx.ViewData("SortName", SortName)
	return mvc.View{Name: "ziyougem/Index/list.html"}
	//return mvc.Response{Text: fmt.Sprintf("%v", err)}
}

// 书籍列表页子页page
// Method:   GET
// resource:http://localhost:8080/List/1/0/0/0/0
func (c *MainMController) GetBy(class int, full int, label int, size int, page int) mvc.Result {
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	var SortName string
	articlelist, pages, err := repositories.QueryList(repositories.DB, class, full, label, size, page)
	if err != nil {
		return mvc.Response{Code: 404}
	}
	config := ServerConfig()
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
	return mvc.View{Name: "ziyougem/Index/list.html"}
	//return mvc.Response{Text: fmt.Sprintf("%v", err)}
}

package controllers

import (
	"html"
	"regexp"

	"github.com/kataras/iris/v12"

	//	"html/template"
	"math"
	. "zyg/datamodels"

	"github.com/dchest/captcha"

	//	"github.com/kataras/iris"
	"time"
	"zyg/repositories"
	"zyg/service"

	valid "github.com/asaskevich/govalidator"
	"github.com/gorilla/securecookie"
	"github.com/kataras/iris/v12/mvc"
)

var (
	// AES only supports key sizes of 16, 24 or 32 bytes.
	// You either need to provide exactly that amount or you derive the key from what you type in.
	hashKey  = []byte("the-min-and-secret-fash-key-here")
	blockKey = []byte("lot-secret-of-characters-min-too")
	sc       = securecookie.New(hashKey, blockKey)
)

type BookCaseController struct {
	Ctx iris.Context
}

func (c *MainController) PostCollect() mvc.Result {
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	if !c.Ctx.IsAjax() {
		return mvc.Response{Code: 404}
	}
	var Json JsonInsetMess
	err := c.Ctx.ReadJSON(&Json)
	if err != nil {
		return mvc.Response{Object: service.ReAjax("出错啦!", 0, "")}
	}
	result, err := valid.ValidateStruct(Json)
	if err != nil {
		return mvc.Response{Object: service.ReAjax("请检查信息是否有误", 0, "")}
	}
	if !result {
		return mvc.Response{Object: service.ReAjax("请检查信息是否有误", 0, "")}
	}
	id := c.Ctx.GetCookie("ratel_u", iris.CookieDecode(sc.Decode))
	name := c.Ctx.GetCookie("ratel_n", iris.CookieDecode(sc.Decode))
	userinfo, err := service.IsLogin(id, name)
	if err == UErr.NotLogin {
		return mvc.Response{Object: service.ReAjax("请登录!", 0, "")}
	} else if err != nil {
		return mvc.Response{Object: service.ReAjax("出错啦!", 0, "")}
	}
	count, err := service.InputMess(userinfo, Json)
	if err != nil {
		return mvc.Response{Object: service.ReAjax("出错啦!", 0, "")}
	} else if count > 0 {
		return mvc.Response{Object: service.ReAjax("成功!", 1, "")}
	}
	return mvc.Response{Object: service.ReAjax("未知错误", 0, "")}
}
func (c *MainController) PostDelmark() mvc.Result {
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	if !c.Ctx.IsAjax() {
		return mvc.Response{Code: 404}
	}
	var Json JsonDel
	err := c.Ctx.ReadJSON(&Json)
	if err != nil {
		return mvc.Response{Object: service.ReAjax("出错啦!", 0, "")}
	}
	if len(Json.Bookid) <= 0 || len(Json.Bookid) > 15 {
		return mvc.Response{Object: service.ReAjax("请选择书籍!", 0, "")}
	}
	id := c.Ctx.GetCookie("ratel_u", iris.CookieDecode(sc.Decode))
	name := c.Ctx.GetCookie("ratel_n", iris.CookieDecode(sc.Decode))
	userinfo, err := service.IsLogin(id, name)
	if err == UErr.NotLogin {
		return mvc.Response{Object: service.ReAjax("请登录!", 0, "")}
	} else if err != nil {
		return mvc.Response{Object: service.ReAjax("出错啦!", 0, "")}
	}
	count, err := service.Delmark(userinfo.Userid, Json.Bookid)
	if err != nil {
		return mvc.Response{Object: service.ReAjax("出错啦!", 0, "")}
	} else if count > 0 {
		return mvc.Response{Object: service.ReAjax("成功!", 1, "")}
	}
	return mvc.Response{Object: service.ReAjax("未知错误", 0, "")}
}
func (c *MainController) PostMovemark() mvc.Result {
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	if !c.Ctx.IsAjax() {
		return mvc.Response{Code: 404}
	}
	var Json JsonMove
	err := c.Ctx.ReadJSON(&Json)
	if err != nil {
		return mvc.Response{Object: service.ReAjax("出错啦!", 0, "")}
	}
	if Json.Action < 0 && Json.Action > 10 {
		return mvc.Response{Object: service.ReAjax("请选择书架!", 0, "")}
	}
	if len(Json.Bookid) < 1 && len(Json.Bookid) > 15 {
		return mvc.Response{Object: service.ReAjax("请选择书籍!", 0, "")}
	}
	id := c.Ctx.GetCookie("ratel_u", iris.CookieDecode(sc.Decode))
	name := c.Ctx.GetCookie("ratel_n", iris.CookieDecode(sc.Decode))
	userinfo, err := service.IsLogin(id, name)
	if err == UErr.NotLogin {
		return mvc.Response{Object: service.ReAjax("请登录!", 0, "")}
	} else if err != nil {
		return mvc.Response{Object: service.ReAjax("出错啦!", 0, "")}
	}
	count, err := service.Movemark(userinfo.Userid, Json.Action, Json.Bookid)
	if err != nil {
		return mvc.Response{Object: service.ReAjax("出错啦!", 0, "")}
	} else if count > 0 {
		return mvc.Response{Object: service.ReAjax("成功!", 1, "")}
	}
	return mvc.Response{Object: service.ReAjax("未知错误", 0, "")}

}
func (c *MainController) PostAddmark() mvc.Result {
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	if !c.Ctx.IsAjax() {
		return mvc.Response{Code: 404}
	}
	var Json JsonAdd
	err := c.Ctx.ReadJSON(&Json)
	if err != nil {
		return mvc.Response{Object: service.ReAjax("出错啦!", 0, "")}
	}
	result, err := valid.ValidateStruct(Json)
	if err != nil {
		return mvc.Response{Object: service.ReAjax("添加失败", 0, "")}
	}
	if !result {
		return mvc.Response{Object: service.ReAjax("添加失败", 0, "")}
	}
	id := c.Ctx.GetCookie("ratel_u", iris.CookieDecode(sc.Decode))
	name := c.Ctx.GetCookie("ratel_n", iris.CookieDecode(sc.Decode))
	userinfo, err := service.IsLogin(id, name)

	if err == UErr.NotLogin {
		return mvc.Response{Object: service.ReAjax("请登录!", 0, "")}
	} else if err != nil {
		return mvc.Response{Object: service.ReAjax("出错啦!", 0, "")}
	}

	isok, err := service.WriteBookMark(userinfo, Json.Articleid, Json.Chapterid)
	if err == UErr.BookMarkSame {
		return mvc.Response{Object: service.ReAjax("加签成功", 1, "")}
	} else if err == UErr.BookcaseFull {
		return mvc.Response{Object: service.ReAjax("书架已满", 0, "")}
	} else if err != nil {
		return mvc.Response{Object: service.ReAjax("出错啦!", 0, "")}
	}
	if isok > 0 {
		return mvc.Response{Object: service.ReAjax("加签成功", 1, "")}
	} else {
		return mvc.Response{Object: service.ReAjax("出错啦!", 0, "")}
	}
}
func (c *MainController) PostRegister() mvc.Result {
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	var userreg UserRegister
	userreg.Username = c.Ctx.FormValue("username")
	userreg.Email = c.Ctx.FormValue("email")
	userreg.Password = c.Ctx.FormValue("password")
	userreg.Repassword = c.Ctx.FormValue("repassword")
	userreg.Sex = c.Ctx.FormValue("sex")
	userreg.Code = c.Ctx.FormValue("code")
	userreg.Vid = c.Ctx.FormValue("vid")
	result, err := valid.ValidateStruct(userreg)
	if err != nil {
		return mvc.Response{Object: service.ReAjax("信息填写有误", 0, "")}
	}
	if !result {
		return mvc.Response{Object: service.ReAjax("信息填写有误", 0, "")}
	}
	var re = regexp.MustCompile(`\/|\~|\!|\@|\#|\\$|\%|\^|\&|\*|\(|\)|\+|\{|\}|\:|\<|\>|\?|\[|\]|\,|\.|\/|\;|\'|\x60|\-|\=|\\\|\|`)
	if len(re.FindStringIndex(userreg.Username)) > 0 {
		return mvc.Response{Object: service.ReAjax("用户名存在特殊字符", 0, "")}
	}
	userreg.Username = html.EscapeString(userreg.Username)
	if !captcha.VerifyString(userreg.Vid, userreg.Code) {
		return mvc.Response{Object: service.ReAjax("验证码错误", 0, "")}
	}
	userreg.Password = service.GetMd5(userreg.Password)
	userreg.Repassword = service.GetMd5(userreg.Repassword)
	if userreg.Password != userreg.Repassword {
		return mvc.Response{Object: service.ReAjax("两次密码不一样", 0, "")}
	}
	_, err = service.WriteUser(userreg)
	if err != nil {
		return mvc.Response{Object: service.ReAjax("未知错误", 1, "")}
	}
	return mvc.Response{Object: service.ReAjax("成功啦", 1, "/login")}
}
func (c *MainController) PostLogin() mvc.Result {
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	var (
		username = c.Ctx.FormValue("username")
		password = c.Ctx.FormValue("password")
		code     = c.Ctx.FormValue("code")
		auto     = c.Ctx.FormValue("auto")
		vid      = c.Ctx.FormValue("vid")
	)
	if !captcha.VerifyString(vid, code) {
		return mvc.Response{Object: service.ReAjax("验证码错误", 0, "")}
	}
	var re = regexp.MustCompile(`\/|\~|\!|\@|\#|\\$|\%|\^|\&|\*|\(|\)|\+|\{|\}|\:|\<|\>|\?|\[|\]|\,|\.|\/|\;|\'|\x60|\-|\=|\\\|\|`)
	if len(re.FindStringIndex(username)) > 0 {
		return mvc.Response{Object: service.ReAjax("用户名存在特殊字符", 0, "")}
	}
	username = html.EscapeString(username)
	pass := service.GetMd5(password)
	userinfo, err := repositories.QueryUserByName(repositories.DB, username, pass)
	if err != nil {
		return mvc.Response{Object: service.ReAjax("账号或密码错误", 0, "")}
	}
	date := time.Hour * 24
	if auto == "yes" {
		date = date * 30
	}
	c.Ctx.SetCookieKV("ratel_n", userinfo.Username, iris.CookieEncode(sc.Encode), iris.CookieExpires(date))
	c.Ctx.SetCookieKV("ratel_u", userinfo.Userid, iris.CookieEncode(sc.Encode), iris.CookieExpires(date))
	repositories.SessSet(userinfo, time.Hour*24)
	//c.Ctx.SetCookieKV("n", username, iris.CookieEncode(sc.Encode), iris.CookieExpires(time.Hour*24*30))
	return mvc.Response{Object: service.ReAjax("登录成功", 1, "/bookcase")}
	//返回格式{"info":"\u9a8c\u8bc1\u7801\u9519\u8bef\uff01","status":0,"url":""}
}
func (c *MainController) GetOutlogin() mvc.Result {
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	c.Ctx.RemoveCookie("ratel_u")
	c.Ctx.RemoveCookie("ratel_n")
	return mvc.Response{Path: "/login"}
}
func (c *MainController) GetLogin() mvc.Result {
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	id := c.Ctx.GetCookie("ratel_u", iris.CookieDecode(sc.Decode))
	name := c.Ctx.GetCookie("ratel_n", iris.CookieDecode(sc.Decode))
	_, err := service.IsLoginGet(id, name)
	if err == UErr.ServiceUnavailable {
		return mvc.Response{Code: 503}
	} else if err == nil { //id 获取到的name与session里的name 不符合
		return mvc.Response{Path: "/bookcase"}
	}
	captchaId := struct {
		CaptchaId string
	}{
		captcha.NewLen(4),
	}
	c.Ctx.ViewData("Captcha", captchaId)
	c.Ctx.ViewData("WebInfo", Web)
	return mvc.View{Name: "ziyouge/User/login.html"}
}
func (c *MainController) GetRegister() mvc.Result {
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	id := c.Ctx.GetCookie("ratel_u", iris.CookieDecode(sc.Decode))
	name := c.Ctx.GetCookie("ratel_n", iris.CookieDecode(sc.Decode))
	_, err := service.IsLoginGet(id, name)
	if err == UErr.ServiceUnavailable {
		return mvc.Response{Code: 503}
	} else if err == nil { //id 获取到的name与session里的name 不符合
		return mvc.Response{Path: "/bookcase"}
	}
	captchaId := struct {
		CaptchaId string
	}{
		captcha.NewLen(4),
	}
	c.Ctx.ViewData("Captcha", captchaId)
	c.Ctx.ViewData("WebInfo", Web)
	return mvc.View{Name: "ziyouge/User/register.html"}
}

//func (c *MainController) GetBookcase() mvc.Result {
//	var (
//		pageInfo ListPages
//		pno      int = 15
//		class    int = 0
//		page     int = 0
//		sort     int = 0
//	)
//	id := c.Ctx.GetCookie("ratel_u", iris.CookieDecode(sc.Decode))
//	name := c.Ctx.GetCookie("ratel_n", iris.CookieDecode(sc.Decode))
//	userinfo, err := service.IsLoginGet(id, name)
//	if err == UErr.ServiceUnavailable {
//		return mvc.Response{Code: 503}
//	} else if err == UErr.NotLogin { //id 获取到的name与session里的name 不符合
//		return mvc.Response{Path: "/login"}
//	}

//	///////
//	c.Ctx.ViewData("WebInfo", ServerConfig())
//	c.Ctx.ViewData("UserInfo", userinfo)
//	n, err := repositories.CountBookcaseById(repositories.DB, userinfo.Userid, 0)
//	if err != nil {
//		return mvc.Response{Code: 503}
//	}
//	c.Ctx.ViewData("BookcaseCount", n)
//	pageInfo.TotalPages = int64(math.Ceil(float64(n) / float64(pno)))
//	//	if page <= 0 {
//	//		page = 1
//	//	}
//	pageInfo.CurrPage = int64(page)
//	pageInfo.Pages = service.GetUpDownPages(int(pageInfo.CurrPage), int(pageInfo.TotalPages), pno)
//	pageInfo.Sort = sort
//	pageInfo.Class = class

//	result, err := repositories.QueryBookcaseByIdSortClass(repositories.DB, userinfo.Userid, 0, 0, 0, pno)
//	if err != nil {
//		return mvc.Response{Code: 503}
//	}
//	for x, y := range result {
//		result[x].LastTime = time.Unix(int64(y.Lastvisit), 0)
//	}
//	c.Ctx.ViewData("BookcaseList", result)
//	shelf := []BookShelf{{Name: "默认书架", Id: 0}, {Name: "书架一", Id: 1}, {Name: "书架二", Id: 2}, {Name: "书架三", Id: 3}, {Name: "书架四", Id: 4}}
//	c.Ctx.ViewData("BookShelf", shelf)
//	c.Ctx.ViewData("PageInfo", pageInfo)
//	//return fmt.Sprintf("找到啦%+v", userinfo)
//	return mvc.View{Name: "ziyouge/User/bookcase.html"}
//}

func (c *MainController) customBookcase(sort int, class int, page int) mvc.Result {
	if c.Ctx.ClientSupportsGzip() {
		c.Ctx.Gzip(true)
	}
	var (
		pageInfo ListPages
	)
	if class > Scon.BookCaseClassMax {
		return mvc.Response{Code: 503}
	}
	id := c.Ctx.GetCookie("ratel_u", iris.CookieDecode(sc.Decode))
	name := c.Ctx.GetCookie("ratel_n", iris.CookieDecode(sc.Decode))
	userinfo, err := service.IsLoginGet(id, name)
	if err == UErr.ServiceUnavailable {
		return mvc.Response{Code: 503}
	} else if err == UErr.NotLogin { //id 获取到的name与session里的name 不符合
		return mvc.Response{Path: "/login"}
	}

	n, err := repositories.CountBookcaseById(repositories.DB, userinfo.Userid, class)
	if err != nil {
		return mvc.Response{Code: 503}
	}

	pageInfo.TotalPages = int64(math.Ceil(float64(n) / float64(Scon.BookCasePageItemMax)))
	if page <= 0 {
		page = 1
	}
	pageInfo.CurrPage = int64(page)
	pageInfo.Pages = service.GetUpDownPages(int(pageInfo.CurrPage), int(pageInfo.TotalPages), Scon.BookCasePageItemMax)
	pageInfo.Sort = sort
	pageInfo.Class = class
	result, err := repositories.QueryBookcaseByIdSortClass(repositories.DB, userinfo.Userid, sort, class, page-1, Scon.BookCasePageItemMax)
	if err != nil {
		return mvc.Response{Code: 503}
	}
	for x, y := range result {
		result[x].LastTime = time.Unix(int64(y.Lastvisit), 0)
		result[x].ImgUrl = service.GetImgUrl(y.Articleid, y.Imgflag, Web.UrlImg)
	}

	shelf := []BookShelf{{Name: "默认书架", Id: 0},
		{Name: "一号书架", Id: 1},
		{Name: "二号书架", Id: 2},
		{Name: "三号书架", Id: 3},
		{Name: "四号书架", Id: 4},
		{Name: "五号书架", Id: 5},
		{Name: "六号书架", Id: 6},
		{Name: "七号书架", Id: 7},
		{Name: "八号书架", Id: 8},
		{Name: "九号书架", Id: 9},
		{Name: "十号书架", Id: 10}}
	c.Ctx.ViewData("BookcaseList", result)
	c.Ctx.ViewData("WebInfo", Web)
	c.Ctx.ViewData("UserInfo", userinfo)
	c.Ctx.ViewData("BookcaseCount", n)
	c.Ctx.ViewData("BookShelf", shelf)
	c.Ctx.ViewData("Pages", pageInfo)
	//return fmt.Sprintf("找到啦%+v", userinfo)
	return mvc.View{Name: "ziyouge/User/bookcase.html"}
}

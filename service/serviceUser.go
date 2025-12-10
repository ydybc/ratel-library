package service

import (
	//	"fmt"
	"html"

	//	"errors"
	//	"errors"
	"database/sql"
	//"fmt"
	"strconv"
	"time"
	. "zyg/datamodels"
	"zyg/repositories"

	valid "github.com/asaskevich/govalidator"
	"github.com/go-redis/redis"
)

func InputMess(user UserInfo, Json JsonInsetMess) (int64, error) {
	var mess BookErrMessage
	mess.Title = html.EscapeString(Json.Title)
	mess.Content = html.EscapeString(Json.Content)
	mess.Bookid = Json.Articleid
	mess.Attachsig = 0
	mess.Chapterid = Json.Chapterid
	mess.Postdate = time.Now().Unix()
	id, err := repositories.InputErrMessage(repositories.DB, user, mess)
	if err != nil {
		return id, err
	}
	return id, err
}
func Delmark(id string, aid []int) (int64, error) {
	count, err := repositories.DeleteBookcaseByIdAids(repositories.DB, id, aid)
	return count, err
}
func Movemark(id string, ation int, aid []int) (int64, error) {
	count, err := repositories.UpdataBookcaseClassByIdAids(repositories.DB, id, ation, aid)
	return count, err
}
func WriteUser(reg UserRegister) (int64, error) {
	index, err := repositories.InputUser(repositories.DB, reg)
	return index, err
}
func LookupUserOnly(reg UserRegister) (bool, string, error) {
	isEmailok, err := repositories.CountUserByEmail(repositories.DB, reg.Email)
	if err != nil {
		return true, "未知错误", err
	}
	if isEmailok == 0 {
		isUnameok, err := repositories.CountUserByUname(repositories.DB, reg.Username)
		if err != nil {
			return true, "未知错误", err
		}
		if isUnameok == 0 {
			return false, "", nil
		}
		return true, "用户名已经被注册", nil
	}
	return true, "email已经被注册", nil

}
func WriteBookMark(user UserInfo, aid int, cid int) (int64, error) {
	var (
		p    bool = true
		isok int64
	)
	if cid == 0 {
		p = false
	}
	bookcaseallcount, err := repositories.CountBookcaseAllById(repositories.DB, user.Userid)
	if err != nil {
		return isok, err
	}
	if bookcaseallcount > Scon.BookcassAllCount {
		return isok, UErr.BookcaseFull
	}
	bookcaseinfo, err := repositories.QueryBookcaseByIdAid(repositories.DB, user.Userid, aid)
	if err == sql.ErrNoRows {

	} else if err != nil {
		return isok, err

	}
	//fmt.Printf("\n%d,%d,%v,%v", aid, cid, user, bookcaseinfo)
	if cid != 0 && bookcaseinfo.Chapterid == cid {
		return isok, UErr.BookMarkSame
	}
	bookcase, err := GetBookMarkInfo(user, aid, cid)
	if err != nil {
		return isok, err
	}
	//fmt.Printf("\nbookcase%v", bookcase)
	switch bookcaseinfo.Caseid {
	case 0:
		isok, err = repositories.InputBookcase(repositories.DB, bookcase, p)
	default:
		isok, err = repositories.UpdataBookcase(repositories.DB, bookcase, bookcaseinfo.Caseid)
	}
	return isok, err

}
func GetBookMarkInfo(user UserInfo, aid int, cid int) (BookCaseData, error) {
	var bookcase BookCaseData
	switch cid {
	case 0:
		article, err := repositories.QueryArticleById(repositories.DB, aid)
		if err != nil {
			return bookcase, err
		}
		bookcase.Articleid = article.Articleid
		bookcase.Articlename = article.Articlename
		bookcase.Chapterorder = 0
	default:
		chapter, err := repositories.QueryChapterByCid(repositories.DB, cid)
		if err != nil {
			return bookcase, err
		}
		bookcase.Articleid = chapter.Articleid
		bookcase.Articlename = chapter.Articlename
		bookcase.Chapterorder = chapter.Chapterorder
		bookcase.Chapterid = chapter.Chapterid
		bookcase.Chaptername = chapter.Chaptername
	}
	uid, err := strconv.Atoi(user.Userid)
	if err != nil {
		return bookcase, err
	}

	bookcase.Userid = uid
	bookcase.Joindate = int(time.Now().Unix())
	bookcase.Lastvisit = int(time.Now().Unix())
	bookcase.Username = user.Username
	return bookcase, nil
}
func IsLogin(id string, name string) (UserInfo, error) {
	var userinfo UserInfo
	if id == "" && valid.IsInt(id) {
		return userinfo, UErr.NotLogin
	}
	userinfo, err := repositories.SessGet(id) //获取redis中sess
	if err == redis.Nil {                     //获取数据为空
		userinfo, err = Autologin(id, name) //进行数据库比对
		if err != nil {
			return userinfo, UErr.NotLogin
		}
		if userinfo.Userid == id && userinfo.Username == name { //比对成功返回续租sess err
			repositories.SessSet(userinfo, CacheTime.RedisSessExpiry)
			return userinfo, nil
		}
	} else if err != nil { //出现错误
		return userinfo, UErr.ServiceUnavailable
	}
	if name != userinfo.Username { //现有用户名和sess内用户名不匹配
		return userinfo, UErr.NotLogin
	}
	return userinfo, nil //sess匹配返回
}
func IsLoginGet(id string, name string) (UserInfo, error) {
	var userinfo UserInfo
	if id != "" && valid.IsInt(id) {
		userinfo, err := repositories.SessGet(id) ////获取redis中sess
		if err == redis.Nil {                     //获取数据为空
			userinfo, err = Autologin(id, name) //进行数据库比对
			if err != nil {                     //数据查询不能出错
				return userinfo, UErr.NotLogin //返回查询错误
			}
			if userinfo.Userid == id && userinfo.Username == name { //比对查询结果集
				repositories.SessSet(userinfo, Scon.RedisSessExpiry) //续租sess
				return userinfo, nil
			}
		} else if err != nil {
			return userinfo, UErr.ServiceUnavailable
		}
		if name != userinfo.Username {
			return userinfo, UErr.NotLogin
		}
		return userinfo, nil //sess匹配返回
	} else {
		return userinfo, UErr.NotLogin
	}

}
func Autologin(id string, name string) (UserInfo, error) {
	userinfo, err := repositories.QueryUserByIdName(repositories.DB, id, name)
	if err != nil {
		return userinfo, err
	}
	return userinfo, err
}

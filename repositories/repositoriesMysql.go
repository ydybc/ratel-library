package repositories

//package repositories

import (
	//	"database/sql"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	//	"errors"
	//	"fmt"

	"github.com/jmoiron/sqlx"

	. "zyg/datamodels"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sqlx.DB

func init() {
	var err error
	Dsn := Scon.MysqlLink
	DB, err = sqlx.Open("mysql", Dsn)
	if err != nil {
		panic(err)
		return
	}
	DB.SetConnMaxLifetime(time.Second * 599)
	DB.SetMaxIdleConns(100)
	//defer DBW.Db.Close()
}
func CountUserByEmail(DB *sqlx.DB, email string) (int64, error) {
	var count int64
	row := DB.QueryRowx("select count(uid) from jieqi_system_users where email = ?", email)
	err := row.Scan(&count)
	if err != nil {
		return count, err
	}
	return count, err
}
func CountUserByUname(DB *sqlx.DB, uname string) (int64, error) {
	var count int64
	row := DB.QueryRowx("select count(uid) from jieqi_system_users where uname = ?", uname)
	err := row.Scan(&count)
	if err != nil {
		return count, err
	}
	return count, err
}
func QueryUserByName(DB *sqlx.DB, name string, pass string) (UserInfo, error) {
	users := UserData{}
	userinfo := UserInfo{}
	row := DB.QueryRowx("select * from jieqi_system_users where uname = ? and pass = ?", name, pass)
	if err := row.Err(); err != nil {
		return userinfo, err
	}
	err := row.StructScan(&users)
	if err != nil {
		return userinfo, err
	}
	userinfo.Userid = strconv.Itoa(users.Uid)
	userinfo.Username = users.Name
	return userinfo, err
}
func QueryUserByIdName(DB *sqlx.DB, id string, name string) (UserInfo, error) {
	users := UserData{}
	userinfo := UserInfo{}
	row := DB.QueryRowx("select * from jieqi_system_users where uid = ? and uname = ?", id, name)
	if err := row.Err(); err != nil {
		return userinfo, err
	}
	err := row.StructScan(&users)
	if err != nil {
		return userinfo, err
	}
	userinfo.Userid = strconv.Itoa(users.Uid)
	userinfo.Username = users.Name
	return userinfo, err
}

func QueryBookcaseByIdSortClass(DB *sqlx.DB, id string, sort int, class int, page int, pages int) ([]BookCaseResultData, error) {
	var (
		order, where string
	)
	resultSet := make([]BookCaseResultData, 0)
	result := BookCaseResultData{}
	position := page * pages
	field := "jieqi_article_bookcase.articleid," +
		"jieqi_article_bookcase.articlename," +
		"jieqi_article_bookcase.chapterid," +
		"jieqi_article_bookcase.chaptername," +
		"jieqi_article_bookcase.lastvisit," +
		"jieqi_article_bookcase.chapterorder," +
		"jieqi_article_bookcase.classid," +
		"jieqi_article_article.imgflag as imgflag," +
		"jieqi_article_article.author as author," +
		"jieqi_article_article.fullflag as fullflag," +
		"jieqi_article_article.chapters as bookchapterorder," +
		"jieqi_article_article.lastchapterid as booklastchapterid," +
		"jieqi_article_article.lastchapter as booklastchapter"
	from := "jieqi_article_article,jieqi_article_bookcase"
	if class == 100 {
		where = fmt.Sprintf("jieqi_article_bookcase.articleid=jieqi_article_article.articleid and jieqi_article_bookcase.userid=%s", id)
	} else {
		where = fmt.Sprintf("jieqi_article_bookcase.articleid=jieqi_article_article.articleid and jieqi_article_bookcase.userid=%s and jieqi_article_bookcase.classid = %d", id, class)
	}

	if sort == 0 {
		order = "by lastvisit desc"
	} else {
		order = "by lastupdate desc"
	}
	limit := fmt.Sprintf("%d,%d", position, pages)
	sql := fmt.Sprintf("select %s from %s where %s order %s limit %s", field, from, where, order, limit)
	rows, err := DB.Queryx(sql)
	defer rows.Close()
	if err != nil {
		return resultSet, err
	}
	for rows.Next() {
		err = rows.StructScan(&result)
		if err != nil {
			return resultSet, err
		}
		resultSet = append(resultSet, result)
	}

	return resultSet, err

}
func CountBookcaseById(DB *sqlx.DB, id string, class int) (int, error) {
	var (
		count       int
		from, where string
	)
	from = "jieqi_article_article,jieqi_article_bookcase"
	if class > 10 {
		where = fmt.Sprintf("jieqi_article_bookcase.articleid=jieqi_article_article.articleid and jieqi_article_bookcase.userid=%s", id)
	} else {
		where = fmt.Sprintf("jieqi_article_bookcase.articleid=jieqi_article_article.articleid and jieqi_article_bookcase.userid=%s and jieqi_article_bookcase.classid=%d", id, class)
	}
	sql := fmt.Sprintf("select count(*) from %s where %s", from, where)
	row := DB.QueryRowx(sql)
	err := row.Scan(&count)
	if err != nil {
		return count, err
	}
	return count, err
}
func CountBookcaseAllById(DB *sqlx.DB, id string) (int, error) {
	var (
		count int
	)
	sql := "select count(*) from jieqi_article_bookcase where userid = ?"
	row := DB.QueryRowx(sql, id)
	err := row.Scan(&count)
	if err != nil {
		return count, err
	}
	return count, err
}
func QueryBookcaseByIdAid(DB *sqlx.DB, id string, aid int) (BookCaseData, error) {
	var (
		bookcase BookCaseData
		row      *sqlx.Row
	)
	row = DB.QueryRowx("select * from jieqi_article_bookcase where userid = ? and articleid = ?", id, aid)
	if err := row.Err(); err != nil {
		return bookcase, err
	}
	err := row.StructScan(&bookcase)
	if err != nil {
		return bookcase, err
	}
	return bookcase, nil
}
func DeleteBookcaseByIdAids(DB *sqlx.DB, id string, aid []int) (int64, error) {
	var aidstring string
	for x, y := range aid {
		if x == 0 {
			aidstring = fmt.Sprintf("%d", y)
		} else {
			aidstring = fmt.Sprintf("%s,%d", aidstring, y)
		}
	}
	sql := fmt.Sprintf("delete from jieqi_article_bookcase where articleid in (%s) and userid = ?", aidstring)
	result := DB.MustExec(sql, id)
	count, err := result.RowsAffected()
	if err != nil {
		return count, err
	}
	return count, nil

}
func UpdataBookcaseClassByIdAids(DB *sqlx.DB, id string, class int, aid []int) (int64, error) {
	var aidstring string
	for x, y := range aid {
		if x == 0 {
			aidstring = fmt.Sprintf("%d", y)
		} else {
			aidstring = fmt.Sprintf("%s,%d", aidstring, y)
		}
	}

	sql := fmt.Sprintf("update jieqi_article_bookcase set classid = ? where userid =? and articleid in (%s)", aidstring)
	result := DB.MustExec(sql, class, id)
	count, err := result.RowsAffected()
	if err != nil {
		return count, err
	}
	return count, nil
}
func InputErrMessage(DB *sqlx.DB, user UserInfo, mess BookErrMessage) (int64, error) {
	sql := "insert into jieqi_system_message (postdate,fromid,fromname,toid,toname,title,content)values(?,?,?,?,?,?,?)"
	//date := time.Now().Unix()
	con := fmt.Sprintf("Id:%d,Cid,%d\n%s", mess.Bookid, mess.Chapterid, mess.Content)
	result := DB.MustExec(sql, mess.Postdate, user.Userid, user.Username, 1, "admin", mess.Title, con)
	messageid, err := result.LastInsertId()
	if err != nil {
		return messageid, err
	}
	return messageid, nil
}
func InputUser(DB *sqlx.DB, user UserRegister) (int64, error) {
	date := time.Now().Unix()
	sql := "insert into jieqi_system_users (uname,name,email,sex,pass,groupid,regdate,lastlogin,sign,intro,setting,badges)values(?,?,?,?,?,?,?,?,?,?,?,?)"
	result := DB.MustExec(sql, user.Username, user.Username, user.Email, user.Sex, user.Password, 3, date, date, "", "", "", "")
	messageid, err := result.LastInsertId()
	if err != nil {
		return messageid, err
	}
	return messageid, nil

}
func InputBookcase(DB *sqlx.DB, bookcase BookCaseData, p bool) (int64, error) {
	var (
		column, value string
		caseid        int64
	)
	if p {
		column = "articleid,articlename,classid,userid,username,chapterid,chaptername,chapterorder,joindate,lastvisit"
		value = ":articleid,:articlename,:classid,:userid,:username,:chapterid,:chaptername,:chapterorder,:joindate,:lastvisit"
	} else {
		column = "articleid,articlename,classid,userid,username,joindate,lastvisit"
		value = ":articleid,:articlename,:classid,:userid,:username,:joindate,:lastvisit"
	}
	sql := fmt.Sprintf("insert into jieqi_article_bookcase (%s) values (%s)", column, value)
	result, err := DB.NamedExec(sql, &bookcase)
	if err != nil {
		return caseid, err
	}
	caseid, err = result.LastInsertId()
	if err != nil {
		return caseid, err
	}
	return caseid, nil
}
func UpdataBookcase(DB *sqlx.DB, bookcase BookCaseData, caseid int) (int64, error) {
	var count int64
	result := DB.MustExec("update  jieqi_article_bookcase set chapterid=?,chaptername=?,chapterorder=?,lastvisit=? where articleid=? and caseid = ?",
		bookcase.Chapterid, bookcase.Chaptername, bookcase.Chapterorder, bookcase.Lastvisit, bookcase.Articleid, caseid)
	count, err := result.RowsAffected()
	if err != nil {
		return count, err
	}
	return count, nil
}
func QueryArticleById(DB *sqlx.DB, Id int) (ArticleStruct, error) {
	article := ArticleStruct{}
	row := DB.QueryRowx("select * from jieqi_article_article where articleid = ?", Id)

	if err := row.Err(); err != nil {
		return article, err
	}
	err := row.StructScan(&article)
	if err != nil {
		return article, err
	}
	return article, err
}
func QueryChapterById(Db *sqlx.DB, Id int) ([]ChapterSeruct, error) {
	chapterSet := make([]ChapterSeruct, 0)
	chapter := ChapterSeruct{}
	rows, err := Db.Queryx("select * from jieqi_article_chapter where articleid = ? and chaptertype <> 1", Id)
	defer rows.Close()
	if err != nil {
		return chapterSet, err
	}
	for rows.Next() {
		err = rows.StructScan(&chapter)
		if err != nil {
			return chapterSet, err
		}
		chapterSet = append(chapterSet, chapter)
	}
	return chapterSet, err
}
func QueryChapterByCid(Db *sqlx.DB, Cid int) (ChapterSeruct, error) {
	var chapter ChapterSeruct
	row := DB.QueryRowx("select * from jieqi_article_chapter where chaptertype <> 1 and chapterid = ?", Cid)
	if err := row.Err(); err != nil {
		return chapter, err
	}
	err := row.StructScan(&chapter)
	if err != nil {
		return chapter, err
	}
	return chapter, err
}
func QueryArticleLastupdate(Db *sqlx.DB, n int) ([]ArticleStruct, error) {
	articleSet := make([]ArticleStruct, 0)
	article := ArticleStruct{}
	rows, err := DB.Queryx("select * from jieqi_article_article order by lastupdate desc limit ?", n)
	defer rows.Close()
	//err := row.StructScan(&article)
	if err != nil {
		return articleSet, err
	}
	for rows.Next() {
		err = rows.StructScan(&article)
		if err != nil {
			return articleSet, err
		}
		articleSet = append(articleSet, article)
	}
	return articleSet, err
}
func CountArticleAll(DB *sqlx.DB) (int64, error) {
	var count int64
	row := DB.QueryRowx("select count(articleid) from jieqi_article_article")
	err := row.Scan(&count)
	if err != nil {
		return count, err
	}
	return count, err
}
func QueryBySort(Db *sqlx.DB, index int, max int) ([]ArticleStruct, error) {
	articleSet := make([]ArticleStruct, 0)
	article := ArticleStruct{}
	rows, err := DB.Queryx("select articleid,lastupdate from jieqi_article_article limit ?,?", index, max)
	defer rows.Close()
	//err := row.StructScan(&article)
	if err != nil {
		return articleSet, err
	}
	for rows.Next() {
		err = rows.StructScan(&article)
		if err != nil {
			return articleSet, err
		}
		articleSet = append(articleSet, article)
	}
	return articleSet, err
}
func QueryArticleByIdSet(Db *sqlx.DB, listId []int) ([]ArticleStruct, error) {
	articleSet := make([]ArticleStruct, 0)
	article := ArticleStruct{}
	str := ""
	for k, v := range listId {
		if k == 0 {
			str = fmt.Sprintf("'%d'", v)
		} else {
			str = fmt.Sprintf("%s,'%d'", str, v)
		}
	}
	sql := fmt.Sprintf("select * from jieqi_article_article where articleid in (%s)", str)
	rows, err := DB.Queryx(sql)
	defer rows.Close()
	//err := row.StructScan(&article)
	if err != nil {
		return articleSet, err
	}
	for rows.Next() {
		err = rows.StructScan(&article)
		if err != nil {
			return articleSet, err
		}
		articleSet = append(articleSet, article)
	}
	return articleSet, err
}
func QueryList(Db *sqlx.DB, class int, full int, label int, size int, page int) ([]ArticleStruct, ListPages, error) {
	var (
		condition                                 ListCndition
		sql, sqlclass, sqlfull, sqllabel, sqlsize string
		articleList                               []ArticleStruct
		article                                   ArticleStruct
		maxItem                                   int64
		pageInfo                                  ListPages
		pageItem                                  int = 20
	)
	condition.Class = Web.Class
	condition.Full = []string{"", "fullflag = 0", "fullflag = 1"}
	condition.Label = []string{"lastupdate", "monthvote", "monthvisit", "monthvote", "postdate"}
	condition.Size = []string{"", "size<100000", "size>100000 and size<300000", "size>300000 and size<500000", "size>500000 and size<1000000", "size>1000000 and size<2000000", "size >2000000"}

	if class > 0 || full > 0 || size > 0 {
		sql = "from jieqi_article_article where"
	} else {
		sql = "from jieqi_article_article"
	}
	if class < len(condition.Class) {
		if class > 0 {
			sqlclass = fmt.Sprintf("sortid = %d", class)
			sql = fmt.Sprintf("%s %s", sql, sqlclass)
		}
	} else {
		return articleList, pageInfo, errors.New("NotCondition")
	}
	if full <= 2 {
		if full > 0 {
			sqlfull = condition.Full[full]
			if class > 0 {
				sql = fmt.Sprintf("%s and %s", sql, sqlfull)
			} else {
				sql = fmt.Sprintf("%s %s", sql, sqlfull)
			}
		}
	} else {
		return articleList, pageInfo, errors.New("NotCondition")
	}
	if size < len(condition.Size) {
		if size > 0 {
			sqlsize = condition.Size[size]
			if class > 0 || full > 0 {
				sql = fmt.Sprintf("%s and %s", sql, sqlsize)
			} else {
				sql = fmt.Sprintf("%s %s", sql, sqlsize)
			}
		}

	} else {
		return articleList, pageInfo, errors.New("NotCondition")
	}
	if label < len(condition.Label) {
		sqllabel = fmt.Sprintf("order by %s desc", condition.Label[label])
		sql = fmt.Sprintf("%s %s", sql, sqllabel)
	} else {
		return articleList, pageInfo, errors.New("NotCondition")
	}
	countsql := fmt.Sprintf("select count(*) %s", sql)
	sql = fmt.Sprintf("select * %s", sql)
	//fmt.Printf("\nsql:%s\ncountsql:%s", sql, countsql)
	ros := DB.QueryRowx(countsql)
	if err := ros.Err(); err != nil {
		return articleList, pageInfo, err
	}
	err := ros.Scan(&maxItem)
	if err != nil {
		return articleList, pageInfo, err
	}
	pageInfo.Pages = make([]int, 0)
	if page <= 0 {
		page = 1

	}
	pageInfo.Pages = append(pageInfo.Pages, page)
	pageInfo.CurrPage = int64(page)
	pageInfo.TotalPages = int64(math.Ceil(float64(maxItem) / float64(pageItem)))
	pageInfo.LastPage = pageInfo.TotalPages
	if page > int(pageInfo.TotalPages)+1 {
		return articleList, pageInfo, errors.New("NotCondition")
	}
	if page >= 1 {
		for i := 1; i < 5; i++ {
			if page-i <= 0 {
				break
			}
			pageInfo.Pages = append(pageInfo.Pages, page-i)
		}
		for i := 1; i < 5; i++ {
			if page+i > int(pageInfo.TotalPages) {
				break
			}
			pageInfo.Pages = append(pageInfo.Pages, page+i)
		}
	}
	pageInfo.FirstPage = int64(pageInfo.Pages[0])
	sort.Ints(pageInfo.Pages)
	page--
	sql = fmt.Sprintf("%s limit %d,%d", sql, page*pageItem, pageItem)
	rows, err := DB.Queryx(sql)
	if err != nil {
		return articleList, pageInfo, err
	}
	for rows.Next() {
		err = rows.StructScan(&article)
		if err != nil {
			return articleList, pageInfo, err
		}
		articleList = append(articleList, article)
	}
	defer rows.Close()
	return articleList, pageInfo, nil
}

package datamodels

import (
	"errors"
	//	"html/template"
	"sort"
	"time"
)

type BookShelf struct {
	Name string
	Id   int
}
type AjaxData struct {
	Info   string
	Status int
	Url    string
}

type UserInfo struct {
	Username string
	Userid   string
}
type ListPages struct {
	Pages      []int
	CurrPage   int64
	FirstPage  int64
	LastPage   int64
	TotalPages int64
	PrevPage   int64
	NextPage   int64
	Sort       int
	Class      int
	Full       int
	Label      int
	Size       int
	ClassName  []string
	FullName   []string
	LabelName  []string
	SizeName   []string
}
type ListCndition struct {
	Class []string
	Full  []string
	Label []string
	Size  []string
}
type BookState struct {
	NotFound           error
	ServiceUnavailable error
}

type BookData struct {
	Article     ArticleStruct
	Chapter     []ChapterSeruct
	Webinfo     WebInfo
	BackChapter []ChapterSeruct
	Cache       bool
}
type BookDataChapter struct {
	Article        ArticleStruct
	Chapter        []ChapterSeruct
	Webinfo        WebInfo
	BackChapter    []ChapterSeruct
	CurrentChapter ChapterRelevant
}

type WebInfo struct {
	Name   string
	Url    string
	UrlM   string
	Class  []string
	UrlImg string
}
type ChapterRelevant struct {
	DownChapterId         int
	DownChapterName       string
	CurrentChapterContent string
	CurrentChapterName    string
	CurrentChapterId      int
	UpChapterId           int
	UpChapterName         string
}
type ArticleStructSet struct {
	List []ArticleStruct
}
type ArticleStruct struct {
	Articleid      int `db:"articleid"`
	Siteid         int `db:"siteid"`
	Postdate       int `db:"postdate"`
	Lastupdate     int `db:"lastupdate"`
	LastTime       time.Time
	Makezipdate    int    `db:"makezipdate"`
	Articlename    string `db:"articlename"`
	Keywords       string `db:"keywords"`
	Initial        string `db:"initial"`
	Authorid       int    `db:"authorid"`
	Author         string `db:"author"`
	Posterid       int    `db:"posterid"`
	Poster         string `db:"poster"`
	Agentid        int    `db:"agentid"`
	Agent          string `db:"agent"`
	Sortid         int    `db:"sortid"`
	Class          string
	Typeid         int    `db:"typeid"`
	Intro          string `db:"intro"`
	Notice         string `db:"notice"`
	Setting        string `db:"setting"`
	Lastvolumeid   int    `db:"lastvolumeid"`
	Lastvolume     string `db:"lastvolume"`
	Lastchapterid  int    `db:"lastchapterid"`
	Lastchapter    string `db:"lastchapter"`
	Chapters       int    `db:"chapters"`
	Size           int    `db:"size"`
	Lastvisit      int    `db:"lastvisit"`
	Dayvisit       int    `db:"dayvisit"`
	Weekvisit      int    `db:"weekvisit"`
	Monthvisit     int    `db:"monthvisit"`
	Allvisit       int    `db:"allvisit"`
	Lastvote       int    `db:"lastvote"`
	Dayvote        int    `db:"dayvote"`
	Weekvote       int    `db:"weekvote"`
	Monthvote      int    `db:"monthvote"`
	Allvote        int    `db:"allvote"`
	Vipvotetime    int    `db:"vipvotetime"`
	Vipvotenow     int    `db:"vipvotenow"`
	Vipvotepreview int    `db:"vipvotepreview"`
	Goodnum        int    `db:"goodnum"`
	Hoodnum        int    `db:"hoodnum"`
	Badnum         int    `db:"badnum"`
	Toptime        int    `db:"toptime"`
	Saleprice      int    `db:"saleprice"`
	Salenum        int    `db:"salenum"`
	Totalcost      int    `db:"totalcost"`
	Articletype    int    `db:"articletype"`
	Permission     int    `db:"permission"`
	Firstflag      int    `db:"firstflag"`
	Fullflag       int    `db:"fullflag"`
	Imgflag        int    `db:"imgflag"`
	ImgUrl         string
	Power          int `db:"power"`
	Display        int `db:"display"`
}
type ChapterSeruct struct {
	Chapterid    int
	Siteid       int
	Articleid    int
	Articlename  string
	Weedid       string
	Volumeid     int
	Posterid     int
	Poster       string
	Postdate     int
	Lastupdate   int
	Chaptername  string
	Chapterorder int
	Size         int
	Saleprice    int
	Salenum      int
	Totalcost    int
	Attachment   string
	Isvip        int
	Chaptertype  int
	Power        int
	Display      int
}
type BookListPageJson struct {
	List []ChapterPageJson
	Page string
}
type ChapterPageJson struct {
	Articleid   int
	Chapterid   int
	Chaptername string
}
type BookcaseSeruct struct {
	Caseid       int
	Articleid    int
	Articlename  string
	Classid      int
	Userid       int
	Username     string
	Chapterid    int
	Chaptername  string
	Chapterorder int
	Joindate     int
	Lastvisit    int
	Flag         int
}
type UserErr struct {
	UorPErr            error
	BookMarkSame       error
	NoRow              error
	NotLogin           error
	ServiceUnavailable error
	LoginTrue          error
	SessExpireRelet    error
	BookcaseFull       error
}
type JsonSearchArticle struct {
	Author   string
	Name     string
	Id       int
	Cid      int
	Cname    string
	Intro    string
	Fullflag int
	Imgurl   string
	Sortid   int
	Class    string
}
type JsonSearchAll struct {
	List   []JsonSearchArticle
	Pages  int
	Page   int
	Status int
}
type ChapterWrapper struct {
	List []ChapterSeruct
	By   func(p, q *ChapterSeruct) bool
}
type SortBy func(p, q *ChapterSeruct) bool
type Sconfig struct {
	WeedMaster          string
	WeedFiler           []string
	WeedTimeout         time.Duration
	RedisHost           string
	RedisPass           string
	MysqlLink           string
	RedisSess           int
	RedisBook           int
	RedisSessExpiry     time.Duration
	RedisBookExpiry     time.Duration
	Tag                 string
	BookcassAllCount    int
	BookCaseClassMax    int
	BookCasePageItemMax int
	BookListPageItemMax int
	BookRecSet1         []string
	BookRecSet2         []string
	EsHost              string
	EsIndex             string
}
type PagesState struct {
	ErrInfo     string
	Err         string
	ErrSynopsis string
	Web         WebInfo
}
type JsonSoapi struct {
	Q string `valid:"stringlength(1|20),required"`
	P int    `valid:"int,required"`
}

var (
	BState BookState
	UErr   UserErr
	Scon   Sconfig
	Web    WebInfo
)

func init() {
	Scon.BookcassAllCount = 200
	Scon.BookCaseClassMax = 10
	Scon.BookCasePageItemMax = 15
	Scon.BookListPageItemMax = 20
	Scon.BookRecSet1 = []string{}
	Scon.BookRecSet2 = []string{}
	Scon.WeedMaster = "192.168.44.143:9444"
	Scon.WeedFiler = []string{"192.168.44.143:9888"}
	Scon.WeedTimeout = 5 * time.Minute
	Scon.RedisHost = "192.168.44.143:6379"
	Scon.RedisPass = "tWY9U52j"
	Scon.MysqlLink = "yxroot:J6Yrkp3uBH6XVkZV@tcp(192.168.44.143:3306)/yxdata?charset=utf8"
	Scon.RedisSess = 4
	Scon.RedisBook = 1
	Scon.Tag = "y_"
	//小时
	Scon.RedisSessExpiry = 24 * time.Hour
	//秒
	Scon.RedisBookExpiry = 3600 * time.Second
	//--------------
	//	Scon.BookcassAllCount = 200
	//	Scon.WeedMaster = "192.168.1.111:9333"
	//	Scon.WeedFiler = []string{"192.168.1.111:9888"}
	//	Scon.WeedTimeout = 5 * time.Minute
	//	Scon.RedisHost = "192.168.1.111:6379"
	//	Scon.RedisPass = ""
	//	Scon.MysqlLink = "root:@tcp(127.0.0.1:3306)/jieqi?charset=utf8"
	//	Scon.RedisSess = 4
	//	Scon.RedisBook = 1
	//	Scon.Tag = "y_"
	//	//小时
	//	Scon.RedisSessExpiry = 24 * time.Hour
	//	//秒
	//	Scon.RedisBookExpiry = 3600 * time.Second
	//---------------
	/////////////
	BState.NotFound = errors.New("404")
	BState.ServiceUnavailable = errors.New("503")
	UErr.UorPErr = errors.New("user Or pass Error")
	UErr.BookMarkSame = errors.New("mark same")
	UErr.NoRow = errors.New("sql: no rows in result set")
	UErr.NotLogin = errors.New("NotLogin")
	UErr.ServiceUnavailable = errors.New("503")
	UErr.LoginTrue = errors.New("LoginTrue")
	UErr.SessExpireRelet = errors.New("SessExpireRelet")
	UErr.BookcaseFull = errors.New("BookcaseFull")
	///////////////////////
}
func (book ChapterWrapper) Len() int { // 重写 Len() 方法
	return len(book.List)
}
func (book ChapterWrapper) Swap(i, j int) { // 重写 Swap() 方法
	book.List[i], book.List[j] = book.List[j], book.List[i]
}
func (book ChapterWrapper) Less(i, j int) bool { // 重写 Less() 方法
	return book.By(&book.List[i], &book.List[j])
}
func SortPerson(List []ChapterSeruct, By SortBy) { // SortPerson 方法
	sort.Sort(ChapterWrapper{List, By})
}
func ServerConfig() WebInfo {
	return Web
}

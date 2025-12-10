package datamodels

import (
	"time"
)

type UserRegister struct {
	Username   string `valid:"stringlength(2|14),required"`
	Email      string `valid:"email,required""`
	Password   string `valid:"stringlength(6|16),required"`
	Repassword string `valid:"stringlength(6|16),required"`
	Sex        string `valid:"int"`
	Code       string `valid:"stringlength(4|4),int,required"`
	Vid        string `valid:"required"`
}

type BookCaseData struct {
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
type BookCaseResultData struct {
	Articleid         int
	Author            string
	Articlename       string
	Chapterid         int
	Chaptername       string
	Lastvisit         int
	Chapterorder      int
	Bookchapterorder  int
	Booklastchapterid int
	Booklastchapter   string
	LastTime          time.Time
	ImgUrl            string
	Imgflag           int
	Fullflag          int
	Classid           int
}
type UserData struct {
	Uid        int
	Siteid     int
	Uname      string
	Name       string
	Pass       string
	Groupid    int
	Regdate    int
	Initial    string
	Sex        string
	Email      string
	Url        string
	Avatar     int
	Workid     int
	Qq         string
	Icq        string
	Msn        string
	Mobile     string
	Sign       string
	Setting    string
	Intro      string
	Badges     string
	Lastlogin  int
	Showsign   int
	Viewemail  int
	Notifymode int
	Adminemail int
	Monthscore int
	Weekscore  int
	Dayscore   int
	Lastscore  int
	Experience int
	Score      int
	Egold      int
	Esilver    int
	Credit     int
	Goodnum    int
	Badnum     int
	Isvip      int
	Overtime   int
	State      int
}
type JsonMove struct {
	Bookid []int
	Action int `valid:"int,required"`
}
type JsonDel struct {
	Bookid []int
}
type JsonAdd struct {
	Articleid int `valid:"int,required"`
	Chapterid int
}
type JsonInsetMess struct {
	Articleid int    `valid:"int,required"`
	Content   string `valid:"stringlength(0|300),required"`
	Title     string `valid:"stringlength(2|60),required"`
	Chapterid int
}
type BookErrMessage struct {
	Bookid    int
	Chapterid int
	Title     string
	Content   string
	Postdate  int64
	Toid      int
	Toname    string
	Attachsig int
}

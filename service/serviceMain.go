package service

import (
	"database/sql"
	"errors"
	"fmt"

	//	"fmt"
	"math"

	//	"strings"
	"strconv"
	"time"

	//	"fmt"
	. "zyg/datamodels"
	"zyg/repositories"

	//	log "github.com/cihub/seelog"
	"github.com/go-redis/redis"
)

var (
	CacheTime cacheTimes
	RecIdList []int
)

type cacheTimes struct {
	RedisBookExpiry time.Duration
	RedisSessExpiry time.Duration
}

func GetBook(Id int) (BookData, error) {
	Idstring := strconv.Itoa(Id)
	book, err := repositories.GetRedisBookData(Idstring)
	if err == redis.Nil {
		book, err = GetBookCombo(Id)
	} else if err != nil {
		//err = errors.New("503")
		return book, err
	}
	book.Article.ImgUrl = GetImgUrl(book.Article.Articleid, book.Article.Imgflag, Web.UrlImg)
	//err = errors.New("200")
	return book, err
}
func GetBookChapter(Id int, Cid int) (BookDataChapter, error) {
	var bookDataChapter BookDataChapter
	Idstring := strconv.Itoa(Id)
	book, err := repositories.GetRedisBookData(Idstring)
	if err == redis.Nil {
		book, err = GetBookCombo(Id)
		if err == sql.ErrNoRows {
			return bookDataChapter, BState.NotFound
		} else if err != nil {
			return bookDataChapter, err
		}
		book.Cache = false
	} else if err != nil {
		return bookDataChapter, err
	}
	chapterRelevant, err := LookupChapterRelevant(book, Cid)
	if err == BState.NotFound {
		if book.Cache == true {
			book, err = GetBookCombo(Id)
			chapterRelevant, err = LookupChapterRelevant(book, Cid)
			if err != nil {
				return bookDataChapter, err
			}
		} else {
			return bookDataChapter, err
		}
	} else if err != nil {
		return bookDataChapter, err
	}
	bookDataChapter.Article = book.Article
	bookDataChapter.Chapter = book.Chapter
	bookDataChapter.Webinfo = book.Webinfo
	bookDataChapter.CurrentChapter = chapterRelevant
	return bookDataChapter, err

}
func LookupChapterRelevant(book BookData, Cid int) (ChapterRelevant, error) {
	var (
		chapterRelevant ChapterRelevant
		err             error
	)
	for x, y := range book.Chapter {
		if y.Chapterid == Cid {
			chapterRelevant.CurrentChapterContent, err = repositories.GetSeaweedVolumeUrl(y.Articleid, y.Chapterid, y.Weedid)
			//chapterRelevant.CurrentChapterContent, err = repositories.GetChapterContent(y.Weedid)
			if err != nil {
				return chapterRelevant, err
			}
			//chapterRelevant.CurrentChapterContent = unescaped(Content)
			chapterRelevant.CurrentChapterName = y.Chaptername
			chapterRelevant.CurrentChapterId = y.Chapterid
			if x > 0 {
				chapterRelevant.UpChapterId = book.Chapter[x-1].Chapterid
				chapterRelevant.UpChapterName = book.Chapter[x-1].Chaptername
			} else {
				chapterRelevant.UpChapterId = 0
			}
			if x < len(book.Chapter)-1 {
				chapterRelevant.DownChapterId = book.Chapter[x+1].Chapterid
				chapterRelevant.DownChapterName = book.Chapter[x+1].Chaptername
			} else {
				chapterRelevant.DownChapterId = 0
			}
			return chapterRelevant, err
		}
	}
	return chapterRelevant, BState.NotFound
}
func GetBookCombo(Id int) (BookData, error) {
	var (
		chaptersort ChapterWrapper
		book        BookData
		err         error
	)
	Idstring := strconv.Itoa(Id)
	book.Webinfo = ServerConfig()
	//获取书籍信息
	book.Article, err = repositories.QueryArticleById(repositories.DB, Id)
	if err == sql.ErrNoRows {
		return book, BState.NotFound
	} else if err != nil {

		return book, err
	}
	//获取章节信息
	chaptersort.List, err = repositories.QueryChapterById(repositories.DB, Id)
	if err != nil {
		//err = errors.New("503")
		return book, err
	}
	//章节进行排序
	SortPerson(chaptersort.List, func(p, q *ChapterSeruct) bool {
		return p.Chapterorder < q.Chapterorder // Name 递增排序
	})
	//转换更新时间
	book.Article.LastTime = time.Unix(int64(book.Article.Lastupdate), 0)
	//转换分类
	book.Article.Class = book.Webinfo.Class[book.Article.Sortid]
	//获取倒叙章节
	book.BackChapter = GetFlashbackChapter(chaptersort.List, 6)
	//获取imgurl
	book.Article.ImgUrl = GetImgUrl(Id, book.Article.Imgflag, Web.UrlImg)
	book.Chapter = chaptersort.List
	//写入redis
	book.Cache = true
	repositories.WriteRedisBookData(Idstring, book, CacheTime.RedisBookExpiry)
	return book, err
}
func GetLastUpdate(n int) ([]ArticleStruct, error) {
	var articleSet ArticleStructSet
	name := fmt.Sprintf("lastupdate%d", n)
	articleSet, err := repositories.GetRedisRecBookData(name)
	if err == redis.Nil {
		articleSet.List, err = repositories.QueryArticleLastupdate(repositories.DB, n)
		if err != nil {
			return articleSet.List, err
		}
		for x, y := range articleSet.List {
			articleSet.List[x].ImgUrl = GetImgUrl(y.Articleid, y.Imgflag, Web.UrlImg)
		}
		repositories.WriteRedisRecBookData(name, articleSet, CacheTime.RedisBookExpiry)
	} else if err != nil {
		//err = errors.New("503")
		return articleSet.List, err
	}
	//err = errors.New("200")
	return articleSet.List, err
}
func GetRecId(n int) ([]ArticleStruct, error) {
	var (
		listId     []int
		articleSet ArticleStructSet
	)
	max := len(RecIdList)
	name := fmt.Sprintf("rec%d", n)
	articleSet, err := repositories.GetRedisRecBookData(name)
	if err == redis.Nil {
		if n > max {
			articleSet.List, err = repositories.QueryArticleLastupdate(repositories.DB, n)
			if err != nil {
				return articleSet.List, err
			}
			//return res, err
		}
		idSer := GenerateRandomNumber(0, max, n)
		for _, v := range idSer {
			listId = append(listId, RecIdList[v])
		}
		articleSet.List, err = repositories.QueryArticleByIdSet(repositories.DB, listId)
		if err != nil {
			return articleSet.List, err
		}
		for x, y := range articleSet.List {
			articleSet.List[x].ImgUrl = GetImgUrl(y.Articleid, y.Imgflag, Web.UrlImg)
		}
		repositories.WriteRedisRecBookData(name, articleSet, CacheTime.RedisBookExpiry)
	} else if err != nil {
		//err = errors.New("503")
		return articleSet.List, err
	}
	//err = errors.New("200")
	return articleSet.List, err

}
func GetChapterListPage(chapter []ChapterSeruct, prepage int, page int) ([]ChapterPageJson, ListPages, error) {
	list := make([]ChapterPageJson, 0)
	pages := ListPages{Pages: []int{}}
	nums := len(chapter)
	pageMax := int(math.Ceil(float64(nums) / float64(prepage)))
	if page > pageMax {
		return list, pages, errors.New("没有此页")
	}

	for i := 1; i <= pageMax; i++ {
		pages.Pages = append(pages.Pages, i)
	}
	pages.TotalPages = int64(pageMax)
	pages.CurrPage = int64(page)
	if page-1 <= 0 {
		pages.PrevPage = 0
	} else {
		pages.PrevPage = int64(page - 1)
	}
	if page+1 > pageMax {
		pages.NextPage = 0
	} else {
		pages.NextPage = int64(page + 1)
	}
	page--
	for k, v := range chapter {
		if k >= (page*prepage)+prepage {
			break
		}
		if k >= page*prepage {
			a := ChapterPageJson{Articleid: v.Articleid, Chapterid: v.Chapterid, Chaptername: v.Chaptername}
			list = append(list, a)
		}

	}
	return list, pages, nil
}

func SearchByES(index string, class string, querytext string, size int, from int) (JsonSearchAll, error) {
	var results JsonSearchAll
	res, err := repositories.QueryByAll(index, class, querytext, size, from)
	if err != nil {
		return results, err
	}
	results.Page = res.Page
	results.Pages = res.Pages
	results.Status = res.Status
	if res.Pages == 0 {
		return results, err
	}
	results.List = make([]JsonSearchArticle, 0)
	for _, v := range res.Data {
		var result JsonSearchArticle
		result.Author = v.Author
		result.Cid = v.Cid
		if v.Imgflag {
			result.Imgurl = GetImgUrl(v.Id, 1, Web.UrlImg)
		} else {
			result.Imgurl = GetImgUrl(v.Id, 0, Web.UrlImg)
		}
		result.Class = Web.Class[v.Sortid]
		if v.Fullflag {
			result.Fullflag = 1
		} else {
			result.Fullflag = 0
		}
		result.Id = v.Id
		result.Intro = v.Intro
		result.Name = v.Name
		result.Sortid = v.Sortid
		result.Cname = v.Cname
		results.List = append(results.List, result)
	}
	return results, nil
}
func ReAjax(info string, status int, url string) AjaxData {
	errinfo := AjaxData{
		Info:   info,
		Status: status,
		Url:    url,
	}
	return errinfo
}

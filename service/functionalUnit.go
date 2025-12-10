package service

import (
	"math/rand"
	"strconv"
	"time"

	//	"errors"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"sort"
	. "zyg/datamodels"
)

func GetFlashbackChapter(chapter []ChapterSeruct, n int) []ChapterSeruct {
	back := make([]ChapterSeruct, 0)
	m := len(chapter) - 1
	for i := m; i >= 0; i-- {
		if i == m-n {
			break
		}
		back = append(back, chapter[i])

	}
	return back
}

func GetImgUrl(id int, p int, imgurl string) string {
	var url string
	if p > 0 {
		if imgurl != "" {
			url = fmt.Sprintf("%s/%d/%d/%ds.jpg", imgurl, id/1000, id, id)
		} else {
			url = fmt.Sprintf("/%d/%d/%ds.jpg", id/1000, id, id)
		}

	} else {
		url = imgurl + "/nocover.jpg"
	}
	return url
}
func unescaped(x string) interface{} {
	return template.HTML(x)
}
func GetMd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))                    // 需要加密的字符串为
	return hex.EncodeToString(h.Sum(nil)) // 输出加密结果
}
func GetUpDownPages(index int, max int, num int) []int {
	pnum := make([]int, 0)
	if index == 0 {
		index = 1
	}
	if index > max {
		return pnum
	}
	pnum = append(pnum, index)
	for i := 1; i < num; i++ {
		if index-i <= 0 {
			break
		}
		pnum = append(pnum, index-i)
	}
	for i := 1; i < num; i++ {
		if index+i > int(max) {
			break
		}
		pnum = append(pnum, index+i)
	}
	sort.Ints(pnum)
	return pnum
}

//生成count个[start,end)结束的不重复的随机数
func GenerateRandomNumber(start int, end int, count int) []int {
	//范围检查
	if end < start || (end-start) < count {
		return nil
	}

	//存放结果的slice
	nums := make([]int, 0)
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		//生成随机数
		num := r.Intn((end - start)) + start

		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, num)
		}
	}

	return nums
}
func StrToInt(a string, b string) (int, int, error) {
	var aid, bid int
	aid, err := strconv.Atoi(a)
	if err != nil {
		//fmt.Println(err)
		return aid, bid, err
	}
	bid, err = strconv.Atoi(b)
	if err != nil {
		//fmt.Println(err)
		return aid, bid, err
	}
	return aid, bid, err
}
func QQreminder(articleid int, articlename string) string {
	return fmt.Sprintf("http://qzs.qq.com/snsapp/app/bee/widget/open.htm#content=《%s》好像已经更新了，赶快去看一看吧！&time=%v 9:00&advance=0&url=http://www.yixuanju.com/book/%d",
		articlename, time.Now().Format("2006-1-2"), articleid)
}

package repositories

//package repositories

import (
	"errors"
	"fmt"

	//	. "zyg/datamodels"
	"math/rand"
	"regexp"
	"strings"
	"time"

	//	"time"
	"unicode/utf8"

	"github.com/BelieveR44/goseaweedfs"
	"github.com/axgle/mahonia"
	"github.com/kirinlabs/HttpRequest"
)

var (
	Sw       *goseaweedfs.Seaweed
	SwConfig weedConfig
)

type weedConfig struct {
	Master  string
	Filer   []string
	Volume  []string
	Timeout time.Duration
}

func init() {
	Sw = goseaweedfs.NewSeaweed("http", SwConfig.Master, SwConfig.Filer, 2*1024*1024, SwConfig.Timeout)
	//defer DBW.Db.Close()
}
//20.8.26 去掉filers的依赖
func GetChapterContent(Weedid string) (string, error) {
	val, err := Sw.LookupFileID(Weedid, nil, true)
	if err != nil {
		return "", err
	}
	//fmt.Println(val)
	body, err := GetSeaweedUrl(val)
	if err != nil {
		return body, err
	}
	if strings.Count(body, "")-1 < 2 {
		err = errors.New(fmt.Sprintf("内容长度过短WeedId:%v", Weedid))
		//return body, errors.New("内容长度过短/WeedId")
	}
	fmt.Println(body)
	txt := contentFmt(body)
	return txt, err
}
func GetSeaweedUrl(url string) (string, error) {
	req := HttpRequest.NewRequest()
	req.SetTimeout(5)
	req.DisableKeepAlives(false)
	res, err := req.Get(url, nil)
	if err != nil {
		//log.Println(err)
		return "", err
	}
	body, err := res.Body()
	if err != nil {
		//log.Println(err)
		return string(body), err
	}
	if res.StatusCode() != 200 && res.StatusCode() != 304 {
		err = errors.New(fmt.Sprintf("状态码不符合Url:%v,Code:%v", url, res.StatusCode()))
		//return string(body), err
	}
	return string(body), err
}

//func GetFilerUrl(id, cid int) (string, error) {
//	if len(Sw.Filers) == 0 {
//		return "", errors.New("filers为空")
//	}
//	url := fmt.Sprintf("http://%s/%d/%d.txt", Scon.WeedFiler[0], id, cid)
//	resp, err := http.Get(url)
//	if err != nil {
//		return "", err
//	}

//}
//func ArrangeContent(con string) (string, error) {

//}
func GetSeaweedVolumeUrl(id, cid int, wid string) (string, error) {
	if len(Sw.Filers) == 0 {
		return "", errors.New("filers为空")
	}
	rand.Seed(time.Now().UnixNano())
	x := rand.Intn(1)
	url := fmt.Sprintf("%s/%s.txt", SwConfig.Volume[x], wid)
	req := HttpRequest.NewRequest()
	req.SetTimeout(5)
	req.DisableKeepAlives(false)
	res, err := req.Get(url, nil)
	if err != nil {
		//log.Println(err)
		return "", err
	}
	defer res.Response().Body.Close()
	if res.StatusCode() != 200 && res.StatusCode() != 304 {
		err = errors.New(fmt.Sprintf("状态码不符合Url:%v,Code:%v", url, res.StatusCode()))
		return "", err
	}
	body, err := res.Body()
	if err != nil {
		return string(body), err
	}
	if strings.Count(string(body), "")-1 < 2 {
		err = errors.New(fmt.Sprintf("内容长度过短cId:%d", cid))
		return "", err
	}
	if !utf8.Valid(body) {
		body = ConvertToByte(string(body), "gbk", "utf8")
	}
	var re = regexp.MustCompile(`\<`)
	txt := re.ReplaceAllString(string(body), "&lt;")
	re = regexp.MustCompile(`\>`)
	txt = re.ReplaceAllString(txt, "&gt;")
	re, _ = regexp.Compile(`\r\n`)
	txt = re.ReplaceAllString(txt, "<br />")
	//	re, _ = regexp.Compile(`\s\s`)
	//	txt = re.ReplaceAllString(txt, "&nbsp;&nbsp;")
	//fmt.Printf("\n%v", txt)
	return txt, nil
}
func GetSeaweedFilerUrl(id, cid int) (string, error) {
	if len(Sw.Filers) == 0 {
		return "", errors.New("filers为空")
	}
	url := fmt.Sprintf("http://%s/%d/%d/%d.txt", SwConfig.Filer[0], id/1000, id, cid)
	req := HttpRequest.NewRequest()
	req.SetTimeout(5)
	req.DisableKeepAlives(false)
	res, err := req.Get(url, nil)
	if err != nil {
		//log.Println(err)
		return "", err
	}
	defer res.Response().Body.Close()
	if res.StatusCode() != 200 && res.StatusCode() != 304 {
		err = errors.New(fmt.Sprintf("状态码不符合Url:%v,Code:%v", url, res.StatusCode()))
		return "", err
	}
	body, err := res.Body()
	if err != nil {
		return string(body), err
	}
	if strings.Count(string(body), "")-1 < 2 {
		err = errors.New(fmt.Sprintf("内容长度过短cId:%d", cid))
		return "", err
	}
	if !utf8.Valid(body) {
		body = ConvertToByte(string(body), "gbk", "utf8")
	}
	var re = regexp.MustCompile(`\<`)
	txt := re.ReplaceAllString(string(body), "&lt;")
	re = regexp.MustCompile(`\>`)
	txt = re.ReplaceAllString(txt, "&gt;")
	re, _ = regexp.Compile(`\r\n`)
	txt = re.ReplaceAllString(txt, "<br />")
	//	re, _ = regexp.Compile(`\s\s`)
	//	txt = re.ReplaceAllString(txt, "&nbsp;&nbsp;")
	//fmt.Printf("\n%v", txt)
	return txt, nil
}
func ConvertToByte(src string, srcCode string, targetCode string) []byte {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(targetCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)

	return cdata
}
func contentFmt(body string)string{
	txtbyte := []byte(body)
	if !utf8.Valid([]byte(body)) {
		txtbyte = ConvertToByte(body, "gbk", "utf8")
	}
	var re = regexp.MustCompile(`\<`)
	txt := re.ReplaceAllString(string(txtbyte), "&lt;")
	re = regexp.MustCompile(`\>`)
	txt = re.ReplaceAllString(txt, "&gt;")
	re, _ = regexp.Compile(`    `)
	txt = re.ReplaceAllString(txt, "&nbsp;&nbsp;&nbsp;&nbsp;")
	re, _ = regexp.Compile(`\r\n`)
	txt = re.ReplaceAllString(txt, "<br />")
	return txt
}
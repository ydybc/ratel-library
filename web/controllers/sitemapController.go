package controllers

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"regexp"
	"strconv"
	"strings"
	"time"
	"zyg/service"

	"github.com/kataras/iris/v12/context"
)

//type urlset struct {
//	Url []SitemapDoc `xml:"url"`
//}
type urlset struct {
	Url []SitemapDoc `xml:"url"`
}
type sitemapindex struct {
	Sitemap []SitemapindexDoc `xml:"sitemap"`
}

type SitemapDoc struct {
	Loc     string `xml:"loc"`
	Lastmod string `xml:"lastmod"`
	Mobile string `xml:"mobile"`
}
type MobileAttr struct {
	mobile string `xml:",attr"`
}
type SitemapindexDoc struct {
	Loc string `xml:"loc"`
}

var xmlconfig context.XML

func init() {
	xmlconfig.Prefix = `<?xml version="1.0" encoding="utf-8"?>`
}

/*
func (c *MainController) GetSitemapIndexAll() {
	sitemapAll := sitemapindex{}
	page, err := service.SitemapCount()
	//fmt.Println(page)
	if err != nil {
		c.Ctx.StatusCode(503)
		return
	}
	for i := 0; i <= int(page); i++ {
		sitemapAll.Sitemap = append(sitemapAll.Sitemap,
			SitemapindexDoc{fmt.Sprintf("http://www.ziyouge.com/sitemap/index/%d", i+1)})
	}
	c.Ctx.XML(sitemapAll, xmlconfig)
}



func (c *MainController) GetSitemapIndexBy(index int) {
	sitemap := urlset{}
	articleSet, err := service.SitemapQuery(index)
	if err != nil {
		c.Ctx.StatusCode(503)
		return
	}
	for _, v := range articleSet {
		sitemap.Url = append(sitemap.Url, SitemapDoc{fmt.Sprintf("http://www.ziyouge.com/book/%d", v.Articleid),
			time.Unix(int64(v.Lastupdate), 0).Format("2006-01-02")})
		sitemap.Url = append(sitemap.Url, SitemapDoc{fmt.Sprintf("http://m.ziyouge.com/novel/%d", v.Articleid),
			time.Unix(int64(v.Lastupdate), 0).Format("2006-01-02")})
	}

	c.Ctx.XML(sitemap, xmlconfig)
}
*/
func GetSitemapIndexMain(ctx iris.Context) {
	index := ctx.Params().Get("index")
	index = strings.Replace(index, ".xml", "", -1)
	if index == "all" {
		// 响应客户端
		sitemapAll := sitemapindex{}
		page, err := service.SitemapCount()
		//fmt.Println(page)
		if err != nil {
			ctx.StatusCode(503)
			return
		}
		for i := 0; i <= int(page); i++ {
			sitemapAll.Sitemap = append(sitemapAll.Sitemap,
				SitemapindexDoc{fmt.Sprintf("http://www.ziyouge.com/sitemap/index/%d.xml", i+1)})
		}
		ctx.XML(sitemapAll, xmlconfig)
	}else{
		indexInt, err := strconv.Atoi(index)
		if err != nil {
			ctx.StatusCode(503)
		}else if indexInt > 0 && indexInt <= 9999 {
			sitemap := urlset{}
			articleSet, err := service.SitemapQuery(indexInt)
			if err != nil {
				ctx.StatusCode(503)
				return
			}
			for _, v := range articleSet {
				sitemap.Url = append(sitemap.Url, SitemapDoc{fmt.Sprintf("http://www.ziyouge.com/book/%d", v.Articleid),
					time.Unix(int64(v.Lastupdate), 0).Format("2006-01-02"),
				""})
				//sitemap.Url = append(sitemap.Url, SitemapDoc{fmt.Sprintf("http://m.ziyouge.com/novel/%d", v.Articleid),
				//	time.Unix(int64(v.Lastupdate), 0).Format("2006-01-02")})
			}
		//ctx.XML(sitemap,xmlconfig)
			data, err := xml.Marshal(sitemap)
			if err != nil {
				fmt.Println(err)
				ctx.StatusCode(404)
			}else{
				doc2,err:=MobileXmlREG(data)
				if err != nil {
					fmt.Println(err)
					ctx.StatusCode(404)
				}else{
					ctx.WriteGzip(doc2)
				}
			}
		} else {
			ctx.StatusCode(404)
		}
	}

}
func MobileXmlREG(doc []byte)([]byte,error){
	docstring := string(doc)
	re := regexp.MustCompile(`(?m)<mobile>(.+?)<\/mobile>`)
	re2 := regexp.MustCompile(`(?m)<mobile><\/mobile>`)
	substitution:="<mobile:mobile type=\"$1\"/>"
	//substitution:="<mobile mobile=\"fdfd\"/>"
	substitution2:=""
	if(re2.Match(doc)){
		docstring=re.ReplaceAllString(docstring, substitution2)
	}else if(re.Match(doc)){
		//re.ReplaceAll()
		docstring=re.ReplaceAllString(docstring, substitution)
	}else{
		const fd = "为空"
		return []byte(""),errors.New(fd)
	}

	docstring="<?xml version=\"1.0\" encoding=\"utf-8\"?>"+docstring
	return []byte(docstring),nil
}
func GetSitemapIndexSon(ctx iris.Context) {
	index := ctx.Params().Get("index")
	index = strings.Replace(index, ".xml", "", -1)
	if index == "all" {
		sitemapAll := sitemapindex{}
		page, err := service.SitemapCount()
		//fmt.Println(page)
		if err != nil {
			ctx.StatusCode(503)
			return
		}
		for i := 0; i <= int(page); i++ {
			sitemapAll.Sitemap = append(sitemapAll.Sitemap,
				SitemapindexDoc{fmt.Sprintf("http://m.ziyouge.com/sitemap/index/%d.xml", i+1)})
		}
		ctx.XML(sitemapAll, xmlconfig)
	}else {
		indexInt, err := strconv.Atoi(index)
		if err != nil {
			ctx.StatusCode(503)
		}else if indexInt > 0 && indexInt <= 9999 {
			sitemap := urlset{}
			articleSet, err := service.SitemapQuery(indexInt)
			if err != nil {
				ctx.StatusCode(503)
				return
			}
			for _, v := range articleSet {
				sitemap.Url = append(sitemap.Url, SitemapDoc{fmt.Sprintf("http://m.ziyouge.com/novel/%d", v.Articleid),
					time.Unix(int64(v.Lastupdate), 0).Format("2006-01-02"),
					"mobile"})
			}
			//ctx.XML(sitemap, xmlconfig)
			data, err := xml.Marshal(sitemap)
			if err != nil {
				fmt.Println(err)
				ctx.StatusCode(404)
			}else{
				doc2,err:=MobileXmlREG(data)
				if err != nil {
					fmt.Println(err)
					ctx.StatusCode(404)
				}else{
					ctx.WriteGzip(doc2)
				}
			}
		} else {
			ctx.StatusCode(404)
		}
	}
}

/*
func (c *MainMController) GetSitemapIndexAll() {
	sitemapAll := sitemapindex{}
	page, err := service.SitemapCount()
	//fmt.Println(page)
	if err != nil {
		c.Ctx.StatusCode(503)
		return
	}
	for i := 0; i <= int(page); i++ {
		sitemapAll.Sitemap = append(sitemapAll.Sitemap,
			SitemapindexDoc{fmt.Sprintf("http://m.ziyouge.com/sitemap/index/%d", i+1)})
	}
	c.Ctx.XML(sitemapAll, xmlconfig)
}
func (c *MainMController) GetSitemapIndexBy(index int) {
	sitemap := urlset{}
	articleSet, err := service.SitemapQuery(index)
	if err != nil {
		c.Ctx.StatusCode(503)
		return
	}
	for _, v := range articleSet {
		sitemap.Url = append(sitemap.Url, SitemapDoc{fmt.Sprintf("http://m.ziyouge.com/novel/%d", v.Articleid),
			time.Unix(int64(v.Lastupdate), 0).Format("2006-01-02")})
	}
	c.Ctx.XML(sitemap, xmlconfig)
}
*/

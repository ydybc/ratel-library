package repositories

//package repositories

import (
	//	"fmt"
	"html"

	//	"fmt"
	"math"

	//	"database/sql"
	"context"
	"encoding/json"

	//	"fmt"

	"gopkg.in/olivere/elastic.v5"
)

var EsClient *elastic.Client
var Host = "http://127.0.0.1:9200/"

type SearchArticle struct {
	Author   string `json:"author"`
	Name     string `json:"name"`
	Id       int    `json:"id"`
	Cid      int    `json:"cid"`
	Cname    string `json:"cname"`
	Intro    string `json:"intro"`
	Fullflag bool
	Imgflag  bool
	Sortid   int
}
type SearchAll struct {
	Data   []SearchArticle
	Pages  int
	Page   int
	Status int
}

func QueryByAll(index string, class string, querytext string, size int, from int) (SearchAll, error) {
	var Data SearchAll
	querytext = html.EscapeString(querytext)
	query := elastic.NewBoolQuery()
	//query = query.Should(elastic.NewMatchQuery("aname", querytext).Boost(2))
	query = query.Should(
		elastic.NewMatchQuery("name", querytext).Boost(2),
		elastic.NewMatchQuery("author", querytext).Boost(1),
		elastic.NewMatchQuery("cname", querytext).Boost(0.5),
		elastic.NewMatchQuery("intro", querytext).Boost(0.5))
	//	data, err := json.MarshalIndent(src, "", "  ")
	//	if err != nil {
	//		panic(err)
	//	}
	//fmt.Println(string(data))
	Data.Page = from
	if from <= 0 {
		Data.Pages = 0
		Data.Status = 0
		return Data, nil
	}
	q := EsClient.Search(index)
	q = q.Size(size)
	q = q.From((from - 1) * size)
	q = q.Type(class)
	q = q.Query(query)
	res, err := q.Do(context.Background())
	//res, err := EsClient.Search(index).Size(size).From(from).Type(class).Query(query).Do(context.Background())
	if err != nil {
		return Data, err
	}
	if res.Hits.TotalHits <= 0 {
		Data.Pages = 0
		Data.Status = 0
		return Data, nil
	}
	Data.Pages = int(math.Ceil(float64(res.Hits.TotalHits) / float64(size)))
	Data.Data = make([]SearchArticle, 0)
	for _, v := range res.Hits.Hits {
		var film SearchArticle
		if err := json.Unmarshal(*v.Source, &film); err != nil {
			return Data, err
		}
		Data.Data = append(Data.Data, film)
	}
	return Data, nil
}

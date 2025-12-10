package repositories

//package repositories

import (
	//	"zyg/datasource"
	"testing"
)

func TestQuerySearch(t *testing.T) {
	res, err := QuerySearch("zyg", "article", "不太懂", 3, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

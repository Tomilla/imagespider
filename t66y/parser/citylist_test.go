package parser

import (
	"io/ioutil"
	"testing"
)

func TestParseCityList(t *testing.T) {
	//contents, err := fetcher.Fetch("http://www.zhenai.com/zhenghun")
	contents, err := ioutil.ReadFile(
		"citylist_test_data.html")
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%s", contents)
	result := ParseTopicList(contents)
	const resultSize = 470
	expectedUrls := []string{
		"http://www.t66y.com/zhenghun/aba",
		"http://www.t66y.com/zhenghun/akesu",
		"http://www.t66y.com/zhenghun/alashanmeng",
	}

	expectedCities := []string{
		"City: 阿坝", "City: 阿克苏", "City: 阿拉善盟",
	}

	for i, url := range expectedUrls {
		if result.Requests[i].Url != url {
			t.Errorf("expected url #%d: %s; but "+
				"was %s", i, url, result.Requests[i].Url)
		}
	}
	for i, city := range expectedCities {
		if result.Items[i].(string) != city {
			t.Errorf("expected City #%d: %s; but "+
				"was %s", i, city, result.Items[i])
		}
	}

	if len(result.Items) != resultSize {
		t.Errorf("result should have %d "+
			"requests; but had %d", resultSize, len(result.Items))
	}

	if len(result.Requests) != resultSize {
		t.Errorf("result should have %d "+
			"requests; but had %d", resultSize, len(result.Requests))
	}
}

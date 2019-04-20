package fetcher

import (
	"bufio"
	"fmt"
	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/net"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
)

type Fetcher interface {
	Fetch(url string) ([]byte, error)
}

func Fetch(link string) ([]byte, error) {
	//@@@@@@@@@@@@@@@@@@@@@@

	client := net.NewClient()
	req, err := http.NewRequest("GET", link, nil)
	res, err := client.Do(req)

	//res, err := http.Get(link)
	//@@@@@@@@@@@@@@@@@@@@@@@@
	if err != nil {
		//panic(err)
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		log.Println("Error: status code", res.StatusCode)
		return nil, fmt.Errorf("wrong status code: %d", res.StatusCode)
	}
	defer res.Body.Close()

	bodyReader := bufio.NewReader(res.Body)
	e := determinEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)

}

func determinEncoding(r *bufio.Reader) encoding.Encoding {

	bytes, err := r.Peek(1024)
	if err != nil {
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}

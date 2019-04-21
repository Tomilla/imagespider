package fetcher

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/wuxiangzhou2010/imagespider/net"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type Fetcher interface {
	Fetch(url string) ([]byte, error)
}

func Fetch(link string) ([]byte, error) {
	//@@@@@@@@@@@@@@@@@@@@@@

	client := net.NewClient(true)
	req, err := http.NewRequest("GET", link, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36 t66y_com")
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

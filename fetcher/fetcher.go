package fetcher

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Fetcher interface {
	Fetch(url string) ([]byte, error)
}

func Fetch(link string) ([]byte, error) {
	//@@@@@@@@@@@@@@@@@@@@@@
	proxyStr := "socks5://localhost:1080"
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		log.Println(err)
	}
	tr := &http.Transport{ //解决x509: certificate signed by unknown authority
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyURL(proxyURL),
	}
	client := &http.Client{
		Timeout:   15 * time.Second,
		Transport: tr, //解决x509: certificate signed by unknown authority
	}
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

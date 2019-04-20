package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/wuxiangzhou2010/imagespider/model"
)

var C *Config
var StartPages = []string{
	//"http://t66y.com/thread0806.php?fid=8",  // 新时代
	//"http://t66y.com/thread0806.php?fid=16", //达盖尔
	"http://t66y.com/thread0806.php?fid=21", //下载区
}

type Config struct {
	sync.RWMutex
	Image        ImageConfig `json:"image"`
	StartPages   []string    `json:"startPages"`
	PageLimit    int         `json:"pageLimit"`
	Net          Net         `json:"net"`
	MameLenLimit int         `json:"nameLenLimit"`
}

func NewConfig() *Config {
	return LoadConfig()
}
func (c *Config) GetNameLenLimit() int {
	c.RLock()
	defer c.RUnlock()

	return c.MameLenLimit
}
func (c *Config) GetImageWorkerCount() int {
	return c.Image.WorkerCount
}

func (c *Config) GetStartPages() []string {
	c.RLock()
	defer c.RUnlock()

	return c.StartPages
}
func (c *Config) GetPageLimit() int {
	c.RLock()
	defer c.RUnlock()

	return c.PageLimit
}

func (c *Config) GetImageConfig() *ImageConfig {
	c.RLock()
	defer c.RUnlock()

	return &c.Image
}

func (c *Config) GetImageChan() chan model.Topic {
	c.RLock()
	defer c.RUnlock()

	return c.Image.ImageChan
}
func (c *Config) SetImageChan(ch chan model.Topic) {
	c.Lock()
	defer c.Unlock()

	c.Image.ImageChan = ch
}

func getConfigFileName() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path.Join(wd, "config.json")

}

func (c *Config) GetProxyURL() string {
	return c.Net.ProxyURL
}

func (c *Config) GetNetTimeOut() int {
	return c.Net.TimeOut

}

// LoadConfig, load config
func LoadConfig() (c *Config) {

	filename := getConfigFileName()
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err.Error())
	}
	newReader := removeComment(f)
	defer f.Close()

	jsonParser := json.NewDecoder(newReader)
	jsonParser.Decode(&c)

	//PrintConfig(c)
	return
}

//reference: https://stackoverflow.com/questions/12682405/strip-out-c-style-comments-from-a-byte

func removeComment(reader io.Reader) (newReader io.Reader) {

	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	s := string(bs)
	fmt.Println("before ", s)
	re1 := regexp.MustCompile(`(?im)^\s+\/\/.*$`) // 整行注释

	s = re1.ReplaceAllString(s, "")
	fmt.Println("after1 ", s)
	re2 := regexp.MustCompile(`\/\/[^"\[\]]+\n`) // 行末
	s = re2.ReplaceAllString(s, "")
	fmt.Println("after2 ", s)
	newReader = strings.NewReader(s)
	return

}

func PrintConfig(cfg *Config) {

	fmt.Printf("%+v\n", cfg)
}

func init() {
	C = NewConfig() // default config
}

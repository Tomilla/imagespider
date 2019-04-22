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

type Config struct {
	sync.RWMutex
	Image        ImageConfig `json:"image"`
	Init         Init        `jsonL"init"`
	Net          Net         `json:"net"`
	MameLenLimit int         `json:"nameLenLimit"`
	Engine       Engine      `json:"engine"`
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

	return c.Init.Seeds
}
func (c *Config) GetPageLimit() int {
	c.RLock()
	defer c.RUnlock()

	return c.Init.TopicPerPage
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
func (c *Config) GetEngineWorkerCount() int {
	c.Lock()
	defer c.Unlock()

	return c.Engine.WorkerCount
}
func (c *Config) GetEngineElasticUrl() string {
	c.Lock()
	defer c.Unlock()

	return c.Engine.ElasticUrl
}
func (c *Config) GetStartPageNum() int {
	c.Lock()
	defer c.Unlock()

	return c.Init.StartPageNum
}

func (c *Config) GetEndPageNum() int {
	c.Lock()
	defer c.Unlock()

	return c.Init.EndPageNum
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
	//fmt.Println("before ", s)
	re1 := regexp.MustCompile(`(?im)^\s+\/\/.*$`) // 整行注释

	s = re1.ReplaceAllString(s, "")
	//fmt.Println("after1 ", s)
	re2 := regexp.MustCompile(`\/\/[^"\[\]]+\n`) // 行末
	s = re2.ReplaceAllString(s, "")
	//fmt.Println("after2 ", s)
	newReader = strings.NewReader(s)
	return

}

func PrintConfig(cfg *Config) {

	fmt.Printf("%+v\n", cfg)
}

func init() {
	C = NewConfig() // default config
}

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/model"
)

var C *Config
var StartPages = []string{
	//"http://t66y.com/thread0806.php?fid=8",  // 新时代
	//"http://t66y.com/thread0806.php?fid=16", //达盖尔
	"http://t66y.com/thread0806.php?fid=21", //下载区
}

type Config struct {
	sync.RWMutex
	Image      ImageConfig `json:"image"`
	StartPages []string    `json:"startPages"`
	PageLimit  int         `json:"pageLimit"`
	ProxyURL   string      `json:"proxyUrl"`
}

func NewConfig() *Config {
	return LoadConfig()
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
	return c.ProxyURL
}

func LoadConfig() (c *Config) {

	filename := getConfigFileName()
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer f.Close()

	jsonParser := json.NewDecoder(f)
	jsonParser.Decode(&c)

	PrintConfig(c)
	return
}

func PrintConfig(cfg *Config) {
	//log.Print(
	//	"\n##### config ########\n",
	//	"rate\t\t\t: ", cfg.Rate,
	//	"\nendpoint\t\t: ", cfg.Endpoint,
	//	"\ntxBuffer\t\t: ", cfg.RawTxBuffer,
	//	"\nsignedTxBuffer\t\t: ", cfg.SignedTxBuffer,
	//	"\nlast\t\t\t: ", cfg.Last,
	//	"\n#######################\n",
	//)
	fmt.Printf("%+v\n", cfg)
}

func init() {
	C = NewConfig()
}

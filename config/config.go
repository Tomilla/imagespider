package config

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "os"
    "path"
    "runtime"
    "sync"
    "time"

    _ "github.com/go-sql-driver/mysql"

    "github.com/Tomilla/imagespider/glog"
    "github.com/Tomilla/imagespider/model"
    "github.com/Tomilla/imagespider/util"
)

var (
    C       *Config
    DB      *sql.DB
    BaseDir string
)

const NilParser = "NilParser"

type Config struct {
    sync.RWMutex
    Image        ImageConfig `json:"image"`
    Init         Init        `json:"init"`
    Log          Log         `json:"log"`
    Db           DBConfig    `json:"db"`
    Net          Net         `json:"net"`
    MameLenLimit int         `json:"nameLenLimit"`
    Engine       Engine      `json:"engine"`
    Elastic      Elastic
}

type InRange struct {
    min int
    max int
}

func (r InRange) GetRange() (int, int) {
    return r.min, r.max
}

func NewConfig() *Config {
    return LoadConfig()
}

func (c *Config) SetElasticChan(ch chan model.Topic) {
    c.Lock()
    defer c.Unlock()
    c.Elastic.topicChan = ch
}

func (c *Config) GetElasticChan() chan model.Topic {
    c.Lock()
    defer c.Unlock()
    return c.Elastic.topicChan
}

func (c *Config) GetNameLenLimit() int {
    c.RLock()
    defer c.RUnlock()

    return c.MameLenLimit
}

func (c *Config) GetSleepRange() (int, int) {
    r := c.Init.SleepRange
    return r[0], r[1]
}

func (c *Config) GetReplyRange() (int, int) {
    r := c.Init.ReplyRange
    return r[0], r[1]
}

func (c *Config) GetPageLimit() (int, int) {
    c.RLock()
    defer c.RUnlock()

    r := c.Init.PageRange
    return r[0], r[1]
}

func (c *Config) GetImageWorkerCount() int {
    return c.Image.WorkerCount
}

func (c *Config) GetStartPages() []string {
    c.RLock()
    defer c.RUnlock()

    return c.Init.Realms
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

func (c *Config) SetImageHungryChan(ch chan bool) {
    c.Lock()
    defer c.Unlock()

    c.Image.HungryChan = ch
}

func (c *Config) GetImageHungryChan() chan bool {
    c.Lock()
    defer c.Unlock()

    return c.Image.HungryChan
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

func getConfigFileName() string {
    var (
        wd  string
        err error
    )
    if BaseDir != "" {
        wd = BaseDir
    } else {
        wd, err = os.Getwd()
        if err != nil {
            panic(err)
        }
    }
    // load local config first

    var finalPath string
    finalPath = path.Join(wd, "local_config.json")
    if glog.CheckPathExists(finalPath) {
        return finalPath
    } else {
        finalPath = path.Join(wd, "config.json")
        if !glog.CheckPathExists(finalPath) {
            log.Panicln("cannot find config file")
        }
        return finalPath
    }
}

func (c *Config) GetProxyURL() string {
    return c.Net.ProxyURL
}

func (c *Config) GetNetTimeOut() int {
    return c.Net.TimeOut

}

func (c *Config) GetLogPath() string {
    return c.Log.Path
}

func (c *Config) GetDbEngine() string {
    return c.Db.Engine
}

func (c *Config) GetDbDSN() string {
    return c.Db.DSN
}

func (c *Config) GetDbMaxOpenConns() int {
    return c.Db.MaxOpenConns
}

func (c *Config) GetDbMaxIdleConns() int {
    return c.Db.MaxIdleConns
}

func (c *Config) GetDbConnMaxLifetime() time.Duration {
    return time.Duration(c.Db.ConnMaxLifetime)
}

// LoadConfig, load config
//noinspection GoUnhandledErrorResult
func LoadConfig() (c *Config) {

    filename := getConfigFileName()
    f, err := os.Open(filename)
    if f == nil {
        panic("illegal file")
    }
    if err != nil {
        fmt.Println(err.Error())
    }
    newReader := util.RemoveComment(f)
    defer f.Close()

    jsonParser := json.NewDecoder(newReader)
    jsonParser.Decode(&c)

    PrintConfig(c)
    return
}

func PrintConfig(cfg *Config) {
    fmt.Printf("%+v\n", cfg)
}

func init() {
    var err error
    rand.Seed(time.Now().UnixNano())

    _, _curDir, _, ok := runtime.Caller(0)
    if !ok {
        panic("No caller information")
    }
    BaseDir = path.Dir(path.Dir(_curDir))

    C = NewConfig() // default config

    DB, err = sql.Open(C.GetDbEngine(), C.GetDbDSN())
    if err != nil {
        log.Panicf("cannot connect to database, since: %v", err)
    }

    DB.SetMaxOpenConns(C.GetDbMaxOpenConns())
    DB.SetMaxIdleConns(C.GetDbMaxIdleConns())
    DB.SetConnMaxLifetime(C.GetDbConnMaxLifetime())
}

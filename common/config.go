package common

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "os"
    "path"
    "runtime"
    "strconv"
    "sync"
    "time"

    _ "github.com/go-sql-driver/mysql"

    "github.com/go-redis/redis"

    "github.com/Tomilla/imagespider/common/model"
    "github.com/Tomilla/imagespider/glog"
)

var (
    C          *Config
    DB         *sql.DB
    Redis      *redis.Client
    L          *glog.Logger
    TopicEnum  *model.TopicPersist
    BaseDir    string
    EndChan    chan os.Signal
    LocalSaved = make(map[string]string)
)

const NilParser = "NilParser"

type Config struct {
    sync.RWMutex
    Image        ImageConfig `json:"image"`
    Init         Init        `json:"init"`
    Log          Log         `json:"log"`
    Db           DBConfig    `json:"db"`
    Redis        RedisConfig `json:"redis"`
    Net          Net         `json:"net"`
    Limit        Limit       `json:"limit"`
    MameLenLimit int         `json:"nameLenLimit"`
    Engine       Engine      `json:"engine"`
    Elastic      Elastic
}

type PostStatusType uint

const (
    PostAdded PostStatusType = iota
    PostContentFailParsed
    PostImgFailParsed
    PostImgPartParsed
    PostImgAllParsed
    PostImgFailDownloaded
    PostImgPartDownloaded
    PostImgAllDownloaded
    PostDone
)

var postStatusStrings = []string{
    "Post was added",
    "The content of post failed to parse",
    "The images within post failed to parse",
    "Some images within post failed to parse",
    "All images within post have been parsed",
    "The images within post failed to download",
    "Some images within post failed to download",
    "All images within post have been downloaded",
    "Post Persist was Done",
}

type PostEnum interface {
    Name() string
    Ordinal() int
    String() string
    Values() *[]string
}

func (pst PostStatusType) Name() string {
    return postStatusStrings[pst]
}

func (pst PostStatusType) String() string {
    return strconv.Itoa(pst.Ordinal())
}

func (pst PostStatusType) Ordinal() int {
    return int(pst)
}

func (pst PostStatusType) Values() *[]string {
    return &postStatusStrings
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

func (c *Config) GetPostNameLenLimit() int {
    return c.Limit.PostNameLenLimit
}

func (c *Config) GetImagePathLenLimit() int {
    return c.Limit.ImagePathLenLimit
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

func (c *Config) GetLimitConfig() *Limit {
    c.RLock()
    defer c.RUnlock()
    return &c.Limit
}

func (c *Config) SetLimitConfig(postNameLenLimit int, imgPathLenLimit int) {
    c.RLock()
    defer c.RUnlock()
    c.Limit.PostNameLenLimit = postNameLenLimit
    c.Limit.ImagePathLenLimit = imgPathLenLimit
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

func (c *Config) GetShowDownloadProgress() bool {
    return c.Log.ShowDownloadProgress
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
    newReader := RemoveComment(f)
    defer f.Close()

    jsonParser := json.NewDecoder(newReader)
    jsonParser.Decode(&c)

    PrintConfig(c)
    return
}

func PrintConfig(cfg *Config) {
    fmt.Printf("%+v\n", cfg)
}

func (c Config) NewRedisClient() *redis.Client {
    addr := c.Redis.Host + ":" + strconv.Itoa(c.Redis.Port)
    client := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: c.Redis.Password,
        DB:       c.Redis.DB, // use default DB
    })

    pong, err := client.Ping().Result()
    if err != nil {
        return nil
    }
    fmt.Println(pong, err)
    // Output: PONG <nil>
    return client
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

    L = glog.NewStdoutLogger()
    Redis = C.NewRedisClient()
    TopicEnum = &model.TopicPersist{
        CountReply:           "cntReply",
        CountImage:           "cntImage",
        CountDownloadedImage: "cntImgDL",
        Status:               "status",
        Name:                 "name",
        Key:                  "key",
        Url:                  "url",
        FailedImages:         "urlFailedImages",
    }
    EndChan = make(chan os.Signal, 1)
}

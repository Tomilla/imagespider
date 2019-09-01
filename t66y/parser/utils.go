package parser

import (
    "bytes"
    "fmt"
    netUrl "net/url"
    "os"
    "path"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"

    "github.com/PuerkitoBio/goquery"
    "github.com/go-redis/redis"

    "github.com/Tomilla/imagespider/common"
    "github.com/Tomilla/imagespider/engine"
    "github.com/Tomilla/imagespider/glog"
    "github.com/Tomilla/imagespider/util"
)

func NormalizeName(s string) string {
    // s = strings.Trim(s, "[]")
    // fmt.Println("before -- > ", s)
    limit := common.C.GetPostNameLenLimit()
    // remove quote
    s = quoteRe.ReplaceAllString(s, `$1`)
    // remove punctuation
    s = punctuationRe.ReplaceAllString(s, `$1`)
    // replace full characters to half characters
    s = halfToFullRe.ReplaceAllString(s, `[$1]`)
    if strings.Contains(s, `/`) { // 去除名字中的反斜杠
        s = strings.Replace(s, `/`, ``, -1)
    }
    // tags := tagBracketRe.FindAllString(s, -1)
    // fmt.Printf("tags: %v\n", tags)
    s = tagBracketRe.ReplaceAllString(s, "")
    s = nonTagBracketRe.ReplaceAllString(s, "$1")
    // replace non Chinese or non English world
    s = nonChineseEnglishRe.ReplaceAllString(s, " ")
    s = whiteSpaceInsideRe.ReplaceAllString(s, "$1 _ $2")
    s = whiteSpaceRe.ReplaceAllString(s, "")

    characters := []rune(s)

    if len(characters) > limit {
        characters = characters[:limit]
    }
    return string(characters)
}

func ExtractImageUrls(content interface{}, matches chan string) {
    switch content.(type) {
    case *goquery.Selection:
        content := content.(*goquery.Selection)
        found := content.Find("input[type=image]")
        for _, node := range found.Nodes {
            inputTag := goquery.NewDocumentFromNode(node)
            val, exist := inputTag.Attr("data-src")
            if !exist {
                val, exist = inputTag.Attr("src")
                if !exist {
                    html, err := inputTag.Html()
                    if err != nil {
                        common.L.Debug("It's not a valid image tag.")
                    } else {
                        common.L.Debugf("It's not a valid image tag.\n", html)
                    }
                }
            }
            if viidiiImageRe.MatchString(val) {
                continue
            }
            if postImageRe.MatchString(val) {
                matches <- val
            }
        }
        break
    case string:
        content := content.(string)
        doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
        if err != nil {
            common.L.Debugf("Cannot create document from content: %v\n", content)
        }
        ExtractImageUrls(doc, matches)
        break
    case []byte:
        content := content.([]byte)
        doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
        if err != nil {
            common.L.Debugf("Cannot create document from content: %v\n", content)
        }
        ExtractImageUrls(doc, matches)
    default:
        break
    }
    close(matches)
}

func RemoveExt(p string) string {
    for i := len(p) - 1; i >= 0 && p[i] != '/'; i-- {
        if p[i] == '.' {
            return p[:i]
        }
    }
    return p
}

// Notice: the htm_data(for web) or htm_mob(for mobile) will
// move the above prefix and replace '\' with "_"
// if includeExt is false, the extension name will be removed
// for example: htm_data/1980/16/2401237.html -> 1980_16_2401237
func NormalizePostUrl(url string, includeExt bool) string {
    var _path string
    if strings.HasPrefix(url, "/") {
        _path = url
    } else {
        u, err := netUrl.Parse(url)
        if err != nil {
            return url
        }
        _path = u.Path
    }
    if !includeExt {
        _path = RemoveExt(_path)
    }

    return strings.Trim(postPathRe.ReplaceAllString(_path, "_"), "_")
}

func SimplifyPostUrl(url string) string {
    return NormalizePostUrl(url, false)
}

func SyncRedisAndLocalDatum(attr map[string]string, p engine.Post) bool {
    // attr from redis
    attrCntReply, ok := attr[common.TopicEnum.CountReply]
    if !ok {
        return false
    }
    title, ok := attr[common.TopicEnum.Name]
    if !ok {
        return false
    }
    cntReply, err := strconv.Atoi(attrCntReply)
    if err != nil {
        return false
    }
    normalTitle := NormalizeName(title)
    if cntReply < p.CountReply {
        dirNameOld, ok := common.LocalSaved[normalTitle]
        if !ok {
            dirNameOld = fmt.Sprintf(FileNameFormat, cntReply, p.CountImage, NormalizeName(title))
        }
        dirNameNew := fmt.Sprintf(FileNameFormat, p.CountReply, p.CountImage, normalTitle)
        src := path.Join(ImageDir, dirNameOld)
        dest := path.Join(ImageDir, dirNameNew)
        if glog.CheckPathExists(src) {
            err = os.Rename(src, dest)
            common.L.Infof("[Rename]: from %v to %v", src, dest)
            if err != nil {
                return false
            }
        }
    }
    return true
}

func LoadLocalSavedPosts(_path string) error {
    replyAndImagePrefix := regexp.MustCompile(`^\[\d+\w]\[\d+\w]`)
    err := filepath.Walk(_path, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            return nil
        }
        fileName := info.Name()
        if replyAndImagePrefix.MatchString(fileName) {
            fileNameNew := replyAndImagePrefix.ReplaceAllString(fileName, "")
            common.LocalSaved[fileNameNew] = fileName
            common.L.Info(fileNameNew)
        }
        return nil
    })
    if err != nil {
        return err
    }

    return nil
}

func GetLocalArchivedPostHTMLs(_path string) error {
    err := filepath.Walk(_path, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        fileName := info.Name()
        common.L.Info(fileName)
        name, ok := util.GetRegexNamedGroupMapping(postArchiveRe, fileName)["Name"]
        if ok {
            common.Redis.HSet(name, common.TopicEnum.Status, common.PostDone.Ordinal())
        }
        return nil
    })
    if err != nil {
        return err
    }
    return nil
}

func repairRedisKeys(keyPattern string, f func(string) string) error {
    keys, err := common.Redis.Keys("19*_*_*").Result()
    if err != nil {
        return err
    }
    for _, key := range keys {
        newKey := f(key)
        _, err := common.Redis.Rename(key, newKey).Result()
        if err != nil {
            if err == redis.Nil {
                common.L.Errorf("non-exist key: %v", key)
            } else {
                common.L.Errorf("rename failed(%v -> %v): %v", key, newKey, err)
            }
        }
    }
    return nil
}

func init() {
    // load archive from local
    // err := GetLocalArchivedPostHTMLs(common.C.GetLogPath())
    // common.L.Info(err)
    // yearMonth := regexp.MustCompile(`^\d{4}_`)

    // fix key format compatibility
    // err := repairRedisKeys("19*_*_*", func (key string) string {
    //     return yearMonth.ReplaceAllString(key, "")
    // })
    // if err != nil {
    //     return
    // }

    err := LoadLocalSavedPosts(common.C.GetImageConfig().Path)
    if err != nil {
        common.L.Info(err)
    }
}

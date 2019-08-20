package parser

import (
    "bytes"
    netUrl "net/url"
    "os"
    "path/filepath"
    "strings"

    "github.com/PuerkitoBio/goquery"

    "github.com/Tomilla/imagespider/common"
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

func GetLocalArchivedPosts(_path string) error {
    err := filepath.Walk(_path, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return filepath.SkipDir
        }
        fileName := info.Name()
        postArchiveRe.MatchString(fileName)

        return nil
    })
    if err != nil {
        return err
    }
    return nil
}

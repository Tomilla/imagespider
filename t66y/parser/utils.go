package parser

import (
    "bytes"
    "strings"

    "github.com/PuerkitoBio/goquery"

    "github.com/Tomilla/imagespider/config"
)

func NormalizeName(s string) string {
    // s = strings.Trim(s, "[]")
    // fmt.Println("before -- > ", s)
    limit := config.C.GetPostNameLenLimit()
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
                        config.L.Debug("It's not a valid image tag.")
                    } else {
                        config.L.Debugf("It's not a valid image tag.\n", html)
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
            config.L.Debugf("Cannot create document from content: %v\n", content)
        }
        ExtractImageUrls(doc, matches)
        break
    case []byte:
        content := content.([]byte)
        doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
        if err != nil {
            config.L.Debugf("Cannot create document from content: %v\n", content)
        }
        ExtractImageUrls(doc, matches)
    default:
        break
    }
    close(matches)
}

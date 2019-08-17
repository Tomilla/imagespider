package parser

import (
    "fmt"
    "regexp"

    "github.com/Tomilla/imagespider/collections/set"
)

// from post_list
var (
    postUrlRe        = regexp.MustCompile(`htm_data/\d+/\d+/\d+\.html`)
    countImageRe     = regexp.MustCompile(`[\[［](\d+).?[］\]]`)
    ignoredPostColor = set.New("red", "blue", "orange")
    includePostColor = set.New("green")
)

// from post
var (
    // template
    tmpChineseEnglish = `[%v\p{Han}\p{Latin}0-9_-]`
    chineseEnglish    = fmt.Sprintf(tmpChineseEnglish, ``)
    nonChineseEnglish = fmt.Sprintf(tmpChineseEnglish, `^`)
    // imageRe           = regexp.MustCompile(`(?i)(data-src|data-link|src)=['"](http[s]?://[^'"]+(jpg|png|jpeg|gif))['"]`)
    // titleRe             = regexp.MustCompile(`<title>([^>]+)(\s+-\s*\S+\s*\|\s*\S+\s*-\s*\S+\s*)</title>`)
    postImageRe         = regexp.MustCompile(`(?:http|https):[^\s]*?(?:jpg|jpeg|png|gif)`)
    viidiiImageRe       = regexp.MustCompile(`(?:http|https)://[^\s]+viidii[^\s]+(?:jpg|jpeg|png|gif)`)
    quoteRe             = regexp.MustCompile(`(?:['"‘“])(.*?)(?:['"’”])`)
    punctuationRe       = regexp.MustCompile(`(?:[(（{｛])(.*?)(?:\[\)）}｝])`)
    halfToFullRe        = regexp.MustCompile(`(?:[［「【『〖])(.*?)(?:[］」】』〗])`)
    tagBracketRe        = regexp.MustCompile(`\[.{0,6}]`)
    nonTagBracketRe     = regexp.MustCompile(`\[(.*?)]`)
    whiteSpaceRe        = regexp.MustCompile(`\s+`)
    whiteSpaceInsideRe  = regexp.MustCompile(fmt.Sprintf(`(%v)\s+(%v)`, chineseEnglish, chineseEnglish))
    nonChineseEnglishRe = regexp.MustCompile(fmt.Sprintf(`%v+`, nonChineseEnglish))
    PathSeparatorRe     = regexp.MustCompile(`[\\/]`)
)

const (
    HOSTNAME = "http://t66y.com/"
)

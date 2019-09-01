package parser

import (
    "fmt"
    "regexp"

    "github.com/Tomilla/imagespider/collections/set"
    "github.com/Tomilla/imagespider/common"
)

// from post_list
var (
    postUrlRe        = regexp.MustCompile(`htm_data/\d+/\d+/\d+\.html`)
    countImageRe     = regexp.MustCompile(`[\[［](\d+).?[］\]]`)
    ignoredPostColor = set.New("red", "blue", "orange")
    includePostColor = set.New("green")
    ValidWebPageExt  = set.New("html", "htm")
    FinishPostStatus = set.New(common.PostImgAllDownloaded.String(), common.PostDone.String())
    ImageDir         = common.C.GetImageConfig().Path
)

// from post
var (
    // template
    tmpChineseEnglish = `[%v\p{Han}\p{Latin}0-9_-]`
    chineseEnglish    = fmt.Sprintf(tmpChineseEnglish, ``)
    nonChineseEnglish = fmt.Sprintf(tmpChineseEnglish, `^`)
    // imageRe           = regexp.MustCompile(`(?i)(data-src|data-link|src)=['"](http[s]?://[^'"]+(jpg|png|jpeg|gif))['"]`)
    // titleRe             = regexp.MustCompile(`<title>([^>]+)(\s+-\s*\S+\s*\|\s*\S+\s*-\s*\S+\s*)</title>`)
    postPathRe    = regexp.MustCompile(`(?i)_+|(?:html?_(?:data|mob)/\d{4}|[&=?]|[/\\])`)
    postArchiveRe = regexp.MustCompile(`(?i)(?P<Name>\d+_\d+_\d+)(?P<Ext>\.[a-z]+)?`)
    postImageRe   = regexp.MustCompile(`(?i)(?:http|https):[^\s]*?(?:jpg|jpeg|png|gif)`)
    viidiiImageRe = regexp.MustCompile(`(?i)(?:http|https)://[^\s]+viidii[^\s]+(?:jpg|jpeg|png|gif)`)
    quoteRe       = regexp.MustCompile(`(?:['"‘“])(.*?)(?:['"’”])`)
    punctuationRe = regexp.MustCompile(`(?:[(（{｛])(.*?)(?:\[\)）}｝])`)
    halfToFullRe  = regexp.MustCompile(`(?:[［「【『〖])(.*?)(?:[］」】』〗])`)
    tagBracketRe    = regexp.MustCompile(`\[.{0,6}]`)
    nonTagBracketRe = regexp.MustCompile(`\[(.*?)]`)
    whiteSpaceRe        = regexp.MustCompile(`\s+`)
    whiteSpaceInsideRe  = regexp.MustCompile(fmt.Sprintf(`(%v)\s+(%v)`, chineseEnglish, chineseEnglish))
    nonChineseEnglishRe = regexp.MustCompile(fmt.Sprintf(`%v+`, nonChineseEnglish))
    // PathSeparatorRe     = regexp.MustCompile(`[\\/]`)
)

const (
    HOSTNAME          = "http://t66y.com/"
    DefaultWebPageExt = ".html"
    FileNameFormat    = "[%03vR][%03vP]%v" // reply, picture_count, and title
)

package image

import (
    "testing"

    "github.com/stretchr/testify/assert"

    "github.com/Tomilla/imagespider/common"
)

func TestDownloader_GetFileName(t *testing.T) {
    limit := common.C.GetLimitConfig()
    oldPostNameLenLimit := limit.PostNameLenLimit
    oldImgPathLenLimit := limit.ImagePathLenLimit

    common.C.SetLimitConfig(60, 25)
    d := NewDownloader(common.C.GetImageConfig(), common.C.GetLimitConfig())
    var sourceAndExcept = [][]string{
        {"https://www.skeimg.com/u/20190803/15403922.jpg", "00年小母狗还需要一些K9 ", "000_15403922.jpg"},
        {"https://www.touimg.com/u/20190727/06262236.gif", "熟女陪我666", "001_06262236.gif"},
        {"https://www.privacypic.com/images/2019/07/29/IMG_20181024_1341321232c8813f705aae.jpg",
            "偷情少妇_有腰窝_大屁股_胸大活好不粘人", "002_341321232c8813f705aae.jpg"},
        {"https://www.privacypic.com/images/2019/07/29/IMG_20181024_134135dfe0291a8e570962.jpg",
            "偷情少妇_有腰窝_大屁股_胸大活好不粘人", "003_34135dfe0291a8e570962.jpg"},
        {"https://www.privacypic.com/images/2019/07/29/IMG_20181024_1341321232c8813f705aae.jpg",
            "偷情少妇_有腰窝_大屁股_胸大活好不粘人", "004_341321232c8813f705aae.jpg"},
        {"https://www.privacypic.com/images/2019/07/29/IMG_20181024_12431629373124823lj2lj24jl2j4l.jpg",
            "偷情少妇_有腰窝_大屁股_胸大活好不粘人", "005_73124823lj2lj24jl2j4l.jpg"},
    }
    for i, item := range sourceAndExcept {
        src, name, except := item[0], item[1], item[2]
        dest := d.GetFileName("", name, src, i)
        assert.Equal(t, dest, except)
    }

    // restore config
    common.C.SetLimitConfig(oldPostNameLenLimit, oldImgPathLenLimit)
}

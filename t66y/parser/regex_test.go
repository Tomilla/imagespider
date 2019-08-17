package parser

import (
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestNormalizeName(t *testing.T) {
    var (
        expected = [][]string{
            {"[原创][啪照出品]00年小母狗'还需要'一些K9[17P]",
                "00年小母狗还需要一些K9"},
            {"[原创][98A]熟女陪我\"666\"[60P]",
                "熟女陪我666"},
            {"[原创]偷情少妇，有‘腰窝’，大‘屁股’，胸大活好不粘人[27P]",
                "偷情少妇_有腰窝_大屁股_胸大活好不粘人"},
            {"[原创]98年的“小处女”［13P］",
                "98年的小处女"},
            {"[原创投稿][画家'洋洋']一线天画家洋洋晚间\"野外\"外拍，与向日葵共舞[露脸][50P]",
                "一线天画家洋洋晚间野外外拍_与向日葵共舞"},
            {"[原创] 开裆‘黑丝’剃逼毛，没毛的“小母狗”似乎更加淫荡了。[42P]",
                "开裆黑丝剃逼毛_没毛的小母狗似乎更加淫荡了"},
            {"蜜丝的原创，路灯下捆绑露出，骚蜜丝穿着丝袜高跟，全裸被完全绑在路灯下，虽然是晚上，但远处的汽车灯光还是让人心跳不已[16P]",
                "蜜丝的原创_路灯下捆绑露出_骚蜜丝穿着丝袜高跟_全裸被完全绑在路灯下_虽然是晚上_但远处的汽车灯光还是让人心跳不已"},
            {"[原创投稿][白袜袜格罗丫][嫩逼少女格罗丫挺着傲人的大乳说想你们啦，防毒面具是她专属的标识][23P]",
                "嫩逼少女格罗丫挺着傲人的大乳说想你们啦_防毒面具是她专属的标识"},
        }
    )
    for _, exp := range expected {
        src, dest := exp[0], exp[1]
        replaced := NormalizeName(src)
        // t.Log(src)
        // t.Log(dest)
        // t.Log(replaced)
        assert.Equal(t, dest, replaced)
    }
}

func TestCleanPathRe(t *testing.T) {
    var (
        expected = [][]string{
            {"/htm_data/1908/16/3604380.html", "1908_16_3604380.html"},
            {"/htm_data/1908/16/3597335.html", "1908_16_3597335.html"},
            {"/htm_data/1908/16/3598134.html", "1908_16_3598134.html"},
            {"/htm_data/1908/16/3618060.html", "1908_16_3618060.html"},
            {"/htm_data/1908/15/3618059.html", "1908_15_3618059.html"},
        }
    )

    for _, exp := range expected {
        src, dest := exp[0], exp[1]
        replaced := strings.Trim(postPathRe.ReplaceAllString(src, "_"), "_")
        assert.Equal(t, dest, replaced)
    }
}

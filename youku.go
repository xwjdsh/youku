package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	address  = flag.String("a", "http://v.youku.com/v_show/id_XNTAzMDQ5NTM2.html", "vedio address")
	format   = flag.String("f", "bz", "vedio format")
	isMore   = flag.Bool("m", false, "is play list")
	savePath = flag.String("s", ".", "save path")
)
var nameMap = map[string]string{
	"标准": "bz",
	"高清": "gq",
	"超清": "cq",
}

var checkMap = map[string]int{
	"bz": 3,
	"gq": 2,
	"cq": 1,
}

var orderMap = map[int]string{
	3: "bz",
	2: "gq",
	1: "cq",
}

func main() {

	flag.Parse()

	if *address == "" {
		panic("必须设置视频路径!")
	}

	if _, ok := checkMap[*format]; !ok {
		panic("视频格式只能设置为'bz'(标准)/'gq'(高清)/'cq'(超清)")
	}

	if *savePath == "" {
		panic("必须设置保存路径!")
	}
	if _, err := os.Stat(*savePath); err != nil && os.MkdirAll(*savePath, os.ModePerm) != nil {
		panic("指定保存路径错误!")
	}

	if *isMore {
		//videos, err := getPlayList()
		//if err != nil {
		//fmt.Println(err.Error())
		//return
		//}
		//for i, video := range videos {

		//}

	} else {
		links, name, err := getLinksAndName(*address)
		fmt.Println(links, name)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	}

}

func getPlayList() ([]string, error) {
	doc, err := goquery.NewDocument(*address)
	if err != nil {
		return nil, err
	}
	playlist := doc.Find("div .items .item .sn").Map(func(i int, s *goquery.Selection) string {
		return s.AttrOr("href", "")
	})
	if len(playlist) == 0 {
		return nil, errors.New("wrong address,can't get the vedio list!")
	}
	return playlist, nil

}

func getLinksAndName(address string) ([]string, string, error) {
	formatMap := make(map[string][]string)

	params := url.Values{}
	params.Add("url", address)
	resp, err := http.PostForm("http://www.shokdown.com/parse.php", params)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "", err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		fmt.Println(err.Error())
		return nil, "", err
	}
	linkAddrs := []string{}
	doc.Find("#main a").Each(func(i int, s *goquery.Selection) {
		linkAddr := s.AttrOr("href", "")
		if strings.HasPrefix(linkAddr, "http") {
			linkAddrs = append(linkAddrs, linkAddr)
		}
	})
	if len(linkAddrs) == 0 {
		return nil, "", errors.New("不能解析出下载地址!")
	}

	index := 0

	doc.Find("#main font[color='red']").Each(func(i int, s *goquery.Selection) {
		if v, ok := nameMap[s.Text()]; ok {
			num, _ := strconv.Atoi(s.Next().Text())
			if index < len(linkAddrs) {
				formatMap[v] = linkAddrs[index : index+num]
				index += num
			}
		}
	})

	name := doc.Find("td[align='left']").First().Map(func(i int, s *goquery.Selection) string {
		s.Find("b").Remove()
		return strings.TrimSpace(s.Text())
	})

	for i := checkMap[*format]; i > 0; i-- {
		*format = orderMap[i]
		if _, ok := formatMap[*format]; ok {
			break
		}
	}
	return formatMap[*format], name[0], nil
}

func download(path string) error {
	links, name, err := getLinksAndName(path)
	if err != nil {
		return err
	}
	fmt.Println(links, name)
	return nil

}

package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	address = flag.String("a", "http://v.youku.com/v_show/id_XNTAzMDQ5NTM2.html", "vedio address")
	format  = flag.String("f", "bz", "vedio format")
	isMore  = flag.Bool("m", false, "is play list")
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
		panic("must set vedio download address!")
	}
	if _, ok := checkMap[*format]; !ok {
		panic("set the correct format,only 'bz'/'gq'/'cq'")
	}

	if *isMore {

	} else {

	}

	//links, name, err := getLinksAndName(*address)
	//fmt.Println(links, name)
	//if err != nil {
	//fmt.Println(err.Error())
	//return
	//}

}

func getPlayList() ([]string, error) {
	doc, err := goquery.NewDocument(*address)
	if err != nil {
		return nil, err
	}
	playlist:=doc.Find("div .items .item .sn").Map(func(i int, s *goquery.Selection) string {
		return s.AttrOr("href", "")
	})
	if len(playlist)==0{
		reutrn nil,errors.New("")
		
	}

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
			formatMap[v] = linkAddrs[index : index+num]
			index += num
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

func download(path, name string) {

}

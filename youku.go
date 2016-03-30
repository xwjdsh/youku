package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	address = flag.String("a", "http://v.youku.com/v_show/id_XNTAzMDQ5NTM2.html", "vedio address")
	format  = flag.String("f", "bz", "vedio format")
)
var nameMap = map[string]string{
	"标准": "bz",
	"高清": "gq",
	"超清": "cq",
}

func main() {

	flag.Parse()
	if *address == "" {
		panic("must set vedio download address!")
	}
	if *format != "bz" && *format != "gq" && *format != "cq" {
		panic("set the correct format,only 'bz'/'gq'/'cq'")
	}
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("can't write log to logfile!")
	}
	defer f.Close()
	log.SetOutput(f)

	formatMap := make(map[string][]string)

	params := url.Values{}
	params.Add("url", *address)
	resp, err := http.PostForm("http://www.shokdown.com/parse.php", params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	linkAddrs := []string{}
	doc.Find("#main a").Each(func(i int, s *goquery.Selection) {
		linkAddr := s.AttrOr("href", "")
		if strings.HasPrefix(linkAddr, "http") {
			linkAddrs = append(linkAddrs, linkAddr)
		}
	})
	index := 0

	doc.Find("#main font[color='red']").Each(func(i int, s *goquery.Selection) {
		if v, ok := nameMap[s.Text()]; ok {
			num, _ := strconv.Atoi(s.Next().Text())
			formatMap[v] = linkAddrs[index : index+num]
			index += num
		}
	})

	//doc.Find("td[align='left'] b").Remove()
	tmp := doc.Find("td[align='left']").First().Map(func(i int, s *goquery.Selection) string {
		s.Find("b").Remove()
		return strings.TrimSpace(s.Text())
	})
	fmt.Println(tmp)

	fmt.Println(len(linkAddrs))
	for k, v := range formatMap {
		for i, j := range v {
			if k == "gq" {
				res, _ := http.Get(j)
				file, _ := os.Create("test.flv")
				io.Copy(file, res.Body)
				fmt.Println("下载完成！")
				fmt.Println(k, i)
				return
			}
		}
	}

}

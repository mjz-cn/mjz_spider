package proxy

import (
	"fmt"
	"strings"
	"net/http"
	"image/gif"
	"image/png"
	"path"
	"os"

    "github.com/PuerkitoBio/goquery"
    "github.com/astaxie/beego/logs"

	"mjz_spider/utils"
	"mjz_spider/config"
	"mjz_spider/models"
)

var log *logs.BeeLogger

type incloakMap struct {
	TypeMap map[string]int8
	AnonMap map[string]int8
}

var inMap = incloakMap{
	TypeMap: map[string]int8{
		"h": 1,
		"2": 2,
		"3": 3,
		"4": 4,
	},
	AnonMap: map[string]int8{
		"1": 0,
		"2": 1,
		"3": 2,
		"4": 3,
	},
}

type incloak struct {
	Name string
	SiteUrl string
	StartUrl string
	DataPath string

	proxyUrl string
	client *http.Client
}

func newIncloak() SpiderHandler {
	s := &incloak{
		Name: "incloak",
		SiteUrl: "http://incloak.com/",
		StartUrl: "http://incloak.com/proxy-list",
		proxyUrl: "socks5://127.0.0.1:1080",
		DataPath: config.GlobalConfig.ExecuteDir + "/data/incloak/",
	}
	s.Init()

	return s
}

func (this *incloak) Init() {
	this.client, _ = utils.NewClient(this.proxyUrl)
	os.MkdirAll(this.DataPath, 0777)
}

func (this *incloak) ocr(portPath string) (port string, err error) {
	cmdPath := config.GlobalConfig.Tesseract
	var args = []string{
		portPath,
        "stdout",
        "digits",
	}
	port, err = utils.ExecuteCmd(cmdPath, args)
	os.Remove(portPath)
	if err != nil {
		return "", err
	}

	port = strings.Trim(port, "\n ")
	// 根据测试，Tesseract很容易把8识别成5，而port端口中，8占据大多数
	// 因此将结果中的5替换成8
	port = strings.Replace(port, "5", "8", -1)
	
	return port, nil
}

func (this *incloak) getDoc(url string) (*goquery.Document, error) {
	var (
		doc *goquery.Document
		resp *http.Response
		err error
	)
	if resp, err = this.client.Get(url); err != nil {
		return doc, err
	}
	// defer resp.Body.Close()
	if doc, err = goquery.NewDocumentFromResponse(resp); err != nil {
		return doc, err
	}
	return doc, nil
}

/**
 * 从一个url中解析出数据，并保存到数据库中
 */
func (this *incloak) parse(url, country, typ, anon string) (err error) {
	var doc *goquery.Document

	if doc, err = this.getDoc(url); err != nil {
		return err
	}

	doc.Find(".pl tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		host := s.Find(".tdl").Text()
		port_src, _ := s.Find("img").Attr("src")
		resp, err1 := this.client.Get(this.SiteUrl + port_src)
		if err1 != nil {
			log.Warn(err1.Error())
			return
		}
		defer resp.Body.Close()
		img, _ := gif.Decode(resp.Body)
		_, name := path.Split(port_src)
		port_path := this.DataPath + name + ".png"
		fd, _ := os.Create(port_path)

		if err = png.Encode(fd, img); err != nil {
			log.Warn(err.Error())
			return
		}
		port := ""
		if port, err = this.ocr(port_path); err != nil {
			log.Warn(err.Error())
			return
		}
		log.Info("%s, %s", host, port)
		models.SaveProxy(host, port, inMap.TypeMap[typ], inMap.AnonMap[anon], country)
	})
	return nil
}

/**
 * 国家，代理类型，匿名类型
 * 通过选择三个不同的参数，获取最多的proxy
 */
func (this *incloak) parseParams() (countryArr []string, err error){
	var doc *goquery.Document

	if doc, err = this.getDoc(this.StartUrl); err != nil {
		return nil, err
	}

	doc.Find("#country option").Each(func(i int, s *goquery.Selection) {
		if code, exists := s.Attr("value"); exists {
			countryArr = append(countryArr, code)
		}
	})

	return countryArr, nil
}

func (this *incloak) Run() {
	typeArr := []string{"h", "s", "4", "5"}
	anonAtrr := []string{"1", "2", "3", "4"}
	countryArr, err := this.parseParams()
	if err != nil {
		log.Warn(err.Error())
		return
	}

	url_fmt := "http://incloak.com/proxy-list/?country=%s&type=%s&anon=%s"
	for _, country := range countryArr {
		// if country[0:1] <= "M" {
		// 	continue
		// }
		for _, typ := range typeArr {
			for _, anon := range anonAtrr {
				url := fmt.Sprintf(url_fmt, country, typ, anon)
				log.Info(url)

				if err := this.parse(url, country, typ, anon); err != nil {
					log.Warn(err.Error())
				}
			}
		}
	}
}

func init() {
	log = logs.NewLogger(1000)
	log.SetLogger("console", "") 
	Register(newIncloak)
}
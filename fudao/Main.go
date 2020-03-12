package main

import (
	"database/sql"
	_ "database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	f, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("读取配置文件失败")
	}
	var config map[string]interface{}
	jsonErr := json.Unmarshal([]byte(f), &config)
	if jsonErr != nil {
		fmt.Println("转换 json 失败")
	}
	//fmt.Println(config["dbhost"])
	//fmt.Println(config["dbusername"])
	//fmt.Println(config["dbpwd"])

	db, err := sql.Open("mysql",  config["dbusername"].(string) + ":"+ config["dbpwd"].(string)  + "@/" + config["dbname"].(string))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	insSql, err := db.Prepare("INSERT INTO `fudao` (`id`, `ke_id`, `grade`, `class`, `title`, `url`, `teacher`, `price`) VALUES (NULL, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer insSql.Close()

	c := colly.NewCollector()
	colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36")



	c2 := c.Clone()
	c3 := c.Clone()
	// 抓取系统课数据
	// 系统课有 n 个子课程组成
	c4 := c.Clone()
	// 单课程页面信息抓取
	c5 := c.Clone()

	// 找到所有年级的 url
	c.OnHTML("div[class=filter-content-wrap]>div[class=grade-content]>div[class=grade-area]>a", func(element *colly.HTMLElement) {
		url := element.Request.AbsoluteURL(element.Attr("href"))
		// 年级
		element.Response.Ctx.Put("grade", element.Text)
		_ = c2.Request("GET", url, nil, element.Response.Ctx, nil)
	})

	// 找到每个年级的所有课程学科
	c2.OnHTML("ul[class=subject-list]>li[class=subject-item]>a", func(element *colly.HTMLElement) {

		if element.Text == "全部" {
			return
		}

		url := element.Request.AbsoluteURL(element.Attr("href"))
		// 学科
		element.Response.Ctx.Put("class", element.Text)
		fmt.Println(url)
		_ = c3.Request("GET", url, nil, element.Response.Ctx, nil)
	})

	// 每个学科页面
	// 年级,学科,url
	c3.OnHTML("section[class=subject-page]", func(element *colly.HTMLElement) {

		// 系统课
		subSystem := element.DOM.Find("section[class=sub-system-ctn]>ul>li")
		systemUrl := subSystem.Find("a")
		url,exists := systemUrl.Attr("href")
		if exists == true {
			// 系统课主课程标题
			title := subSystem.Find("div[class=subject-course--content]>h2").Text()
			element.Response.Ctx.Put("title", title)
			_ = c4.Request("GET", url, nil, element.Response.Ctx, nil)
		}

		// 专题课
		subSubject := element.DOM.Find("section[class=sub-subject-ctn]>ul>li")
		subjectUrl := subSubject.Find("a")
		url,exists = subjectUrl.Attr("href")
		if exists == true {
			id := strings.Replace(url, "/pc/course.html?course_id=", "", -1)
			element.Response.Ctx.Put("id", id)
			url := element.Request.AbsoluteURL(element.Attr("href"))
			_ = c5.Request("GET", url, nil, element.Response.Ctx, nil)
		}
	})

	// 系统课引导页
	// https://fudao.qq.com/subject/6005/subject_system/str_sys_course_pkg_info_1_6005_8_0
	c4.OnHTML("div[class=sys-pkg-ct]>li>a", func(element *colly.HTMLElement) {
		// 课程 url
		url := element.Request.AbsoluteURL(element.Attr("href"))

		// 课程 id
		id := strings.Replace(element.Attr("href"), "/pc/course.html?course_id=", "", -1)
		element.Response.Ctx.Put("id", id)

		_ = c5.Request("GET", url, nil, element.Response.Ctx, nil)
	})

	// 课程详情页
	c5.OnHTML("div[id=react-body]", func(element *colly.HTMLElement) {
		// title
		title := element.DOM.Find("div[class=fixed-title]>h1>span[class=tt-word]").Text()
		fmt.Println("课程标题:" ,title)
		// teacher
		teacherName := element.DOM.Find("ul[class=teacherList]>li:nth-child(1)>div[class=teacherContent]>div[class=caption]>p").Text()
		teacherName = strings.ReplaceAll(teacherName, "授课老师：", "")
		fmt.Println("教师名称:", teacherName)
		// id
		id := element.Response.Ctx.Get("id")
		fmt.Println("id:", id)
		// price
		price := element.DOM.Find("div[class=tt-price-wrap]").Text()
		price = strings.ReplaceAll(price, "¥", "")
		priceInt := 0
		if price == "免费" {
			priceInt = 0
		}else{
			priceInt, _ = strconv.Atoi(price)
		}
		fmt.Println("价格:", priceInt)
		// 年级
		grade := element.Response.Ctx.Get("grade")
		fmt.Println("年级",grade)
		// 学科
		class := element.Response.Ctx.Get("class")
		fmt.Println("学科",class)

		_, err = insSql.Exec(id, grade, class, title, element.Request.URL.String(), teacherName, priceInt)
		if err != nil {
			panic(err.Error())
		}
	})


	_ = c.Visit("https://fudao.qq.com/")
}
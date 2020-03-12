package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
)

func getDbData() string {
	db, err := sql.Open("mysql", "root:123@/fudao")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM fudao")
	if err != nil {
		panic(err.Error())
	}

	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}


	var group []map[string]interface{}

	for rows.Next() {

		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		var value string
		m := make(map[string]interface{})
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			m[columns[i]] = value

			fmt.Println(columns[i], ": ", value)
		}
		group = append(group, m)
		fmt.Println("-----------------------------------")
	}
	if err = rows.Err(); err != nil {
		panic(err.Error())
	}

	return listToJSON(group)
}

func listToJSON(s []map[string]interface{}) string {
	data, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	// fmt.Println(data)
	return string(data)
}


func handler1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json") //返回数据格式是json

	fmt.Println("接受请求")

	s, _ := ioutil.ReadAll(r.Body) //把  body 内容读入字符串 s
	fmt.Println(string(s)) //在返回页面中显示内容。

	_, _ = fmt.Fprintf(w, getDbData())
}

func main() {
	ports := []string{":8887"}
	for _, v := range ports {
		go func(port string) {
			mux := http.NewServeMux()
			mux.HandleFunc("/", handler1)
			_ = http.ListenAndServe(port, mux)
		}(v)
	}
	select {}
}

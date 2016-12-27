package db

import (
	"github.com/astaxie/beego"
	"fmt"
	"github.com/xormplus/xorm"
	"strings"
	"encoding/json"
	"strconv"
	"reflect"
)

type Db struct {
	Db            *xorm.Engine
	FilterSession *xorm.Session
}

var DB *Db
//
func dsn() string {
	db_user := beego.AppConfig.String("db_user")
	db_pass := beego.AppConfig.String("db_pass")
	db_host := beego.AppConfig.String("db_host")
	db_port := beego.AppConfig.String("db_port")
	db_name := beego.AppConfig.String("db_name")
	dsn := db_user + ":" + db_pass + "@tcp(" + db_host + ":" + db_port + ")/" + db_name + "?charset=utf8"
	return dsn
}
//初始化
func Init() {
	var err error
	DB = new(Db)
	DB.Db, err = xorm.NewEngine("mysql", dsn())
	if err != nil {
		fmt.Println("NewEngine", err)
	}
	DB.Db.ShowSQL(true)
}
// NewDb create new db
func NewDb() *xorm.Engine {
	return DB.Db
}

type QuerySession struct {
	Session *xorm.Session
}

var Query *QuerySession

func Filter(where map[string]interface{}) *xorm.Session {
	db := DB.Db
	Query=new(QuerySession)
	i := 1
	for k, v := range where {
		//fmt.Println(k, v, reflect.TypeOf(v))
		//fmt.Println("?号个数为", strings.Count(k, "?"))
		QuestionMarkCount := strings.Count(k, "?")
		isEmpty := false
		isMap := false
		arrCount := 0
		str := ""
		var arr []string
		switch v.(type) {
		case string:
			//是字符时做的事情
			isEmpty = v == ""
		case int:

		//是整数时做的事情
		case []string :
			isMap = true
			arr = v.([]string)
			arrCount = len(arr)
			isEmpty = arrCount == 0
			for j, val := range arr {
				if j > 0 {
					str += ","
				}
				str += val
			}
		case []int :
			isMap = true
			arrInt := v.([]int)
			arrCount = len(arrInt)
			isEmpty = arrCount == 0
			for j, val := range arrInt {
				if j > 0 {
					str += ","
				}
				str += strconv.Itoa(val)
			}
		}
		if QuestionMarkCount == 0 && isEmpty {
			FilterWhereAnd(db, i, k, "")
		} else if QuestionMarkCount == 0 && !isEmpty {
			//是数组
			if (isMap) {

				FilterWhereAnd(db, i, k, str)
			} else {
				//不是数组
				FilterWhereAnd(db, i, k + " = ?", v)
			}
		} else if QuestionMarkCount == 1 && isEmpty {
			//值为空字符串,不是数组
			FilterWhereAnd(db, i, k, "''")
		} else if QuestionMarkCount == 1 && !isEmpty {
			//是数组
			if isMap {
				fmt.Println("ArrToStr_key", k)
				fmt.Println("ArrToStr", str)
				if arrCount > 1 {
					new_q := ""
					for z := 1; z <= arrCount; z++ {
						if z > 1 {
							new_q += ","
						}
						new_q += "?"
					}
					str2 := strings.Replace(k, "?", new_q, -1)
					fmt.Println("ArrToStr", str)
					fmt.Println("arr", arr)
					//var inter =arr
					inter := make([]interface{}, arrCount)
					for y, x := range arr {
						inter[y] = x
					}
					FilterWhereAnd(db, i, str2, inter...)
				} else {
					fmt.Println("22222", str)
					FilterWhereAnd(db, i, k, str)
				}

			} else {
				//不是数组
				//不是数组，有值
				FilterWhereAnd(db, i, k, v)
			}
		} else if QuestionMarkCount > 1 && isEmpty {
			//不是数组，空值
			FilterWhereAnd(db, i, k, "")
		} else if QuestionMarkCount > 1 && !isEmpty && isMap {
			//问号 与  数组相同时
			if QuestionMarkCount == arrCount {
				//不是数组
				FilterWhereAnd(db, i, k, v)
			} else {
				//问号 与  数组不同时
				FilterWhereAnd(db, i, k, str)
			}
		} else {
			fmt.Println("其他还没有收录")
		}
		i++
	}
	return Query.Session
}
func FilterWhereAnd(db *xorm.Engine, i int, key string, value ...interface{}) {
	fmt.Println("key", key)
	fmt.Println("value", value)
	fmt.Println("TypeOf", reflect.TypeOf(value))
	if i == 1 {

		Query.Session = db.Where(key, value...)
	} else {
		Query.Session = Query.Session.And(key, value...)
	}
}
func ArrToStr(val interface{}) string {
	str, _ := json.Marshal(val)
	str2 := string(str)
	fmt.Println("ArrToStr", str2)
	str2 = strings.Replace(str2, "{", "", -1)
	str2 = strings.Replace(str2, "}", "", -1)
	return str2
}
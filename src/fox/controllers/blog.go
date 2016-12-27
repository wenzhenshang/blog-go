package controllers

import (
	"errors"
	"strconv"
	"strings"

	"regexp"
	"fox/service/blog"
	"fmt"
)

type BlogController struct {
	BaseNoLoginController
}


// GetOne ...
// @Title Get One
// @Description get Blog by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Blog
// @Failure 403 :id is empty
// @router /article/:id [get]
func (c *BlogController) Get() {
	idStr := c.Ctx.Input.Param(":id")
	var ser *blog.Blog
	var err error
	var read map[string]interface{}
	if ok, _ := regexp.Match(`^\d+$`, []byte(idStr)); ok {
		id, _ := strconv.Atoi(idStr)
		read, err = ser.Read(id)
	} else {
		read, err = ser.ReadByUrlRewrite(idStr)
	}
	if err != nil {
		c.Error(err.Error())
		return
	} else {
		c.Data["info"] = read["Blog"]
		c.Data["statistics"] = read["Statistics"]
		c.Data["TimeAdd"] = read["TimeAdd"]
		c.Data["Content"] = read["Content"]
	}
	c.TplName = "blog/get.html"
}

// GetAll ...
// @Title Get All
// @Description get Blog
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Blog
// @Failure 403
// @router / [get]
func (c *BlogController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int = 20

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt("limit"); err == nil {
		limit = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}
	page,_:=c.GetInt("page")
	data,err:=blog.GetAllBlog(query, fields, sortby, order, page, limit)
	if err != nil {
		//c.Data["data"] = err.Error()
		fmt.Println(err.Error())
	} else {
		c.Data["data"] = data
	}
	c.TplName = "blog/index.html"
	//fmt.Println("==========")
	//db:=db.NewDb()
	//bb:=make([]model.Blog,0)
	//err=db.Where("blog_id =?",54).Find(&bb)
	//fmt.Println(err)
	//fmt.Println(bb)
	//fmt.Println("==========")
	//b2:=new(model.Blog)
	//var ok bool
	//ok,err=db.Get(b2)
	//fmt.Println(err)
	//fmt.Println(ok)
	//fmt.Println(b2)
	//fmt.Println("==========")
	//b22:=new(model.Blog)
	//err=db.Find(&b22)
	//fmt.Println(err)
	//fmt.Println(b22)
	//where :=make(map[string]interface{})
	//where["blog_id in (?)"]=[]string{"53","54"}
	////where["blog_id in (?)"]=[]int{53,54}
	//o:=db.Filter(where).OrderBy("blog_id asc")
	//err=o.Find(&bb)
	//fmt.Println(err)
	//fmt.Println(bb)
	q:=make(map[string]interface{})
	q["content like ?"]="%jpeg%"
	var b *blog.Blog
	pag,err:=b.GetAll(q,[]string{},"blog_id",page,20)
	fmt.Println(err)
	fmt.Println(pag)

}
package goshua

import (
	"encoding/json"

	"errors"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

// oprations for Accounts
type BaseCRUDController struct {
	beego.Controller
	Model CRUDModel
}

func (c *BaseCRUDController) URLMapping() {
	c.Mapping("Post", c.HandleError(c.Post))
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

func (c *BaseCRUDController) HandleError(handler func() (interface{}, error)) {
	v, err := handler()
	if err != nil {
		c.Data["json"] = err.Error()
		if(c.Ctx.Output.Status == 0)
			c.Ctx.Output.SetStatus(500)
	}

	c.Data["json"] = v
	
	if(c.Ctx.Output.Status == 0)
		c.Ctx.Output.SetStatus(200)

	c.ServeJSON()
}

// @Title Post
// @Description create Accounts
// @Param	body		body 	models.Account	true		"body for Accounts content"
// @Success 201 {int} models.Account
// @Failure 403 body is empty
// @router / [post]
func (c *BaseCRUDController) Post() (v interface{}, err error) {
	v = c.Model.New(0)

	err = json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		return
	}

	err = c.BeforeCreate(v)
	if err != nil {
		return
	}

	_, err := c.Model.Add(v)
	if err != nil {
		c.Ctx.Output.SetStatus(409)
		return
	}

	err != c.AfterCreate(v)
	if err != nil {
		return
	}

	c.Ctx.Output.SetStatus(201)
}

// @Title Get
// @Description get Accounts by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Account
// @Failure 403 :id is empty
// @router /:id [get]
func (c *BaseCRUDController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := c.Model.GetById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// @Title Get All
// @Description get Accounts
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Account
// @Failure 403
// @router / [get]
func (c *BaseCRUDController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query map[string]string = make(map[string]string)
	var limit int64 = 10
	var offset int64 = 0

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
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

	l, err := c.Model.GetAll(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// @Title Update
// @Description update the Accounts
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Account	true		"body for Accounts content"
// @Success 200 {object} models.Account
// @Failure 403 :id is not int
// @router /:id [put]
func (c *BaseCRUDController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := c.Model.New(id)
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, v); err == nil {
		if err := c.Model.UpdateById(v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// @Title Delete
// @Description delete the Accounts
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *BaseCRUDController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := c.Model.Delete(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *BaseCRUDController) Options() {
	c.Ctx.ResponseWriter.WriteHeader(200)
}

//Observers

func (c *BaseCRUDController) BeforeCreate(v interface{}) (err error) {}
func (c *BaseCRUDController) AfterCreate(v interface{}) (err error)  {}

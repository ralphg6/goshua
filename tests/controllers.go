package tests

import (
	"encoding/json"
	"fmt"

	"git.labbs.com.br/sandman/run_accounts/models"

	"github.com/astaxie/beego/context"

	"github.com/ralphg6/goshua/controllers"
)

type AccountsController struct {
	controllers.BaseCRUDController
}

func (c *AccountsController) Init(ctx *context.Context, controllerName, actionName string, app interface{}) {
	c.Model = NewAccountModel()
	c.BaseCRUDController.Init(ctx, controllerName, actionName, app)
}

func (c *AccountsController) URLMapping() {
	fmt.Println("teste")
	c.Mapping("Post", c.Post)
}

// @Title Post
// @Description create Accounts
// @Param	body		body 	models.Account	true		"body for Accounts content"
// @Success 201 {int} models.Account
// @Failure 403 body is empty
// @router / [post]
func (c *AccountsController) Post() {
	var v models.Account
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddAccount(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

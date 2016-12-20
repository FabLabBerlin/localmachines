// /api/categories
package categories

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/categories"
	"github.com/astaxie/beego"
	"strings"
)

type Controller struct {
	controllers.Controller
}

// @Title GetAll
// @Description Get all categories
// @Success 200
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router / [get]
func (c *Controller) GetAll() {
	cs, err := categories.GetAll()
	if err != nil {
		c.CustomAbort(500, "Failed to get all categories")
	}

	c.Data["json"] = cs
	c.ServeJSON()
}

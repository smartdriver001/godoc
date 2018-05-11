package controllers

import (
	"github.com/astaxie/beego"
	"showdoc/consts"
	"showdoc/models"
)

// 我的项目
type ItemController struct {
	beego.Controller
}

// @Title  删除item
// @Description 分类列表
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.catalogs
// @router /delete [post]
func (this *ItemController) Delete() {
	json := consts.Json{}
	id,_ := this.GetInt("item_id")
	uid := this.GetSession(consts.SESSION_UID)
	if uid == nil {
		this.Abort("403")
	} else {
		ret,item := models.GetOneItem(id)
		if ret {
			err := item.Delete()
			if err == nil {
				json.Set(0,"删除成功")
			} else {
				json.Set(404,"删除失败")
			}
		} else {
			json.Set(404,"数据不存在")
		}
		this.Data["json"] = json.VendorOk()

		this.ServeJSON()
	}
}

// @Title item info
// @Description item info
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.item
// @router /info [post]
func (u *ItemController) Info() {

	json := consts.Json{}
	id,_ := u.GetInt("item_id")
	uid := u.GetSession(consts.SESSION_UID)

	_,item := models.GetOneItem(id)
	itemInfo := item.GetItemInfo()

	//是否登录
	if uid == nil {
		itemInfo.IsLogin = false
	}  else {
		itemInfo.IsLogin = true
	}

	//是否可以编辑
	if item.UserId == uid {
		itemInfo.ItemCreator = true
		itemInfo.ItemPermn = true
	}


	json.SetData(itemInfo)
	u.Data["json"] = json.VendorOk()
	u.ServeJSON()


}

// @Title MyList
// @Description mylist
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.item
// @router /myList [get]
func (u *ItemController) MyList() {
	json := consts.Json{}

	uid ,err := consts.IsLogin(u.Controller)
	if err != nil {
		u.Abort("403")
	} else {

		myItem := models.GetMyItem(uid)

		json.SetData(myItem)
		u.Data["json"] = json.VendorOk()
		u.ServeJSON()
	}


}


// @Title add item
// @Description add item
// @Param   item_type     formData    int   true        "项目类型 1常规项目  2单页项目"
// @Param   item_name     formData    string  true        "项目名称"
// @Param   password     formData    string  false        "查看密码"
// @Param   item_description     formData    string  false        "项目描述"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router /add [post]
func (u *ItemController) Add() {

	var err error
	var json consts.Json
	var userId int
	uid := u.GetSession(consts.SESSION_UID)
	if uid == nil {
		json.Set(10000, "用户未登录")
		u.Data["json"] = json.VendorError()
		u.ServeJSON()
		return
	} else {
		value, ok := uid.(int64)
		if ok {
			userId = int(value)
		}
	}

	item_type,_ := u.GetInt("item_type")
	item_name := u.GetString("item_name")
	password := u.GetString("password")
	item_description := u.GetString("item_description")

	var item models.Item
	item.Title = item_name
	item.Type = item_type
	item.Password = password
	item.Description = item_description
	item.UserId = userId
	item.Id,err = item.Create()

	if err != nil {
		json.Set(500, err.Error())
		u.Data["json"] = json.VendorError()
		u.ServeJSON()
		return
	}
	json.Set(0,"成功")
	u.Data["json"] = json.VendorOk()
	u.ServeJSON()
}

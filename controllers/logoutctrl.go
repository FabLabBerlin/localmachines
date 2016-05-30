package controllers

type LogoutController struct {
	Controller
}

func (this *LogoutController) Get() {
	this.DestroySession()
	this.Redirect("/machines/#/login", 302)
}

package userctrls

import (
	"github.com/astaxie/beego"
)

type OAuth2Controller struct {
	Controller
}

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// @Title LoginOAuth2
// @Description Logs user into the system by using OAuth2
// @Param	grant_type		body 	string	true		"Only password supported at the moment"
// @Param	client_id		body 	string 	true 		"Client Id"
// @Param	client_secret	body 	string 	false 		"Client Secret"
// @Param	scope			body 	string 	false 		"Scope"
// @Param	username		body 	string 	true 		"user's username"
// @Param	password		body 	string 	true 		"user's password"
// @Param	location		body 	int 	true 		"Location Id"
// @Success 200 {object} models.LoginResponse
// @Failure 401 Failed to authenticate
// @router /login [post]
func (this *OAuth2Controller) Login() {
	locId, err := this.GetInt64("location")
	if err != nil {
		beego.Error("get location:", err)
		this.CustomAbort(400, "Bad Request")
	}

	if this.GetString("grant_type") != "password" {
		this.CustomAbort(400, "Only grant_type password supported")
	}

	username := this.GetString("username")
	pw := this.GetString("password")

	_, unregisteredAtLocation, err := login(locId, false, username, pw)
	if err != nil {
		beego.Error("login:", err)
		this.Abort("401")
	}
	if unregisteredAtLocation {
		this.CustomAbort(400, "EASY LAB user but not registered (accepted terms) at requested location")
	}

	this.Data["json"] = AccessTokenResponse{
		AccessToken: "foo",
		ExpiresIn:   3600,
	}
	this.ServeJSON()
}

package userctrls

type OAuth2Controller struct {
	Controller
}

// @Title LoginOAuth2
// @Description Logs user into the system by using OAuth2
// @Param	grant_type		body 	string	true		"Only password supported at the moment"
// @Param	client_id		body 	string 	true 		"Client ID"
// @Param	client_secret	body 	string 	true 		"Client Secret"
// @Param	scope			body 	string 	true 		"Scope"
// @Param	username		body 	string 	true 		"user's username"
// @Param	password		body 	string 	true 		"user's password"
// @Success 200 {object} models.LoginResponse
// @Failure 401 Failed to authenticate
// @router /login [post]
func (this *OAuth2Controller) Login() {
	this.ServeJSON()
}

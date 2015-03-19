package controllers

type ActivationController struct {
	Controller
}

// @Title Get
// @Description Get activation by activation ID
// @Param	aid		path 	int	true		"Activation ID"
// @Success 200 {object} models.Activation
// @Failure	403	Failed to get activation
// @Failure	401	Not authorized
// @router /:aid [get]
func (this *ActivationController) Get() {

}

// @Title GetAll
// @Description Get all active activations
// @Success 200 {object} models.Activation
// @Failure	403	Failed to get all users
// @router / [get]
func (this *ActivationController) GetActive() {

}

// @Title Create
// @Description Create new activation
// @Param	body		body 	models.Activation	true		"Activation body"
// @Success 201 {object} models.ActivationCreateResponse
// @Failure	403	Failed to create activation
// @Failure 401 Not authorized
// @router / [post]
func (this *ActivationController) Create() {

}

// @Title Close
// @Description Close running activation
// @Param	aid		path 	int	true		"Activation ID"
// @Success 200 {object} models.StatusResponse
// @Failure	403	Failed to close activation
// @Failure 401 Not authorized
// @router /:aid [put]
func (this *ActivationController) Close() {

}
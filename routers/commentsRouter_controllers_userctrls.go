package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:ForgotPassword"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:ForgotPassword"],
		beego.ControllerComments{
			Method: "ForgotPassword",
			Router: `/forgot_password`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:ForgotPassword"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:ForgotPassword"],
		beego.ControllerComments{
			Method: "CheckPhone",
			Router: `/forgot_password/phone`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:ForgotPassword"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:ForgotPassword"],
		beego.ControllerComments{
			Method: "ResetPassword",
			Router: `/forgot_password/reset`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:OAuth2Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:OAuth2Controller"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/oauth2/login`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserDashboardController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserDashboardController"],
		beego.ControllerComments{
			Method: "GetDashboard",
			Router: `/:uid/dashboard`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserDashboardController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserDashboardController"],
		beego.ControllerComments{
			Method: "LP",
			Router: `/:uid/dashboard/lp`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserDashboardController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserDashboardController"],
		beego.ControllerComments{
			Method: "WS",
			Router: `/:uid/dashboard/ws`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserLocationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserLocationsController"],
		beego.ControllerComments{
			Method: "GetUserLocations",
			Router: `/:uid/locations`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserLocationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserLocationsController"],
		beego.ControllerComments{
			Method: "PostUserLocation",
			Router: `/:uid/locations/:lid`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserLocationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserLocationsController"],
		beego.ControllerComments{
			Method: "PutUserLocation",
			Router: `/:uid/locations/:lid`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserLocationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserLocationsController"],
		beego.ControllerComments{
			Method: "DeleteUserLocation",
			Router: `/:uid/locations/:lid`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserMembershipsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserMembershipsController"],
		beego.ControllerComments{
			Method: "PostUserMemberships",
			Router: `/:uid/memberships`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserMembershipsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserMembershipsController"],
		beego.ControllerComments{
			Method: "GetUserMemberships",
			Router: `/:uid/memberships`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserMembershipsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserMembershipsController"],
		beego.ControllerComments{
			Method: "PutUserMembership",
			Router: `/:uid/memberships/:umid`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserMembershipsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserMembershipsController"],
		beego.ControllerComments{
			Method: "DeleteUserMembership",
			Router: `/:uid/memberships/:umid`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserPermissionsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserPermissionsController"],
		beego.ControllerComments{
			Method: "CreateUserPermission",
			Router: `/:uid/permissions`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserPermissionsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserPermissionsController"],
		beego.ControllerComments{
			Method: "DeleteUserPermission",
			Router: `/:uid/permissions`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserPermissionsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserPermissionsController"],
		beego.ControllerComments{
			Method: "UpdateUserPermissions",
			Router: `/:uid/permissions`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserPermissionsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserPermissionsController"],
		beego.ControllerComments{
			Method: "GetUserPermissions",
			Router: `/:uid/permissions`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserPermissionsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserPermissionsController"],
		beego.ControllerComments{
			Method: "GetUserMachinePermissions",
			Router: `/:uid/machinepermissions`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			Method: "Logout",
			Router: `/logout`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			Method: "GetCurrentUser",
			Router: `/current`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			Method: "Signup",
			Router: `/signup`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:uid`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:uid`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			Method: "GetUserMachines",
			Router: `/:uid/machines`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			Method: "GetUserBill",
			Router: `/:uid/bill`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			Method: "GetUserNames",
			Router: `/names`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			Method: "PostUserPassword",
			Router: `/:uid/password`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			Method: "GetByNfcId",
			Router: `/by_nfc_id`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			Method: "PostNfcId",
			Router: `/:uid/nfc_id`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

}

package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:ForgotPassword"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:ForgotPassword"],
		beego.ControllerComments{
			"ForgotPassword",
			`/forgot_password`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:ForgotPassword"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:ForgotPassword"],
		beego.ControllerComments{
			"CheckPhone",
			`/forgot_password/phone`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:ForgotPassword"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:ForgotPassword"],
		beego.ControllerComments{
			"ResetPassword",
			`/forgot_password/reset`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:OAuth2Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:OAuth2Controller"],
		beego.ControllerComments{
			"Login",
			`/oauth2/login`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserDashboardController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserDashboardController"],
		beego.ControllerComments{
			"GetDashboard",
			`/:uid/dashboard`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserDashboardController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserDashboardController"],
		beego.ControllerComments{
			"LP",
			`/:uid/dashboard/lp`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserDashboardController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserDashboardController"],
		beego.ControllerComments{
			"WS",
			`/:uid/dashboard/ws`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserLocationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserLocationsController"],
		beego.ControllerComments{
			"GetUserLocations",
			`/:uid/locations`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserLocationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserLocationsController"],
		beego.ControllerComments{
			"PostUserLocation",
			`/:uid/locations/:lid`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserLocationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserLocationsController"],
		beego.ControllerComments{
			"PutUserLocation",
			`/:uid/locations/:lid`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserLocationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserLocationsController"],
		beego.ControllerComments{
			"DeleteUserLocation",
			`/:uid/locations/:lid`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserMembershipsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserMembershipsController"],
		beego.ControllerComments{
			"PostUserMemberships",
			`/:uid/memberships`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserMembershipsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserMembershipsController"],
		beego.ControllerComments{
			"GetUserMemberships",
			`/:uid/memberships`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserMembershipsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserMembershipsController"],
		beego.ControllerComments{
			"PutUserMembership",
			`/:uid/memberships/:umid`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserMembershipsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserMembershipsController"],
		beego.ControllerComments{
			"DeleteUserMembership",
			`/:uid/memberships/:umid`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserPermissionsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserPermissionsController"],
		beego.ControllerComments{
			"CreateUserPermission",
			`/:uid/permissions`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserPermissionsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserPermissionsController"],
		beego.ControllerComments{
			"DeleteUserPermission",
			`/:uid/permissions`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserPermissionsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserPermissionsController"],
		beego.ControllerComments{
			"UpdateUserPermissions",
			`/:uid/permissions`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserPermissionsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserPermissionsController"],
		beego.ControllerComments{
			"GetUserPermissions",
			`/:uid/permissions`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserPermissionsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UserPermissionsController"],
		beego.ControllerComments{
			"GetUserMachinePermissions",
			`/:uid/machinepermissions`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"GetByNfcId",
			`/by_nfc_id`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"PostNfcId",
			`/:uid/nfc_id`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"Login",
			`/login`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"Logout",
			`/logout`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"GetCurrentUser",
			`/current`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"Signup",
			`/signup`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"Get",
			`/:uid`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"Put",
			`/:uid`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"GetUserMachines",
			`/:uid/machines`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"GetUserBill",
			`/:uid/bill`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"GetUserNames",
			`/names`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"PostUserPassword",
			`/:uid/password`,
			[]string{"post"},
			nil})

}

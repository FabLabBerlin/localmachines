package routers

import (
	"github.com/astaxie/beego"
)

func init() {
	
	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"Login",
			`/login`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"LoginUid",
			`/loginuid`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"Logout",
			`/logout`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"Signup",
			`/signup`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"Get",
			`/:uid`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"Put",
			`/:uid`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"GetUserMachines",
			`/:uid/machines`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"GetUserBill",
			`/:uid/bill`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"GetUserNames",
			`/names`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"PostUserPassword",
			`/:uid/password`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UsersController"],
		beego.ControllerComments{
			"UpdateNfcUid",
			`/:uid/nfcuid`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UserMembershipsController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UserMembershipsController"],
		beego.ControllerComments{
			"PostUserMemberships",
			`/:uid/memberships`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UserMembershipsController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UserMembershipsController"],
		beego.ControllerComments{
			"GetUserMemberships",
			`/:uid/memberships`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UserMembershipsController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UserMembershipsController"],
		beego.ControllerComments{
			"DeleteUserMembership",
			`/:uid/memberships/:umid`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UserMembershipsController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UserMembershipsController"],
		beego.ControllerComments{
			"PutUserMembership",
			`/:uid/memberships/:umid`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UserPermissionsController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UserPermissionsController"],
		beego.ControllerComments{
			"CreateUserPermission",
			`/:uid/permissions`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UserPermissionsController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UserPermissionsController"],
		beego.ControllerComments{
			"DeleteUserPermission",
			`/:uid/permissions`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UserPermissionsController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UserPermissionsController"],
		beego.ControllerComments{
			"UpdateUserPermissions",
			`/:uid/permissions`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UserPermissionsController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UserPermissionsController"],
		beego.ControllerComments{
			"GetUserPermissions",
			`/:uid/permissions`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UserPermissionsController"] = append(beego.GlobalControllerRouter["github.com/kr15h/fabsmith/controllers/userctrls:UserPermissionsController"],
		beego.ControllerComments{
			"GetUserMachinePermissions",
			`/:uid/machinepermissions`,
			[]string{"get"},
			nil})

}

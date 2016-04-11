/*
user_roles declaration.

Super Admins can basically do anything that is possible within the policies.
Local Admins can do basically anything that is only relevant for a certain
location. Staff cannot visit the Admin Page but they can turn on and off
all machines. Also they can turn on and off maintenance.
*/
package user_roles

const (
	SUPER_ADMIN    = Role("superadmin")
	ADMIN          = Role("admin")
	STAFF          = Role("staff")
	MEMBER         = Role("member")
	API            = Role("api")
	NOT_AFFILIATED = Role("notaffiliated")
)

type Role string

func (r Role) String() string {
	return string(r)
}

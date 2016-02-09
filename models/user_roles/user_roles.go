package user_roles

const (
	SUPER_ADMIN    = Role("superadmin")
	ADMIN          = Role("admin")
	STAFF          = Role("staff")
	MEMBER         = Role("member")
	NOT_AFFILIATED = Role("notaffiliated")
)

type Role string

func (r Role) String() string {
	return string(r)
}

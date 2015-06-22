package userTests

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/kr15h/fabsmith/models"
	. "github.com/kr15h/fabsmith/tests/models"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	ConfigDB()
}

func TestDeleteUser(t *testing.T) {
	Convey("Testing Delete user", t, func() {
		u := models.User{
			FirstName: "test",
			LastName:  "test",
		}
		Convey("Creating User and delete it", func() {
			uc, err := models.CreateUser(&u)

			err = models.DeleteUser(uc)
			So(err, ShouldBeNil)
		})
		Convey("Try to delete non-existing user", func() {
			err := models.DeleteUser(0)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestUserCreate(t *testing.T) {
	Convey("Testing CreateUser", t, func() {
		u := models.User{
			Email:    "test",
			Username: "test",
		}
		Convey("Creating one user into database", func() {
			// Creating user
			uc, err := models.CreateUser(&u)
			defer models.DeleteUser(uc)
			So(err, ShouldBeNil)
			So(uc, ShouldBeGreaterThan, 0)
		})
		Convey("Creating 2 users that are identical into database, should get an error", func() {
			// Creating first user
			uc, err := models.CreateUser(&u)
			defer models.DeleteUser(uc)
			So(err, ShouldBeNil)
			So(uc, ShouldBeGreaterThan, 0)

			// Creating second user
			uc2, err2 := models.CreateUser(&u)
			So(err2, ShouldNotBeNil)
			So(uc2, ShouldEqual, 0)
		})
	})
}

func TestDeleteUserAuth(t *testing.T) {
	Convey("Testing DeleteUserAuth", t, func() {
		u := models.User{
			Username: "test",
		}
		Convey("Creating user with password and delete his Auth", func() {
			uid, _ := models.CreateUser(&u)
			defer models.DeleteUser(uid)
			models.AuthSetPassword(uid, "test")

			err := models.DeleteUserAuth(uid)
			So(err, ShouldBeNil)
		})
		Convey("Delete auth on non-existing user", func() {
			err := models.DeleteUserAuth(0)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestAuthSetPassword(t *testing.T) {
	Convey("Testing AuthenticateUser", t, func() {
		u := models.User{
			Username: "test",
		}
		Convey("Creating a user and setting him a password", func() {
			uid, err := models.CreateUser(&u)
			defer models.DeleteUser(uid)

			err = models.AuthSetPassword(uid, "test")
			defer models.DeleteUserAuth(uid)
			So(err, ShouldBeNil)
		})
		Convey("Try setting password on non-existing user", func() {
			err := models.AuthSetPassword(0, "test")
			So(err, ShouldNotBeNil)
		})
	})
}

func TestAuthenticateUser(t *testing.T) {
	Convey("Testing AuthenticateUser", t, func() {
		u := models.User{
			Username: "test",
		}
		Convey("Creating a user and setting him a password", func() {
			uid, err := models.CreateUser(&u)
			defer models.DeleteUser(uid)

			err = models.AuthSetPassword(uid, "test")
			defer models.DeleteUserAuth(uid)
			So(err, ShouldBeNil)
		})
		Convey("Creating a user with a password and try to authenticate him", func() {
			uid, err := models.CreateUser(&u)
			defer models.DeleteUser(uid)

			err = models.AuthSetPassword(uid, "test")
			defer models.DeleteUserAuth(uid)

			authUID, err := models.AuthenticateUser("test", "test")
			So(authUID, ShouldEqual, uid)
			So(err, ShouldBeNil)
		})
		Convey("Creating a user with a password and try to authenticate with wrong username", func() {
			uid, err := models.CreateUser(&u)
			defer models.DeleteUser(uid)

			err = models.AuthSetPassword(uid, "test")
			defer models.DeleteUserAuth(uid)

			authUID, err := models.AuthenticateUser("wrong", "test")
			So(authUID, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})
		Convey("Creating a user with a password and try to authenticate with wrong password", func() {
			uid, err := models.CreateUser(&u)
			defer models.DeleteUser(uid)

			err = models.AuthSetPassword(uid, "test")
			defer models.DeleteUserAuth(uid)

			authUID, err := models.AuthenticateUser("test", "wrong")
			So(authUID, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestAuthUpdateNfcUid(t *testing.T) {
	Convey("Testing AuthUpdateNfcUid", t, func() {
		u := models.User{
			Username: "test",
		}
		Convey("Creating a user and setting him a NFC UID", func() {
			uid, _ := models.CreateUser(&u)
			defer models.DeleteUser(uid)

			_ = models.AuthSetPassword(uid, "test")
			defer models.DeleteUserAuth(uid)

			err := models.AuthUpdateNfcUid(uid, "123456")
			So(err, ShouldBeNil)
		})
		Convey("Setting NFC UID to non-existing user", func() {
			err := models.AuthUpdateNfcUid(0, "123456")
			So(err, ShouldNotBeNil)
		})
	})
}

func TestAuthenticateUserUID(t *testing.T) {
	Convey("Testing AuthenticateUser", t, func() {
		u := models.User{
			Username: "test",
		}
		Convey("Creating a user with a NFC UID and try to authenticate him", func() {
			uid, err := models.CreateUser(&u)
			defer models.DeleteUser(uid)

			err = models.AuthSetPassword(uid, "test")
			defer models.DeleteUserAuth(uid)

			err = models.AuthUpdateNfcUid(uid, "123456")

			_, authUID, err := models.AuthenticateUserUid("123456")
			So(authUID, ShouldEqual, uid)
			So(err, ShouldBeNil)
		})
		Convey("Creating a user with a NFC UID and try to authenticate him with wrong UID", func() {
			uid, err := models.CreateUser(&u)
			defer models.DeleteUser(uid)

			err = models.AuthSetPassword(uid, "test")
			defer models.DeleteUserAuth(uid)

			err = models.AuthUpdateNfcUid(uid, "123456")

			_, authUID, err := models.AuthenticateUserUid("654321")
			So(authUID, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestGetUser(t *testing.T) {
	Convey("Testing GetUser", t, func() {
		u := models.User{
			Username: "test",
		}
		Convey("Creating user and get it", func() {
			uid, err := models.CreateUser(&u)
			defer models.DeleteUser(uid)

			user, err := models.GetUser(uid)
			So(user.Id, ShouldEqual, uid)
			So(user.Username, ShouldEqual, "test")
			So(err, ShouldBeNil)
		})
		Convey("Try to get non-existing user", func() {
			user, err := models.GetUser(0)
			So(user, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestGetAllUsers(t *testing.T) {
	Convey("Testing GetAllUsers", t, func() {
		u1 := models.User{
			Email:    "u1",
			Username: "u1",
		}
		u2 := models.User{
			Email:    "u2",
			Username: "u2",
		}
		Convey("Getting all users with 0 users in the database", func() {
			users, err := models.GetAllUsers()
			So(len(users), ShouldEqual, 0)
			So(err, ShouldBeNil)
		})
		Convey("Creating 2 users and get them", func() {
			uid1, _ := models.CreateUser(&u1)
			defer models.DeleteUser(uid1)
			uid2, _ := models.CreateUser(&u2)
			defer models.DeleteUser(uid2)

			fmt.Print(uid1, uid2)

			users, err := models.GetAllUsers()
			So(len(users), ShouldEqual, 2)
			So(err, ShouldBeNil)
		})
	})
}

func TestUpdateUser(t *testing.T) {
	Convey("Testing UpdateUser", t, func() {
		u := models.User{
			Username: "test",
		}
		Convey("Creating a user and try to modify FirstName", func() {
			uid, _ := models.CreateUser(&u)
			defer models.DeleteUser(uid)
			user, _ := models.GetUser(uid)

			user.FirstName = "YOLO"
			err := models.UpdateUser(user)
			So(err, ShouldBeNil)
			user, _ = models.GetUser(user.Id)
			So(user.FirstName, ShouldEqual, "YOLO")
		})
		Convey("Try updating non-existing user", func() {
			err := models.UpdateUser(&u)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestCreateUserPermission(t *testing.T) {
	Convey("Testing CreateUserPermission", t, func() {
		u := models.User{
			Username: "test",
		}
		Convey("Creating one UserPermission", func() {
			uid, _ := models.CreateUser(&u)
			defer models.DeleteUser(uid)

			err := models.CreateUserPermission(uid, 0)
			So(err, ShouldBeNil)
		})
	})
}

func TestDeleteUserPermission(t *testing.T) {
	Convey("Testing DeleteUserPermission", t, func() {
		u := models.User{
			Username: "test",
		}
		Convey("Creating one UserPermission and delete it", func() {
			uid, _ := models.CreateUser(&u)
			defer models.DeleteUser(uid)
			models.CreateUserPermission(uid, 0)

			err := models.DeleteUserPermission(uid, 0)
			So(err, ShouldBeNil)
		})
	})
}

func TestUpdateUserPermission(t *testing.T) {
	Convey("Testing UpdateUserPermission", t, func() {
		u := models.User{
			Username: "test",
		}
		Convey("Creating one User and update his permission", func() {
			uid, _ := models.CreateUser(&u)
			defer models.DeleteUser(uid)

			perms := &[]models.Permission{
				models.Permission{UserId: uid, MachineId: 0},
				models.Permission{UserId: uid, MachineId: 1},
				models.Permission{UserId: uid, MachineId: 2},
				models.Permission{UserId: uid, MachineId: 3},
			}
			err := models.UpdateUserPermissions(uid, perms)
			So(err, ShouldBeNil)
		})
	})
}

package modelTest

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/kr15h/fabsmith/models"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	ConfigDB()
}

func TestUsers(t *testing.T) {
	Convey("Testing User model", t, func() {
		Reset(ResetDB)
		Convey("Testing Delete user", func() {
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
		Convey("Testing CreateUser", func() {
			u := models.User{
				Email:    "test",
				Username: "test",
			}
			Convey("Creating one user into database", func() {
				uc, err := models.CreateUser(&u)

				So(err, ShouldBeNil)
				So(uc, ShouldBeGreaterThan, 0)
			})
			Convey("Creating 2 users that are identical into database, should get an error", func() {
				// Creating first user
				uc, err := models.CreateUser(&u)
				uc2, err2 := models.CreateUser(&u)

				So(err, ShouldBeNil)
				So(uc, ShouldBeGreaterThan, 0)
				So(err2, ShouldNotBeNil)
				So(uc2, ShouldEqual, 0)
			})
		})
		Convey("Testing DeleteUserAuth", func() {
			u := models.User{
				Username: "test",
			}
			Convey("Creating user with password and delete his Auth", func() {
				uid, _ := models.CreateUser(&u)
				models.AuthSetPassword(uid, "test")
				err := models.DeleteUserAuth(uid)

				So(err, ShouldBeNil)
			})
			Convey("Delete auth on non-existing user", func() {
				err := models.DeleteUserAuth(0)

				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing AuthSetPassword", func() {
			u := models.User{
				Username: "test",
			}
			Convey("Creating a user and setting him a password", func() {
				uid, err := models.CreateUser(&u)
				err = models.AuthSetPassword(uid, "test")

				So(err, ShouldBeNil)
			})
			Convey("Try setting password on non-existing user", func() {
				err := models.AuthSetPassword(0, "test")

				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing AuthenticateUser", func() {
			u := models.User{
				Username: "test",
			}
			Convey("Creating a user with a password and try to authenticate him", func() {
				uid, err := models.CreateUser(&u)
				err = models.AuthSetPassword(uid, "test")
				authUID, err := models.AuthenticateUser("test", "test")

				So(authUID, ShouldEqual, uid)
				So(err, ShouldBeNil)
			})
			Convey("Creating a user with a password and try to authenticate with wrong username", func() {
				uid, err := models.CreateUser(&u)
				err = models.AuthSetPassword(uid, "test")
				authUID, err := models.AuthenticateUser("wrong", "test")

				So(authUID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
			})
			Convey("Creating a user with a password and try to authenticate with wrong password", func() {
				uid, err := models.CreateUser(&u)
				err = models.AuthSetPassword(uid, "test")
				authUID, err := models.AuthenticateUser("test", "wrong")

				So(authUID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing AuthUpdateNfcUid", func() {
			u := models.User{
				Username: "test",
			}
			Convey("Creating a user and setting him a NFC UID", func() {
				uid, _ := models.CreateUser(&u)
				_ = models.AuthSetPassword(uid, "test")
				err := models.AuthUpdateNfcUid(uid, "123456")

				So(err, ShouldBeNil)
			})
			Convey("Setting NFC UID to non-existing user", func() {
				err := models.AuthUpdateNfcUid(0, "123456")

				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing AuthenticateUserUid", func() {
			u := models.User{
				Username: "test",
			}
			Convey("Creating a user with a NFC UID and try to authenticate him", func() {
				uid, err := models.CreateUser(&u)
				err = models.AuthSetPassword(uid, "test")
				err = models.AuthUpdateNfcUid(uid, "123456")
				_, authUID, err := models.AuthenticateUserUid("123456")

				So(authUID, ShouldEqual, uid)
				So(err, ShouldBeNil)
			})
			Convey("Creating a user with a NFC UID and try to authenticate him with wrong UID", func() {
				uid, err := models.CreateUser(&u)
				err = models.AuthSetPassword(uid, "test")
				err = models.AuthUpdateNfcUid(uid, "123456")
				_, authUID, err := models.AuthenticateUserUid("654321")

				So(authUID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing GetUser", func() {
			u := models.User{
				Username: "test",
			}
			Convey("Creating user and get it", func() {
				uid, err := models.CreateUser(&u)
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
		Convey("Testing GetAllUsers", func() {
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
				models.CreateUser(&u1)
				models.CreateUser(&u2)
				users, err := models.GetAllUsers()

				So(len(users), ShouldEqual, 2)
				So(err, ShouldBeNil)
			})
		})
		Convey("Testing UpdateUser", func() {
			u := models.User{
				Username: "test",
				Email:    "test@example.com",
			}
			Convey("Creating a user and try to modify FirstName", func() {
				uid, _ := models.CreateUser(&u)
				user, _ := models.GetUser(uid)
				user.FirstName = "YOLO"
				err := models.UpdateUser(user)
				user, _ = models.GetUser(user.Id)

				So(err, ShouldBeNil)
				So(user.FirstName, ShouldEqual, "YOLO")
			})
			Convey("Try updating non-existing user", func() {
				err := models.UpdateUser(&u)

				So(err, ShouldNotBeNil)
			})
			Convey("Update user with duplicate username, should return error", func() {
				user1 := models.User{}
				user1.Username = "hesus"
				user1.Email = "hesus@example.com"
				user2 := models.User{}
				user2.Username = "hrist"
				user2.Email = "hrist@example.com"
				models.CreateUser(&user1)
				models.CreateUser(&user2)
				user2.Username = user1.Username
				err := models.UpdateUser(&user2)

				So(err, ShouldNotBeNil)
			})
			Convey("Update user with duplicate email, should return error", func() {
				user1 := models.User{}
				user1.Username = "hesus"
				user1.Email = "hesus@example.com"
				user2 := models.User{}
				user2.Username = "hrist"
				user2.Email = "hrist@example.com"
				models.CreateUser(&user1)
				models.CreateUser(&user2)
				user2.Email = user1.Email
				err := models.UpdateUser(&user2)

				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing CreateUserPermission", func() {
			u := models.User{
				Username: "test",
			}
			Convey("Creating one UserPermission", func() {
				uid, _ := models.CreateUser(&u)
				err := models.CreateUserPermission(uid, 0)

				So(err, ShouldBeNil)
			})
		})
		Convey("Testing DeleteUserPermission", func() {
			u := models.User{
				Username: "test",
			}
			Convey("Creating one UserPermission and delete it", func() {
				uid, _ := models.CreateUser(&u)
				models.CreateUserPermission(uid, 0)
				err := models.DeleteUserPermission(uid, 0)

				So(err, ShouldBeNil)
			})
		})
		Convey("Testing UpdateUserPermission", func() {
			u := models.User{
				Username: "test",
			}
			Convey("Creating one User and update his permission", func() {
				uid, _ := models.CreateUser(&u)
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
	})
}

package modelTest

import (
	"testing"
	"time"

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
			Convey("Creating user and deleting it should be no problem", func() {
				user := models.User{}
				user.Email = "user@example.com"
				uid, err := models.CreateUser(&user)
				So(err, ShouldBeNil)
				So(uid, ShouldBeGreaterThan, 0)

				err = models.DeleteUser(uid)
				So(err, ShouldBeNil)
			})
			Convey("Try to delete non-existing user", func() {
				err := models.DeleteUser(0)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing CreateUser", func() {
			Convey("Creating user with invalid email should not work", func() {
				user := models.User{}
				user.Username = "hackenberg"
				user.Email = "hackingbackintime"
				_, err := models.CreateUser(&user)
				So(err, ShouldNotBeNil)

				user.Email = ""
				_, err = models.CreateUser(&user)
				So(err, ShouldNotBeNil)

				user.Email = "hacking@backintime"
				_, err = models.CreateUser(&user)
				So(err, ShouldNotBeNil)

				user.Email = "hacking@backintime."
				_, err = models.CreateUser(&user)
				So(err, ShouldNotBeNil)
			})
			Convey("Creating user with invalid email should not work #2", func() {
				user := models.User{}
				user.Username = "hackenberg"
				user.Email = ".Aackingbackintime@-email.com"
				_, err := models.CreateUser(&user)
				So(err, ShouldNotBeNil)
				user.Email = ""

				_, err = models.CreateUser(&user)
				So(err, ShouldNotBeNil)

				user.Email = "hacking@backintime"
				_, err = models.CreateUser(&user)
				So(err, ShouldNotBeNil)

				user.Email = "hacking@backintime."
				_, err = models.CreateUser(&user)
				So(err, ShouldNotBeNil)
			})
			Convey("Creating user with valid email should work", func() {
				user := models.User{}
				user.Username = "hackerman"
				user.Email = "Hacking@bAckin.tIme"
				uid, err := models.CreateUser(&user)

				So(err, ShouldBeNil)
				So(uid, ShouldBeGreaterThan, 0)
			})
			Convey("Creating 2 users that are identical into database, should get an error", func() {
				user := models.User{}
				user.Username = "hackerman"
				user.Email = "hacking@backin.time"
				uid1, err1 := models.CreateUser(&user)
				uid2, err2 := models.CreateUser(&user)

				So(err1, ShouldBeNil)
				So(uid1, ShouldBeGreaterThan, 0)
				So(err2, ShouldNotBeNil)
				So(uid2, ShouldEqual, 0)
			})
			Convey("Create 2 users with identical username, should return error", func() {
				user1 := models.User{}
				user1.Username = "william"
				user1.Email = "william@example.com"
				user2 := models.User{}
				user2.Username = user1.Username
				user2.Email = "ismael@example.com"
				uid, err := models.CreateUser(&user1)
				_, err = models.CreateUser(&user2)

				So(uid, ShouldBeGreaterThan, 0)
				So(err, ShouldNotBeNil)
			})
			Convey("Create 2 users with identical email, should return error", func() {
				user1 := models.User{}
				user1.Username = "william"
				user1.Email = "william@example.com"
				user2 := models.User{}
				user2.Username = "baram"
				user2.Email = user1.Email
				uid, err := models.CreateUser(&user1)
				_, err = models.CreateUser(&user2)

				So(uid, ShouldBeGreaterThan, 0)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing DeleteUserAuth", func() {
			Convey("Creating user with password and delete his Auth", func() {
				user := models.User{}
				user.Email = "user@example.com"
				uid, _ := models.CreateUser(&user)
				So(uid, ShouldBeGreaterThan, 0)

				err := models.AuthSetPassword(uid, "test")
				So(err, ShouldBeNil)

				err = models.DeleteUserAuth(uid)
				So(err, ShouldBeNil)
			})
			Convey("Delete auth on non-existing user", func() {
				err := models.DeleteUserAuth(0)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing AuthSetPassword", func() {
			Convey("Creating a user and setting him a password", func() {
				user := models.User{}
				user.Email = "user@example.com"
				uid, err := models.CreateUser(&user)
				So(uid, ShouldBeGreaterThan, 0)
				So(err, ShouldBeNil)

				err = models.AuthSetPassword(uid, "test")
				So(err, ShouldBeNil)
			})
			Convey("Try setting password on non-existing user", func() {
				err := models.AuthSetPassword(0, "test")

				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing AuthenticateUser", func() {
			Convey("Creating a user with a password and try to authenticate him", func() {
				user := models.User{}
				user.Username = "test"
				user.Email = "user@example.com"
				uid, _ := models.CreateUser(&user)
				models.AuthSetPassword(uid, "test")
				authUID, err := models.AuthenticateUser(user.Username, "test")
				So(authUID, ShouldEqual, uid)
				So(err, ShouldBeNil)
			})
			Convey("Creating a user with a password and try to authenticate with wrong username", func() {
				user := models.User{}
				user.Email = "user@example.com"
				uid, err := models.CreateUser(&user)
				err = models.AuthSetPassword(uid, "test")
				authUID, err := models.AuthenticateUser("wrong", "test")

				So(authUID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
			})
			Convey("Creating a user with a password and try to authenticate with wrong password", func() {
				user := models.User{}
				user.Email = "user@example.com"
				uid, err := models.CreateUser(&user)
				err = models.AuthSetPassword(uid, "test")
				authUID, err := models.AuthenticateUser("test", "wrong")

				So(authUID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing AuthUpdateNfcUid", func() {
			Convey("Creating a user and setting him a NFC UID", func() {
				user := models.User{}
				user.Username = "test"
				user.Email = "user@example.com"
				uid, _ := models.CreateUser(&user)
				_ = models.AuthSetPassword(uid, "test")
				err := models.AuthUpdateNfcUid(uid, "123456")
				So(err, ShouldBeNil)
			})
			Convey("Creating a user and setting him a NFC UID without setting a password", func() {
				user := models.User{}
				user.Username = "test"
				user.Email = "user@example.com"
				uid, _ := models.CreateUser(&user)
				err := models.AuthUpdateNfcUid(uid, "123456")
				So(err, ShouldBeNil)
			})
			Convey("Setting NFC UID to non-existing user", func() {
				err := models.AuthUpdateNfcUid(0, "123456")
				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing AuthenticateUserUid", func() {
			Convey("Creating a user with a NFC UID and try to authenticate him", func() {
				user := models.User{}
				user.Email = "user@example.com"
				uid, err := models.CreateUser(&user)
				err = models.AuthUpdateNfcUid(uid, "123456")
				_, authUID, err := models.AuthenticateUserUid("123456")

				So(authUID, ShouldEqual, uid)
				So(err, ShouldBeNil)
			})
			Convey("Creating a user with a NFC UID and try to authenticate him with wrong UID", func() {
				user := models.User{}
				user.Email = "user@example.com"
				uid, err := models.CreateUser(&user)
				err = models.AuthSetPassword(uid, "test")
				err = models.AuthUpdateNfcUid(uid, "123456")
				_, authUID, err := models.AuthenticateUserUid("654321")

				So(authUID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing GetUser", func() {
			Convey("Creating user and get it", func() {
				u := models.User{}
				u.Username = "test"
				u.Email = "user@example.com"
				uid, _ := models.CreateUser(&u)
				user, err := models.GetUser(uid)

				So(user.Id, ShouldEqual, uid)
				So(user.Username, ShouldEqual, "test")
				So(err, ShouldBeNil)

				Convey("The created time should be close to now", func() {
					So(user.Created, ShouldHappenWithin,
						time.Duration(1)*time.Second, time.Now())
				})
			})
			Convey("Try to get non-existing user", func() {
				user, err := models.GetUser(0)

				So(user, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing GetAllUsers", func() {
			u1 := models.User{
				Email:    "u1@example.com",
				Username: "u1",
			}
			u2 := models.User{
				Email:    "u2@example.com",
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
			Convey("Updating user with invalid email should not work", func() {
				user := models.User{}
				user.Username = "hackerman"
				user.Email = "hacker@inspacefeels.good"
				user.Id, _ = models.CreateUser(&user)

				user.Email = ""
				err := models.UpdateUser(&user)
				So(err, ShouldNotBeNil)

				user.Email = "invalidemail"
				err = models.UpdateUser(&user)
				So(err, ShouldNotBeNil)

				user.Email = "invalid@email"
				err = models.UpdateUser(&user)
				So(err, ShouldNotBeNil)

				user.Email = "invalid@email."
				err = models.UpdateUser(&user)
				So(err, ShouldNotBeNil)
			})
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

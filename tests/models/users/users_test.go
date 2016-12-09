package users

import (
	"fmt"
	"testing"
	"time"

	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/user_permissions"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	setup.ConfigDB()
}

func TestUsers(t *testing.T) {
	Convey("Testing User model", t, func() {
		Reset(setup.ResetDB)
		Convey("Testing CreateUser", func() {
			Convey("Creating user with invalid email should not work", func() {
				user := users.User{}
				user.Username = "hackenberg"
				user.Email = "hackingbackintime"
				_, err := users.CreateUser(&user)
				So(err, ShouldNotBeNil)

				user.Email = ""
				_, err = users.CreateUser(&user)
				So(err, ShouldNotBeNil)

				user.Email = "hacking@backintime"
				_, err = users.CreateUser(&user)
				So(err, ShouldNotBeNil)

				user.Email = "hacking@backintime."
				_, err = users.CreateUser(&user)
				So(err, ShouldNotBeNil)
			})
			Convey("Creating user with invalid email should not work #2", func() {
				user := users.User{}
				user.Username = "hackenberg"
				user.Email = ".Aackingbackintime@-email.com"
				_, err := users.CreateUser(&user)
				So(err, ShouldNotBeNil)
				user.Email = ""

				_, err = users.CreateUser(&user)
				So(err, ShouldNotBeNil)

				user.Email = "hacking@backintime"
				_, err = users.CreateUser(&user)
				So(err, ShouldNotBeNil)

				user.Email = "hacking@backintime."
				_, err = users.CreateUser(&user)
				So(err, ShouldNotBeNil)
			})
			Convey("Creating user with valid email should work", func() {
				user := users.User{}
				user.Username = "hackerman"
				user.Email = "Hacking@bAckin.tIme"
				uid, err := users.CreateUser(&user)

				So(err, ShouldBeNil)
				So(uid, ShouldBeGreaterThan, 0)
			})
			Convey("Creating 2 users that are identical into database, should get an error", func() {
				user := users.User{}
				user.Username = "hackerman"
				user.Email = "hacking@backin.time"
				uid1, err1 := users.CreateUser(&user)
				uid2, err2 := users.CreateUser(&user)

				So(err1, ShouldBeNil)
				So(uid1, ShouldBeGreaterThan, 0)
				So(err2, ShouldNotBeNil)
				So(uid2, ShouldEqual, 0)
			})
			Convey("Create 2 users with identical username, should return error", func() {
				user1 := users.User{}
				user1.Username = "william"
				user1.Email = "william@example.com"
				user2 := users.User{}
				user2.Username = user1.Username
				user2.Email = "ismael@example.com"
				uid, err := users.CreateUser(&user1)
				_, err = users.CreateUser(&user2)

				So(uid, ShouldBeGreaterThan, 0)
				So(err, ShouldNotBeNil)
			})
			Convey("Create 2 users with identical email, should return error", func() {
				user1 := users.User{}
				user1.Username = "william"
				user1.Email = "william@example.com"
				user2 := users.User{}
				user2.Username = "baram"
				user2.Email = user1.Email
				uid, err := users.CreateUser(&user1)
				_, err = users.CreateUser(&user2)

				So(uid, ShouldBeGreaterThan, 0)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing DeleteUserAuth", func() {
			Convey("Creating user with password and delete his Auth", func() {
				user := users.User{}
				user.Email = "user@example.com"
				uid, _ := users.CreateUser(&user)
				So(uid, ShouldBeGreaterThan, 0)

				err := users.AuthSetPassword(uid, "test")
				So(err, ShouldBeNil)

				err = users.DeleteUserAuth(uid)
				So(err, ShouldBeNil)
			})
			Convey("Delete auth on non-existing user", func() {
				err := users.DeleteUserAuth(0)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing AuthSetPassword", func() {
			Convey("Creating a user and setting him a password", func() {
				user := users.User{}
				user.Email = "user@example.com"
				uid, err := users.CreateUser(&user)
				So(uid, ShouldBeGreaterThan, 0)
				So(err, ShouldBeNil)

				err = users.AuthSetPassword(uid, "test")
				So(err, ShouldBeNil)
			})
			Convey("Try setting password on non-existing user", func() {
				err := users.AuthSetPassword(0, "test")

				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing AuthenticateUser", func() {
			Convey("Creating a user with a password and try to authenticate him", func() {
				user := users.User{}
				user.Username = "test"
				user.Email = "user@example.com"
				uid, _ := users.CreateUser(&user)
				users.AuthSetPassword(uid, "test")
				authUID, err := users.AuthenticateUser(user.Username, "test")
				So(authUID, ShouldEqual, uid)
				So(err, ShouldBeNil)
			})
			Convey("Creating a user with a password and try to authenticate with wrong username", func() {
				user := users.User{}
				user.Email = "user@example.com"
				uid, err := users.CreateUser(&user)
				err = users.AuthSetPassword(uid, "test")
				authUID, err := users.AuthenticateUser("wrong", "test")

				So(authUID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
			})
			Convey("Creating a user with a password and try to authenticate with wrong password", func() {
				user := users.User{}
				user.Email = "user@example.com"
				uid, err := users.CreateUser(&user)
				err = users.AuthSetPassword(uid, "test")
				authUID, err := users.AuthenticateUser("test", "wrong")

				So(authUID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing GetUser", func() {
			Convey("Creating user and get it", func() {
				u := users.User{}
				u.Username = "test"
				u.Email = "user@example.com"
				uid, _ := users.CreateUser(&u)
				user, err := users.GetUser(uid)

				So(user.Id, ShouldEqual, uid)
				So(user.Username, ShouldEqual, "test")
				So(err, ShouldBeNil)

				Convey("The created time should be close to now", func() {
					So(user.Created, ShouldHappenWithin,
						time.Duration(1)*time.Second, time.Now())
				})
			})
			Convey("Try to get non-existing user", func() {
				user, err := users.GetUser(0)

				So(user, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing GetAllUsers", func() {
			u1 := users.User{
				Email:    "u1@example.com",
				Username: "u1",
			}
			u2 := users.User{
				Email:    "u2@example.com",
				Username: "u2",
			}
			Convey("Getting all users with 0 users in the database", func() {
				users, err := users.GetAllUsersAt(1)

				So(len(users), ShouldEqual, 0)
				So(err, ShouldBeNil)
			})
			Convey("Creating 2 users and get them", func() {
				id1, err := users.CreateUser(&u1)
				if err != nil {
					panic(err.Error())
				}
				id2, err := users.CreateUser(&u2)
				if err != nil {
					panic(err.Error())
				}
				_, err = user_locations.Create(&user_locations.UserLocation{
					UserId:     id1,
					LocationId: 1,
				})
				if err != nil {
					panic(err.Error())
				}
				_, err = user_locations.Create(&user_locations.UserLocation{
					UserId:     id2,
					LocationId: 1,
				})
				if err != nil {
					panic(err.Error())
				}
				users, err := users.GetAllUsersAt(1)
				if err != nil {
					panic(err.Error())
				}

				So(len(users), ShouldEqual, 2)
				So(err, ShouldBeNil)
			})
		})
		Convey("Testing CheckEmail", func() {
			validAddresses := []string{
				"a@fablab.berlin",
				"piet@fun.mobi",
				"jimbo+foo@gmail.com",
			}
			u := users.User{}
			for _, addr := range validAddresses {
				u.Email = addr
				So(u.CheckEmail(), ShouldBeNil)
			}
		})
		Convey("Testing Update", func() {
			u := users.User{
				Username: "test",
				Email:    "test@example.com",
			}
			Convey("Updating user with invalid email should not work", func() {
				user := users.User{}
				user.Username = "hackerman"
				user.Email = "hacker@inspacefeels.good"
				user.Id, _ = users.CreateUser(&user)

				user.Email = ""
				err := user.Update()
				So(err, ShouldNotBeNil)

				user.Email = "invalidemail"
				err = user.Update()
				So(err, ShouldNotBeNil)

				user.Email = "invalid@email"
				err = user.Update()
				So(err, ShouldNotBeNil)

				user.Email = "invalid@email."
				err = user.Update()
				So(err, ShouldNotBeNil)
			})
			Convey("Creating a user and try to modify FirstName", func() {
				uid, _ := users.CreateUser(&u)
				user, _ := users.GetUser(uid)
				user.FirstName = "YOLO"
				user.LastName = "Maijer"
				err := user.Update()
				user, _ = users.GetUser(user.Id)

				So(err, ShouldBeNil)
				So(user.FirstName, ShouldEqual, "YOLO")
			})
			Convey("Try updating non-existing user", func() {
				err := u.Update()

				So(err, ShouldNotBeNil)
			})
			Convey("Update user with duplicate username, should return error", func() {
				user1 := users.User{}
				user1.Username = "hesus"
				user1.Email = "hesus@example.com"
				user2 := users.User{}
				user2.Username = "hrist"
				user2.Email = "hrist@example.com"
				users.CreateUser(&user1)
				users.CreateUser(&user2)
				user2.Username = user1.Username
				err := user2.Update()

				So(err, ShouldNotBeNil)
			})
			Convey("Update user with duplicate email, should return error", func() {
				user1 := users.User{}
				user1.Username = "hesus"
				user1.Email = "hesus@example.com"
				user2 := users.User{}
				user2.Username = "hrist"
				user2.Email = "hrist@example.com"
				users.CreateUser(&user1)
				users.CreateUser(&user2)
				user2.Email = user1.Email
				err := user2.Update()

				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing CreateUserPermission", func() {
			u := users.User{
				Username: "test",
			}
			Convey("Creating one UserPermission", func() {
				uid, _ := users.CreateUser(&u)
				err := user_permissions.Create(uid, 0)

				So(err, ShouldBeNil)
			})
		})
		Convey("Testing DeleteUserPermission", func() {
			u := users.User{
				Username: "test",
			}
			Convey("Creating one UserPermission and delete it", func() {
				uid, _ := users.CreateUser(&u)
				user_permissions.Create(uid, 0)
				err := user_permissions.Delete(uid, 0)

				So(err, ShouldBeNil)
			})
		})
		Convey("Testing UpdateUserPermission", func() {
			u := users.User{
				Username: "test",
				Email:    "joe@example.com",
			}
			Convey("Creating one User and update his permission", func() {
				uid, _ := users.CreateUser(&u)
				perms := &[]user_permissions.Permission{
					{UserId: uid, MachineId: 0},
					{UserId: uid, MachineId: 1},
					{UserId: uid, MachineId: 2},
					{UserId: uid, MachineId: 3},
				}
				err := user_permissions.Update(uid, 1, perms)

				So(err, ShouldBeNil)
			})
		})
		Convey("Testing that default row limit of 1000 doesn't get in the way", func() {
			uids := make(map[int64]struct{})
			n := 2000
			for i := 0; i < n; i++ {
				u := &users.User{
					Username: fmt.Sprintf("user_%v", i),
				}
				u.Email = u.Username + "@example.com"
				uid, err := users.CreateUser(u)
				if err != nil {
					panic(err.Error())
				}
				uids[uid] = struct{}{}
				ul := &user_locations.UserLocation{
					LocationId: 1,
					UserId:     uid,
				}
				if _, err := user_locations.Create(ul); err != nil {
					panic(err.Error())
				}
			}
			So(len(uids), ShouldEqual, n)
			us, err := users.GetAllUsersAt(1)
			if err != nil {
				panic(err.Error())
			}
			So(n, ShouldBeLessThanOrEqualTo, len(us))
			for _, u := range us {
				delete(uids, u.Id)
			}
			So(len(uids), ShouldEqual, 0)
		})
	})
}

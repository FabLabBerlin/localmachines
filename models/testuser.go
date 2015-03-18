package models

import (
	"errors"
	"strconv"
	"time"
)

var (
	TestUserList map[string]*TestUser
)

func init() {
	TestUserList = make(map[string]*TestUser)
	u := TestUser{"user_11111", "astaxie", "11111", TestProfile{"male", 20, "Singapore", "astaxie@gmail.com"}}
	TestUserList["user_11111"] = &u
}

type TestUser struct {
	Id       string
	Username string
	Password string
	Profile  TestProfile
}

type TestProfile struct {
	Gender  string
	Age     int
	Address string
	Email   string
}

func AddUser(u TestUser) string {
	u.Id = "user_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	TestUserList[u.Id] = &u
	return u.Id
}

/*
func GetUser(uid string) (u *TestUser, err error) {
	if u, ok := TestUserList[uid]; ok {
		return u, nil
	}
	return nil, errors.New("User not exists")
}
*/

func GetAllUsers() map[string]*TestUser {
	return TestUserList
}

func UpdateUser(uid string, uu *TestUser) (a *TestUser, err error) {
	if u, ok := TestUserList[uid]; ok {
		if uu.Username != "" {
			u.Username = uu.Username
		}
		if uu.Password != "" {
			u.Password = uu.Password
		}
		if uu.Profile.Age != 0 {
			u.Profile.Age = uu.Profile.Age
		}
		if uu.Profile.Address != "" {
			u.Profile.Address = uu.Profile.Address
		}
		if uu.Profile.Gender != "" {
			u.Profile.Gender = uu.Profile.Gender
		}
		if uu.Profile.Email != "" {
			u.Profile.Email = uu.Profile.Email
		}
		return u, nil
	}
	return nil, errors.New("User Not Exist")
}

func Login(username, password string) bool {
	for _, u := range TestUserList {
		if u.Username == username && u.Password == password {
			return true
		}
	}
	return false
}

func DeleteUser(uid string) {
	delete(TestUserList, uid)
}

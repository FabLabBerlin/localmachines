package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
	"golang.org/x/crypto/scrypt"
	"io"
)

// cf. http://stackoverflow.com/a/23039768/485185
const (
    PW_SALT_BYTES = 32
    PW_HASH_BYTES = 64
)

type LoginController struct {
	Controller
}

type UserDataResponse struct {
	Status    string
	UserId    int
	Username  string
	FirstName string
	LastName  string
	Email     string
	Admin     bool
	Staff     bool
	Member    bool
}

// Override our custom root controller's Prepare method as it is checking
// if we are logged in and we don't want that here at this point
func (this *LoginController) Prepare() {
	beego.Info("Skipping login check")
}

// Log in user, handle API /login request
func (this *LoginController) Login() {

	// Attempt to get stored session username
	sessUsername := this.GetSession(SESSION_FIELD_NAME_USERNAME)

	if sessUsername == nil {

		// If not set, user is not logged in
		if this.isUserValid() {

			// If user is valid, log in, save username in session
			reqUsername := this.GetString(REQUEST_FIELD_NAME_USERNAME)
			this.SetSession(SESSION_FIELD_NAME_USERNAME, reqUsername)

			// Save the user ID in session as well
			userId, err := this.getUserId(reqUsername)
			if err != nil {
				if beego.AppConfig.String("runmode") == "dev" {
					panic("User valid, but could not get user ID")
				}
				// This is really strange -
				// respond with error in case of bad spirits...
				this.serveErrorResponse("Invalid username or password")
			}
			this.SetSession(SESSION_FIELD_NAME_USER_ID, userId)
			beego.Info("User", reqUsername, "successfully logged in")

			this.serveUserData(userId, "ok")

		} else {
			// If user not valid, respond with error
			this.serveErrorResponse("Invalid username or password")
			beego.Info("Failed to authenticate user")
		}
	} else {

		// Get stored session user ID
		var userId int = this.GetSession(SESSION_FIELD_NAME_USER_ID).(int)
		this.serveUserData(userId, "logged")
	}
}

func (this *LoginController) getUserId(username string) (int, error) {
	userModel := models.User{}
	userModel.Username = username
	beego.Trace("Attempt to get user id for username ", username)
	o := orm.NewOrm()
	err := o.Read(&userModel, "Username")
	if err != nil {
		beego.Error("Could not get user ID with username", username, ":", err)
		return int(0), err
	}
	return userModel.Id, nil
}

func (this *LoginController) isPasswordValid(username, password string) (bool, error) {
	authModel := models.Auth{}
	o := orm.NewOrm()
	err := o.Raw("SELECT hash, salt FROM auth INNER JOIN user ON auth.user_id = user.id WHERE user.username = ?",
		username).QueryRow(&authModel)
	if err != nil {
		beego.Error("Could not read into AuthModel:", err)
		return false, err
	}
	authModelSalt, err := hex.DecodeString(authModel.Salt)
	if err != nil {
		beego.Error("Could not decode authModel.Salt:", err)
		return false, err
	}
	hash, err := this.hash(password, authModelSalt)
	if err != nil {
		beego.Error("Could not calculate hash:")
		return false, err
	}
	beego.Info("calculated hash:", hex.EncodeToString(hash))
	return hex.EncodeToString(hash) == authModel.Hash, err
}

func (this *LoginController) hash(password string, salt []byte) ([]byte, error) {
    hash, err := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, PW_HASH_BYTES)
    if err != nil {
        return []byte{}, err
    }
	return hash, nil
}

func (this *LoginController) createSalt() ([]byte, error) {
	salt := make([]byte, PW_SALT_BYTES)
	_, err := io.ReadFull(rand.Reader, salt)
	return salt, err
}

func (this *LoginController) getUserData(userId int) (*models.User, error) {

	// Create userModel struct for user data storage
	userModel := models.User{Id: userId}

	// Fill user model with data from database
	o := orm.NewOrm()
	beego.Info("Getting user data...")
	// Could be replaced with plain Read, but this gets only the data we need for now
	err := o.Raw("SELECT username, first_name, last_name, email FROM user WHERE id = ?",
		userId).QueryRow(&userModel)
	if err != nil {
		beego.Error("There was an error while getting user data:", err)
		return &userModel, err
	}

	// Return the user model to someone else
	beego.Info("Success getting user data")
	return &userModel, nil

}

func (this *LoginController) isUserValid() bool {
	// Get request variables
	username := this.GetString("username")
	password := this.GetString("password")

	beego.Trace("POST username:", username)
	beego.Trace("POST password:", password)

	valid, err := this.isPasswordValid(username, password)
	return valid && err == nil
}

// Serves user data on successful login or when user is already logged in
func (this *LoginController) serveUserData(userId int, status string) {

	userModel, err := this.getUserData(userId)

	if err != nil {

		beego.Error("Failed to get user data:", err)

	} else {

		// We need to get user roles from database as well
		userRolesModel := models.UserRoles{
			UserId: userId,
			Admin:  false,
			Staff:  false,
			Member: false}
		o := orm.NewOrm()
		err = o.Read(&userRolesModel)

		if err != nil {
			beego.Error("Could not get user roles for user ID", userId, ", error:", err)
			// We can continue here as if we can't get user roles, we have none probably
		}

		// Fill out response object
		userDataResponse := UserDataResponse{
			Status:    status,
			UserId:    userModel.Id,
			Username:  userModel.Username,
			FirstName: userModel.FirstName,
			LastName:  userModel.LastName,
			Email:     userModel.Email,
			Admin:     userRolesModel.Admin,
			Staff:     userRolesModel.Staff,
			Member:    userRolesModel.Member}

		this.Data["json"] = userDataResponse
		this.ServeJson()

	}
}

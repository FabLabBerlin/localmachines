package models

// This is being returned when client awaits a JSON message like
// {"Status":"ok"}
type StatusResponse struct {
	Status string
}

// Used to return
// {"Status":"error", "Message":"Error message, what went wrong"}
type ErrorResponse struct {
	Status  string
	Message string
}

// Used to return user ID on successful login
// {"Status":"ok", 		"UserId":2}
// {"Status":"logged", 	"UserId":2}
type LoginResponse struct {
	Status     string
	UserId     int64
	LocationId int64
}

// Used to return activation ID after it has been created
// {"ActivationId":39}
type ActivationCreateResponse struct {
	ActivationId int64
}

// Used to return machine ID after a machine has been created
// {"MachineId":101}
type MachineCreatedResponse struct {
	MachineId int64
}

type UserNamesResponse struct {
	Users []UserNameResponse
}

// Used to return only the full name of an user
// {"UserId": 2, "FirstName": "Milov", "LastName": "Milovič"}
type UserNameResponse struct {
	UserId    int64
	FirstName string
	LastName  string
}

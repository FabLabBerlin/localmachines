// This file contains only specific response types that don't need to be registered
package models

// This is being returned when client awaits a JSON message like
// {"Status":"ok"}
type StatusResponse struct {
	Status 	string
}

// Used to return
// {"Status":"error", "Message":"Error message, what went wrong"}
type ErrorResponse struct {
	Status 	string
	Message string
}

// Used to return user ID on successful login
// {"Status":"ok", 		"UserId":2}
// {"Status":"logged", 	"UserId":2}
type LoginResponse struct {
	Status 	string
	UserId  int
}
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go/test/data"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

type UsersHandlers struct {
}

func GetUsersHandlerFunc() *UsersHandlers {
	return &UsersHandlers{}
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
}

var collection = data.Open_Collection(data.GetMongoDBFunctions().Mongo_Connect(), "USERS")

func (userhandler *UsersHandlers) AddUser(rw http.ResponseWriter, r *http.Request) {

	user_details := r.Context().Value(KeyUser{}).(data.User)

	// if email is valid
	if !data.Isvalid(user_details.Email) {
		http.Error(rw, "the given email is invalid.", http.StatusBadRequest)
		return
	}

	// checking if user email or phone no already registered
	ctx := context.TODO()
	email_count, err := collection.CountDocuments(ctx, bson.M{"email": user_details.Email})

	if err != nil {
		http.Error(rw, fmt.Sprintf("Error while retriving user count : %v", err), http.StatusBadRequest)
		return
	}
	phone_count, err := collection.CountDocuments(ctx, bson.M{"phone": user_details.Phone})
	if err != nil {
		http.Error(rw, fmt.Sprintf("Error while retriving user count : %v", err), http.StatusBadRequest)
		return
	}

	// if count of either is more than one then return error
	if email_count > 0 || phone_count > 0 {
		http.Error(rw, "Email or phone no. already registered.", http.StatusBadRequest)
		return
	}
	fmt.Println("Add User request queued.")

	// adding to user db
	accessToken, refreshToken, err := data.AddUser(&user_details)

	// error handling if error occured while adding user
	if err != nil {
		http.Error(rw, fmt.Sprintf("%v", err), http.StatusInternalServerError)
	}
	response := &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(response)

}

func (userhandler *UsersHandlers) Login(rw http.ResponseWriter, r *http.Request) {

	// check if request contains username and password

	userObj := &data.UserLogin{}

	err := json.NewDecoder(r.Body).Decode(&userObj)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Unable to unmarshal the given json : %v", err), http.StatusBadRequest)
		return
	}

	err = userObj.LoginUserValidator()
	if err != nil {
		http.Error(rw, fmt.Sprintf("Validatin error : %v", err), http.StatusBadRequest)
		return
	}
	fmt.Println("User login request received with validation checks-true")
	accessToken, refreshToken, err := data.LoginUser(userObj)

	if err != nil {
		http.Error(rw, fmt.Sprintf("Error : %v", err), http.StatusBadRequest)
		return
	}
	response := &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(response)

}

type KeyUser struct{}

func (userhandler *UsersHandlers) MiddlewaresUserHandlers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// creating a instance of our user model and then decoding json to something json understands.
		user_obj := &data.User{}
		err := json.NewDecoder(r.Body).Decode(&user_obj)

		if err != nil {
			http.Error(rw, fmt.Sprintf("Unable to unmarshal the given json : %v", err), http.StatusBadRequest)
			return
		}
		// CAN ADD VALIDATION AND AUTHENTICATION HERE IN FUTURE

		// checking struct validation
		err = user_obj.UserValidator()
		if err != nil {
			http.Error(rw, fmt.Sprintf("Validation error : %v", err), http.StatusBadRequest)
			return
		}

		// adding the user object in context for the coming functions to read it from there
		ctx := context.WithValue(r.Context(), KeyUser{}, *user_obj)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)

	})
}

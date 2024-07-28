package data

import (
	"context"
	"encoding/json"
	"fmt"
	"go/test/helpers"
	"io"
	"net/mail"
	"reflect"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	FirstName    string             `json:"firstname" validate:"required,min=2,max=100"`
	LastName     string             `json:"lastname"`
	Password     string             `json:"password" validate:"required,min=6,max=100"`
	Email        string             `json:"email" validate:"required,email"`
	Phone        string             `json:"phone"`
	Role         string             `json:"role" validate:"required,oneof=ADMIN CUSTOMER MANAGER BARISTA TESTER AUDITOR"`
	UserID       int                `json:"user_id"`
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
	CreatedOn    time.Time          `json:"created_on"`
	UpdatedOn    time.Time          `json:"updated_on"`
}

// used for schema check in login
type UserLogin struct {
	Username string `json:"username" validate:"required,username"`
	Password string `json:"password" validate:"required"`
}

var collection = Open_Collection(GetMongoDBFunctions().Mongo_Connect(), "USERS")

func Get_User(user *User) {
	ctx := context.TODO()
	// var userDef User

	data, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Print("Error occured while getting user")
	}
	fmt.Print(data)
}

func (user *User) FromJSON(data io.Reader) error {
	decoded_data := json.NewDecoder(data)
	fmt.Println(decoded_data)

	val := reflect.TypeOf(decoded_data.Decode(user))
	fmt.Println("changes", val)
	return decoded_data.Decode(user)
}

// validatin registration for signup
func (user *User) UserValidator() error {

	validate := validator.New()
	validate.RegisterValidation("email", CustomUserValidationFunc)
	return validate.Struct(user)
}

// validatin registration for login

func (loginUser *UserLogin) LoginUserValidator() error {

	validate := validator.New()
	validate.RegisterValidation("username", CustomUserValidationFunc)
	return validate.Struct(loginUser)
}
func CustomUserValidationFunc(fl validator.FieldLevel) bool {

	reg := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	matches := reg.FindAllString(fl.Field().String(), -1) // -1 returns all matches from the string

	return len(matches) == 1
}
func (user *User) ToJSON(w io.Writer) error {
	newEncoder := json.NewEncoder(w)
	print(newEncoder.Encode(user))
	return newEncoder.Encode(user)
}
func Isvalid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
func AddUser(user_details *User) (accessToken string, refreshToken string, err error) {

	// hashing password with brypt
	hashed_pass, err := generateHashedPass(user_details.Password)
	if err != nil {
		return "", "", err
	}
	user_details.Password = hashed_pass

	// giving a created, updated and new OID to mongo DB user.
	user_details.CreatedOn, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user_details.UpdatedOn, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user_details.ID = primitive.NewObjectID()

	// create operation is same for all routes as its just addding, so using mongo functin. Note get operation is different for all types of routes in mongo so create a function.
	all_func := GetMongoDBFunctions()
	max_id, err := getNextUser()

	// getNextUser always returns -1 or err in case of an error occured
	if err != nil || max_id == -1 {
		return "", "", err
	}
	// adding max_id with one in both cases if no doc exists then first doc has 1 else it increments it.
	user_details.UserID = max_id + 1
	all_func.Insert_one_mdb("USERS", &user_details)

	full_name := user_details.FirstName + " " + user_details.LastName
	accessToken, refreshToken, err = helpers.GenerateToken(user_details.Email, full_name, user_details.Role, user_details.UserID)

	if err != nil {
		return "", "", err
	}
	err = UpdateAll(user_details.UserID, accessToken, refreshToken)

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil

}

func LoginUser(user_details *UserLogin) (string, string, error) {

	userFromDB, err := Get_UserByEmailorID(user_details.Username, -1)
	// fmt.Println(UserSchema)
	if err != nil {
		return "", "", fmt.Errorf("error:%v", err)
	}

	hashedPass := userFromDB.Password
	isSame := MatchPassword(hashedPass, string(user_details.Password))
	if !isSame {
		return "", "", fmt.Errorf("wrong Password")
	}
	fmt.Print(userFromDB)
	full_name := userFromDB.FirstName + " " + userFromDB.LastName
	accessToken, refreshToken, err := helpers.GenerateToken(userFromDB.Email, full_name, userFromDB.Role, userFromDB.UserID)
	if err != nil {
		return "", "", fmt.Errorf("error generating token %v", err)
	}

	err = UpdateAll(userFromDB.UserID, accessToken, refreshToken)

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// NOTE: need to change the mongo  code to remove redundancy

func getNextUser() (int, error) {
	// query to get the maximum user id given to the user.
	pipeline := mongo.Pipeline{
		{{"$group", bson.D{{"_id", nil}, {"max_userid", bson.D{{"$max", "$userid"}}}}}},
		{{"$project", bson.D{{"_id", 0}, {"max_userid", 1}}}},
	}

	// Execute aggregation pipeline
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		fmt.Printf("error executing aggregation pipeline: %v", err)
		return -1, err
	}
	// defer cursor.Close(context.TODO())
	// convert the incoming cursor to result struct that has int max user is.
	var result struct {
		MaxUserID int `bson:"max_userid"`
	}
	if cursor.Next(context.Background()) {
		err := cursor.Decode(&result)
		if err != nil {
			return -1, fmt.Errorf("error occured while converting mongo db response.%v", err)
		}
		// return max user id if present.
		return result.MaxUserID, nil
	} else {
		fmt.Println("No documents found")
		return 0, nil
	}
}

// hashing functinos
func generateHashedPass(plaintext string) (string, error) {

	plainbytes := []byte(plaintext)

	// it taked bytes as plain text so we convert it and send min cost. hashing time is calculated by 2^cost,
	// higher cost higher the hashing time as more number of iterations taken as no of iterations = cost.
	cipherText, err := bcrypt.GenerateFromPassword(plainbytes, bcrypt.MinCost)

	if err != nil {
		return "", fmt.Errorf("error occured while hashing %v", err)
	}
	return string(cipherText), nil
}

func MatchPassword(cipherText string, plaintext string) bool {
	// complaring returns err if error is nil  then matches
	err := bcrypt.CompareHashAndPassword([]byte(cipherText), []byte(plaintext))
	return err == nil
}

func UpdateAll(user_id int, accessToken string, refreshToken string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"userid": user_id}
	current_time, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	fmt.Println(reflect.TypeOf(current_time))
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	update := bson.M{
		"$set": bson.M{
			"Token":        accessToken,
			"refreshToken": refreshToken,
			"updatedon":    current_time,
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no user found with the given ID")
	}
	// Updating values

	return nil
}

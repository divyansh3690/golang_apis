package data

import (
	"fmt"
	"testing"
)

func TestCheckValidation(t *testing.T) {
	p := &Product{Name: "Divaynsh ", Price: 1.7, SKU: "absxh-as-as"}

	err := p.Validator()

	if err != nil {
		t.Fatal(err)
	}

}

func TestBearerToken(t *testing.T) {

}

// test case for user data
func TestGetUser(t *testing.T) {
	email := "divyansh3690@gmail.com"
	user_id := 0

	user_test, err := Get_UserByEmailorID(email, user_id)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(user_test)
}

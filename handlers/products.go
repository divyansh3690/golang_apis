package handlers

import (
	"encoding/json"
	"fmt"
	"go/test/data"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
)

type Product struct {
}

func GetProductsHandlerfunc() *Product {
	return &Product{}
}

func (prod *Product) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		prod.GetReqProd(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		prod.AddProduct(rw, r)
		return
	}

	// two things to keep in mind in put you get id form url and body that you wanna change

	if r.Method == http.MethodPut {
		print("Inside put")

		regexInpt := "([0-9]+)"
		reg := regexp.MustCompile(regexInpt)
		all_string_match := reg.FindAllStringSubmatch(r.URL.Path, -1)

		fmt.Println(all_string_match)

		if len(all_string_match) != 1 {
			fmt.Println("Error unmarshing the url path as more than one id")

			http.Error(rw, "Invalid url path", http.StatusBadRequest)
		}

		if len(all_string_match[0]) != 2 {
			fmt.Println(all_string_match[0])
			fmt.Println(all_string_match[1])
			fmt.Println("Error unmarshing the url path as more than one id")

			http.Error(rw, "Invalid url path", http.StatusBadRequest)
		}

		idProd := all_string_match[0][1]
		id, err := strconv.Atoi(idProd)
		fmt.Print(id)
		if err != nil {
			fmt.Println("Error as couldnt convert int id to string in PUT")
		}
		prod.UpdateProduct(id, rw, r)
		return

	}

	// handling other requests
	rw.WriteHeader(http.StatusMethodNotAllowed)

}

// HANDLES GET REQ AND CONVERTS DATA TO JSON
func (prod *Product) GetReqProd(rw http.ResponseWriter, r *http.Request) {
	dt := data.GetProducts()
	fmt.Println(reflect.TypeOf(dt))
	value, err := json.Marshal(dt)
	if err != nil {
		fmt.Print("Error occurede while converting data to json", err)
	}
	rw.Write(value)
}

//  Adds data into the database

func (prod *Product) AddProduct(rw http.ResponseWriter, r *http.Request) {

	// we need to create obj of our interface as
	prod_obj := &data.Product{}

	err := prod_obj.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Unable to unmarshal the given json : %v", err), http.StatusBadRequest)

	}
	data.AddProduct(prod_obj)
	fmt.Printf("DATA: %#v", prod_obj)
}

func (prod *Product) UpdateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	prod_obj := &data.Product{}
	err := prod_obj.FromJSON(r.Body)
	if err != nil {
		fmt.Println("Error occured, unable to unmarshal json")
	}
	err2 := data.UpdateProd(id, prod_obj)
	if err2 != nil {
		fmt.Print(err2)
	}
}

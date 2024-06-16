package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go/test/data"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Product struct {
}

func GetProductsHandlerfunc() *Product {
	return &Product{}
}

// NOTE WE need to keep this signature of (rw, request ) as any method handling request should have the signature of SERVE HTTP
// HANDLES GET REQ AND CONVERTS DATA TO JSON
func (prod *Product) GetReqProd(rw http.ResponseWriter, r *http.Request) {
	dt := data.GetProducts()
	// fmt.Println(reflect.TypeOf(dt))
	value, err := json.Marshal(dt)
	if err != nil {
		fmt.Print("Error occurede while converting data to json", err)
	}
	rw.Write(value)
}

//  Adds data into the database

func (prod *Product) AddProduct(rw http.ResponseWriter, r *http.Request) {

	// we need to create obj of our interface as

	prod_obj := r.Context().Value(KeyProduct{}).(data.Product)

	data.AddProduct(&prod_obj)
	fmt.Printf("DATA: %#v", prod_obj)
}

func (prod *Product) UpdateProduct(rw http.ResponseWriter, r *http.Request) {

	id_str := mux.Vars(r)
	// mux vars gives us variables passed
	id, err := strconv.Atoi(id_str["id"])
	if err != nil {
		http.Error(rw, "Unable to filer out id from url ", http.StatusBadGateway)
	}
	prod_obj := r.Context().Value(KeyProduct{}).(data.Product)

	err2 := data.UpdateProd(id, &prod_obj)
	if err2 != nil {
		fmt.Print(err2)
	}
}

func (prod *Product) RemoveProduct(rw http.ResponseWriter, r *http.Request) {
	id_str := mux.Vars(r)
	id, err := strconv.Atoi(id_str["id"])
	if err != nil {
		fmt.Println("Error occured while filtering id from request url", err)
	}
	err = data.RemoveProd(id)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
	}
}

type KeyProduct struct{}

func (prod *Product) MiddlewaresHandlers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		prod_obj := &data.Product{}

		err := prod_obj.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, fmt.Sprintf("Unable to unmarshal the given json : %v", err), http.StatusBadRequest)
			return
		}
		fmt.Println("VSALE:,", prod_obj)
		err = prod_obj.Validator()
		if err != nil {
			http.Error(rw, fmt.Sprintf("Unable to unmarshal the given json in validation: %v", err), http.StatusBadRequest)
			return
		}
		fmt.Println(prod_obj)
		ctx := context.WithValue(r.Context(), KeyProduct{}, *prod_obj)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)

	})
}

package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"cost" validate:"min=1.0"`
	SKU         string  `json:"category" `
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

func (p *Product) ToJSON(w io.Writer) error {
	newEncoder := json.NewEncoder(w)
	print(newEncoder.Encode(p))
	return newEncoder.Encode(p)
}

func (p *Product) Validator() error {

	validate := validator.New()
	validate.RegisterValidation("sku", customValidationFunc)
	return validate.Struct(p)
}

func customValidationFunc(fl validator.FieldLevel) bool {

	reg := regexp.MustCompile("[a-z]+-[a-z]+-[a-z]")
	matches := reg.FindAllString(fl.Field().String(), -1) // -1 returns all matches from the string

	return len(matches) == 1

}

// This will be of type product because from json takes io writer which is a body format and converts to a format Product that i will add to next index of my array.
func (p *Product) FromJSON(data io.Reader) error {
	decoded_data := json.NewDecoder(data)
	fmt.Println("Inside decoding")
	return decoded_data.Decode(p)

}

func GetProducts() []*Product {
	fmt.Println("HERE")
	prod_mongo_db := GetMongoDBFunctions()
	prod_list, err := prod_mongo_db.GetAll_mdb("products_data")
	if err != nil {
		fmt.Println(err)
	}
	return prod_list
}

func AddProduct(p *Product) {
	p.ID = getNextProd()
	fmt.Println(p)
	// productList = append(productList, p)
	prod_mongo_db := GetMongoDBFunctions()
	prod_mongo_db.Insert_one_mdb("products_data", p)

}

func UpdateProd(id int, p *Product) error {

	err := GetMongoDBFunctions().Update_one_mdb("products_data", p, id)
	if err != nil {
		return err
	}
	return nil

}

func RemoveProd(id int) error {
	err := GetMongoDBFunctions().Delete_one_mdb("products_data", id)
	if err != nil {
		return err
	}

	return nil

}

// UNUSED AS WE NEED TO DEFINE A HANDLER FOR GET PROD BY ID

func GetProdByID(id int) (p *Product, err error) {

	_, filteredProd, err := GetMongoDBFunctions().GetProdByID_mdb("products_data", id)
	if err != nil {
		return nil, fmt.Errorf("product not found or out of bound")
	}
	return filteredProd, nil
}

func getNextProd() int {
	// last_prod := productList[len(productList)-1]

	newProd_list, err := GetMongoDBFunctions().GetAll_mdb("products_data")
	if err != nil {
		fmt.Println("Error occured while getting the collection:", err)
	}
	var maxID = 0
	for _, prod := range newProd_list {

		if prod.ID > maxID {
			maxID = prod.ID
		}
	}

	return maxID + 1
}

// var productList = []*Product{}

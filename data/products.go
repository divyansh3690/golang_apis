package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"cost" validate:"min=1.0"`
	SKU         string  `json:"category" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
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
	return productList
}

func AddProduct(p *Product) {
	p.ID = getNextProd()

	productList = append(productList, p)

}

func UpdateProd(id int, p *Product) error {
	_, index, err := getProdByID(id)
	if err != nil {
		return err
	}
	fmt.Print(index)

	// p.ID = id
	productList[index] = p

	return nil

}

func RemoveProd(id int) error {
	_, index, err := getProdByID(id)
	if err != nil {
		return fmt.Errorf("error ocurred while finding the product: %w", err)
	}

	productList[index] = nil
	if index == 0 {
		fmt.Print(index)
		productList = productList[index+1:]
	} else if index > 0 && index < len(productList) {
		productList = append(productList[:index], productList[index+1:]...)
	} else {
		productList = productList[:index]
	}

	return nil

}

func getProdByID(id int) (p *Product, index int, err error) {
	for index, product := range productList {
		// fmt.Println(product, '\n', product.ID)
		if product.ID == id {

			return product, index, nil
		}
	}
	return nil, -1, fmt.Errorf("product not found or out of bound")
}

func getNextProd() int {
	last_prod := productList[len(productList)-1]
	return last_prod.ID + 1
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}

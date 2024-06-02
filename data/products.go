package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"cost"`
	SKU         string  `json:"category"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
	newEncoder := json.NewEncoder(w)

	return newEncoder.Encode(p)
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

func getProdByID(id int) (p *Product, index int, err error) {
	for index, product := range productList {
		// fmt.Println(product, '\n', product.ID)
		if product.ID == id {

			return product, index, nil
		}
	}
	return nil, -1, fmt.Errorf("NOT FOUND OR OUT OF BOUND")
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

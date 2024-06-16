package data

import "testing"

func TestCheckValidation(t *testing.T) {
	p := &Product{Name: "Divaynsh ", Price: 1.7, SKU: "absxh-as-as"}

	err := p.Validator()

	if err != nil {
		t.Fatal(err)
	}

}

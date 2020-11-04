package data

import "testing"

func TestCheckValidation(t *testing.T) {
	p := &Product{}
	p.Name = "band"
	p.Price = 10
	p.SKU = "abc-qwe-asdgf"

	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}

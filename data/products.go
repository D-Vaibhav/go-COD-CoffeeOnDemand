package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreratedOn  string  `json:"-"` //wont appear for JSON
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

func (p *Product) FromJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(p)
}

func (p *Product) Validate() error {
	// creating new validator object
	validate := validator.New()

	// to validate some custom input filed text
	validate.RegisterValidation("sku", validateSKU)

	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	// the valid format of SKU is: akbkabc-adkckba-abcjhavacfa
	re := regexp.MustCompile(`[a-z]-+[a-z]+-[a-z]+`)
	matchedString := re.FindAllString(fl.Field().String(), -1)

	if len(matchedString) != 1 {
		return false
	}
	return true
}

type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(p)
}

// we'll using this function to get all of the products
func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}
	p.ID = id
	productList[pos] = p
	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found")

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrProductNotFound
}

func getNextID() int {
	lastProduct := productList[len(productList)-1]
	return lastProduct.ID + 1
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Latte coffee",
		Price:       2.45,
		SKU:         "lat123",
		CreratedOn:  time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Espresso coffee",
		Price:       1.99,
		SKU:         "exp123",
		CreratedOn:  time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}

package product

import (
	"fmt"
	"net/http"

	"alwis.dev/selectify/internal/model"
)

func (c *controller) GetProduct(w http.ResponseWriter, r *http.Request, s *model.LoggedInSession) {
	fmt.Println("GetProduct")
	fmt.Println(s.IpAddress)
	fmt.Println(*s.DeviceId)
	return
}

func (c *controller) AddNewProduct(w http.ResponseWriter, r *http.Request, s *model.LoggedInSession) {
	fmt.Println("PostProduct")
}

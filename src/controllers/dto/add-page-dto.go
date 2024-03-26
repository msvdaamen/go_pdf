package dto

import (
	"encoding/json"
	"net/http"
)

type AddPageDto struct {
	Pdf []byte `form:"pdf"`
}

func ValidateAddPage(r *http.Request) (addPageDto AddPageDto, err error) {
	err = json.NewDecoder(r.Body).Decode(&addPageDto)
	if err != nil {
		return
	}
	return
}

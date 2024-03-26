package dto

import (
	"encoding/json"
	"net/http"
)

type MergePdfsDto struct {
	Pdfs [][]byte `json:"pdfs"`
}

func ValidateMergePdfs(r *http.Request) (mergePdfsDto *MergePdfsDto, err error) {
	err = json.NewDecoder(r.Body).Decode(mergePdfsDto)
	if err != nil {
		return
	}
	return
}

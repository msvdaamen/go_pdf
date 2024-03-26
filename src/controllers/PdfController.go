package controllers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	pdf2 "print_com/src/modules/pdf"
	"strconv"
)

func Create(w http.ResponseWriter, r *http.Request) {
	pdf := pdf2.Create()
	pdf.AddPage()
	pdf.SetLineWidth(2)
	pdf.SetLineType("dashed")
	pdf.Line(10, 30, 585, 30)
	file, err := pdf.ToFile()
	if err != nil {
		return
	}
	defer file.Close()
	http.ServeFile(w, r, file.Name())
}

func AddPage(w http.ResponseWriter, r *http.Request) {
	fileBytes, err := getFileBytesFromRequest(r)
	if err != nil {
		return
	}
	pdf := pdf2.Create()
	rs := io.ReadSeeker(bytes.NewReader(fileBytes))

	pageNumber := 1
	for {
		errT := pdf.ImportPage(&rs, pageNumber)
		if errT != nil {
			break
		}
		pageNumber++
	}
	pdf.AddPage()
	pdf.SetLineWidth(5)
	pdf.SetLineType("dotted")
	pdf.Line(10, 30, 585, 30)

	file, err := pdf.ToFile()
	if err != nil {
		return
	}
	defer file.Close()
	http.ServeFile(w, r, file.Name())
}

func RemovePage(w http.ResponseWriter, r *http.Request) {
	pageNumberToRemove, err := strconv.Atoi(r.PathValue("pageNumber"))
	if err != nil {
		return
	}
	fileBytes, err := getFileBytesFromRequest(r)
	pdf := pdf2.Create()
	rs := io.ReadSeeker(bytes.NewReader(fileBytes))

	pageNumber := 1
	for {
		if pageNumber == pageNumberToRemove {
			pageNumber++
			continue
		}
		errT := pdf.ImportPage(&rs, pageNumber)
		if errT != nil {
			break
		}
		pageNumber++
	}
	file, err := pdf.ToFile()
	if err != nil {
		return
	}
	defer file.Close()
	http.ServeFile(w, r, file.Name())
}

func ReorderPages(w http.ResponseWriter, r *http.Request) {
	pageFrom, err := strconv.Atoi(r.PostFormValue("from"))
	pageTo, err := strconv.Atoi(r.PostFormValue("to"))
	if err != nil {
		return
	}
	fileBytes, err := getFileBytesFromRequest(r)
	pdf := pdf2.Create()
	rs := io.ReadSeeker(bytes.NewReader(fileBytes))
	var pages []int
	pageNumber := 1
	for {
		template, err := pdf.GetPage(&rs, pageNumber)
		pages = append(pages, template)
		if err != nil {
			break
		}
		pageNumber++
	}
	newPageOrder := make([]int, len(pages)-1)
	copy(newPageOrder, pages)
	newPageOrder[pageTo-1] = pages[pageFrom-1]
	newPageOrder[pageFrom-1] = pages[pageTo-1]
	for _, page := range newPageOrder {
		pdf.UseImportedTemplate(page)
	}
	file, err := pdf.ToFile()
	if err != nil {
		return
	}
	defer file.Close()
	http.ServeFile(w, r, file.Name())
}

func MergePdfs(w http.ResponseWriter, r *http.Request) {
	fileBytesSlice, err := getFilesBytesFromRequest(r)
	if err != nil {
		return
	}
	pdf := pdf2.Create()
	for _, fileBytes := range fileBytesSlice {
		rs := io.ReadSeeker(bytes.NewReader(fileBytes))
		pageNumber := 1
		for {
			err := pdf.ImportPage(&rs, pageNumber)
			if err != nil {
				break
			}
			pageNumber++
		}
	}
	file, err := pdf.ToFile()
	if err != nil {
		return
	}
	defer file.Close()
	http.ServeFile(w, r, file.Name())
}

func SplitPdf(w http.ResponseWriter, r *http.Request) {
	fileBytes, err := getFileBytesFromRequest(r)
	if err != nil {
		return
	}

	rs := io.ReadSeeker(bytes.NewReader(fileBytes))
	archive, err := os.Create("archive.zip")
	if err != nil {
		return
	}
	defer archive.Close()
	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()
	pageNumber := 1
	for {
		err := func() (err error) {
			pdf := pdf2.Create()
			err = pdf.ImportPage(&rs, pageNumber)
			if err != nil {
				return
			}
			pageNumber++
			file, err := pdf.ToFile()
			if err != nil {
				return
			}
			fileInfo, err := file.Stat()
			if err != nil {
				return
			}
			header, err := zip.FileInfoHeader(fileInfo)
			header.Name = fmt.Sprintf("pdf-%v.pdf", pageNumber-1)
			if err != nil {
				return
			}
			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return
			}
			_, err = io.Copy(writer, file)
			if err != nil {
				return
			}
			defer file.Close()
			return
		}()
		if err != nil {
			break
		}
	}
	w.Header().Set("Content-Type", "application/zip")
	http.ServeFile(w, r, archive.Name())
}

func getFileBytesFromRequest(r *http.Request) (fileBytes []byte, err error) {
	//fileType = r.PostFormValue("type")
	file, _, err := r.FormFile("pdf")
	if err != nil {
		return
	}
	fileBytes, err = io.ReadAll(file)
	if err != nil {
		return
	}
	//detectedFileType := http.DetectContentType(fileBytes)
	//fmt.Println(detectedFileType)
	defer file.Close()
	return
}

func getFilesBytesFromRequest(r *http.Request) (fileBytesSlice [][]byte, err error) {
	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		return
	}
	files := r.MultipartForm.File["pdfs"]
	for _, multipart := range files {
		func() {
			file, err := multipart.Open()
			fileBytes, err := io.ReadAll(file)
			if err != nil {
				return
			}
			fileBytesSlice = append(fileBytesSlice, fileBytes)
			defer file.Close()
		}()
	}
	return
}

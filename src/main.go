package main

import (
	"log"
	"net/http"
	"os"
	"print_com/src/controllers"
)

func main() {

	router := http.NewServeMux()

	router.HandleFunc("GET /", func(writer http.ResponseWriter, request *http.Request) {
		//fmt.Fprintf(writer, "Hello world")
		zipFile, err := os.OpenFile("archive.zip", os.O_RDONLY, 0666)
		if err != nil {
			panic(err)
		}
		writer.Header().Set("Content-Type", "application/zip")
		http.ServeFile(writer, request, zipFile.Name())
	})
	router.HandleFunc("POST /pdf/merge", controllers.MergePdfs)
	router.HandleFunc("POST /pdf/{fileId}", controllers.Create)
	router.HandleFunc("POST /pdf/{fileId}/add-page", controllers.AddPage)
	router.HandleFunc("POST /pdf/{fileId}/reorder-page", controllers.ReorderPages)
	router.HandleFunc("POST /pdf/{fileId}/split", controllers.SplitPdf)
	router.HandleFunc("DELETE /pdf/{fileId}/remove-page/{pageNumber}", controllers.RemovePage)
	err := http.ListenAndServe(":8080", router)
	log.Fatal(err)
}

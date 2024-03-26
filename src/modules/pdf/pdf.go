package pdf

import (
	"github.com/pkg/errors"
	"github.com/signintech/gopdf"
	"io"
	"os"
)

type Pdf struct {
	pdf gopdf.GoPdf
}

func (pdf *Pdf) AddPage() {
	pdf.pdf.AddPage()
}

func (pdf *Pdf) ImportPage(sourceStream *io.ReadSeeker, pageNumber int) (err error) {
	template, err := pdf.GetPage(sourceStream, pageNumber)
	if err != nil {
		return
	}
	pdf.AddPage()
	pdf.pdf.UseImportedTemplate(template, 0, 0, 595, 842)
	return
}

func (pdf *Pdf) UseImportedTemplate(template int) {
	pdf.AddPage()
	pdf.pdf.UseImportedTemplate(template, 0, 0, 595, 842)
}

func (pdf *Pdf) SetLineWidth(width float64) {
	pdf.pdf.SetLineWidth(width)
}

func (pdf *Pdf) SetLineType(lineType string) {
	pdf.pdf.SetLineType(lineType)
}

func (pdf *Pdf) Line(x1, y1, x2, y2 float64) {
	pdf.pdf.Line(x1, y1, x2, y2)
}

func (pdf *Pdf) ToFile() (file *os.File, err error) {
	file, err = os.CreateTemp("", "temp")
	if err != nil {
		return
	}
	err = pdf.pdf.WritePdf(file.Name())
	if err != nil {
		return
	}
	return
}

/*
*
Wrapper function to catch the panic from the gopdf library and return error.
*/
func (pdf *Pdf) GetPage(sourceStream *io.ReadSeeker, pageNumber int) (page int, err error) {
	// Defer a function that will recover from any panic
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Invalid page number")
		}
	}()
	page = pdf.pdf.ImportPageStream(sourceStream, pageNumber, "/MediaBox")

	return
}

func Create() *Pdf {
	pdf := Pdf{pdf: gopdf.GoPdf{}}
	pdf.pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	return &pdf
}

package main

import (
	"github.com/jung-kurt/gofpdf"
	"golang.org/x/exp/errors/fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type book struct {
	name                string
	urlTemplate         string
	imageType           string
	firstPage, lastPage int
}

var books = []book{
	{"I_volume_Weber_Wellstein", "http://www.mathesis.ru/books/weber13/mathesis_weber13_%03d.jpg", "jpg", 1, 696},
	{"II_volume_Weber_Wellstein", "http://www.mathesis.ru/books/weber22/mathesis_weber22_%03d.jpg", "jpg", 3, 382},
	//{"II_volume_Trigonometry_Weber_Wellstein", "http://www.mathesis.ru/books/weber31/mathesis_weber31_%03d.jpg", "jpg", 5, 334},
}

func main() {
	for _, book := range books {
		setUpDataDir()
		download(book)
		exportToPdf(book)
	}
}

func setUpDataDir() {
	err := os.RemoveAll("data")
	if err != nil {
		log.Fatalf("Data directory could not be removed: %s", err)
	}
	err = os.Mkdir("data", 0777)
	if err != nil {
		log.Fatalf("Could not create data directory: %s", err)
	}
}

func download(b book) {
	for pageNum := b.firstPage; pageNum <= b.lastPage; pageNum++ {
		pageUrl := fmt.Sprintf(b.urlTemplate, pageNum)
		resp, err := http.Get(pageUrl)
		if err != nil {
			log.Fatalf("Can't get page (%d) %s!", pageNum, pageUrl)
		}
		defer resp.Body.Close()

		fileName := fmt.Sprintf("data/%03d.%s", pageNum, b.imageType)
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatalf("Error while creating the file %s: %s", fileName, err)
		}
		defer file.Close()

		io.Copy(file, resp.Body)

		log.Println(pageUrl)
	}
}

func exportToPdf(b book) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	for pageNum := b.firstPage; pageNum <= b.lastPage; pageNum++ {
		pdf.AddPage()
		pageWidth, pageHeight := pdf.GetPageSize()
		imageOptions := gofpdf.ImageOptions{"", true, true}
		pdf.ImageOptions(fmt.Sprintf("data/%03d.%s", pageNum, b.imageType), 0, 0, pageWidth, pageHeight, false, imageOptions, 0, "")
	}

	pdfFileName := fmt.Sprintf("%s.pdf", b.name)
	err := pdf.OutputFileAndClose(pdfFileName)
	if err != nil {
		log.Fatalf("Error while saving PDF file %s: %s", pdfFileName, err)
	}
}

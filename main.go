package main

import (
	"log"

	"github.com/exanubes/url-shortener/internal/drivers"
)

func main() {
	driver := drivers.NewHttpDriver()

	if err := driver.Run(); err != nil {
		log.Fatal()
	}
}

package main

import (
	"fmt"
	"github.com/thomas-bouvier/palette-extractor"
)

func main() {
	// Creating the extractor object
	extractor := extractor.NewExtractor("../resources/img1.png", 10)

	// Displaying the top 5 dominant colors of the image
	fmt.Println(extractor.GetPalette(5))
}
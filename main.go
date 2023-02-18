package main

import (
	cr "crypto/rand"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	"log"
	"math/big"
	"os"
)

func main() {
	reader, err := os.Open("testdata/photo.jpg")
	outName := "testdata/output.jpg"
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := m.Bounds()
	outImg := image.NewRGBA(bounds)
	out, err := os.Create(outName)
	defer out.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fuckItUp(m, outImg)
	jpeg.Encode(out, outImg, nil)
}

type ColorGenerator struct {
	inToOut  map[color.Color]color.Color
	outFound map[color.Color]bool
	rand     func() uint8
}

func NewColorGenerator() ColorGenerator {
	return ColorGenerator{
		inToOut:  make(map[color.Color]color.Color),
		outFound: make(map[color.Color]bool),
		rand:     genCryptoRng(256),
	}
}

func (c ColorGenerator) GetOutputColor(in color.Color) color.Color {
	if val, ok := c.inToOut[in]; ok {
		return val
	}
	// r, g, b, a := in.RGBA()
	var newColor color.Color
	for {
		newColor = c.randColor()
		if !c.outFound[newColor] {
			c.outFound[newColor] = true
			c.inToOut[in] = newColor
			return newColor
		}

	}
}

func (c ColorGenerator) randColor() color.Color {
	r := c.rand()
	g := c.rand()
	b := c.rand()
	a := c.rand()
	out := color.RGBA{r, g, b, a}
	return out
}

// only goes up to 255, so no need to return int64
func genCryptoRng(nOpts int64) func() uint8 {
	max := big.NewInt(nOpts)
	return func() uint8 {
		bn, err := cr.Int(cr.Reader, max)
		if err != nil {
			panic(err)
		}
		return uint8(bn.Int64())
	}
}

func fuckItUp(img image.Image, out *image.RGBA) {
	cg := NewColorGenerator()
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			c := cg.GetOutputColor((img.At(x, y)))
			out.Set(x, y, c)
		}
	}
}

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
)

type ComplexLine struct {
	index int
	line  []complex128
}

type IntegerLine struct {
	index int
	line  []int
}

func main() {
	dx, dy := 1920, 1080
	niter := 200
	start, end := bounds(-0.1-0.9i, 0.1, dx, dy)
	complexGrid := complexArray(start, end, dx, dy)
	integerGrid := mandelbrot(complexGrid, niter)

	rect := image.Rect(0, 0, dx, dy)
	m := image.NewGray(rect)

	for y_i, y := range integerGrid {
		for x_i, x := range y {
			grayVal := uint16(65535 * float64(x) / float64(niter))
			m.Set(x_i, y_i, color.Gray16{grayVal})
		}
	}

	saveImage(m)
	fmt.Println("Saved image")
}

func bounds(center complex128, span float64, dx int, dy int) (complex128, complex128) {
	spanX := span * (float64(dx) / float64(dy))
	start := complex(real(center)-spanX, imag(center)-span)
	end := complex(real(center)+spanX, imag(center)+span)
	return start, end
}

func saveImage(img image.Image) {
	file, err := os.OpenFile("image.png", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic("Failed to open image file")
	}
	png.Encode(file, img)
	file.Close()
}

func mandelbrot(complexGrid [][]complex128, niter int) [][]int {
	c := make(chan IntegerLine)

	for index, line := range complexGrid {
		complexLine := ComplexLine{index, line}
		go complexLine.mandelbrot(c, niter)
	}

	integerGrid := make([][]int, len(complexGrid))

	for i := 0; i < len(complexGrid); i++ {
		integerLine := <-c
		integerGrid[integerLine.index] = integerLine.line
	}

	return integerGrid
}

func (line ComplexLine) mandelbrot(c chan IntegerLine, niter int) {
	out := IntegerLine{
		index: line.index,
		line:  make([]int, len(line.line)),
	}

	for elementIndex, complexVal := range line.line {
		z := 0 + 0i
		resultIteration := 0
		for iteration := 0; iteration < niter; iteration++ {
			if cmplx.Abs(z) > 2 {
				resultIteration = iteration
				break
			}
			z = cmplx.Pow(z, 2) + complexVal
		}
		out.line[elementIndex] = resultIteration
	}

	c <- out
}

func complexArray(start complex128, end complex128, dx int, dy int) [][]complex128 {
	startY := imag(start)
	endY := imag(end)
	startX := real(start)
	endX := real(end)

	incrementY := (endY - startY) / float64(dy)
	incrementX := (endX - startX) / float64(dx)

	out := make([][]complex128, dy)

	for y := 0; y < dy; y++ {
		out[y] = make([]complex128, dx)
		for x := 0; x < dx; x++ {
			r := startX + (incrementX * float64(x))
			im := startY + (incrementY * float64(y))
			out[y][x] = complex(r, im)
		}
	}

	return out
}

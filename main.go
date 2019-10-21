package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/cmplx"
	"net/http"
	"os/exec"
	"runtime"
)

var (
	red   color.Color = color.RGBA{255, 30, 0, 255}
	blue  color.Color = color.RGBA{0, 0, 255, 255}
	green color.Color = color.RGBA{0, 255, 0, 255}
	black color.Color = color.RGBA{0, 0, 0, 255}
	white color.Color = color.RGBA{255, 255, 255, 255}
)
var (
	height int = 4096
	weight int = 2160
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Hello)")
	fmt.Println("Proc Num:", runtime.NumCPU())

	url := fmt.Sprintf("http://localhost:%d/", 2234)
	fmt.Printf("Opening %s...\n", url)
	if err := open(url); err != nil {
		fmt.Println("Auto-open failed:", err)
		fmt.Printf("Open %s in your browser.\n", url)
	}
	http.HandleFunc("/", drawHtmlHandle)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 2234), nil))

}

func drawHtmlHandle(w http.ResponseWriter, r *http.Request) {
	m := image.NewRGBA(image.Rect(0, 0, height, weight))

	//draw.Draw(m, m.Bounds(), &image.Uniform{black}, image.Point{}, draw.Src)
	//for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
	//	m.Set(x, m.Bounds().Max.Y/2, white) // to change a single pixel
	//}
	//for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
	//	m.Set(m.Bounds().Max.X/2, y, white) // to change a single pixel
	//}

	fractal2(m)
	//Fractal(m)
	//Line(0, 0, m.Bounds().Max.X, 0, m, red)
	//Line(m.Bounds().Max.X, 0, 0, 0, m, red)
	//Line(0, 0, 0, m.Bounds().Max.Y-1, m, blue)
	//Line(m.Bounds().Max.X-1, 0, m.Bounds().Max.X-1, m.Bounds().Max.Y-1, m, green)
	//Line(0, m.Bounds().Max.Y-1, m.Bounds().Max.X-1, m.Bounds().Max.Y-1, m, white)
	//Circle(m.Bounds().Max.X/2, m.Bounds().Max.Y/2, 100, m, green)
	//jpeg.Encode(w, m , &jpeg.Options{90})
	png.Encode(w, m)
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var args []string
	args = append(args, url)

	return exec.Command("xdg-open", args...).Start()
}

func Line(x0, y0, x1, y1 int, img *image.RGBA, clr color.Color) { //Bresenham's line algorithm
	steep := false
	if math.Abs(float64(x0)-float64(x1)) < math.Abs(float64(y0)-float64(y1)) {
		y0, x0 = x0, y0
		y1, x1 = x1, y1
		steep = true
	}
	if x0 > x1 { // make it left-to-right
		//x0, y0 = x1, y1
		tmp := x0
		x0 = x1
		x1 = tmp
		tmp = y1
		y0 = y1
		y1 = tmp
	}
	dx := x1 - x0
	dy := y1 - y0
	derror := math.Abs(float64(dy) / float64(dx))
	var error float64
	y := y0
	for x := x0; x <= x1; x++ {
		if steep {
			img.Set(y, x, clr)
		} else {
			img.Set(x, y, clr)
		}
		error += derror

		if error > .5 {
			if y1 > y0 {
				y++
			} else {
				y--
			}
			error -= 1.
		}
	}
}

func Circle(x1, y1, radius int, img *image.RGBA, clr color.Color) { //Bresenham's line algorithm
	x := 0
	y := radius
	delta := 1 - 2*radius
	error := 0
	for y >= 0 {
		img.Set(x1+x, y1+y, clr)
		img.Set(x1+x, y1-y, clr)
		img.Set(x1-x, y1+y, clr)
		img.Set(x1-x, y1-y, clr)
		error = 2*(delta+y) - 1
		if (delta < 0) && (error <= 0) {
			x++
			delta += 2*x + 1
			continue
		}
		if (delta > 0) && (error > 0) {
			y--
			delta -= 2*y + 1
			continue
		}
		x++
		y--
		delta += 2 * (x - y)
	}
}

func fractal(img *image.RGBA) {
	dx := img.Bounds().Max.X
	dy := img.Bounds().Max.Y
	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			z := mandelbrot(complex(
				(float64(y)/float64(dy)*3 - 2.15),
				(float64(x)/float64(dx)*3 - 1.5),
			))
			img.Set(x, y, color.RGBA{-z ^ 2, -z ^ 2, -z ^ 2, 255})
		}
	}
}

func mandelbrot(in complex128) uint8 {
	n := in
	for i := uint8(0) + 255; i > 0; i-- {
		if cmplx.Abs(n) > 2 {
			return i
		}
		n = cmplx.Pow(n, complex(2, 0)) + in
	}
	return 255
}

func mandelb(x0, y0 float64, iter int) int {
	x := x0
	y := y0
	for i := 0; i < iter; i++ {
		real2 := x * x
		imag2 := y * y
		if (real2 + imag2) > 4.0 {
			return i
		}
		y = 2*x*y + y0
		x = real2 - imag2 + x0
	}
	return iter
}

func fractal2(img *image.RGBA) {
	dx := height
	dy := weight
	userX := -0.794591379577363
	userY := 0.16093921135504
	zoom := int64(9990999900)
	iterations := 255
	xShift := float64(dx / 2)
	yShift := float64(dy / 2)

	for v := 0; v < dy; v++ {
		for u := 0; u < dx; u++ {
			x := float64(u) - xShift
			y := (float64(v) * -1) + float64(dy) - yShift
			x = x + userX*float64(zoom)
			y = y + userY*float64(zoom)
			x = x / float64(zoom)
			y = y / float64(zoom)

			level := mandelb(x, y, iterations)
			if level == iterations {
				img.Set(u, v, black)
			} else {
				clr := 255 - uint8(level*255/iterations)
				img.Set(u, v, color.RGBA{clr, clr, clr, 255})
			}

		}
	}
}

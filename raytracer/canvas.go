package raytracer

import rl "github.com/gen2brain/raylib-go/raylib"

type Canvas struct {
	Pixels []rl.Color
	Width  int32
	Height int32
	View   View
}

func (c *Canvas) PutPixel(x, y int32, color rl.Color) {
	cX := c.Width/2 + x
	cY := c.Height/2 - y

	if cX < 0 || cX >= c.Width || cY < 0 || cY >= c.Height {
		return
	}

	c.Pixels[cY*c.Width+cX] = color
}

func (c *Canvas) CanvasToViewport(x, y int32) Vec3 {
	return Vec3{
		X: float64(x) * c.View.X / float64(c.Width),
		Y: float64(y) * c.View.Y / float64(c.Height),
		Z: c.View.D,
	}
}

type View struct {
	X, Y, D float64
}

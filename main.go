package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const MAX_INF int32 = 1_000_000_000

func main() {
	const (
		screenWidth  int32 = 1000
		screenHeight int32 = 1000
	)

	rl.InitWindow(screenWidth, screenHeight, "go converted = raylib dynamic texture")

	width := screenWidth
	height := screenHeight

	pixels := make([]rl.Color, width*height)

	c := canvas{
		pixels: pixels,
		width:  width,
		height: height,
		view:   View{x: 1, y: 1, d: 1},
	}

	checkedIm := rl.GenImageColor(int(width), int(height), rl.White)

	checked := rl.LoadTextureFromImage(checkedIm)
	rl.UpdateTexture(checked, c.pixels)

	rl.UnloadImage(checkedIm)

	rl.SetTargetFPS(24)

	posX := screenWidth/2 - checked.Width/2
	posY := screenHeight/2 - checked.Height/2

	camera_pos := rl.Vector3{
		X: 0,
		Y: 0,
		Z: 0,
	}
	spheres := []Sphere{
		{
			center: rl.Vector3{X: 0, Y: -1, Z: 3},
			radius: 1,
			color:  rl.Red,
		},
		{
			center: rl.Vector3{X: 2, Y: 0, Z: 4},
			radius: 1,
			color:  rl.Blue,
		},
		{
			center: rl.Vector3{X: -2, Y: 0, Z: 4},
			radius: 1,
			color:  rl.Green,
		},
		{
			center: rl.Vector3{X: 0, Y: -5001, Z: 0},
			radius: 5000,
			color:  rl.Brown,
		},
	}

	for !rl.WindowShouldClose() {
		if rl.IsKeyDown(rl.KeyW) {
			spheres[2].center.Z += 1 * rl.GetFrameTime()
		}
		if rl.IsKeyDown(rl.KeyS) {
			spheres[2].center.Z -= 1 * rl.GetFrameTime()
		}
		if rl.IsKeyDown(rl.KeyA) {
			spheres[2].center.X -= 1 * rl.GetFrameTime()
		}
		if rl.IsKeyDown(rl.KeyD) {
			spheres[2].center.X += 1 * rl.GetFrameTime()
		}

		for x := -c.width / 2; x < c.width/2; x++ {
			for y := -c.height / 2; y < c.height/2; y++ {
				direction := c.CanvasToViewport(x, y)
				color := TraceRay(camera_pos, direction, 1, float32(MAX_INF), spheres)
				c.PutPixel(x, y, color)
			}
		}

		rl.UpdateTexture(checked, c.pixels)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawTexture(checked, posX, posY, rl.White)
		rl.DrawFPS(10, 10)

		rl.EndDrawing()

		if rl.IsKeyDown(rl.KeyQ) {
			rl.UnloadTexture(checked)
			rl.CloseWindow()
			break
		}
	}

	rl.UnloadTexture(checked)
	rl.CloseWindow()
}

func TraceRay(O, D rl.Vector3, t_min, t_max float32, spheres []Sphere) rl.Color {
	closest_t := float32(MAX_INF)
	var closest_sphere Sphere

	for _, sphere := range spheres {
		t1, t2 := IntersectRaySphere(O, D, sphere)
		if t1 < closest_t && t_min < t1 && t1 < t_max {
			closest_t = t1
			closest_sphere = sphere
		}
		if t2 < closest_t && t_min < t2 && t2 < t_max {
			closest_t = t2
			closest_sphere = sphere
		}
	}

	if closest_sphere.radius == 0 {
		return rl.White
	}

	return closest_sphere.color
}

func IntersectRaySphere(O, D rl.Vector3, sphere Sphere) (float32, float32) {
	r := float32(sphere.radius)
	CO := rl.Vector3{
		X: O.X - sphere.center.X,
		Y: O.Y - sphere.center.Y,
		Z: O.Z - sphere.center.Z,
	}

	a := D.X*D.X + D.Y*D.Y + D.Z*D.Z
	b := 2 * (CO.X*D.X + CO.Y*D.Y + CO.Z*D.Z)
	c := (CO.X*CO.X + CO.Y*CO.Y + CO.Z*CO.Z) - (r * r)

	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return float32(MAX_INF), float32(MAX_INF)
	}

	t1 := (-b + float32(math.Sqrt(float64(discriminant)))) / (2 * a)
	t2 := (-b - float32(math.Sqrt(float64(discriminant)))) / (2 * a)

	return t1, t2
}

type canvas struct {
	pixels []rl.Color
	width  int32
	height int32
	view   View
}

func (c *canvas) PutPixel(x, y int32, color rl.Color) {
	cX := c.width/2 + x
	cY := c.height/2 - y

	if cX < 0 || cX >= c.width || cY < 0 || cY >= c.height {
		return
	}

	c.pixels[cY*c.width+cX] = color
}

func (c *canvas) CanvasToViewport(x, y int32) rl.Vector3 {
	return rl.Vector3{
		X: float32(x) * c.view.x / float32(c.width),
		Y: float32(y) * c.view.y / float32(c.height),
		Z: float32(c.view.d),
	}
}

type View struct {
	x, y, d float32
}

type Sphere struct {
	center rl.Vector3
	radius int32
	color  rl.Color
}

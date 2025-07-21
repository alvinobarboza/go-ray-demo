package main

import (
	"go-ray-demo/raytracer"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	const (
		screenWidth  int32 = 1000
		screenHeight int32 = 1000
	)

	rl.InitWindow(screenWidth, screenHeight, "go raytracer - raylib screen texture")

	width := screenWidth
	height := screenHeight

	pixels := make([]rl.Color, width*height)

	c := raytracer.Canvas{
		Pixels: pixels,
		Width:  width,
		Height: height,
		View:   raytracer.View{X: 1, Y: 1, D: 1},
	}

	checkedIm := rl.GenImageColor(int(width), int(height), rl.White)

	checked := rl.LoadTextureFromImage(checkedIm)
	rl.UpdateTexture(checked, c.Pixels)

	rl.UnloadImage(checkedIm)

	rl.SetTargetFPS(2)

	posX := screenWidth/2 - checked.Width/2
	posY := screenHeight/2 - checked.Height/2

	camera := raytracer.Camera{
		Position: rl.Vector3{
			X: 3,
			Y: 0,
			Z: 1,
		},
		Direction: [][]float32{
			{0.7071, 0, -0.7071},
			{0, 1, 0},
			{0.7071, 0, 0.7071},
		},
	}
	spheres := []raytracer.Sphere{
		{
			Center:     rl.Vector3{X: 0, Y: -1, Z: 3},
			Radius:     1,
			Color:      rl.Red,
			Specular:   500,
			Reflective: 0.2,
		},
		{
			Center:     rl.Vector3{X: 2, Y: 0, Z: 4},
			Radius:     1,
			Color:      rl.Blue,
			Specular:   500,
			Reflective: 0.3,
		},
		{
			Center:     rl.Vector3{X: -2, Y: 0, Z: 4},
			Radius:     1,
			Color:      rl.Green,
			Specular:   10,
			Reflective: 0.4,
		},
		{
			Center:     rl.Vector3{X: 0, Y: -501, Z: 0},
			Radius:     500,
			Color:      rl.Brown,
			Specular:   1000,
			Reflective: 0.5,
		},
	}

	lights := []raytracer.Ligths{
		{
			TypeL:     raytracer.AMBIENT,
			Intensity: 0.2,
		},
		{
			TypeL:     raytracer.POINT,
			Intensity: 0.6,
			Position:  rl.Vector3{X: 2, Y: 1, Z: 0},
		},
		{
			TypeL:     raytracer.DIRECTIONAL,
			Intensity: 0.2,
			Direction: rl.Vector3{X: 1, Y: 4, Z: 4},
		},
	}

	for !rl.WindowShouldClose() {
		if rl.IsKeyDown(rl.KeyW) {
			camera.Position.Z += 1 * rl.GetFrameTime()
		}
		if rl.IsKeyDown(rl.KeyS) {
			camera.Position.Z -= 1 * rl.GetFrameTime()
		}
		if rl.IsKeyDown(rl.KeyA) {
			camera.Position.X -= 1 * rl.GetFrameTime()
		}
		if rl.IsKeyDown(rl.KeyD) {
			camera.Position.X += 1 * rl.GetFrameTime()
		}

		for x := -c.Width / 2; x < c.Width/2; x++ {
			for y := -c.Height / 2; y < c.Height/2; y++ {
				direction := c.CanvasToViewport(x, y)
				newDirection := raytracer.MatrixMultiplication(camera.Direction, direction)

				color := raytracer.TraceRay(
					camera.Position,
					newDirection,
					c.View.D,
					float32(raytracer.MAX_INF),
					spheres,
					lights,
					raytracer.MAX_RECURSION,
				)
				c.PutPixel(x, y, color)
			}
		}

		rl.UpdateTexture(checked, c.Pixels)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawTexture(checked, posX, posY, rl.White)
		rl.DrawFPS(10, 10)

		rl.EndDrawing()

		if rl.IsKeyDown(rl.KeyQ) {
			break
		}
	}

	rl.UnloadTexture(checked)
	rl.CloseWindow()
}

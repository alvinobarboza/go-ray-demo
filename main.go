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

	rl.InitWindow(screenWidth, screenHeight, "go converted = raylib dynamic texture")

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

	rl.SetTargetFPS(24)

	posX := screenWidth/2 - checked.Width/2
	posY := screenHeight/2 - checked.Height/2

	camera_pos := rl.Vector3{
		X: 0,
		Y: 0,
		Z: 0,
	}
	spheres := []raytracer.Sphere{
		{
			Center:   rl.Vector3{X: 0, Y: -1, Z: 3},
			Radius:   1,
			Color:    rl.Red,
			Specular: 200,
		},
		{
			Center:   rl.Vector3{X: 2, Y: 0, Z: 4},
			Radius:   1,
			Color:    rl.Blue,
			Specular: 10,
		},
		{
			Center:   rl.Vector3{X: -2, Y: 0, Z: 4},
			Radius:   1,
			Color:    rl.Green,
			Specular: 1000,
		},
		{
			Center:   rl.Vector3{X: 0, Y: -5001, Z: 0},
			Radius:   5000,
			Color:    rl.Brown,
			Specular: 2000,
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
			spheres[2].Center.Z += 1 * rl.GetFrameTime()
		}
		if rl.IsKeyDown(rl.KeyS) {
			spheres[2].Center.Z -= 1 * rl.GetFrameTime()
		}
		if rl.IsKeyDown(rl.KeyA) {
			spheres[2].Center.X -= 1 * rl.GetFrameTime()
		}
		if rl.IsKeyDown(rl.KeyD) {
			spheres[2].Center.X += 1 * rl.GetFrameTime()
		}

		for x := -c.Width / 2; x < c.Width/2; x++ {
			for y := -c.Height / 2; y < c.Height/2; y++ {
				direction := c.CanvasToViewport(x, y)
				color := raytracer.TraceRay(camera_pos, direction, 1, float32(raytracer.MAX_INF), spheres, lights)
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
			rl.UnloadTexture(checked)
			rl.CloseWindow()
			break
		}
	}

	rl.UnloadTexture(checked)
	rl.CloseWindow()
}

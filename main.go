package main

import (
	"fmt"
	"go-ray-demo/raytracer"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	const (
		screenWidth  int32 = 800
		screenHeight int32 = 450

		moveSpeed float32 = 5
		turnSpeed float32 = 80
	)

	rl.InitWindow(screenWidth, screenHeight, "go raytracer - raylib screen texture")

	width := screenWidth
	height := screenHeight

	pixels := make([]rl.Color, width*height)

	c := raytracer.Canvas{
		Pixels: pixels,
		Width:  width,
		Height: height,
		View:   raytracer.View{X: 1, Y: .6, D: 1},
	}

	checkedIm := rl.GenImageColor(int(width), int(height), rl.White)

	checked := rl.LoadTextureFromImage(checkedIm)
	rl.UpdateTexture(checked, c.Pixels)

	rl.UnloadImage(checkedIm)

	rl.SetTargetFPS(30)

	posX := screenWidth/2 - checked.Width/2
	posY := screenHeight/2 - checked.Height/2

	camera := raytracer.Camera{
		Position: rl.Vector3{
			X: 3,
			Y: 0,
			Z: 1,
		},
		Rotation: rl.Vector3{
			X: 0,
			Y: 90,
			Z: 0,
		},
		Direction: rl.Vector3{
			X: 0,
			Y: 0,
			Z: 1,
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
			Reflective: 0.001,
		},
		{
			Center:     rl.Vector3{X: -2, Y: 0, Z: 4},
			Radius:     1,
			Color:      rl.Green,
			Specular:   10,
			Reflective: 0.1,
		},
		{
			Center:     rl.Vector3{X: -2, Y: 0, Z: -4},
			Radius:     1,
			Color:      rl.Green,
			Specular:   10,
			Reflective: 0.2,
		},
		{
			Center:     rl.Vector3{X: 0, Y: -501, Z: 0},
			Radius:     500,
			Color:      rl.DarkGreen,
			Specular:   200,
			Reflective: 0.1,
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

	doTask := func() func(task func(), c bool) {
		tasks := make(chan func())

		for range 4 {
			go func() {
				for t := range tasks {
					t()
				}
			}()
		}

		return func(task func(), c bool) {
			if c {
				close(tasks)
				return
			}
			tasks <- task
		}
	}()

	var wg sync.WaitGroup

	for !rl.WindowShouldClose() {
		if rl.IsKeyDown(rl.KeyW) {
			camera.MoveForward(moveSpeed * rl.GetFrameTime())
		}
		if rl.IsKeyDown(rl.KeyS) {
			camera.MoveBackward(moveSpeed * rl.GetFrameTime())
		}
		if rl.IsKeyDown(rl.KeyA) {
			camera.MoveLeft(moveSpeed * rl.GetFrameTime())
		}
		if rl.IsKeyDown(rl.KeyD) {
			camera.MoveRight(moveSpeed * rl.GetFrameTime())
		}

		if rl.IsKeyDown(rl.KeyRight) {
			camera.Rotation.Y -= turnSpeed * rl.GetFrameTime()
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			camera.Rotation.Y += turnSpeed * rl.GetFrameTime()
		}
		if rl.IsKeyDown(rl.KeyUp) {
			camera.Rotation.X += turnSpeed * rl.GetFrameTime()
		}
		if rl.IsKeyDown(rl.KeyDown) {
			camera.Rotation.X -= turnSpeed * rl.GetFrameTime()
		}

		startW, endW := -c.Width/2, c.Width/2
		startH, endH := -c.Height/2, c.Height/2

		listT := [][]int32{
			{startW, 0, 0, endH},
			{startW, 0, startH, 0},
			{0, endW, 0, endH},
			{0, endW, startH, 0},
		}

		for _, i := range listT {
			wg.Add(1)
			doTask(func() {
				for x := i[0]; x < i[1]; x++ {
					for y := i[2]; y < i[3]; y++ {
						direction := c.CanvasToViewport(x, y)
						newDirection := raytracer.RotateXYZ(camera.Rotation, direction)

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
				wg.Done()
			}, false)
		}
		wg.Wait()

		rl.UpdateTexture(checked, c.Pixels)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawTexture(checked, posX, posY, rl.White)
		rl.DrawFPS(10, 10)
		rl.DrawText(
			fmt.Sprintf("Cam-> \nX:%02f Y:%02f Z:%02f", camera.Position.X, camera.Position.Y, camera.Position.Z),
			10, 50, 15, rl.Red)

		rl.EndDrawing()

		if rl.IsKeyDown(rl.KeyQ) {
			break
		}
	}

	rl.UnloadTexture(checked)
	rl.CloseWindow()
}

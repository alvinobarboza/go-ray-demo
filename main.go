package main

import (
	"fmt"
	"go-ray-demo/raytracer"
	"runtime"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	const (
		screenWidth  int32 = 1000
		screenHeight int32 = 1000

		moveSpeed = 4
		turnSpeed = 70
	)
	// Dynamically get CPUs to spawn a
	// reasonable number of goroutines
	threads := runtime.NumCPU()
	if threads > 3 {
		threads -= 2
	}

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

	rl.SetTargetFPS(30)

	posX := screenWidth/2 - checked.Width/2
	posY := screenHeight/2 - checked.Height/2

	camera := raytracer.Camera{
		Position: raytracer.Vec3{
			X: 0,
			Y: 0,
			Z: 0,
		},
		Rotation: raytracer.Vec3{
			X: 0,
			Y: 0,
			Z: 0,
		},
		Direction: raytracer.Vec3{
			X: 0,
			Y: 0,
			Z: 1,
		},
	}
	spheres := []raytracer.Sphere{
		{
			Center:     raytracer.Vec3{X: 0, Y: -1, Z: 3},
			Radius:     1,
			Color:      rl.Red,
			Specular:   500,
			Reflective: 0.2,
		},
		{
			Center:     raytracer.Vec3{X: 2, Y: 0, Z: 4},
			Radius:     1,
			Color:      rl.Blue,
			Specular:   500,
			Reflective: 0.001,
		},
		{
			Center:     raytracer.Vec3{X: -2, Y: 0, Z: 4},
			Radius:     1,
			Color:      rl.Green,
			Specular:   10,
			Reflective: 0.1,
		},
		{
			Center:          raytracer.Vec3{X: -.5, Y: 0, Z: 2},
			Radius:          .4,
			Color:           rl.Blue,
			Specular:        200,
			Reflective:      0.2,
			Opacity:         0.9,
			RefractionIndex: 1.33,
		},
		{
			Center:     raytracer.Vec3{X: 0, Y: -501, Z: 0},
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
			Position:  raytracer.Vec3{X: 2, Y: 1, Z: 0},
		},
		{
			TypeL:     raytracer.DIRECTIONAL,
			Intensity: 0.2,
			Direction: raytracer.Vec3{X: 1, Y: 4, Z: 4},
		},
	}

	// Before main loop, spaw "threads" amount of goroutines that will listen in a channel for tasks
	// This return a func which accepts a func as task to send in the chan in its scope
	doTask := func() func(task func(), c bool) {
		tasks := make(chan func())

		for range threads {
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

	xs, xe, ys, ye := -c.Width/2, c.Width/2, -c.Height/2, c.Height/2

	columnRange := make([]int32, 0)
	for i := xs; i < xe; i++ {
		columnRange = append(columnRange, i)
	}

	rowRange := make([]int32, 0)
	for j := ys; j < ye; j++ {
		rowRange = append(rowRange, j)
	}

	// Dynamically choose size of minor matrix to run in the task
	// Trying to maximize goroutine use, as a too small window
	// would make too many calls to "chan", and too large window
	// would waste goroutine in idle
	matrixSides := int(width * height / int32(threads) / width)

	// Precompute the task that will run in the frame update
	task := func(csx, cex, csy, cey int) func() {
		return func() {
			for xx := csx; xx < cex; xx++ {
				for yy := csy; yy < cey; yy++ {
					direction := c.CanvasToViewport(columnRange[xx], rowRange[yy])
					newDirection := direction.RotateXYZ(camera.Rotation)

					color := raytracer.TraceRay(
						camera.Position,
						newDirection,
						c.View.D,
						raytracer.MaxInf,
						spheres,
						lights,
						raytracer.MaxRecursion,
					)
					c.PutPixel(columnRange[xx], rowRange[yy], color)
				}
			}
			wg.Done()
		}
	}

	for !rl.WindowShouldClose() {
		if rl.IsKeyDown(rl.KeyU) {
			spheres[3].Center.Z += 2 * float64(rl.GetFrameTime())
		}
		if rl.IsKeyDown(rl.KeyJ) {
			spheres[3].Center.Z -= 2 * float64(rl.GetFrameTime())
		}
		if rl.IsKeyDown(rl.KeyH) {
			spheres[3].Center.X -= 2 * float64(rl.GetFrameTime())
		}
		if rl.IsKeyDown(rl.KeyK) {
			spheres[3].Center.X += 2 * float64(rl.GetFrameTime())
		}

		if rl.IsKeyDown(rl.KeyW) {
			camera.MoveForward(moveSpeed * float64(rl.GetFrameTime()))
		}
		if rl.IsKeyDown(rl.KeyS) {
			camera.MoveBackward(moveSpeed * float64(rl.GetFrameTime()))
		}
		if rl.IsKeyDown(rl.KeyA) {
			camera.MoveLeft(moveSpeed * float64(rl.GetFrameTime()))
		}
		if rl.IsKeyDown(rl.KeyD) {
			camera.MoveRight(moveSpeed * float64(rl.GetFrameTime()))
		}
		if rl.IsKeyDown(rl.KeySpace) {
			camera.Position.Y += moveSpeed * float64(rl.GetFrameTime())
		}
		if rl.IsKeyDown(rl.KeyLeftControl) {
			camera.Position.Y -= moveSpeed * float64(rl.GetFrameTime())
		}

		if rl.IsKeyDown(rl.KeyRight) {
			camera.Rotation.Y -= turnSpeed * float64(rl.GetFrameTime())
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			camera.Rotation.Y += turnSpeed * float64(rl.GetFrameTime())
		}
		if rl.IsKeyDown(rl.KeyUp) {
			camera.Rotation.X += turnSpeed * float64(rl.GetFrameTime())
			if camera.Rotation.X >= 90 {
				camera.Rotation.X = 89
			}
		}
		if rl.IsKeyDown(rl.KeyDown) && camera.Rotation.X < 90 {
			camera.Rotation.X -= turnSpeed * float64(rl.GetFrameTime())
			if camera.Rotation.X <= -90 {
				camera.Rotation.X = -89
			}
		}

		cStartX := 0
		cEndX := 0
		cX := 0

		for i := range columnRange {
			cX++
			if cX == matrixSides || i+1 == len(columnRange) {
				cStartX = cEndX
				cEndX += cX

				cStartY := 0
				cEndY := 0
				cY := 0
				for j := range rowRange {
					cY++
					if cY == matrixSides {
						cStartY = cEndY
						cEndY += cY
						cX = 0
						cY = 0

						wg.Add(1)
						doTask(task(cStartX, cEndX, cStartY, cEndY), false)
						continue
					}
					if j+1 == len(rowRange) {
						cStartY = cEndY
						cEndY += cY
						cX = 0
						cY = 0
						wg.Add(1)
						doTask(task(cStartX, cEndX, cStartY, cEndY), false)
						continue
					}
				}
			}
		}

		// Can run without a waiting, but has too many tearing in the image,
		// as the pixel array will have old pixel data..
		// The tasks are run in parallel, so they will fill
		// pixel data when finished, not in this loop
		// wait, prevents that
		wg.Wait()

		rl.UpdateTexture(checked, c.Pixels)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawTexture(checked, posX, posY, rl.White)
		rl.DrawFPS(10, 10)
		rl.DrawText(
			fmt.Sprintf("Cam-> \nX:%01f \nY:%01f \nZ:%01f", camera.Position.X, camera.Position.Y, camera.Position.Z),
			10, 30, 20, rl.White)
		rl.DrawText("Move: A/W/S/D\nControl Camera: UP/DOWN/LEFT/RIGHT",
			10, 120, 20, rl.White)
		rl.DrawText("I'm a bit too lazy to make \nit work with the mouse...",
			10, 160, 10, rl.White)

		rl.EndDrawing()

		if rl.IsKeyDown(rl.KeyQ) {
			break
		}
	}

	rl.UnloadTexture(checked)
	rl.CloseWindow()
}

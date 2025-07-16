package raytracer

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func VecDot(v1, v2 rl.Vector3) float32 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

func VecLen(v rl.Vector3) float32 {
	return float32(math.Sqrt(float64(VecDot(v, v))))
}

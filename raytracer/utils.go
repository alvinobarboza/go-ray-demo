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

func MatrixMultiplication(m [][]float32, v rl.Vector3) rl.Vector3 {
	result := []float32{0, 0, 0}
	vec := []float32{v.X, v.Y, v.Z}

	for i := range 3 {
		for j := range 3 {
			result[i] += vec[j] * m[i][j]
		}
	}
	return rl.Vector3{X: result[0], Y: result[1], Z: result[2]}
}

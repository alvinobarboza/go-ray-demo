package raytracer

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	TAU        float32 = 2 * math.Pi
	DEG_TO_RAD         = TAU / 360
)

func VecDot(v1, v2 rl.Vector3) float32 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

func VecLen(v rl.Vector3) float32 {
	return float32(math.Sqrt(float64(VecDot(v, v))))
}

func MatrixMultiplication(m []float32, v rl.Vector3) rl.Vector3 {
	result := []float32{0, 0, 0}
	vec := []float32{v.X, v.Y, v.Z}

	length := 3

	for h := range length {
		for w := range length {
			result[h] += vec[w] * m[length*h+w]
		}
	}
	return rl.Vector3{X: result[0], Y: result[1], Z: result[2]}
}

func RotateXYZ(angle rl.Vector3, v rl.Vector3) rl.Vector3 {
	cosa := float32(math.Cos(-float64(angle.X * DEG_TO_RAD)))
	sina := float32(math.Sin(-float64(angle.X * DEG_TO_RAD)))

	cosb := float32(math.Cos(-float64(angle.Y * DEG_TO_RAD)))
	sinb := float32(math.Sin(-float64(angle.Y * DEG_TO_RAD)))

	cosga := float32(math.Cos(-float64(angle.Z * DEG_TO_RAD)))
	singa := float32(math.Sin(-float64(angle.Z * DEG_TO_RAD)))

	// Formula for general 3D roation using matrix
	matrix := []float32{
		cosb * cosga, sina*sinb*cosga - cosa*singa, cosa*sinb*cosga + sina*singa,
		cosb * singa, sina*sinb*singa + cosa*cosga, cosa*sinb*singa - sina*cosga,
		-sinb, sina * cosb, cosa * cosb,
	}

	value := MatrixMultiplication(matrix, v)

	return value
}

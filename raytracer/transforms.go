package raytracer

import "math"

type Vec3 struct {
	X, Y, Z float64
}

func (v Vec3) VecDot() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vec3) VecLen() float64 {
	return math.Sqrt(v.VecDot())
}

func (v Vec3) VecNormal() Vec3 {
	n := v.VecLen()
	return Vec3{
		X: v.X / n,
		Y: v.Y / n,
		Z: v.Z / n,
	}
}

func (v Vec3) VecMultiply(n float64) Vec3 {
	return Vec3{
		X: v.X * n,
		Y: v.Y * n,
		Z: v.Z * n,
	}
}

func (v Vec3) MatrixMultiplication(m []float64) Vec3 {
	result := []float64{0, 0, 0}
	vec := []float64{v.X, v.Y, v.Z}

	length := 3

	for h := range length {
		for w := range length {
			result[h] += vec[w] * m[length*h+w]
		}
	}
	return Vec3{X: result[0], Y: result[1], Z: result[2]}
}

func (v Vec3) RotateXYZ(angle Vec3) Vec3 {
	cosa := math.Cos(angle.X * -DEG_TO_RAD)
	sina := math.Sin(angle.X * -DEG_TO_RAD)

	cosb := math.Cos(angle.Y * -DEG_TO_RAD)
	sinb := math.Sin(angle.Y * -DEG_TO_RAD)

	cosga := math.Cos(angle.Z * -DEG_TO_RAD)
	singa := math.Sin(angle.Z * -DEG_TO_RAD)

	// Formula for general 3D roation using matrix
	matrix := []float64{
		cosb * cosga, sina*sinb*cosga - cosa*singa, cosa*sinb*cosga + sina*singa,
		cosb * singa, sina*sinb*singa + cosa*cosga, cosa*sinb*singa - sina*cosga,
		-sinb, sina * cosb, cosa * cosb,
	}

	value := v.MatrixMultiplication(matrix)

	return value
}

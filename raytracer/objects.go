package raytracer

import rl "github.com/gen2brain/raylib-go/raylib"

type lightType int

const (
	AMBIENT lightType = iota
	POINT
	DIRECTIONAL
)

type Sphere struct {
	Center   rl.Vector3
	Radius   int32
	Color    rl.Color
	Specular int32
}

type Ligths struct {
	TypeL     lightType
	Intensity float32
	Direction rl.Vector3
	Position  rl.Vector3
}

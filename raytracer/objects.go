package raytracer

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type lightType int

const (
	AMBIENT lightType = iota
	POINT
	DIRECTIONAL
)

type Sphere struct {
	Center          Vec3
	Radius          float64
	Color           rl.Color
	Specular        int32
	Reflective      float64
	Opacity         float64
	RefractionIndex float64
}

type Ligths struct {
	TypeL     lightType
	Intensity float64
	Direction Vec3
	Position  Vec3
}

type Camera struct {
	Position  Vec3
	Rotation  Vec3
	Direction Vec3
}

func (c *Camera) MoveForward(unit float64) {
	direction := RotateXYZ(c.Rotation, c.Direction)
	lenD := VecLen(direction)
	normalDirection := Vec3{
		X: direction.X / lenD,
		Y: direction.Y / lenD,
		Z: direction.Z / lenD,
	}
	c.Position.X += normalDirection.X * unit
	c.Position.Y += normalDirection.Y * unit
	c.Position.Z += normalDirection.Z * unit
}

func (c *Camera) MoveBackward(unit float64) {
	c.MoveForward(-unit)
}

func (c *Camera) MoveLeft(unit float64) {
	direction := RotateXYZ(c.Rotation, c.Direction)
	sideDirection := CrossProdutc(direction, Vec3{Y: 1})

	lenD := VecLen(sideDirection)
	normalDirection := Vec3{
		X: sideDirection.X / lenD,
		Y: sideDirection.Y / lenD,
		Z: sideDirection.Z / lenD,
	}
	c.Position.X += normalDirection.X * unit
	c.Position.Y += normalDirection.Y * unit
	c.Position.Z += normalDirection.Z * unit
}

func (c *Camera) MoveRight(unit float64) {
	c.MoveLeft(-unit)
}

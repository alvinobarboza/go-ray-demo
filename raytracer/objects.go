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
	Center     rl.Vector3
	Radius     int32
	Color      rl.Color
	Specular   int32
	Reflective float32
}

type Ligths struct {
	TypeL     lightType
	Intensity float32
	Direction rl.Vector3
	Position  rl.Vector3
}

type Camera struct {
	Position  rl.Vector3
	Rotation  rl.Vector3
	Direction rl.Vector3
}

func (c *Camera) MoveForward(unit float32) {
	direction := RotateXYZ(c.Rotation, c.Direction)
	lenD := VecLen(direction)
	normalDirection := rl.Vector3{
		X: direction.X / lenD,
		Y: direction.Y / lenD,
		Z: direction.Z / lenD,
	}
	c.Position.X += normalDirection.X * unit
	c.Position.Y += normalDirection.Y * unit
	c.Position.Z += normalDirection.Z * unit
}

func (c *Camera) MoveBackward(unit float32) {
	c.MoveForward(-unit)
}

func (c *Camera) MoveLeft(unit float32) {
	direction := RotateXYZ(c.Rotation, c.Direction)
	sideDirection := CrossProdutc(direction, rl.Vector3{Y: 1})

	lenD := VecLen(sideDirection)
	normalDirection := rl.Vector3{
		X: sideDirection.X / lenD,
		Y: sideDirection.Y / lenD,
		Z: sideDirection.Z / lenD,
	}
	c.Position.X += normalDirection.X * unit
	c.Position.Y += normalDirection.Y * unit
	c.Position.Z += normalDirection.Z * unit
}

func (c *Camera) MoveRight(unit float32) {
	c.MoveLeft(-unit)
}

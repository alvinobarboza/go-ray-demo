package raytracer

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const MAX_INF int32 = 1_000_000_000
const MAX_RECURSION int8 = 2

func TraceRay(O, D rl.Vector3, t_min, t_max float32, spheres []Sphere, lights []Ligths, recursion int8) rl.Color {
	closest_sphere, closest_t := ClosesIntersection(O, D, t_min, t_max, spheres)

	if closest_sphere.Radius == 0 {
		return rl.Black
	}

	point := rl.Vector3{
		X: O.X + closest_t*D.X,
		Y: O.Y + closest_t*D.Y,
		Z: O.Z + closest_t*D.Z,
	}

	normal := rl.Vector3{
		X: point.X - closest_sphere.Center.X,
		Y: point.Y - closest_sphere.Center.Y,
		Z: point.Z - closest_sphere.Center.Z,
	}

	l_normal := VecLen(normal)

	if l_normal > 0 {
		normal.X = normal.X / l_normal
		normal.Y = normal.Y / l_normal
		normal.Z = normal.Z / l_normal
	}

	objToCam := rl.Vector3{
		X: D.X * -1,
		Y: D.Y * -1,
		Z: D.Z * -1,
	}

	i := ComputeLighting(point, normal, objToCam, lights, spheres, closest_sphere.Specular)

	closest_sphere.Color.R = uint8(float32(closest_sphere.Color.R) * i)
	closest_sphere.Color.G = uint8(float32(closest_sphere.Color.G) * i)
	closest_sphere.Color.B = uint8(float32(closest_sphere.Color.B) * i)

	local_color := closest_sphere.Color

	if recursion <= 0 || closest_sphere.Reflective <= 0 {
		return local_color
	}

	r := closest_sphere.Reflective

	reflected := ReflectRay(objToCam, normal)

	reflected_color := TraceRay(point, reflected, 0.001, float32(MAX_INF), spheres, lights, recursion-1)

	reflected_color.R = uint8(float32(local_color.R)*(1-r) + float32(reflected_color.R)*r)
	reflected_color.G = uint8(float32(local_color.G)*(1-r) + float32(reflected_color.G)*r)
	reflected_color.B = uint8(float32(local_color.B)*(1-r) + float32(reflected_color.B)*r)

	return reflected_color
}

func ClosesIntersection(O, D rl.Vector3, t_min, t_max float32, spheres []Sphere) (Sphere, float32) {
	closest_t := float32(MAX_INF)
	var closest_sphere Sphere

	for _, sphere := range spheres {
		t1, t2 := IntersectRaySphere(O, D, sphere)
		if t1 < closest_t && t_min < t1 && t1 < t_max {
			closest_t = t1
			closest_sphere = sphere
		}
		if t2 < closest_t && t_min < t2 && t2 < t_max {
			closest_t = t2
			closest_sphere = sphere
		}
	}
	return closest_sphere, closest_t
}

func IntersectRaySphere(O, D rl.Vector3, sphere Sphere) (float32, float32) {
	r := float32(sphere.Radius)
	CO := rl.Vector3{
		X: O.X - sphere.Center.X,
		Y: O.Y - sphere.Center.Y,
		Z: O.Z - sphere.Center.Z,
	}

	a := VecDot(D, D)
	b := 2 * (VecDot(CO, D))
	c := (VecDot(CO, CO)) - (r * r)

	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return float32(MAX_INF), float32(MAX_INF)
	}

	t1 := (-b + float32(math.Sqrt(float64(discriminant)))) / (2 * a)
	t2 := (-b - float32(math.Sqrt(float64(discriminant)))) / (2 * a)

	return t1, t2
}

func ComputeLighting(point, normal, objToCam rl.Vector3, lights []Ligths, spheres []Sphere, s int32) float32 {
	var i float32

	for _, light := range lights {
		if light.TypeL == AMBIENT {
			i += light.Intensity
		} else {
			L := rl.Vector3{}
			if light.TypeL == POINT {
				L.X = light.Position.X - point.X
				L.Y = light.Position.Y - point.Y
				L.Z = light.Position.Z - point.Z
			} else {
				L = light.Direction
			}

			// Shadow
			shadow_sphere, _ := ClosesIntersection(point, L, 0.001, float32(MAX_INF), spheres)
			if shadow_sphere.Radius != 0 {
				continue
			}

			// Deffuse
			n_dot_l := VecDot(normal, L)
			if n_dot_l > 0 {
				length_normal := VecLen(normal)
				length_L := VecLen(L)
				i += light.Intensity * n_dot_l / (length_normal * length_L)
			}

			// Specular
			if s != -1 {
				reflected := ReflectRay(L, normal)
				r_dot_oc := VecDot(reflected, objToCam)
				if r_dot_oc > 0 {
					length_reflected := VecLen(reflected)
					length_objToCam := VecLen(objToCam)
					i += light.Intensity * float32(math.Pow(float64(r_dot_oc/(length_reflected*length_objToCam)), float64(s)))
				}
			}
		}
	}

	if i > 1 {
		i = 1
	}

	return i
}

func ReflectRay(ray, normal rl.Vector3) rl.Vector3 {
	r_dot_n := VecDot(ray, normal)
	return rl.Vector3{
		X: 2*normal.X*r_dot_n - ray.X,
		Y: 2*normal.Y*r_dot_n - ray.Y,
		Z: 2*normal.Z*r_dot_n - ray.Z,
	}
}

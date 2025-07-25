package raytracer

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const MaxInf = 1_000_000_000.0
const MaxRecursion int8 = 3

func TraceRay(origin, ray Vec3, t_min, t_max float64, spheres []Sphere, lights []Ligths, recursion int8) rl.Color {
	closest_sphere, closest_t := ClosesIntersection(origin, ray, t_min, t_max, spheres)

	if closest_sphere.Radius == 0 {
		return rl.Gray
	}

	point := Vec3{
		X: origin.X + closest_t*ray.X,
		Y: origin.Y + closest_t*ray.Y,
		Z: origin.Z + closest_t*ray.Z,
	}

	normal := Vec3{
		X: point.X - closest_sphere.Center.X,
		Y: point.Y - closest_sphere.Center.Y,
		Z: point.Z - closest_sphere.Center.Z,
	}

	l_normal := normal.VecLen()

	if l_normal > 0 {
		normal.X = normal.X / l_normal
		normal.Y = normal.Y / l_normal
		normal.Z = normal.Z / l_normal
	}

	objToCam := ray.VecMultiply(-1)

	i := ComputeLighting(point, normal, objToCam, lights, spheres, closest_sphere.Specular)

	closest_sphere.Color.R = uint8(float64(closest_sphere.Color.R) * i)
	closest_sphere.Color.G = uint8(float64(closest_sphere.Color.G) * i)
	closest_sphere.Color.B = uint8(float64(closest_sphere.Color.B) * i)

	local_color := closest_sphere.Color

	if closest_sphere.Opacity > 0 {
		o := closest_sphere.Opacity

		angleRay := RayAngleFromNormal(ray, normal)
		refracted := Refraction(ray, normal, angleRay, closest_sphere.RefractionIndex)
		transparentColor := TraceRay(point, refracted, t_min, t_max, spheres, lights, recursion-1)

		local_color.R = uint8(float64(local_color.R)*(1-o) + float64(transparentColor.R)*o)
		local_color.G = uint8(float64(local_color.G)*(1-o) + float64(transparentColor.G)*o)
		local_color.B = uint8(float64(local_color.B)*(1-o) + float64(transparentColor.B)*o)
	}

	if recursion <= 0 || closest_sphere.Reflective <= 0 {
		return local_color
	}

	r := closest_sphere.Reflective

	reflected := ReflectRay(objToCam, normal)

	reflected_color := TraceRay(point, reflected, 0.001, MaxInf, spheres, lights, recursion-1)

	reflected_color.R = uint8(float64(local_color.R)*(1-r) + float64(reflected_color.R)*r)
	reflected_color.G = uint8(float64(local_color.G)*(1-r) + float64(reflected_color.G)*r)
	reflected_color.B = uint8(float64(local_color.B)*(1-r) + float64(reflected_color.B)*r)

	return reflected_color
}

func ClosesIntersection(origin, ray Vec3, t_min, t_max float64, spheres []Sphere) (Sphere, float64) {
	closest_t := MaxInf
	var closest_sphere Sphere

	for _, sphere := range spheres {
		t1, t2 := IntersectRaySphere(origin, ray, sphere)
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

func IntersectRaySphere(origin, ray Vec3, sphere Sphere) (float64, float64) {
	r := sphere.Radius
	CO := Vec3{
		X: origin.X - sphere.Center.X,
		Y: origin.Y - sphere.Center.Y,
		Z: origin.Z - sphere.Center.Z,
	}

	a := ray.VecDot()
	b := 2 * VecDot(CO, ray)
	c := CO.VecDot() - (r * r)

	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return MaxInf, MaxInf
	}

	t1 := (-b + math.Sqrt(discriminant)) / (2 * a)
	t2 := (-b - math.Sqrt(discriminant)) / (2 * a)

	return t1, t2
}

func ComputeLighting(point, normal, objToCam Vec3, lights []Ligths, spheres []Sphere, s int32) float64 {
	var i float64

	for _, light := range lights {
		if light.TypeL == AMBIENT {
			i += light.Intensity
		} else {
			L := Vec3{}
			if light.TypeL == POINT {
				L.X = light.Position.X - point.X
				L.Y = light.Position.Y - point.Y
				L.Z = light.Position.Z - point.Z
			} else {
				L = light.Direction
			}

			// Shadow
			shadow_sphere, _ := ClosesIntersection(point, L, 0.001, MaxInf, spheres)
			if shadow_sphere.Radius != 0 {
				continue
			}

			// Deffuse
			n_dot_l := VecDot(normal, L)
			if n_dot_l > 0 {
				length_normal := normal.VecLen()
				length_L := L.VecLen()
				i += light.Intensity * n_dot_l / (length_normal * length_L)
			}

			// Specular
			if s != -1 {
				reflected := ReflectRay(L, normal)
				r_dot_oc := VecDot(reflected, objToCam)
				if r_dot_oc > 0 {
					length_reflected := reflected.VecLen()
					length_objToCam := objToCam.VecLen()
					i += light.Intensity * math.Pow(r_dot_oc/(length_reflected*length_objToCam), float64(s))
				}
			}
		}
	}

	if i > 1 {
		i = 1
	}

	return i
}

func ReflectRay(ray, normal Vec3) Vec3 {
	r_dot_n := VecDot(ray, normal)
	return Vec3{
		X: 2*normal.X*r_dot_n - ray.X,
		Y: 2*normal.Y*r_dot_n - ray.Y,
		Z: 2*normal.Z*r_dot_n - ray.Z,
	}
}

// Formula for rotation around a arbtrary orthognal vector(both vector must be normalized)
// u = ortho vector = cross product of ray and normal
// x = vector to rotate = ray
// angle = radian
// newx = u * (u dot x) + cos(angle) * (u cross x) cross u + sin(angle)*(u cross x)
func Refraction(ray, normal Vec3, angleRay, refractionIndex float64) Vec3 {
	angleIndex := math.Asin(math.Sin(angleRay) / refractionIndex)

	crossRayNormal := CrossProdutc(normal, ray)
	if crossRayNormal.X != 0 && crossRayNormal.Y != 0 && crossRayNormal.Z != 0 {
		crossRayNormal = crossRayNormal.VecNormal()
	}
	crossRayCross := CrossProdutc(crossRayNormal, ray)
	c1 := crossRayNormal.VecMultiply(VecDot(crossRayNormal, ray))
	c2 := CrossProdutc(
		crossRayCross.VecMultiply(
			math.Cos(angleIndex),
		),
		crossRayNormal,
	)
	c3 := crossRayCross.VecMultiply(math.Sin(angleIndex))

	c1c2 := VecAdd(c1, c2)
	c1c2c3 := VecAdd(c1c2, c3)
	if c1c2c3.X == 0 && c1c2c3.Y == 0 && c1c2c3.Z == 0 {
		return c1c2c3
	}

	return c1c2c3.VecNormal()
}

// Find angle between two vectors
// angle = cos(angle) = (u dot v) / (length u * length v)
func RayAngleFromNormal(ray, normal Vec3) float64 {
	rayDotNormal := VecDot(ray, normal)
	lenRay := ray.VecLen()
	lenNormal := normal.VecLen()

	angleRay := math.Acos(rayDotNormal / (lenRay * lenNormal))
	return angleRay
}

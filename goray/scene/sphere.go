package scene

import (
	"math"
	vec "github.com/uri-nyx/tinyraytracer/goray/raytracing")

// A sphere is a hittable object, defined by a center (point) and a radius.
type Sphere struct {
	Center vec.Point3d
	Radius float64
	Material Material
}

func (s Sphere) Hit(r vec.Ray, t_min, t_max float64, hit *Hit) bool {
	oc := vec.Sub(r.O, s.Center)
	a := r.Dir.LengthSquared()
	half_b := vec.Dot(oc, r.Dir)
	c := oc.LengthSquared() - s.Radius*s.Radius
	discriminant := half_b*half_b - a*c

	if discriminant < 0 { return false }

	discriminant_root := math.Sqrt(discriminant)
	root := (-half_b - discriminant_root) / a

	if root < t_min || t_max < root {
		root = (-half_b + discriminant_root) / a
		if root < t_min || t_max < root {
			return false
		}
	}

	hit.T = root
	hit.P = r.At(root)
	hit.N = vec.Sub(hit.P, s.Center).DivScalar(s.Radius)
	hit.Material = s.Material

	outwardNormal := vec.Sub(hit.P, s.Center).DivScalar(s.Radius)
    hit.SetFaceNormal(r, outwardNormal)

	return true
}
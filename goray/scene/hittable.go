package scene

import vec "github.com/uri-nyx/tinyraytracer/goray/raytracing"

// Hit is a struct that holds information about a hit.
type Hit struct {
	T float64
	P vec.Point3d
	N vec.Vector3
	FrontFace bool

	Material Material
}

func (h *Hit) SetFaceNormal(r vec.Ray, outward vec.Vector3) {
	h.FrontFace = vec.Dot(r.Dir, outward) < 0;
    
	if h.FrontFace {
		h.N = outward
	} else {
		h.N = outward.MulScalar(-1)
	}
}

// Hittable is an interface that defines a hittable object.
type Hittable interface {
	Hit(r vec.Ray, t_min, t_max float64, hit *Hit) bool
}

type HittableList []Hittable

func (l HittableList) Hit(r vec.Ray, t_min, t_max float64, hit *Hit) bool {
	tempHit := Hit{}
	hitAnything := false
	closestSoFar := t_max

	for _, o := range l {
		if o.Hit(r, t_min, closestSoFar, &tempHit) {
			hitAnything = true
			closestSoFar = tempHit.T
			*hit = tempHit
		}
	}
	return hitAnything
}


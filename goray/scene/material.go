package scene

import (
	"math"
	"math/rand"

	vec "github.com/uri-nyx/tinyraytracer/goray/raytracing"
)

type Material interface {
	Scatter(r vec.Ray, h *Hit, attenuation *vec.Color, scattered *vec.Ray) bool
}

type Lambertian struct {
	Albedo vec.Color
}

func (l Lambertian) Scatter(r vec.Ray, h *Hit, attenuation *vec.Color, scattered *vec.Ray) bool {
	scatterDir := vec.Add(h.N, vec.RandomUnitV3())
	
	if scatterDir.NearZero() {
		scatterDir = h.N
	}

	*scattered = vec.NewRay(h.P, scatterDir)
	*attenuation = l.Albedo

	return true
}

type Metal struct {
	Albedo vec.Color
	Fuzz float64
}

func (m Metal) Scatter(r vec.Ray, h *Hit, attenuation *vec.Color, scattered *vec.Ray) bool {
	reflected := vec.Reflect(r.Dir.Unit(), h.N);
	*scattered = vec.NewRay(h.P, vec.Add(reflected, vec.RandomV3InUnitSphere().MulScalar(m.Fuzz)));
	*attenuation = m.Albedo;
	return (vec.Dot(scattered.Dir, h.N) > 0);
}

type Dielectric struct {
	RefractiveIndex float64
}

func reflectance(cosine, refractiveIndex float64) float64 {
	r0 := math.Pow((1 - refractiveIndex) / (1 + refractiveIndex), 2)
	return r0 + (1-r0) * math.Pow((1 - cosine),5)
}

func (d Dielectric) Scatter(r vec.Ray, h *Hit, attenuation *vec.Color, scattered *vec.Ray) bool {
	*attenuation = vec.NewColor(1, 1, 1)
	var refractionRatio float64

	if h.FrontFace {
		refractionRatio = 1 / d.RefractiveIndex
	} else {
		refractionRatio = d.RefractiveIndex
	}

	unitDir := r.Dir.Unit()

	cosTheta := math.Min(vec.Dot(unitDir.MulScalar(-1), h.N), 1.0)
	sinTheta := math.Sqrt(1.0 - cosTheta * cosTheta)
	cannotRefract := refractionRatio * sinTheta > 1.0
	
	var direction vec.Vector3 
	if cannotRefract || reflectance(cosTheta, refractionRatio) > rand.Float64() {
		direction = vec.Reflect(unitDir, h.N)
	} else {
		direction = vec.Refract(unitDir, h.N, refractionRatio)
	}

	*scattered = vec.NewRay(h.P, direction)

	return true
}
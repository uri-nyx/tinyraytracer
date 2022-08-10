package raytracing

import (
	"math"
	"math/rand"
)

// V3 defines a 3d set
type V3 struct{
	x, y, z float64
}

type Point3d = V3
type Vector3 = V3

// NewV3 returns a new initialized 3d vector
func NewV3(x, y, z float64) V3 {
	return V3{x, y, z}
}

func RandomV3() V3 {
	return V3{rand.Float64(), rand.Float64(), rand.Float64()}
}

func RandomV3Range(min, max float64) V3 {
	return V3{min + (max-min)*rand.Float64(), min + (max-min)*rand.Float64(), min + (max-min)*rand.Float64()}
}

func RandomV3InUnitSphere() V3 {
	for {
		p := RandomV3Range(-1, 1)
		if p.LengthSquared() >= 1 {
			continue
		}

		return p
	}
}

 func RandomInUnitDisk() V3 {
    for {
        p := NewV3(rand.Float64(), rand.Float64(), 0);
        if p.LengthSquared() >= 1 {
			continue
		}
		
        return p;
    }
}

func RandomUnitV3() V3 {
	return RandomV3InUnitSphere().Unit()
}

func (v V3) NearZero() bool {
	s := 1e-8
	return math.Abs(v.x) < s && math.Abs(v.y) < s && math.Abs(v.z) < s 
}

func (v V3) X() float64 {
	return v.x
}

func (v V3) Y() float64 {
	return v.y
}

func (v V3) Z() float64 {
	return v.z
}


// MulScalar multiplies each component of a vector vy a scalar
func (v V3) MulScalar(s float64) V3 {
	return NewV3(v.x * s, v.y * s, v.z * s)
}

// DivScalar divides each component of a vector vy a scalar
func (v V3) DivScalar(s float64) V3 {
	return NewV3(v.x / s, v.y / s, v.z / s)
}

// AddScalar adds each component of a vector to a scalar
func (v V3) AddScalar(s float64) V3 {
	return NewV3(v.x + s, v.y + s, v.z + s)
}

// SubScalar subtracts a scalar from each component of a vector 
func (v V3) SubScalar(s float64) V3 {
	return NewV3(v.x - s, v.y - s, v.z - s)
}

// LengthSquared returns the square of the vector's length
func (v V3) LengthSquared() float64 {
	return v.x * v.x + v.y * v.y + v.z * v.z
}

// Length returns the length (magnitude) of the vector
func (v V3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

// Uinit returns the normalized vector for a vector
func (v V3) Unit() V3 {
	return v.DivScalar(v.Length())
}

// Vector Arithmetic

// Adds vectors v and w
func Add(v, w V3) V3 {
	return NewV3(v.x + w.x, v.y + w.y, v.z + w.z) 
}

// Subtracts vector w from v
func Sub(v, w V3) V3 {
	return NewV3(v.x - w.x, v.y - w.y, v.z - w.z)
}

// Multiplies each component of v by its counterpart on w
func Mul(v, w V3) V3 {
	return NewV3(v.x * w.x, v.y * w.y, v.z * w.z)
}

// Returns the dot product of v·w
func Dot(v, w V3) float64 {
	return v.x*w.x + v.y*w.y + v.z*w.z
}

// Returns the cross product of v×w
func Cross(v, w V3) V3 {
	return NewV3(v.y*w.z - v.z*w.y, v.z*w.x - v.x*w.z, v.x*w.y - v.y*w.x)
}

// COPILOT DID THIS v

// Returns the angle between vectors v and w
func Angle(v, w V3) float64 {
	return math.Acos(Dot(v.Unit(), w.Unit()))
}

// Returns the vector projection of v onto w
func Project(v, w V3) V3 {
	return w.MulScalar(Dot(v, w)/Dot(w, w))
}

// Returns the vector rejection of v onto w
func Reject(v, w V3) V3 {
	return Sub(v, Project(v, w))
}

// Returns the vector reflection of v about w
func ReflectCopilot(v, w V3) V3 {
	return Sub(v, Project(v, w).MulScalar(2))
}

func Reflect(v, w V3) V3 {//v - 2*dot(v,n)*n;
	return Sub(v, w.MulScalar(2 * Dot(v, w)))
}

// Returns the vector refraction of v about w
func RefractCopilot(v, w V3, eta float64) V3 {
	return Sub(v, Project(v, w).MulScalar(2)).MulScalar(eta).AddScalar(1 - eta) // v - project(v,w)*2*eta + (1-eta)
}

func Refract(v, w V3, eta_over_eta float64) V3 {
	cos_theta := math.Min(Dot(v.MulScalar(-1), w), 1)
	r_out_perp := Add(v, w.MulScalar(cos_theta)).MulScalar(eta_over_eta)
	r_out_parallel := w.MulScalar(- math.Sqrt(math.Abs(1 - r_out_perp.LengthSquared())))
	return Add(r_out_perp, r_out_parallel)
}
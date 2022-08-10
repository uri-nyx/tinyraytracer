package raytracing

type Ray struct{
	O Point3d
	Dir Vector3
}

func NewRay(o Point3d, dir Vector3) Ray {
	return Ray{o, dir.Unit()}
}

func (r Ray) At(t float64) V3 {
	return Add(r.O, r.Dir.MulScalar(t))
}
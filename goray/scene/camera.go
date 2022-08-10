package scene

import (
	"math"
	vec "github.com/uri-nyx/tinyraytracer/goray/raytracing"
)

type Camera struct {
	Image *Image

	View_h, View_w, Focal, lensRadius float64

	origin                                  vec.Point3d
	horizontal, vertical, lower_left_corner vec.Vector3
	u, v, w vec.Vector3

}

func NewCamera(image *Image, vfov, aperture, focus float64, lookfrom, lookat, vup vec.V3) Camera {
	theta := toRadians(vfov)
	h := math.Tan(theta/2)
	view_h := 2 * h
	view_w := view_h * image.aspect_ratio



	w := vec.Sub(lookfrom, lookat).Unit()
	u := vec.Cross(vup, w).Unit()
	v := vec.Cross(w, u)

	origin := lookfrom
	horizontal := u.MulScalar(focus * view_w);
	vertical := v.MulScalar(focus * view_h);
	lower_left_corner := vec.Sub(vec.Sub(vec.Sub(origin, horizontal.DivScalar(2)), vertical.DivScalar(2)), w.MulScalar(focus))


	lensRadius := aperture / 2;

	return Camera{
		Image:  image,
		View_h: view_h, View_w: view_w, Focal: 1, lensRadius: lensRadius,
		origin: origin,
		horizontal: horizontal, vertical: vertical, lower_left_corner: lower_left_corner,
		v: v, u: u, w: w}
}

func (c Camera) GetRay(s, t float64) vec.Ray {
	rd := vec.RandomInUnitDisk().MulScalar(c.lensRadius );
    offset := vec.Add(c.u.MulScalar(rd.X()), c.v.MulScalar(rd.Y()))

	return vec.NewRay(vec.Add(c.origin, offset),
	 	   vec.Sub(vec.Sub(vec.Add(vec.Add(c.lower_left_corner, c.horizontal.MulScalar(s)), c.vertical.MulScalar(t)), c.origin), offset))
}

func toRadians(degree float64) float64 {
	return degree * math.Pi / 180
}
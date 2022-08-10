package main

import (
	"log"
	"math"
	"math/rand"
	"runtime"

	ppm "github.com/uri-nyx/tinyraytracer/goray/ppm"
	vec "github.com/uri-nyx/tinyraytracer/goray/raytracing"
	scene "github.com/uri-nyx/tinyraytracer/goray/scene"
)

const (
	MAX_DEPTH = 50
	GAMMA	 = 2.0
)

func randomRange(min, max float64) float64 {
	// Returns a random real in [min,max).
	return min + (max-min)*rand.Float64()
}

func RayColor(r vec.Ray, world scene.HittableList, depth int) vec.Color {
	hit := scene.Hit{}

	if depth <= 0 {
		return vec.NewColor(0, 0, 0)
	}

	if world.Hit(r, 0.001, math.Inf(+1), &hit) {
		var scattered vec.Ray
		var attenuation vec.Color
		if hit.Material.Scatter(r, &hit, &attenuation, &scattered) {
			return vec.Mul(attenuation, RayColor(scattered, world, depth-1))
		}

		return vec.NewColor(0, 0, 0)
	}

	unit_direction := r.Dir.Unit()
	t := 0.5 * (unit_direction.Y() + 1)

	return vec.Add(vec.NewColor(1, 1, 1).MulScalar(1-t),
		vec.NewColor(0.5, 0.7, 1).MulScalar(t))
}

func Render(imgName string, img scene.Image, cam scene.Camera, world scene.HittableList, samplesPerPixel int) *ppm.PpmImage {
	render := ppm.New(imgName, uint(img.Width), uint(img.Height))

	log.Default().Print("Rendering...\n")

	for j := img.Height - 1; j >= 0; j-- {
		log.Default().Printf("Scanlines remaining: %d\n", j)

		for i := 0; i < img.Width; i++ {

			func() {
				cores := runtime.GOMAXPROCS(0)

				px := make(chan vec.Color)

				for core := 0; core < cores; core++ {

					go func(core int, px chan<-vec.Color) {
						pixelColor := vec.NewColor(0, 0, 0)

						start := (samplesPerPixel / cores) * core
						end := start + (samplesPerPixel / cores)

						for s := start; s < end; s++ {
							u := (float64(i) + rand.Float64()) / float64(img.Width - 1)
							v := (float64(j) + rand.Float64()) / float64(img.Height - 1)
							r := cam.GetRay(u, v)
							pixelColor = vec.Add(pixelColor, RayColor(r, world, MAX_DEPTH) )
						}

						px <- pixelColor

					}(core, px)
				}

				sampledColor := vec.NewColor(0, 0, 0)
				for core := 0; core < cores; core++ {
					sampledColor = vec.Add(sampledColor, <- px)
				}

				render.WritePixel(sampledColor.DenormalizeSampledColor(samplesPerPixel, GAMMA))
			} ()
		}
	}

	log.Default().Printf("Done\n")
	return render
}

func random_scene() scene.HittableList {
    var world scene.HittableList

    ground_material := scene.Lambertian{Albedo: vec.NewColor(0.5, 0.5, 0.5)}
    world = append(world, scene.Sphere{vec.NewV3(0,-1000,0), 1000, ground_material});

    for a := -11; a < 11; a++ {
        for b := -11; b < 11; b++ {
            choose_mat := rand.Float64()
            center := vec.NewV3(float64(a) + 0.9*rand.Float64(), 0.2, float64(b) + 0.9*rand.Float64());

            if vec.Sub(center, vec.NewV3(4, 0.2, 0)).Length() > 0.9 {
                var sphere_material scene.Material;

                if (choose_mat < 0.8) {
                    // diffuse
                    albedo := vec.Mul(vec.RandomV3(), vec.RandomV3())
                    sphere_material = scene.Lambertian{Albedo: albedo}
                    world = append(world, scene.Sphere{center, 0.2, sphere_material})
                } else if (choose_mat < 0.95) {
                    // metal
                    albedo := vec.RandomV3Range(0.5, 1)
                    fuzz := randomRange(0, 0.5)
                    sphere_material = scene.Metal{albedo, fuzz}
                    world = append(world, scene.Sphere{center, 0.2, sphere_material})
                } else {
                    // glass
                    sphere_material = scene.Dielectric{1.5}
					world = append(world, scene.Sphere{center, 0.2, sphere_material})

                }
            }
        }
    }

    material1 := scene.Dielectric{1.5}
    world = append(world, scene.Sphere{vec.NewV3(0, 1, 0), 1.0, material1})

    material2 := scene.Lambertian{vec.NewColor(0.4, 0.2, 0.1)}
    world = append(world, scene.Sphere{vec.NewV3(-4, 1, 0), 1.0, material2})

    material3 := scene.Metal{vec.NewColor(0.7, 0.6, 0.5), 0.0};
    world = append(world, scene.Sphere{vec.NewV3(4, 1, 0), 1.0, material3})

    return world;
}

func main() {
	samples := 500
	img := scene.NewImage(16.0/9.0, 400)

	lookfrom, lookat, vup := vec.NewV3(13,2,3), vec.NewV3(0,0,0), vec.NewV3(0,1,0)
	cam := scene.NewCamera(&img, 20, 0.1, 10, lookfrom, lookat, vup)

	world := random_scene()


	rendered := Render("test.ppm", img, cam, world, samples)
	rendered.Save(".")
	log.Default().Printf("Saved image, %dx%d\n", img.Width, img.Height)
}

package ppm

import "testing"

func TestPpm(t *testing.T) {
		// test ppm package

		if err:= New("test.ppm", 256, 256).Save("."); err == nil {
			t.Errorf("Error saving image: %s", err)
		}
}
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/seiflotfy/cmts"
	"github.com/seiflotfy/count-min-log"
)

func main() {
	size := uint(2097152) // 4194304 bits base + ==> total bits ==> 1 MB
	cm := cmts.NewSketches(16384, 256)
	cl, _ := cml.NewSketch(size/4, 4, 1.00026)

	max := 500000
	now := time.Now()
	expected := make([]uint, max, max)
	zipf := rand.NewZipf(rand.New(rand.NewSource(now.UnixNano())), 1.1, 10.0, uint64(max)-1)

	seen := map[string]bool{}
	for k := uint64(0); len(seen) != max; k++ {
		if k%uint64(max) == 0 {
			fmt.Printf("\rCardinality %06d\t Hits: %d", len(seen), k)
		}
		i := zipf.Uint64()
		expected[i]++
		id := fmt.Sprintf("flow-%05d", i)
		seen[id] = true
		cm.Increment([]byte(id))
		cl.Update([]byte(id))
	}

	for i := range expected {
		// some minor print for easier visuals
		if i == 100 || i-50 == (len(expected)/2)+1 || i == len(expected)-100 {
			fmt.Printf("\n---")
		}

		if (i > len(expected)-100 || i < 100) || (i < (len(expected)/2)+50 && i > (len(expected)/2)-50) {
			// id
			id := fmt.Sprintf("flow-%05d", i)

			// estimation
			est2 := float64(cm.Get([]byte(id)))
			est3 := float64(cl.Query([]byte(id)))

			// error ratio
			ratio2 := 100*est2/float64(expected[i]) - 100
			ratio3 := 100*est3/float64(expected[i]) - 100

			fmt.Printf("\n%s:\t\texpected %d\t\tcmt ~= %.2f%%\tcml ~= %.2f%%", id, expected[i], ratio2, ratio3)
		}
	}
}

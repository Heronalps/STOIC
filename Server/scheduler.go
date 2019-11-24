/*
Scheduler decides where to run WTB inferencing job given Pi bandwidth, number of images
*/

package server

import (
	"fmt"
	"sort"
)

/*
Schedule obtains total times of four scenarios
*/
func Schedule() {
	imageNum := ImageCache()
	totalTimes := GetTotalTime(imageNum)
	fmt.Println(totalTimes)
	// Sort the totalTimes map by key
	keys := make([]float64, 0, len(totalTimes))
	for k := range totalTimes {
		keys = append(keys, k)
	}
	keys = sort.Float64Slice(keys)
	fmt.Printf("The task is scheduled at %s for %f seconds\n", totalTimes[keys[0]], keys[0])
}

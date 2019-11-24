/*
Scheduler decides where to run WTB inferencing job given Pi bandwidth, number of images
*/

package server

import (
	"fmt"
	"os/exec"
	"sort"
)

/*
SelectRunTime select the runtime among four scenarios
*/
func SelectRunTime(imageNum int64) string {
	totalTimes := GetTotalTime(imageNum)
	fmt.Println(totalTimes)
	// Sort the totalTimes map by key
	keys := make([]float64, 0, len(totalTimes))
	for k := range totalTimes {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	fmt.Printf("The task is scheduled at %s for %f seconds\n", totalTimes[keys[0]], keys[0])
	return totalTimes[keys[0]]
}

/*
Schedule obtains total times of four scenarios
*/
func Schedule() {
	imageNum := ImageCache()
	switch runtime := SelectRunTime(imageNum); runtime {
	case "euca":
		//RunOnEuca(imageNum)
	default:
		//RunOnNautilus(runtime, imageNum)
	}
}

/*
RunOnEuca runs the task on mini euca edge cloud with AVX support
*/
func RunOnEuca(imageNum int64) {
	PATH := HomeDir() + "/GPU_Serverless/kubeless/image_clf/inference/local_version/image_clf_inf.py "
	cmdVenv := "source " + HomeDir() + "/GPU_Serverless/venv_avx/bin/activate"
	exec.Command(cmdVenv)
	cmdRun := "python " + PATH + string(imageNum)
	exec.Command(cmdRun)
}

/*
RunOnNautilus runs the task on Nautilus public cloud
*/
func RunOnNautilus(runtime string, imageNum int64) {

}

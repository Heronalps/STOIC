/*
Scheduler decides where to run WTB inferencing job given Pi bandwidth, number of images
*/

package server

import (
	"fmt"
	exec "os/exec"
	"sort"
	"strconv"
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
		fmt.Println("Running on euca...")
		RunOnEuca(imageNum)
	default:
		fmt.Println("Running on Nautilus...")
		RunOnNautilus(runtime, imageNum)
	}
}

/*
RunOnEuca runs the task on mini euca edge cloud with AVX support
*/
func RunOnEuca(imageNum int64) {
	var output []byte
	var err error
	var cmd *exec.Cmd
	repoPATH := HomeDir() + "/GPU_Serverless"

	//Activate the venv in python project repo
	cmdVenv := "source " + "venv_avx/bin/activate"
	cmd = exec.Command("bash", "-c", cmdVenv)
	cmd.Dir = repoPATH
	output, err = cmd.Output()
	if err != nil {
		fmt.Printf("Error activating venv. msg: %s \n", err.Error())
		return
	}
	fmt.Printf("Activated AVX support... %s\n", output)

	// Run WTB image classification task
	FILE := "image_clf_inf.py "
	cmdRun := "python " + FILE + strconv.Itoa(int(imageNum))
	cmd = exec.Command("bash", "-c", cmdRun)
	cmd.Dir = repoPATH + "/kubeless/image_clf/inference/local_version"
	output, err = cmd.Output()
	if err != nil {
		fmt.Printf("Error running task. msg: %s \n", err.Error())
		return
	}
	fmt.Printf("Output of task %s\n", output)
}

/*
RunOnNautilus runs the task on Nautilus public cloud
*/
func RunOnNautilus(runtime string, imageNum int64) {

}

package client

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"

	// oidc authorization
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

/*
Request sends request to Nautilus based on runtime and image number
*/
func Request(runtime string, imageNum int) []byte {
	namespace := "racelab"
	deployment := "image-clf-inf"
	var (
		output    []byte
		isGPUSame bool
		cmd       string
		err       error
		result    []byte
	)
	resultChannel := make(chan []byte)

	fmt.Printf("Making request to Nautilus %s %d \n", runtime, imageNum)
	switch runtime {
	case "cpu":
		isGPUSame, _, err = Deploy(namespace, deployment, 0)
	case "gpu1":
		isGPUSame, _, err = Deploy(namespace, deployment, 1)
	case "gpu2":
		isGPUSame, _, err = Deploy(namespace, deployment, 2)
	}
	if err != nil {
		log.Println(err.Error())
		return result
	}

	// Wait 3 second for deployment to complete
	time.Sleep(3 * time.Second)
	go func() {
		// If the pod is re-deployed,
		// make initial kubeless call to depolyed function to avoid cold start
		if !isGPUSame {
			fmt.Println("Probing deployed kubeless function to avoid cold start ...")
			cmd = "sh ./scripts/invoke_inf.sh " + strconv.Itoa(1)
			output, err = exec.Command("bash", "-c", cmd).Output()
			if err != nil {
				fmt.Println("Error msg : ", err.Error())
			}
			fmt.Println(string(output))
			fmt.Println("Finish Probing ...")
		} else {
			fmt.Println("Using same pod, No probing needed....")
		}

		//make kubeless call to deployed function
		cmd = "sh ./scripts/invoke_inf.sh " + strconv.Itoa(imageNum)
		output, err = exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			fmt.Println("Error msg : ", err.Error())
		}
		fmt.Println(string(output))
		resultChannel <- output
	}()

	result = <-resultChannel
	return result
}

/*
StaticKubeconfig : static kubernetes config
*/
var kubeconfig string

func getKubeConfig() string {
	if kubeconfig != "" {
		return kubeconfig
	}
	if home := homedir.HomeDir(); home != "" {
		flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "service-account"), "(optional) absolute path to the kubeconfig file")
	} else {
		flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	return kubeconfig
}

/*
QueryGPUNum queries the number of GPU in current pod
*/
func QueryGPUNum(namespace string, deployment string) int64 {
	kubeconfig := getKubeConfig()
	// use the current context in kubeconfig
	// This config is credential information of kubernetes
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	result, getErr := deploymentsClient.Get(deployment, metav1.GetOptions{})
	if getErr != nil {
		log.Printf("Failed to get latest version of Deployment %v", getErr)
	}
	numGpu := result.Spec.Template.Spec.Containers[0].Resources.Requests["nvidia.com/gpu"]
	return numGpu.Value()
}

/*
Deploy patches kubeless function on Nautilus based on number of GPU
*/
func Deploy(namespace string, deployment string, NumGPU int64) (bool, float64, error) {
	kubeconfig := getKubeConfig()
	// use the current context in kubeconfig
	// This config is credential information of kubernetes
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	deploymentsClient := clientset.AppsV1().Deployments(namespace)

	fmt.Println("Updating deployment...")
	var (
		prevNumGPU int64
		currNumGPU int64
		timeStamp0 time.Time
		timeStamp1 time.Time
		duration   float64
	)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {

		// Recover from "panic: assignment to entry in nil map"
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Panic : %s \n", r)
			}
		}()

		result, getErr := deploymentsClient.Get(deployment, metav1.GetOptions{})
		if getErr != nil {
			log.Printf("Failed to get latest version of Deployment %v", getErr)
			return getErr
		}

		numGpu := result.Spec.Template.Spec.Containers[0].Resources.Requests["nvidia.com/gpu"]
		prevNumGPU = numGpu.Value()
		fmt.Printf("Current Number of GPU is %v \n", prevNumGPU)
		quant := resource.NewQuantity(NumGPU, resource.DecimalSI)
		result.Spec.Template.Spec.Containers[0].Resources.Limits["nvidia.com/gpu"] = *quant
		result.Spec.Template.Spec.Containers[0].Resources.Requests["nvidia.com/gpu"] = *quant

		timeStamp0 = time.Now()

		_, updateErr := deploymentsClient.Update(result)
		fmt.Printf("Updated Number of GPU is %v \n", quant.Value())
		currNumGPU = quant.Value()

		fmt.Println("Waiting kubeless function to be deployed...")
		time.Sleep(3 * time.Second)
		result, getErr = deploymentsClient.Get(deployment, metav1.GetOptions{})
		for !strings.HasSuffix(result.Status.Conditions[1].Message, "has successfully progressed.") {
			time.Sleep(3 * time.Second)
			result, getErr = deploymentsClient.Get(deployment, metav1.GetOptions{})
			fmt.Printf("Message : %s \n", result.Status.Conditions[1].Message)
		}

		timeStamp1 = time.Now()
		duration = float64(timeStamp1.Sub(timeStamp0))

		fmt.Println("Kubeless function is successfully deployed...")
		if updateErr != nil {
			fmt.Println(updateErr)
		}
		return updateErr
	})

	if retryErr != nil {
		log.Printf("Update failed: %v", retryErr)
	}

	return prevNumGPU == currNumGPU, duration, retryErr
}

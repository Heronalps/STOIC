package client

import (
	"flag"
	"fmt"
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
func Request(runtime string, imageNum int) {
	namespace := "racelab"
	deployment := "image-clf-inf"
	fmt.Printf("Making request to Nautilus %s %d \n", runtime, imageNum)
	switch runtime {
	case "cpu":
		deploy(namespace, deployment, 0)
	case "gpu1":
		deploy(namespace, deployment, 1)
	case "gpu2":
		deploy(namespace, deployment, 2)
	}
	//Make initial kubeless call to depolyed function to avoid cold start
	//Wait a second for deployment to complete
	time.Sleep(1 * time.Second)

	fmt.Println("Probing deployed kubeless function to avoid cold start ...")
	cmd := "sh invoke_inf.sh " + strconv.Itoa(1)
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Println("Error msg : ", err.Error())
	}
	fmt.Println(string(output))
	fmt.Println("Finish Probing ...")

	//make kubeless call to deployed function
	cmd = "sh invoke_inf.sh " + strconv.Itoa(imageNum)
	output, err = exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Println("Error msg : ", err.Error())
	}
	fmt.Println(string(output))
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
		flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	return kubeconfig
}

/*
Deploy patches kubeless function on Nautilus based on number of GPU
*/
func deploy(namespace string, deployment string, NumGPU int64) {
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

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentsClient.Get(deployment, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment %v", getErr))
		}

		numGpu := result.Spec.Template.Spec.Containers[0].Resources.Requests["nvidia.com/gpu"]
		fmt.Printf("Current Number of GPU is %v \n", numGpu.Value())
		quant := resource.NewQuantity(NumGPU, resource.DecimalSI)
		result.Spec.Template.Spec.Containers[0].Resources.Limits["nvidia.com/gpu"] = *quant
		result.Spec.Template.Spec.Containers[0].Resources.Requests["nvidia.com/gpu"] = *quant
		_, updateErr := deploymentsClient.Update(result)
		fmt.Printf("Updated Number of GPU is %v \n", quant.Value())

		fmt.Println("Waiting kubeless function to be deployed...")
		time.Sleep(3 * time.Second)
		result, getErr = deploymentsClient.Get(deployment, metav1.GetOptions{})
		for !strings.HasSuffix(result.Status.Conditions[1].Message, "has successfully progressed.") {
			time.Sleep(3 * time.Second)
			result, getErr = deploymentsClient.Get(deployment, metav1.GetOptions{})
			fmt.Printf("Message : %s \n", result.Status.Conditions[1].Message)
		}
		fmt.Println("Kubeless function is successfully deployed...")
		if updateErr != nil {
			fmt.Println(updateErr)
		}
		return updateErr
	})

	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
}

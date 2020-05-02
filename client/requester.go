package client

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	retrygo "github.com/avast/retry-go"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	// oidc authorization
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

/*
RunOnNautilus sends request to Nautilus based on runtime and image number
return output & processing time
*/
func RunOnNautilus(runtime string, zipPath string, imageNum int, app string, version string, transferTime float64) ([]byte, bool, *TimeLog) {
	var (
		output         []byte
		isGPUSame      bool
		isDeployed     bool
		cmd            string
		err            error
		procTime       float64
		deploymentTime float64
		cmdRun         *exec.Cmd
		retryErr       error
	)
	retryErr = retrygo.Do(
		func() error {
			TransferBatch(zipPath)
			return nil
		},
	)
	if retryErr != nil {
		fmt.Printf("Transfer zip to S3 failed: %v \n", retryErr.Error())
	}

	fmt.Printf("Making request to Nautilus %s %d \n", runtime, imageNum)

	switch runtime {
	case "cpu":
		isGPUSame, isDeployed, deploymentTime, err = Deploy(namespace, RunDeployment, 0, app)
	case "gpu1":
		isGPUSame, isDeployed, deploymentTime, err = Deploy(namespace, RunDeployment, 1, app)
	case "gpu2":
		isGPUSame, isDeployed, deploymentTime, err = Deploy(namespace, RunDeployment, 2, app)
	}
	if err != nil {
		log.Println(err.Error())
		return output, isDeployed, nil
	}
	if !isDeployed {
		log.Println("Runtime cannot be deployed. Reschedule workload...")
		return output, isDeployed, nil
	}

	// Wait 3 second for deployment to complete
	time.Sleep(3 * time.Second)
	retryErr = retrygo.Do(
		func() error {
			// If the pod is re-deployed,
			// make initial kubeless call to depolyed function to avoid cold start
			// Update - After refactoring Deploy function, every invocation will redeploy.
			// Add true condition to always probe
			if true || !isGPUSame {
				fmt.Println("Probing deployed kubeless function to avoid cold start ...")

				// The pseudoPath is a workaround of golang exec incompotence option
				cmd = fmt.Sprintf("%s sh ./scripts/invoke_inf.sh pseudoPath", serviceAccountConfig)
				cmdRun = exec.Command("bash", "-c", cmd)

				if output, err = cmdRun.Output(); err != nil {
					fmt.Println("Error msg : ", err.Error())
					return err
				}
				fmt.Println(string(output))
				fmt.Println("Finish Probing ...")
			} else {
				fmt.Println("Using same pod, No probing needed....")
			}

			// Make kubeless call to deployed function
			// Put presetZipPath for testing purpose
			cmd = fmt.Sprintf("%s sh ./scripts/invoke_inf.sh %s", serviceAccountConfig, presetZipPath)
			cmdRun = exec.Command("bash", "-c", cmd)
			if output, err = cmdRun.Output(); err != nil {
				fmt.Println("Error msg : ", err.Error())
				return err
			}
			procTime = ParseElapsed(output)
			fmt.Println(string(output))
			//resultChannel <- output
			return err
		})

	if retryErr != nil {
		fmt.Printf("Request failed: %v ...", retryErr.Error())
	}

	return output, true, CreateTimeLog(transferTime, deploymentTime, procTime)
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
	var (
		numGpu resource.Quantity
	)

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
		return 0
	}
	if result.Spec.Template.Spec.Containers != nil {
		numGpu = result.Spec.Template.Spec.Containers[0].Resources.Requests["nvidia.com/gpu"]
		return numGpu.Value()
	}

	return 0
}

/*
Deploy function deploys new deployment and patches kubeless function on Nautilus based on number of GPU
*/
func Deploy(namespace string, deployment string, NumGPU int64, app string) (bool, bool, float64, error) {
	var (
		retryErr     error
		deployResult DeployResult
		duration     float64
		prevNumGPU   int64
		currNumGPU   int64
	)

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
	// Create Deployment Start Timestamp
	tsCreateStart := time.Now()
	for !deployResult.Progressed {
		err = CreateDeployment(app, NumGPU)
		if err != nil {
			fmt.Println(err.Error())
		}
		//Async call using channel await to join
		deployResult = <-IsPodReady(deployment, deploymentsClient)
		if deployResult.Timeout {
			break
		}
	}

	// Create Deployment End Timestamp
	tsCreateEnd := time.Now()

	// Change GPU number is executed in the deploy.sh

	// retryErr = retry.RetryOnConflict(retry.DefaultRetry, func() error {

	// 	// Recover from "panic: assignment to entry in nil map"
	// 	defer func() {
	// 		if r := recover(); r != nil {
	// 			fmt.Printf("Panic : %s \n", r)
	// 		}
	// 	}()

	// 	result, getErr := deploymentsClient.Get(deployment, metav1.GetOptions{})
	// 	if getErr != nil {
	// 		log.Printf("Failed to get latest version of Deployment %v", getErr)
	// 		return getErr
	// 	}

	// 	numGpu := result.Spec.Template.Spec.Containers[0].Resources.Requests["nvidia.com/gpu"]
	// 	prevNumGPU = numGpu.Value()
	// 	fmt.Printf("Current Number of GPU is %v \n", prevNumGPU)
	// 	quant := resource.NewQuantity(NumGPU, resource.DecimalSI)
	// 	result.Spec.Template.Spec.Containers[0].Resources.Limits["nvidia.com/gpu"] = *quant
	// 	result.Spec.Template.Spec.Containers[0].Resources.Requests["nvidia.com/gpu"] = *quant

	// 	// Update Start Timestamp
	// 	tsUpdateStart = time.Now()

	// 	_, updateErr := deploymentsClient.Update(result)
	// 	fmt.Printf("Updated Number of GPU is %v \n", quant.Value())
	// 	currNumGPU = quant.Value()

	// 	fmt.Println("Waiting kubeless function to be deployed...")

	// 	//Async call using channel await to join
	// 	_ = <-IsPodReady(deployment, deploymentsClient)
	// 	// Update End Timestamp
	// 	tsUpdateEnd = time.Now()

	// 	fmt.Println("Kubeless function is successfully deployed...")
	// 	if updateErr != nil {
	// 		fmt.Println(updateErr)
	// 	}
	// 	return updateErr
	// })

	// if retryErr != nil {
	// 	log.Printf("Update failed: %v", retryErr.Error())
	// }

	// Convert nanoseconds to seconds; Substract 3 seconds of waiting time
	// duration = (float64(tsCreateEnd.Sub(tsCreateStart))+float64(tsUpdateEnd.Sub(tsUpdateStart)))*1e-9 - 3.0
	duration = (float64(tsCreateEnd.Sub(tsCreateStart))) * 1e-9
	return prevNumGPU == currNumGPU, deployResult.Progressed, duration, retryErr
}

/*
CreateDeployment creates new deployment
*/
func CreateDeployment(app string, NumGPU int64) error {
	var (
		pythonVersion string = "3.6"
		err           error
		cmd           *exec.Cmd
	)
	if app == "image-clf-inf" {
		pythonVersion = "3.6"

		// image-clf-inf37 is for periodically querying deployment time
	} else if app == "image-clf-inf37" {
		pythonVersion = "3.7"
	}
	// Create new deployment
	cmdRun := fmt.Sprintf("%s sh ./scripts/deploy.sh %s %s %d", serviceAccountConfig, app, pythonVersion, NumGPU)
	cmd = exec.Command("bash", "-c", cmdRun)
	fmt.Printf("Creating new deployment of app %s \n", app)

	if _, err = cmd.Output(); err != nil {
		fmt.Printf("Error creating deployment. msg: %s \n", err.Error())
		return err
	}

	return err
}

/*
IsPodReady makes judgement about if a deployment is ready
*/
func IsPodReady(deployment string, deploymentsClient appsv1.DeploymentInterface) <-chan DeployResult {
	r := make(chan DeployResult)
	var (
		result       *v1.Deployment
		getErr       error
		progressed   bool
		timeout      bool
		deployResult DeployResult
	)
	go func() {
		for true {
			result, getErr = deploymentsClient.Get(deployment, metav1.GetOptions{})
			time.Sleep(3 * time.Second)
			if len(result.Status.Conditions) > 1 {
				progressed = strings.HasSuffix(result.Status.Conditions[1].Message, "has successfully progressed.")
				timeout = strings.HasSuffix(result.Status.Conditions[1].Message, "has timed out progressing.")
				fmt.Printf("Message : %s \n", result.Status.Conditions[1].Message)
			}

			if progressed || timeout {
				break
			}
		}
		if getErr != nil {
			fmt.Println(getErr.Error())
		}
		deployResult.Progressed = progressed
		deployResult.Timeout = timeout
		r <- deployResult
	}()

	return r
}

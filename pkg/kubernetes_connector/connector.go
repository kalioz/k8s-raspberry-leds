package kubernetes


import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"sort"
	"strconv"
	"strings"
)

const LabelPriority = "raspberry.color/priority"
const LabelColor = "raspberry.color/color"
const LabelPriorityBadFormat = "Warning - Pod %s - label %s should be an integer (current value : %s)"
const LabelColorMissing = "Error - Pod %s - label "+LabelColor +" is missing"

func GetPodsWithColor(limit int, nodeName string) ([]v1.Pod, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientSet
	clientSet, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err.Error())
	}

	// N.B. do not use the Limit parameter of the ListOption as we sort the data after
	podList, err := clientSet.CoreV1().Pods("").List(	metav1.ListOptions{LabelSelector: "raspberry.color/enabled=true"})

	if err != nil {
		log.Printf("ERROR - could not get pods from API (%s)", err)
		return nil, err
	}

	// remove pods which are not on the same node as this pod
	pods := FilterForNodeName(podList.Items, nodeName)

	// sort pods
	pods = sortPods(pods)

	if len(pods) > limit {
		pods = pods[0:limit]
	}

	return pods, nil
}

func FilterForNodeName(pods []v1.Pod, nodeName string) []v1.Pod {
	output := make([]v1.Pod, 0)
	for _, pod := range pods {
		if pod.Spec.NodeName == nodeName {
			output = append(output, pod)
		}
	}
	return output
}

func sortPods(pods []v1.Pod) []v1.Pod{
	var err error
	sort.Slice(pods, func(i, j int) bool{
		// sort by priority
		priority1 := 0
		priority2 := 0

		if priority, ok := pods[i].Labels[LabelPriority]; ok {
			priority1, err = strconv.Atoi(priority)
			if err != nil {
				log.Printf(LabelPriorityBadFormat, pods[i].Name, LabelPriority, priority)
			}
		}

		if priority, ok := pods[j].Labels[LabelPriority]; ok {
			priority2, err = strconv.Atoi(priority)
			if err != nil {
				log.Printf(LabelPriorityBadFormat, pods[j].Name, LabelPriority, priority)
			}
		}

		if priority1 != priority2 {
			return priority1 > priority2
		}

		// sort by color, alphabetically
		var color1 = ""
		var color2 = ""

		if color, ok := pods[i].Labels[LabelColor]; ok {
			color1 = color
		} else {
			log.Fatalf(LabelColorMissing, pods[i].Name)
		}

		if color, ok := pods[j].Labels[LabelColor]; ok {
			color2 = color
		}else {
			log.Fatalf(LabelColorMissing, pods[j].Name)
		}

		if color1 != color2 {
			if color1 == "" {
				return false
			}
			if color2 == "" {
				return true
			}
			return strings.Compare(color1, color2) < 0
		}

		// in last resort sort by namespace+name
		return strings.Compare(pods[i].Namespace+"_"+pods[i].Name, pods[j].Namespace+"_"+pods[j].Name) < 0
	})

	return pods
}

func testConnectionToKubernetes() {

}

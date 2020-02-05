package main

import (
	"github.com/kalioz/raspberry-leds/pkg/pixel"
	k8s "github.com/kalioz/raspberry-leds/pkg/kubernetes_connector"
	"gopkg.in/go-playground/colors.v1"
	"log"
	"os"
	"time"
)

const NumberOfPixels = 8
const LabelColorBadFormat = "Error - Pod %s - the color label could not be parsed (value: '%s')"
const LedFrameRate = 2000 //GetVariable("LED_FRAME_RATE",20000) // number of changes that can occur in a second to the gpio pin
const DutyCycle = 255 // the max value a pin will be able to attain. max = 255 (uint8)
const UseAbsoluteColors = true // if set to true, the leds will only be able to display 8 colors : R,G,B,RG,RB,GB,RGB and black
const Debug = true

func main() {
	log.Println("Starting application... Adding pixels")
	var pixels [NumberOfPixels]*pixel.Pixel

	//TODO set this in variables
	//pixels[0] = pixel.NewPixel(0,1,2)
	//pixels[1] = pixel.NewPixel(0,1,2)
	//pixels[2] = pixel.NewPixel(0,1,2)
	//pixels[3] = pixel.NewPixel(0,1,2)
	//pixels[4] = pixel.NewPixel(0,1,2)
	//pixels[5] = pixel.NewPixel(0,1,2)
	//pixels[6] = pixel.NewPixel(0,1,2)
	//pixels[7] = pixel.NewPixel(0,1,2)


	var nodeName = getNodeName()

	log.Println("Starting to monitor")
	if UseAbsoluteColors {
		mainLoopAbsoluteColors(pixels[:], nodeName)
	} else {
		//TODO
	}
}

func getNodeName() string {
	if name, ok := os.LookupEnv("NODE_NAME"); ok {
		return name
	}

	log.Println("Warning - environment variable NODE_NAME not set; defaulting to hostname (https://kubernetes.io/docs/tasks/inject-data-application/environment-variable-expose-pod-information/#use-pod-fields-as-values-for-environment-variables)")

	hostname, err := os.Hostname()

	if err != nil {
		log.Fatalf("Could not query hostname - %s", err)
	}

	return hostname
}

func mainLoopAbsoluteColors(pixels []*pixel.Pixel, nodeName string){
	ticker := time.NewTicker(1*time.Second)
	for {
		select {
		case _ = <-ticker.C:
			GetPodsWithColorsAndChangeDisplay(pixels, nodeName)
			for _, pixel := range pixels {
				pixel.Display(0)
			}
		}
	}

}

func parseColor(color string) ( *colors.RGBColor, error) {
	c, err := colors.Parse(color)
	if err != nil {
		return nil, err
	}
	output := c.ToRGB()

	return output, nil
}

/**
* Change the value of the lights to be displayed
* It does not force the lights to be displayed.
 */
func GetPodsWithColorsAndChangeDisplay(pixels []*pixel.Pixel, nodeName string) error {
	pods, err := k8s.GetPodsWithColor(len(pixels), nodeName)
	if err != nil {
		return err
	}

	for i, pixel := range pixels {
		if len(pods) > i {
			color, err := parseColor(pods[i].Labels[k8s.LabelColor])

			if err != nil {
				log.Printf(LabelColorBadFormat, pods[i].Name, pods[i].Labels[k8s.LabelColor])
				pixel.ChangeColor(0,0,0)
			} else {
				pixel.ChangeColor(color.R, color.G, color.B)
			}
		} else {
			pixel.ChangeColor(0,0,0)
		}
	}

	return nil
}

func RefreshPixels(pixels []*pixel.Pixel, t time.Time){
	var cycle uint8
	if UseAbsoluteColors {
		cycle = DutyCycle / 2
	} else {
		cycle = uint8((t.Nanosecond() * LedFrameRate) % DutyCycle)
	}
	for _, pixel := range pixels {
		pixel.Display(uint8(cycle))
	}
}
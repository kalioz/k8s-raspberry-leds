package pixel

import "github.com/stianeikeland/go-rpio"

type Led struct {
	PinNumber int
	Value uint8

	// illustrate the current state of the pin, to reduce system calls.
	currentState bool

	rpio.Pin
}

func NewLed(pinNumber int) *Led {
	led := &Led{Pin: rpio.Pin(pinNumber), PinNumber:pinNumber, Value:0 }

	// by default, stop the led
	led.Output()
	led.Low()
	led.currentState = false

	return led
}

/**
 * change the pin output in function of the current status of the cycle
 * if cycle < led.Value, the led is lit
 * this is a simple Duty Cycle where the period is a fixed value of 255 and the duty cycle is equal to Value.
 */
func (led *Led) Display(cycle uint8){
	if led.Value > cycle {
		if led.currentState == false {
			led.High()
			led.currentState = true
		}
	} else {
		if led.currentState {
			led.Low()
			led.currentState = false
		}
	}
}
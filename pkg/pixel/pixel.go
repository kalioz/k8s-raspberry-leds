package pixel

type Pixel struct {
	Red *Led
	Green *Led
	Blue *Led
}

func NewPixel(redPinNumber int, bluePinNumber int, greenPinNumber int) *Pixel {
	pixel := Pixel{
		Red: NewLed(redPinNumber),
		Green: NewLed(greenPinNumber),
		Blue: NewLed(bluePinNumber),
	}

	return &pixel
}

func (pixel *Pixel) ChangeColor(red uint8, green uint8, blue uint8) {
	pixel.Red.Value = red
	pixel.Green.Value = green
	pixel.Blue.Value = blue
}

// Change the pixel color depending of the current number of cycles.
func (pixel *Pixel) Display(cycle uint8)  {
	pixel.Red.Display(cycle)
	pixel.Green.Display(cycle)
	pixel.Blue.Display(cycle)
}
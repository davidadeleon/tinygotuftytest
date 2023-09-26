package main

import (
	"image/color"
	"machine"
	"time"
	"tinygo.org/x/drivers"
)

const clockHz = 133000000

func main() {
	time.Sleep(5 * time.Second)
	println("Initializing Display")
	display := ST7789{
		cs:                machine.LCD_CS,
		dc:                machine.LCD_DC,
		wr:                machine.LCD_WR,
		rd:                machine.LCD_WR,
		d0:                machine.LCD_DB0,
		bl:                machine.LCD_BACKLIGHT,
		stateMachineIndex: 0,
		dmaChannel:        2,
		width:             320,
		height:            240,
		rotation:          drivers.Rotation0,
	}

	println("Initializing PIO")
	display.pio = machine.PIO0
	println("Configuring PIO")
	display.pio.Configure()
	println("Parallel Init")
	display.ParallelInit()

	// Setup DMA
	println("Setting Up DMA")
	dmaConfig := getDefaultDMAConfig(display.dmaChannel)
	setTransferDataSize(dmaConfig, DMA_SIZE_8)
	setBSwap(dmaConfig, false)
	setDREQ(dmaConfig, display.pio.Device.GetIRQ())
	dmaChannelConfigure(display.dmaChannel, dmaConfig, display.pio.Device.TXF0.Reg, 0, 0, false)

	machine.LCD_RD.High()

	println("Display Common Init")
	display.CommonInit()

	println("Making Screen Blue")
	blue := color.RGBA{255, 255, 255, 255}
	display.FillRectangle(0, 0, 320, 240, blue)
}

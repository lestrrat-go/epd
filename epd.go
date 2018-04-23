package epd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ecc1/gpio"
	"github.com/ecc1/spi"
)

func (cmd Command) Byte() byte {
	return byte(cmd)
}

func New() *EPD {
	var pinBUSY gpio.InputPin
	var pinCS gpio.OutputPin
	var pinDC gpio.OutputPin
	var pinRST gpio.OutputPin

	d, err := spi.Open("/dev/spidev0.0", 2000000, 0)
	if err != nil {
		panic("failed to open SPI device: " + err.Error())
	}
	d.SetMode(0)

	{
		p, err := gpio.Input(defPinBusy, false)
		if err != nil {
			panic("failed to setup BUSY pin: " + err.Error())
		}
		pinBUSY = p
	}

	{
		p, err := gpio.Output(defPinCS, false, false)
		if err != nil {
			panic("failed to setup CS pin: " + err.Error())
		}
		pinCS = p
	}

	{
		p, err := gpio.Output(defPinDC, false, false)
		if err != nil {
			panic("failed to setup DC pin: " + err.Error())
		}
		pinDC = p
	}

	{
		p, err := gpio.Output(defPinRST, false, false)
		if err != nil {
			panic("failed to setup RST pin: " + err.Error())
		}
		pinRST = p
	}

	e := &EPD{
		busy:   pinBUSY,
		cs:     pinCS,
		dc:     pinDC,
		lut:    lutFullUpdate,
		reset:  pinRST,
		spi:    d,
		width:  Width,
		height: Height,
	}

	e.Reinitialize()
	return e
}

func (e *EPD) Reinitialize() {
	log.Printf("Reinitialize")
	e.SendCommand(DriverOutputControl, (e.height-1)&0xff, ((e.height-1)>>8)&0xff, 0x00)
	e.SendCommand(BoosterSoftStartControl, 0xd7, 0xd6, 0x9d)
	e.SendCommand(WriteVCOMRegister, 0xA8)
	e.SendCommand(SetDummyLinePeriod, 0x1A)
	e.SendCommand(SetGateTime, 0x08)
	e.SendCommand(DataEntryModeSetting, 0x03)
	e.SetLUT(e.lut)
}

func (e *EPD) Reset() {
	log.Printf("Reset")
	e.reset.Write(false)
	time.Sleep(200 * time.Millisecond)
	e.reset.Write(true)
	time.Sleep(200 * time.Millisecond)
}

func (e *EPD) SendCommand(cmd Command, args ...byte) {
	e.dc.Write(false)

	var buf bytes.Buffer

	fmt.Fprintf(&buf, "CMD: ")
	scmd := cmd.String()
	if len(scmd) > 15 {
		fmt.Fprintf(&buf, "%-12s...", scmd[:12])
	} else {
		fmt.Fprintf(&buf, "%-15s", scmd)
	}
	fmt.Fprintf(&buf, "(0x%02X)", cmd.Byte())
	if len(args) > 0 {
		fmt.Fprintf(&buf, " ARGS: [")
		for i, arg := range args {
			fmt.Fprintf(&buf, "0x%02X", arg)
			if i < len(args)-1 {
				fmt.Fprintf(&buf, ", ")
			}
		}
		fmt.Fprintf(&buf, "]")
	}
	log.Println(buf.String())
	e.spi.Transfer([]byte{cmd.Byte()})
	e.SendData(args...)
}

// SendData sends data through SPI. Arbitrary number of bytes can be
// passed to this method, and they will be each sent in succession
func (e *EPD) SendData(data ...byte) {
	for _, b := range data {
		e.dc.Write(true)
		e.spi.Transfer([]byte{b})
	}
}

func (e *EPD) SetLUT(lut []byte) {
	e.lut = lut
	e.SendCommand(WriteLUTRegister, lut...)
}

func (e *EPD) SetMemoryArea(startX, startY, endX, endY uint8) {
	log.Printf("SetMemoryArea")
	// x point must be multiple of 8 or the last 3 bits will be ignored
	e.SendCommand(SetRamXAddressStartEndPosition, (startX>>3)&0xFF, (endX>>3)&0xFF)
	e.SendCommand(SetRamYAddressStartEndPosition, startY&0xFF, (startY>>8)&0xFF, endY&0xFF, (endY>>8)&0xFF)
}

func (e *EPD) SetMemoryPointer(x, y uint8) {
	e.SendCommand(SetRamXAddressCounter, (x>>3)&0xFF)
	e.SendCommand(SetRamYAddressCounter, y&0xFF, (y>>8)&0xFF)
	e.WaitUntilIdle(nil)
}

func (e *EPD) WaitUntilIdle(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()

	for {
		b, err := e.busy.Read()
		if err == nil && !b {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
		}
	}
	return errors.New(`not reached`)
}

func (e *EPD) ClearFrameMemory(color byte) {
	e.SetMemoryArea(0, 0, e.width-1, e.height-1)
	e.SetMemoryPointer(0, 0)
	log.Printf("Start writing to RAM")

	// XXX this is not optimal, but makes it easier for debugging
	args := make([]byte, e.width/8*e.height)
	for i := range args {
		args[i] = color
	}
	e.SendCommand(WriteRAM, args...)
}

func (e *EPD) DisplayFrame() {
	e.SendCommand(DisplayUpdateControl2, 0xC4)
	e.SendCommand(MasterActivation)
	e.SendCommand(TerminateFrameReadWrite)
	e.WaitUntilIdle(nil)
}

package pcr2

import (
	"fmt"
)

type Device struct {
	transp Transport
}

func NewDevice(transp Transport) *Device {
	return &Device{
		transp: transp,
	}
}

func (d *Device) Get(name string) (string, error) {
	out, err := d.transp.Write(fmt.Sprintf("get %s\r\n", name))
	if err != nil {
		return "", err
	}
	switch name {
	case "temp":
		i := len(out) - 1
		temp := out[:i] + "." + out[i:] + " °C"
		return temp, nil
	case "mode":
		return modeIntToString(out), nil
	default:
		return out, nil
	}
}

func (d *Device) Set(name, param string) (string, error) {
	return d.transp.Write(fmt.Sprintf("set %s %s\r\n", name, param))
}

func (d *Device) LoraGet(name string) (string, error) {
	return d.transp.Write(fmt.Sprintf("lora get %s\r\n", name))
}

func (d *Device) LoraSet(name, param string) (string, error) {
	return d.transp.Write(fmt.Sprintf("lora set %s %s\r\n", name, param))
}

func (d *Device) Clear() error {
	_, err := d.transp.Write("clear\r\n")
	return err
}

func modeIntToString(in string) string {
	switch in {
	case "0":
		return "[0] Timespan, used to sum up detection and send after a certain time (Sending Interval)"
	case "1":
		return "[1] NotZero, Same as Timespan but w/o sending if counters are 0 (zero)"
	case "2":
		return "[2] Trigger, Send on every events. Events can be filtered with Hold Off setting"
	case "3":
		return "[3] Capacity Alert Mode"
	default:
		return "Unknown mode"
	}
}

type Transport interface {
	Write(in string) (string, error)
}

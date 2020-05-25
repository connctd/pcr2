package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/connctd/pcr2"
	"github.com/urfave/cli/v2"
)

var (
	device *pcr2.Device
)

func main() {
	app := &cli.App{
		Name:  "pcr2",
		Usage: "Configure a PCR2 people counter",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "port",
				Usage: "Serial port",
			},
		},
		Before: openPcr2,
		Action: func(c *cli.Context) error {
			out, err := device.Get("typestr")
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", out)
			return nil
		},

		Commands: []*cli.Command{
			&cli.Command{
				Name: "get",
				Action: func(c *cli.Context) error {
					if c.Args().Len() != 1 {
						return errors.New("get takes only exactly one argument")
					}
					out, err := device.Get(c.Args().First())
					if err != nil {
						return err
					}
					fmt.Printf("%s\n", out)
					return nil
				},
			},

			&cli.Command{
				Name: "set",
				Action: func(c *cli.Context) error {
					if c.Args().Len() < 2 {
						return errors.New("set requires at least two arguments")
					}
					name := c.Args().First()
					param := ""
					for i := 1; i < c.Args().Len(); i++ {
						param += c.Args().Get(i) + " "
					}
					param = strings.TrimRight(param, " ")
					out, err := device.Set(name, param)
					if err != nil {
						return err
					}
					fmt.Printf("%s\n", out)
					return nil
				},
			},

			&cli.Command{
				Name: "clear",
				Action: func(c *cli.Context) error {
					if c.Args().Len() != 0 {
						return errors.New("clear takes no arguments")
					}
					err := device.Clear()
					if err == nil {
						fmt.Printf("Device counters cleared\n")
					}
					return err
				},
			},

			&cli.Command{
				Name: "lora",
				Subcommands: []*cli.Command{
					&cli.Command{
						Name: "get",
						Action: func(c *cli.Context) error {
							if c.Args().Len() != 1 {
								return errors.New("The LoRa parameter name to get needs to be specified")
							}
							name := c.Args().First()
							out, err := device.LoraGet(name)
							if err == nil {
								fmt.Sprintf("%s\n", out)
							}
							return err
						},
					},

					&cli.Command{
						Name: "set",
						Action: func(c *cli.Context) error {
							if c.Args().Len() < 2 {
								return errors.New("The LoRa parameter name to get needs to be specified")
							}
							name := c.Args().First()
							param := ""
							for i := 1; i < c.Args().Len(); i++ {
								param += c.Args().Get(i) + " "
							}
							param = strings.TrimRight(param, " ")
							out, err := device.LoraSet(name, param)
							if err == nil {
								fmt.Sprintf("%s\n", out)
							}
							return err
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func openPcr2(ctx *cli.Context) error {
	serialPort := ctx.String("port")
	if serialPort == "" {
		serialPort = pcr2.DefaultSerialPort
	}
	if serialPort == "" {
		return errors.New("Please specify the serial port to use")
	}
	transp, err := pcr2.Open(serialPort)
	if err != nil {
		return fmt.Errorf("Failed to open serial port transport: %w", err)
	}
	device = pcr2.NewDevice(transp)
	return nil
}

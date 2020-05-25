# PCR2 cli config tool

The [PCR2-IN](https://www.parametric.ch/products/pcr2-in/) is a useful sensor to determine whether persons
or objects move through in a certain direction. Unfortunately the company behind this sensor only delivers
a configuration tool for Microsoft Windows. Since I don't have any Windows machines at home and COVID-19
prevented me from visiting colleagues with Windows machines, I decided to write a small cli configuration tool
which should be fairly cross platform and should also be usable for mass provisioning of these devices.

## Usage Examples

This cli mostly follows the serial port cli described [here](https://www.parametric.ch/docs/pcr2/pcr2_cli_v32x), but adds
some convience over using something like miniterm

* Read type string `pcr2 get typestr`
* Read device temperature`pcr2 get temp`
* Set radar sensitivity `pcr2 set sens 50`
* Set LoRa AppKey `pcr2 lora set appkey <appkey>`

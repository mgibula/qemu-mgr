package qemu

type vendorEntry struct {
	number uint
	name   string
}

type deviceEntry struct {
	number uint
	name   string
}

var vendors = []vendorEntry{
	{0x8086, "Intel"},
}

var devices = []deviceEntry{
	{0x0680, "PCI Bridge"},
	{0x1237, "Intel 82441"},
	{0x2918, "Intel ICH9"},
	{0x2922, "Intel ICH9 (AHCI mode)"},
	{0x2930, "Intel ICH9 SMBus Controller"},
	{0x29c0, "Intel Q35 MCH"},
	{0x7000, "Intel 82371SB"},
	{0x7010, "Intel 82371SB IDE Controller"},
	{0x7020, "Intel 82371SB"},
	{0x7110, "Intel 82371AB"},
}

func GetVendorName(vendor uint) string {
	for _, entry := range vendors {
		if entry.number == vendor {
			return entry.name
		}
	}

	return ""
}

func GetDeviceName(device uint) string {
	for _, entry := range devices {
		if entry.number == device {
			return entry.name
		}
	}

	return ""
}

package mamedattools

type Machine struct {
	Name         string      `xml:"name,attr"`
	SourceFile   string      `xml:"sourcefile,attr,omitempty"`
	IsBios       string      `xml:"isbios,attr,omitempty"`
	CloneOf      string      `xml:"cloneof,attr,omitempty"`
	RomOf        string      `xml:"romof,attr,omitempty"`
	SampleOf     string      `xml:"sampleof,attr,omitempty"`
	Description  string      `xml:"description,omitempty"`
	Year         string      `xml:"year,omitempty"`
	Manufacturer string      `xml:"manufacturer,omitempty"`
	Roms         []Rom       `xml:"rom,omitempty"`
	Disks        []Disk      `xml:"disk,omitempty"`
	Samples      []Sample    `xml:"sample,omitempty"`
	DeviceRef    []DeviceRef `xml:"device_ref,omitempty"`
	Driver       Driver      `xml:"driver,omitempty"`
}

type DeviceRef struct {
	Name string `xml:"name,attr"`
}

type Driver struct {
	Status string `xml:"status,attr,omitempty"`
}

func MachineToGame(machine Machine) Game {
	return Game{
		Name:         machine.Name,
		SourceFile:   machine.SourceFile,
		IsBios:       machine.IsBios,
		CloneOf:      machine.CloneOf,
		RomOf:        machine.RomOf,
		SampleOf:     machine.SampleOf,
		Description:  machine.Description,
		Year:         machine.Year,
		Manufacturer: machine.Manufacturer,
		Roms:         machine.Roms,
		Disks:        machine.Disks,
		Samples:      machine.Samples,
		Driver:       machine.Driver,
	}
}

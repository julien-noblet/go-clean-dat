package mamedattools

import (
	"encoding/xml"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/viper"
)

type Datafile struct {
	Header  Header    `xml:"header"`
	Games   []Game    `xml:"game,omitempty"`
	Machine []Machine `xml:"machine,omitempty"`
}

type Header struct {
	Category string `xml:"category"`
}

type Game struct {
	Name         string   `xml:"name,attr"`
	SourceFile   string   `xml:"sourcefile,attr,omitempty"`
	IsBios       string   `xml:"isbios,attr,omitempty"`
	CloneOf      string   `xml:"cloneof,attr,omitempty"`
	RomOf        string   `xml:"romof,attr,omitempty"`
	SampleOf     string   `xml:"sampleof,attr,omitempty"`
	Description  string   `xml:"description,omitempty"`
	Year         string   `xml:"year,omitempty"`
	Manufacturer string   `xml:"manufacturer,omitempty"`
	Roms         []Rom    `xml:"rom,omitempty"`
	Disks        []Disk   `xml:"disk,omitempty"`
	Samples      []Sample `xml:"sample,omitempty"`
	Driver       Driver   `xml:"driver,omitempty"`
}

type Rom struct {
	Name   string `xml:"name,attr"`
	Size   string `xml:"size,attr,omitempty"`
	Crc    string `xml:"crc,attr,omitempty"`
	Sha1   string `xml:"sha1,attr,omitempty"`
	Merge  string `xml:"merge,attr,omitempty"`
	Status string `xml:"status,attr,omitempty"`
}

type Disk struct {
	Name   string `xml:"name,attr"`
	Sha1   string `xml:"sha1,attr,omitempty"`
	Merge  string `xml:"merge,attr,omitempty"`
	Status string `xml:"status,attr,omitempty"`
}

type Sample struct {
	Name string `xml:"name,attr"`
}

func OpenAndUnmarshal(filename string) Datafile {
	mame := Datafile{}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Unmarshal le fichier xml dans la struct Mame
	err = xml.NewDecoder(file).Decode(&mame)
	if err != nil {
		log.Fatal(err)
	}
	return mame
}

func CleanMame(mame *Datafile) Datafile {
	newMame := Datafile{}
	bioses := Datafile{}

	if len(mame.Machine) != 0 {
		// convert machines to games
		for _, machine := range mame.Machine {
			if viper.GetBool("no-preliminary") && machine.Driver.Status == "preliminary" {
				continue
			}
			if viper.GetBool("no-imperfect") && machine.Driver.Status == "imperfect" || machine.Driver.Status == "preliminary" {
				continue
			}
			newMame.Games = append(newMame.Games, MachineToGame(machine))
		}
		mame.Games = newMame.Games
		mame.Machine = []Machine{}
	}
	// liste les bios
	GetBios(mame, &bioses)
	GetBios(mame, &newMame)

	AddGameWithoutBios(mame, &bioses, &newMame)
	return newMame
}

func InB(game *Game, mameB *Datafile, mamebname string) bool {
	if game.IsBios == "yes" { // on ne compare pas les bios
		return false
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	found := false

	for i, gameB := range mameB.Games {
		wg.Add(1)
		go func(i int, gameB Game) {
			defer wg.Done()
			if IsSameGame(game, &gameB) {
				mutex.Lock()
				removeFormList(mameB.Games, i)
				mutex.Unlock()
				log.Printf("Game %s is in %s as %s", game.Name, mamebname, gameB.Name)
				found = true
			}
		}(i, gameB)
	}

	wg.Wait()
	return found
}

func NotInB(mameA, mameB *Datafile, mamebname string) Datafile {

	// display progress bar
	bar := progressbar.Default(int64(len(mameA.Games)))

	newMame := Datafile{}
	for _, gameA := range mameA.Games {
		if bar.Add(1) != nil {
			log.Fatal("Error with progress bar")
		}
		if gameA.IsBios == "yes" {
			newMame.Games = append(newMame.Games, gameA)
			continue
		}
		if !InB(&gameA, mameB, mamebname) {
			newMame.Games = append(newMame.Games, gameA)
		}
	}
	return newMame

}

func CheckIfRomExists(rom *Rom, bioses *Datafile) bool {
	// check par crc
	for _, bios := range bioses.Games {
		for _, biosRom := range bios.Roms {
			if rom.Crc == biosRom.Crc && rom.Crc != "" || rom.Sha1 == biosRom.Sha1 && rom.Sha1 != "" {
				return true
			}
		}
	}
	return false
}

func GameHAveToBeRemoved(element string, list *[]string) bool {
	for _, e := range *list {
		if element == e {
			return true
		}
	}
	return false
}

func AddGameWithoutBios(mame, bioses, newMame *Datafile) {
	for _, game := range mame.Games {
		// remove Pinball games
		// sourcefile start with "pinball/"
		if viper.GetBool("no-pinball") {
			if len(game.SourceFile) > 8 {
				if game.SourceFile != "" && game.SourceFile[:8] == "pinball/" {
					continue
				}
			}
		}
		// remove devices (sourcefile start with "devices/")
		if viper.GetBool("no-devices") {
			if len(game.SourceFile) > 9 {
				if game.SourceFile != "" && game.SourceFile[:9] == "devices/" {
					continue
				}
			}
		}
		// remove mechanical games
		// sourcefile start with "mechanical/"
		if viper.GetBool("no-mechanical") {
			if len(game.SourceFile) > 12 {
				if game.SourceFile != "" && game.SourceFile[:12] == "mechanical/" {
					continue
				}
			}
		}

		// remove clones
		if viper.GetBool("no-clones") && game.CloneOf != "" {
			continue
		}
		ns := viper.GetStringSlice("no-system")
		if GameHAveToBeRemoved(game.RomOf, &ns) {
			continue
		}
		sf := viper.GetStringSlice("no-sourcefile")
		if GameHAveToBeRemoved(game.SourceFile, &sf) {
			continue
		}
		if game.IsBios != "yes" {
			// si le jeu contient des roms ou des disques
			if len(game.Roms) != 0 || len(game.Disks) != 0 {
				newGame := game
				newGame.Roms = []Rom{}
				// on check si les roms existent dans les bios
				for _, rom := range game.Roms {
					if rom.Status == "nodump" {
						//log.Printf("Rom %s is nodump in %s game", rom.Name, game.Name)
						continue
					}
					if !CheckIfRomExists(&rom, bioses) {
						// si elles n'existent pas on ajoute le jeu à la nouvelle liste
						newGame.Roms = append(newGame.Roms, rom)
						continue
					}
				}
				newMame.Games = append(newMame.Games, newGame)
			}
		}
	}
}

func IsSameGame(game1, game2 *Game) bool {
	if game1.IsBios == "yes" || game2.IsBios == "yes" {
		return false
	}

	// try to find 1 disk in game1 not in game2 or 1 disk in game2 not in game1
	// create an with length = len(game1.Disks) + len(game1.Roms) and fill it with false
	disks1 := len(game1.Disks)
	disks2 := len(game2.Disks)
	if disks1 != disks2 {
		return false
	}

	game1files := make([]bool, disks1+len(game1.Roms))
	game2files := make([]bool, disks2+len(game2.Roms))

	var wg sync.WaitGroup

	nvramRegex := regexp.MustCompile(`.*\.nv(?:[0-9])?$`)
	for i, disk1 := range game1.Disks {
		wg.Add(1)
		go func(i int, disk1 Disk) {
			defer wg.Done()
			for j, disk2 := range game2.Disks {
				if disk1.Sha1 == disk2.Sha1 {
					game1files[i] = true
					game2files[j] = true
					break
				}
			}
		}(i, disk1)
	}
	wg.Wait()

	for _, rom1 := range game1.Roms {
		if nvramRegex.MatchString(rom1.Name) || rom1.Status == "nodump" {
			// remove 1 to the length of the array
			game1files = game1files[:len(game1files)-1]
		}
	}

	for _, rom2 := range game2.Roms {
		if nvramRegex.MatchString(rom2.Name) || rom2.Status == "nodump" {
			// remove 1 to the length of the array
			game2files = game2files[:len(game2files)-1]
		}
	}

	if len(game1files) != len(game2files) {
		return false
	}

	i := disks1
	for _, rom1 := range game1.Roms {
		if nvramRegex.MatchString(rom1.Name) {
			continue
		}
		//		wg.Add(1)
		//		go func(i int, rom1 Rom) {
		//			defer wg.Done()
		j := disks2
		for _, rom2 := range game2.Roms {
			if nvramRegex.MatchString(rom2.Name) {
				continue
			}
			if rom1.Crc == rom2.Crc && rom1.Crc != "" || rom1.Sha1 == rom2.Sha1 && rom1.Sha1 != "" {
				game1files[i] = true
				game2files[j] = true
				break
			}
			j++
		}
		//		}(i, rom1)
		//		i++
		//		wg.Wait()
	}

	// if game1files == game2files then the games are the same
	if len(game1files) != len(game2files) {
		return false
	}
	for i, file := range game1files {
		if file != game2files[i] {
			return false
		}
	}
	for _, f := range game1files {
		if !f {
			return false
		}
	}
	for _, f := range game2files {
		if !f {
			return false
		}
	}
	return true

}

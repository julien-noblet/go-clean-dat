package mamedattools

import "log"

func GetBios(mame *Datafile, bioses *Datafile) {
	for _, game := range mame.Games {
		if game.IsBios == "yes" {
			// si le bios contient des roms ou des disques
			if len(game.Roms) != 0 || len(game.Disks) != 0 {
				log.Printf("Bios: %s", game.Name)
				// et on l'ajoute à la nouvelle liste				// et a la liste des bios
				bioses.Games = append(bioses.Games, game)
			}
		}
	}
}

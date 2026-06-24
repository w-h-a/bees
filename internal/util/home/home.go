package home

import (
	"fmt"
	"path/filepath"
)

// Resolved is the bees home and db file on this device.
type Resolved struct {
	Home string
	DB   string
}

// Resolve derives the bees home and the db file under it from the two relevant
// env vars: beesHome is the BEES_HOME override and userHome is the user's OS home
// directory.
func Resolve(beesHome, userHome string) (Resolved, error) {
	if beesHome == "" && userHome == "" {
		return Resolved{}, fmt.Errorf("refused to resolve bees home: BEES_HOME and user home are both empty")
	}

	dir := beesHome
	if dir == "" {
		dir = filepath.Join(userHome, ".bees")
	}

	return Resolved{
		Home: dir,
		DB:   filepath.Join(dir, "bees.db"),
	}, nil
}

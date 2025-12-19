package directory

import (
	"os"
	"path/filepath"
)

func GetCachePath() string { 
	base, _ := os.UserConfigDir();
	base = filepath.Join(base, "space.stelare.pyxis");
	return base
}

func GetDataPath() string {
	base, _ := os.UserHomeDir();
	base = filepath.Join(base, ".local", "share", "Pyxis");
	return base;
}



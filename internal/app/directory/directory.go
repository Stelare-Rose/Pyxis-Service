package directory

import (
	"os"
	"path/filepath"

	"github.com/Stelare-Rose/Pyxis-Service/internal/app/environment"
)

func GetCachePath() string { 
	base, _ := os.UserConfigDir();
	base = filepath.Join(base, "space.stelare.pyxis");
	return base
}

func GetDataPath() string {
	base, _ := os.UserHomeDir();
	if environment.GetEnvironment() == "dev" {
		base = filepath.Join(base, ".local", "share", "Pyxis-dev");
	} else {
		base = filepath.Join(base, ".local", "share", "Pyxis");
	}
	return base;
}



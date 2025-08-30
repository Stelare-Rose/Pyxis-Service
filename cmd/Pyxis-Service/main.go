package main

import (
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/core"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/database"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/filedetection"
)

func main(){
	end := make(chan bool);
	database.Create();
	core.CleanIndexActiveItems();
	core.IndexTags(nil);
	go filedetection.Start();
	<- end;
}


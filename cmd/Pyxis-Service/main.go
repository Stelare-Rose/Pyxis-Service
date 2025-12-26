package main

import (
	"fmt"

	"github.com/Stelare-Rose/Pyxis-Service/internal/app/core"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/database"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/directory"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/environment"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/filedetection"
)

func main(){
	end := make(chan bool);
	fmt.Println("Currently Running on Environment " + environment.GetEnvironment())
	directory.Check();
	database.Create();
	core.CleanIndexActiveItems();
	core.CleanIndexActiveIdeas();
	core.IndexTags(nil);
	go filedetection.Start();
	<- end;
}


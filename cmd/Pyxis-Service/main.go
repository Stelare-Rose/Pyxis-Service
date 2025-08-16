package main

import (
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/directory"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/filedetection"
)

func main(){
	end := make(chan bool);

	go filedetection.Start();
	go directory.IndexActiveItems();
	
	<- end;
}


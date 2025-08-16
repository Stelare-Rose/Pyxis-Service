package filedetection

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Stelare-Rose/Pyxis-Service/internal/app/directory"
	"github.com/fsnotify/fsnotify"
)

func Start(){
	batch := make(chan struct{});
	end := make(chan bool);
	var events []fsnotify.Event;

	watcher, err := fsnotify.NewWatcher();
	if err != nil {
		log.Fatal(err);
	}
	defer watcher.Close();

	// Debounce
	go func() {
		for {
			<-batch
			fmt.Println("Batched Index.. Debouncing");
			time.Sleep(5 * time.Second);
			
			directory.IndexActiveItems();
			e := events;
			events = nil;

			if e != nil {
				fmt.Println("List of Events:");
				for _, i := range e {
					fmt.Println(i);
				}
			}
		}
	}()

	go func() {
        for {
            select {
            case event, ok := <-watcher.Events:
                if !ok {
                    return
                }
				events = append(events, event);
				select {
					case batch <- struct{}{}:
						fmt.Println("Sent Signal!")
					default:
				}
            case err, ok := <-watcher.Errors:
                if !ok {
                    return
                }
                log.Println("error:", err)
            }
        }
    }()

	home, _ := os.UserHomeDir();
	path := filepath.Join(home, ".local", "share", "Pyxis", "Items", "Active")
	err = watcher.Add(path);
	<- end;
}


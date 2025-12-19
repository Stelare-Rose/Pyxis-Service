package filedetection

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/Stelare-Rose/Pyxis-Service/internal/app/core"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/database"
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

	// Close Watcher if Program is Done
	defer watcher.Close();

	// Debounce
	go func() {
		for {
			<-batch
			fmt.Println("Batched Index.. Debouncing");
			time.Sleep(50 * time.Millisecond);
			start := time.Now();
			e := events;
			events = nil;

			files := make(map[string]struct{});
			if e != nil {
				// I think I was drunk when I wrote this code, I'll refactor it some day
				// TODO: Refactor
				fmt.Println("List of Events:");

				// Process Filesystem Events into Map of Files
				for _, i := range e {
					name := filepath.Clean(i.Name);
					if strings.Contains(name, "/.stfolder") {
						continue;
					}
					if strings.Contains(name, ".syncthing"){
						continue;
					}
					if i.Name == filepath.Join(directory.GetDataPath(), "tags.toml") {
						files["tags.toml"] = struct{}{};
						continue;
					}
					parts := strings.Split(name, string(filepath.Separator));

					if len(parts) >= 2 {
						files[filepath.Join(parts[len(parts) - 2], parts[len(parts) - 1])] = struct {}{};
					} else {
						files[name] = struct{}{};
					}

					fmt.Println(i.String());
				}
				
				// Indexing
				tx := database.StartTransaction();
				fmt.Println(files);
				for key := range files {
					if key == "tags.toml" {
						core.IndexTags(tx);
						continue;
					}
					if strings.HasPrefix(key, "Active/") && strings.HasSuffix(key, ".task"){
						core.IndexActiveItem(key, tx);
						continue;
					}
				}
				database.EndTransaction(tx);
			}
			fmt.Println("Small Indexing completed in", time.Since(start));
		}
	}()

	// Start Watcher in Goroutine
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

	// Defile Paths, Add Paths to Watcher
	path := filepath.Join(directory.GetDataPath());
	err = watcher.Add(path);
	path = filepath.Join(directory.GetDataPath(), "Items", "Active");
	err = watcher.Add(path);
	if err != nil {
		fmt.Println(err);
	}
	<- end;
}


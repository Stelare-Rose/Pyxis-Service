package directory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/file"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/types"
)

func IndexActiveItems(){
	start := time.Now();
	//TODO: Fix For Windows
	itemChanged := false;
	home, err := os.UserHomeDir();
	if err != nil {
		fmt.Println(err);
		return;
	}
	path := filepath.Join(home, ".local", "share", "Pyxis", "Items", "Active")

	baseCache, err := os.UserCacheDir();
	if err != nil {
		fmt.Println(err);
		return;
	}
	cache := filepath.Join(baseCache, "Pyxis", "index-ongoing", "items.json");
	data, _ := os.ReadFile(cache);

	var items types.Items;
	json.Unmarshal(data, &items);

	entries, err := os.ReadDir(path);
	if err != nil {
		fmt.Println(err);
		return;
	}
	for i := 0; i < len(items.Item); i++ {
		file, _ := os.Stat(filepath.Join(path, items.Item[i].Path));	
		if file == nil {
			itemChanged = true;
			fmt.Printf("Removed %v\n", items.Item[i].Name);
			items.Item = remove(i, items.Item);
			i--;
		}
	}
	for i := 0; i < len(items.Item); i++ {
		for j := i + 1; j < len(items.Item); j++ {
			if items.Item[i].Id == items.Item[j].Id {
				itemChanged = true;
				fmt.Printf("Removed %v\n", items.Item[i].Name);
				items.Item = remove(j, items.Item);
				j--;
			}
		}	
	}
	
	for _, e := range entries {
		f := file.Fingerprint(filepath.Join(path, e.Name()));
		index := getItemInIndex(items, e.Name());
		if index == -1 {
			itemChanged = true;
			item := file.ReadActiveFile(filepath.Join(path, e.Name())); 
			items.Item = append(items.Item, item);
			items.Item[len(items.Item) - 1].Fingerprint = f;
			items.Item[len(items.Item) - 1].Path = e.Name();
		} else if f != items.Item[index].Fingerprint {
			itemChanged = true;
			item := file.ReadActiveFile(filepath.Join(path, e.Name()));
			items.Item[index] = item;
			items.Item[index].Fingerprint = f;
			items.Item[index].Path = e.Name();
		} else {
			fmt.Printf("File %v was not changed.\n", e.Name());
		}
	}

	if itemChanged {
		os.WriteFile(filepath.Join(baseCache, "Pyxis", "index-ongoing", "index.lock"), []byte{}, 0666);
		jsonData, _ := json.Marshal(items);
		os.WriteFile(cache, jsonData, 0666);
		os.Remove(filepath.Join(baseCache, "Pyxis", "index-ongoing", "index.lock"));
	}
	fmt.Println("Scan and Indexing completed in", time.Since(start));
}

func getItemInIndex(items types.Items, name string) int {
	for i, data := range items.Item {
		if data.Path == name {
			fmt.Println("Id found match at " + data.Name);
			return i;
		}
	}
	return -1;
}

func remove(index int, arr []types.Item) []types.Item {
	arr[index] = arr[len(arr) - 1];
	return arr[:len(arr) - 1];
}

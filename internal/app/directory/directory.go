package directory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/hashing"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/file"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/types"
)

func IndexActiveItems(){
	//TODO: Fix For Windows
	itemChanged := false;
	home, err := os.UserHomeDir();
	if err != nil {
		fmt.Println(err);
		return;
	}
	path := filepath.Join(home, ".local", "share", "Pyxis", "Items", "Active")

	cache, err := os.UserCacheDir();
	if err != nil {
		fmt.Println(err);
		return;
	}
	cache = filepath.Join(cache, "Pyxis", "index-ongoing", "items.json");
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
			fmt.Printf("Removed %v\n", items.Item[i].Name);
			items.Item = remove(i, items.Item);
			i--;
		}
	}
	for i := 0; i < len(items.Item); i++ {
		for j := i + 1; j < len(items.Item); j++ {
			if items.Item[i].Id == items.Item[j].Id {
			fmt.Printf("Removed %v\n", items.Item[i].Name);
				items.Item = remove(j, items.Item);
				j--;
			}
		}	
	}
	
	for _, e := range entries {
		h := hashing.MetadataHash(filepath.Join(path, e.Name()));
		index := getItemInIndex(items, e.Name());
		if index == -1 {
			itemChanged = true;
			item := file.ReadActiveFile(filepath.Join(path, e.Name())); 
			items.Item = append(items.Item, item);
			items.Item[len(items.Item) - 1].Hash = h;
			items.Item[len(items.Item) - 1].Path = e.Name();
		} else if h != items.Item[index].Hash {
			itemChanged = true;
			item := file.ReadActiveFile(filepath.Join(path, e.Name()));
			items.Item[index] = item;
			items.Item[index].Hash = h;
			items.Item[index].Path = e.Name();
		} else {
			fmt.Printf("File %v was not changed.\n", e.Name());
		}
	}

	if itemChanged {
		jsonData, _ := json.Marshal(items);
		os.WriteFile(cache, jsonData, 0666);
	}
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

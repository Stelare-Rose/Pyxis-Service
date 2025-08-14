package directory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/hashing"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/file"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/types"
)

func IndexActiveItems(){
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
	for _, e := range entries {
		h := hashing.MetadataHash(filepath.Join(path, e.Name()));
		id := strings.Split(e.Name(), "§");
		id = strings.Split(id[1], ".");
		index := getItemInIndex(items, id[0]);
		if index == -1 {
			itemChanged = true;
			item := file.ReadActiveFile(filepath.Join(path, e.Name())); 
			items.Item = append(items.Item, item);
			items.Item[len(items.Item) - 1].Hash = h;

		} else if h != items.Item[index].Hash {
			itemChanged = true;
			item := file.ReadActiveFile(filepath.Join(path, e.Name()));
			items.Item[index] = item;
			items.Item[index].Hash = h;
		} else {
			fmt.Printf("Id %v was not changed.\n", id[0]);
		}
	}

	if itemChanged {
		jsonData, _ := json.Marshal(items);
		os.WriteFile(cache, jsonData, 0666);
	}
}

func getItemInIndex(items types.Items, id string) int {
	for i, data := range items.Item {
		if data.Id == id {
			fmt.Println("Id found match at " + data.Name);
			return i;
		}
	}
	return -1;
}

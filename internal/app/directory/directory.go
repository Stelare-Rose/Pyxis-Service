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

func IndexActiveItems(){
	start := time.Now();
	//TODO: Fix For Windows
	itemChanged := false;
	path := filepath.Join(GetDataPath(), "Items", "Active")
	cache := filepath.Join(GetCachePath(), "index-ongoing", "items.json");
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

	compare := make(map[string]bool);
	for i := 0; i < len(items.Item); i++ {
		_, exists := compare[items.Item[i].Id];
		if exists { 
			fmt.Printf("Removed %v\n", items.Item[i].Name);
			items.Item = remove(i, items.Item);
			i--;
		} else {
			compare[items.Item[i].Id] = true;
		}
	}
	
	lookup := make(map[string]int)
	for i := range items.Item {
		lookup[items.Item[i].Path] = i
	}

	for _, e := range entries {
		f := file.Fingerprint(filepath.Join(path, e.Name()));
		index, exists := lookup[e.Name()];
		if !exists {
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
		os.WriteFile(filepath.Join(GetCachePath(), "index-ongoing", "index.lock"), []byte{}, 0666);
		jsonData, _ := json.Marshal(items);
		os.WriteFile(cache, jsonData, 0666);
		os.Remove(filepath.Join(GetCachePath(), "index-ongoing", "index.lock"));
	}
	fmt.Println("Scan and Indexing completed in", time.Since(start));
}

func remove(index int, arr []types.Item) []types.Item {
	arr[index] = arr[len(arr) - 1];
	return arr[:len(arr) - 1];
}

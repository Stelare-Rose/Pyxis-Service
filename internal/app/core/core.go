package core

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/database"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/directory"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/file"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/types"
)

func CleanIndexActiveItems() {
	start := time.Now();
	
	path := filepath.Join(directory.GetDataPath(), "Items", "Active");
	entries, _:= os.ReadDir(path);

	database.ResetActiveItems()
	
	names := make([]string, 0);
	for _, e := range entries {
		name := filepath.Join("Active", e.Name());
		names = append(names, name);	
	}
	database.VerifyItems(names);
	database.DropUnverifiedItems();

	tx := database.StartTransaction()
	for _, e := range entries {
		fmt.Println(e.Name());
		f, _ := database.QueryItemFingerprintByPath(filepath.Join("Active", e.Name()));
		fmt.Println(f);
		if f == 0 {
			f = file.Fingerprint(filepath.Join(path, e.Name()));
			item := file.ReadActiveFile(filepath.Join(path, e.Name()));
			database.RemoveTagWithTransaction(item.Id, tx);
			for _, tag := range item.Tags {
				database.AddTagWithTransaction(item.Id, tag, tx);
			}
			item.Fingerprint = f;
			item.Path = filepath.Join("Active", e.Name())
			fmt.Println(item);
			database.AddItemWithTransaction(&item, false, tx);
		}
	}
	database.EndTransaction(tx);

	fmt.Println("Scan and Indexing completed in", time.Since(start));
}

func IndexActiveItem(shortPath string, tx *sql.Tx){
	start := time.Now();
	path := filepath.Join(directory.GetDataPath(), "Items", shortPath);
	_, err := os.Stat(path);
	if err != nil {
		fmt.Println("File Missing!");
		database.RemoveItemByPathWithTransaction(shortPath, tx);
		return;
	}
	f, _ := database.QueryItemFingerprintByPath(path);
	nf := file.Fingerprint(path);
	if nf != f {
		item := file.ReadActiveFile(path);
		database.RemoveTagWithTransaction(item.Id, tx);
		fmt.Println(item.Id);
		for _, tag := range item.Tags {
			database.AddTagWithTransaction(item.Id, tag, tx);
		}
		item.Fingerprint = nf;
		item.Path = shortPath;
		fmt.Println(item);
		database.AddItemWithTransaction(&item, false, tx);
	}
	fmt.Println("Single Index completed in", time.Since(start));
}

func IndexTags(tx *sql.Tx){
	endTx := false;
	if tx == nil {
		endTx = true;
		tx = database.StartTransaction();
	}
	database.ResetTags(tx);
	var tags types.Tags;
	_, err := toml.DecodeFile(filepath.Join(directory.GetDataPath(), "tags.toml"), &tags);
	fmt.Println(tags);
	if err != nil {
		fmt.Println(err);
	}

	for id, tag := range tags.Tags {
		database.AddTagsWithTransaction(id, tag.Name, tag.Colors, tx);
	}

	if endTx {
		database.EndTransaction(tx);
	}
}

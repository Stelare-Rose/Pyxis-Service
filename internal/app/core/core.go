package core

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
		if !strings.HasSuffix(e.Name(), ".task") {
			fmt.Printf("Detected Bad Filetype at %s\n", e.Name());
			continue;
		}
		f, _ := database.QueryItemFingerprintByPath(filepath.Join("Active", e.Name()));
		curr := file.Fingerprint(filepath.Join(path, e.Name()));
		if f == 0 || f != curr {
			f = curr;
			item := file.ReadActiveFile(filepath.Join(path, e.Name()));
			database.RemoveTagWithTransaction(item.Id, tx);
			for _, tag := range item.Tags {
				database.AddTagWithTransaction(item.Id, tag, tx);
			}
			item.Fingerprint = f;
			item.Path = filepath.Join("Active", e.Name())
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
		fmt.Printf("File Missing! Removing %s\n", path);
		database.RemoveItemByPathWithTransaction(shortPath, tx);
		return;
	}
	f, _ := database.QueryItemFingerprintByPath(path);
	nf := file.Fingerprint(path);
	if nf != f {
		item := file.ReadActiveFile(path);
		database.RemoveTagWithTransaction(item.Id, tx);
		for _, tag := range item.Tags {
			database.AddTagWithTransaction(item.Id, tag, tx);
		}
		item.Fingerprint = nf;
		item.Path = shortPath;
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
	if err != nil {
		fmt.Println(err);
	}

	for id, tag := range tags.Tags {
		database.AddTagsWithTransaction(id, tag.Tag, tag.Color, tx);
	}

	if endTx {
		database.EndTransaction(tx);
	}
}

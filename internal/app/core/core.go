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

// Full Clean Index
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
			item := file.ReadActiveItemFile(filepath.Join(path, e.Name()));
			database.RemoveItemTag(item.Id, tx);
			if len(item.Id) == 0 {
				break;
			}
			for _, tag := range item.Tags {
				database.AddItemTag(item.Id, tag, tx);
			}
			item.Fingerprint = f;
			item.Path = filepath.Join("Active", e.Name())
			database.AddItem(&item, false, tx);
		}
	}
	database.EndTransaction(tx);

	fmt.Println("Scan and Indexing completed in", time.Since(start));
}

func CleanIndexActiveIdeas() {
	start := time.Now();
	
	path := filepath.Join(directory.GetDataPath(), "Ideas", "Active");
	entries, _:= os.ReadDir(path);

	database.ResetActiveIdeas()
	
	names := make([]string, 0);
	for _, e := range entries {
		name := filepath.Join("Active", e.Name());
		names = append(names, name);	
	}
	database.VerifyIdeas(names);
	database.DropUnverifiedIdeas();

	tx := database.StartTransaction()
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".idea") {
			fmt.Printf("Detected Bad Filetype at %s\n", e.Name());
			continue;
		}
		f, _ := database.QueryIdeaFingerprintByPath(filepath.Join("Active", e.Name()));
		curr := file.Fingerprint(filepath.Join(path, e.Name()));
		if f == 0 || f != curr {
			f = curr;
			item := file.ReadActiveIdeaFile(filepath.Join(path, e.Name()));
			database.RemoveIdeaTag(item.Id, tx);
			if len(item.Id) == 0 {
				break;
			}
			for _, tag := range item.Tags {
				database.AddIdeaTag(item.Id, tag, tx);
			}
			item.Fingerprint = f;
			item.Path = filepath.Join("Active", e.Name())
			database.AddIdea(&item, false, tx);
		}
	}
	database.EndTransaction(tx);

	fmt.Println("Scan and Indexing completed in", time.Since(start));
}

// Single Item Index
func IndexActiveItem(shortPath string, tx *sql.Tx){
	start := time.Now();
	path := filepath.Join(directory.GetDataPath(), "Items", shortPath);
	_, err := os.Stat(path);
	if err != nil {
		fmt.Printf("File Missing! Removing %s\n", path);
		database.RemoveItemByPath(shortPath, tx);
		return;
	}
	f, _ := database.QueryItemFingerprintByPath(path);
	nf := file.Fingerprint(path);
	if nf != f {
		item := file.ReadActiveItemFile(path);
		database.RemoveItemTag(item.Id, tx);
		if len(item.Id) == 0 {
			fmt.Println("Single Index completed in", time.Since(start));
			return;
		}
		for _, tag := range item.Tags {
			database.AddItemTag(item.Id, tag, tx);
		}
		item.Fingerprint = nf;
		item.Path = shortPath;
		database.AddItem(&item, false, tx);
	}
	fmt.Println("Single Index completed in", time.Since(start));
}

func IndexActiveIdea(shortPath string, tx *sql.Tx){
	start := time.Now();
	path := filepath.Join(directory.GetDataPath(), "Ideas", shortPath);
	_, err := os.Stat(path);
	if err != nil {
		fmt.Printf("File Missing! Removing %s\n", path);
		return;
	}	
	f, _ := database.QueryIdeaFingerprintByPath(path);
	nf := file.Fingerprint(path);
	if nf != f {
		item := file.ReadActiveIdeaFile(path);
		database.RemoveIdeaTag(item.Id, tx);
		if len(item.Id) == 0 {
			fmt.Println("Single Index completed in", time.Since(start));
			return;
		}
		for _, tag := range item.Tags {
			database.AddIdeaTag(item.Id, tag, tx);
		}
		item.Fingerprint = nf;
		item.Path = shortPath;
		database.AddIdea(&item, false, tx);
	}
	fmt.Println("Single Index completed in", time.Since(start));
}
// Tags Index
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
		database.AddTags(id, tag.Tag, tag.Color, tx);
	}

	if endTx {
		database.EndTransaction(tx);
	}
}

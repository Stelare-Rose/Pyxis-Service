package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Stelare-Rose/Pyxis-Service/internal/app/directory"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/types"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func Create(){
	versionCode := "0.0.2"
	data, _ := os.ReadFile(filepath.Join(directory.GetCachePath(), "version"));
	if string(data) != versionCode {
		fmt.Println("Regenerating DB");
		os.Remove(filepath.Join(directory.GetCachePath(), "main.db"));
	}

	if db == nil {
		open();
	}

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS items (
			id TEXT NOT NULL PRIMARY KEY,
			type TEXT,
			name TEXT,
			status TEXT,
			endDate TEXT,
			startDate TEXT,
			priorityDate TEXT,
			path TEXT,
			fingerprint int64,
			isArchived bool,
			verified bool
		);
		CREATE TABLE IF NOT EXISTS tags (
			id TEXT PRIMARY KEY,
			name TEXT,
			color TEXT,
			verified bool
		);
		CREATE TABLE IF NOT EXISTS items_tags (
			item_id TEXT,
			tag_id TEXT,
			PRIMARY KEY (item_id, tag_id),
			FOREIGN KEY (item_id) REFERENCES items(id)
		);
		CREATE INDEX IF NOT EXISTS idx_items_path ON items(path);
		CREATE INDEX IF NOT EXISTS idx_items_status_archived ON items(status, isArchived);
		CREATE INDEX IF NOT EXISTS idx_items_priorityDate ON items(priorityDate);
		CREATE INDEX IF NOT EXISTS idx_items_name ON items(name);

		CREATE INDEX IF NOT EXISTS idx_items_tags_item_id ON items_tags(item_id);
		CREATE INDEX IF NOT EXISTS idx_items_tags_tag_id ON items_tags(tag_id);
	`)

	os.WriteFile(filepath.Join(directory.GetCachePath(), "version"), []byte(versionCode), 0666);

	if err != nil{
		fmt.Println(err);
	}
	fmt.Println("Success!");
}

func StartTransaction() *sql.Tx{
	if db == nil {
		open();
	}
	tx, _ := db.Begin();
	return tx;
}

func EndTransaction(tx *sql.Tx){
	tx.Commit();
	os.WriteFile(filepath.Join(directory.GetCachePath(), "fingerprint"), []byte(strconv.FormatInt(time.Now().UnixMilli(), 10)), 0666);
}

func AddItem(Item *types.Item, isArchived bool){
	if db == nil {
		open();
	}
	_, err := db.Exec(
		`INSERT INTO items (id, type, name, status, endDate, startDate, priorityDate, path, fingerprint, isArchived, verified) 
		VALUES (?,?,?,?,?,?,?,?,?,?,?)
		ON CONFLICT(id) DO UPDATE SET
			type = excluded.type,
			name = excluded.name,
			status = excluded.status,
			endDate = excluded.endDate,
			startDate = excluded.startDate,
			priorityDate = excluded.priorityDate,
			path = excluded.path,
			fingerprint = excluded.fingerprint,
			isArchived = excluded.isArchived,
			verified = excluded.verified;
		;`, 
		Item.Id, 
		Item.Type, 
		Item.Name,
		Item.Status,
		Item.EndDate,
		Item.StartDate,
		Item.PriorityDate,
		Item.Path,
		Item.Fingerprint,
		isArchived,
		true);

	if err != nil {
		fmt.Println(err);
	}
}
func AddItemWithTransaction(Item *types.Item, isArchived bool, tx *sql.Tx){
	_, err := tx.Exec(
		`INSERT INTO items (id, type, name, status, endDate, startDate, priorityDate, path, fingerprint, isArchived, verified) 
		VALUES (?,?,?,?,?,?,?,?,?, ?,?)
		ON CONFLICT(id) DO UPDATE SET
			type = excluded.type,
			name = excluded.name,
			status = excluded.status,
			endDate = excluded.endDate,
			startDate = excluded.startDate,
			priorityDate = excluded.priorityDate,
			path = excluded.path,
			fingerprint = excluded.fingerprint,
			isArchived = excluded.isArchived,
			verified = excluded.verified;
		;`, 
		Item.Id, 
		Item.Type, 
		Item.Name,
		Item.Status,
		Item.EndDate,
		Item.StartDate,
		Item.PriorityDate,
		Item.Path,
		Item.Fingerprint,
		isArchived,
		true);

	if err != nil {
		fmt.Println(err);
	}
}

func RemoveTagWithTransaction(id string, tx *sql.Tx){
	_, err := tx.Exec(
		`DELETE FROM items_tags WHERE item_id = ?`, id);
	if err != nil {
		fmt.Println(err);
	}
}

func AddTagWithTransaction(id string, tag_id string, tx *sql.Tx){
	fmt.Println("adding " + tag_id);
	_, err := tx.Exec(
		`INSERT OR IGNORE INTO items_tags (item_id, tag_id) values (?, ?)`, 
		id, tag_id);

	if err != nil {
		fmt.Println(err);
	}
}

func AddTagsWithTransaction(id string, name string, colors []string, tx *sql.Tx){
	color := strings.Join(colors, ",")
	_, err := tx.Exec(
		`
		INSERT INTO tags (id, name, color, verified) values (?, ?, ?, 1) 
		ON CONFLICT (id) DO UPDATE SET
			name = excluded.name,
			color = excluded.color,
			verified = 1
		`,
		id, name, color);
	if err != nil {
		fmt.Println(err);
	}
}


func QueryItemById(id string) (types.Item, error){
	if db == nil {
		open();
	}
	var item types.Item;
	row := db.QueryRow(
		`SELECT * FROM items WHERE id=?;`, id,
	)
	
	err := row.Scan(&item.Id, &item.Type, &item.Name, &item.Status, &item.EndDate, &item.StartDate, &item.PriorityDate, &item.Path, &item.Fingerprint, &item.IsArchived, &item.IsVerified);
	if err != nil {
		fmt.Println(err);
	}

	return item, err;
}

func QueryItemFingerprintByPath(path string) (int64, error){
	if db == nil {
		open();
	}
	var fingerprint int64;
	row := db.QueryRow(
		`SELECT fingerprint FROM items WHERE path=?;`, 
		path);
	
	err := row.Scan(&fingerprint);
	if err != nil {
		fmt.Println(err);
	}

	return fingerprint, err;
}

func RemoveItemByPath(path string){
	if db == nil {
		open();
	}

	db.Exec(`
		DELETE FROM items WHERE path=?;
	`, path);
}
func RemoveItemByPathWithTransaction(path string, tx *sql.Tx){
	if db == nil {
		open();
	}

	fmt.Println("Hello!");
	_, err := tx.Exec(`
		DELETE FROM items WHERE path=?;
	`, path);
	fmt.Println(path);
	if err != nil {
		fmt.Println(err);
	}
}
func ResetTags(tx *sql.Tx) {
	_, err := tx.Exec(`
		UPDATE tags
		SET verified=0
		`,
	);
	if err != nil {
		fmt.Println(err);
	}
}
func ResetActiveItems() {
	if db == nil {
		open();
	}
	_, err := db.Exec(`
		UPDATE items
		SET verified=0
		WHERE isArchived=0;
		`,
	);
	if err != nil {
		fmt.Println(err);
	}
}

func VerifyItems(paths []string){
	if db == nil {
		open();
	}
	tx, _ := db.Begin();
	stmt, _ := tx.Prepare(`
		UPDATE items
		SET verified=1
		WHERE path=?;
	`);

	for _, path := range paths {
		stmt.Exec(path);
	}
	tx.Commit();
}

func DropUnverifiedItems(){
	if db == nil {
		open();
	}
	
	tx, _ := db.Begin();
	tx.Exec(`DELETE FROM items WHERE isArchived=0 AND verified=0`);
	tx.Commit();
}

func open() {
	var err error;
	db, err = sql.Open("sqlite", filepath.Join(directory.GetCachePath(), "main.db"))
	if err != nil {
		fmt.Println(err);
	}
}

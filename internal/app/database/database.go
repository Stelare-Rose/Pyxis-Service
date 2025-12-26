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
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/environment"
	"github.com/Stelare-Rose/Pyxis-Service/internal/app/types"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func Create(){
	versionCode := "dev4-rev1"
	fmt.Println("Pyxis Running on Version Code " + versionCode + "!");

	// Database Check
	data, _ := os.ReadFile(filepath.Join(directory.GetCachePath(), "version"));
	if string(data) != versionCode {
		fmt.Println("Regenerating DB");
		if environment.GetEnvironment() == "dev" {
			os.Remove(filepath.Join(directory.GetCachePath(), "dev.db"));
		} else {
			os.Remove(filepath.Join(directory.GetCachePath(), "main.db"));
		}
	}

	if db == nil {
		open();
	}

	// General Schema
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS items (
			id TEXT NOT NULL PRIMARY KEY,
			type TEXT,
			name TEXT,
			status TEXT,
			endDate TEXT,
			startDate TEXT,
			priorityDate TEXT,
			completedDate TEXT,
			path TEXT,
			fingerprint int64,
			isArchived bool,
			verified bool,
			sortDate AS (
				COALESCE(
					CASE
						WHEN status = 'Done' THEN NULLIF(completedDate, '')
						ELSE NULL
					END,
					NULLIF(startDate, ''),
					NULLIF(endDate, ''),
					NULLIF(priorityDate, '')
				)
			)
		);
		CREATE TABLE IF NOT EXISTS ideas (
			id TEXT NOT NULL PRIMARY KEY,
			name TEXT,
			status TEXT,
			createdDate TEXT,
			priorityDate TEXT,
			completedDate TEXT,
			path TEXT,
			fingerprint int64,
			isArchived bool,
			verified bool,
			sortDate AS (
				COALESCE(
					CASE
						WHEN status = 'Done' THEN NULLIF(completedDate, '')
						ELSE NULL
					END,
					NULLIF(createdDate, ''),
					NULLIF(priorityDate, '')
				)
			)
		);
		CREATE TABLE IF NOT EXISTS tags (
			id TEXT PRIMARY KEY,
			tag TEXT,
			color TEXT,
			verified bool
		);
		CREATE TABLE IF NOT EXISTS items_tags (
			item_id TEXT,
			tag_id TEXT,
			PRIMARY KEY (item_id, tag_id),
			FOREIGN KEY (item_id) REFERENCES items(id)
		);	
		CREATE TABLE IF NOT EXISTS ideas_tags (
			idea_id TEXT,
			tag_id TEXT,
			PRIMARY KEY (idea_id, tag_id),
			FOREIGN KEY (idea_id) REFERENCES ideas(id)
		);
		CREATE INDEX IF NOT EXISTS idx_items_name ON items(name);
		CREATE INDEX IF NOT EXISTS idx_items_path ON items(path);
		CREATE INDEX IF NOT EXISTS idx_items_archived ON items(isArchived);
		CREATE INDEX IF NOT EXISTS idx_items_sortDate ON items(sortDate);

		CREATE INDEX IF NOT EXISTS idx_ideas_name ON ideas(name);
		CREATE INDEX IF NOT EXISTS idx_ideas_path ON ideas(path);
		CREATE INDEX IF NOT EXISTS idx_ideas_archived ON ideas(isArchived);
		CREATE INDEX IF NOT EXISTS idx_ideas_sortDate ON ideas(sortDate);

		CREATE INDEX IF NOT EXISTS idx_items_tags_item_id ON items_tags(item_id);
		CREATE INDEX IF NOT EXISTS idx_items_tags_tag_id ON items_tags(tag_id);
		CREATE INDEX IF NOT EXISTS idx_ideas_tags_idea_id ON ideas_tags(idea_id);
		CREATE INDEX IF NOT EXISTS idx_ideas_tags_tag_id ON ideas_tags(tag_id);
	`)

	// Views
	_, err = db.Exec(`
		CREATE VIEW IF NOT EXISTS ActiveItems AS
		SELECT 
		i.id, 
		i.name, 
		i.type, 
		i.path, 
		i.status, 
		i.endDate, 
		i.startDate, 
		i.priorityDate, 
		i.completedDate,
		i.fingerprint,
		GROUP_CONCAT(t.id || ':' || t.tag || ':' || t.color, ';') AS tags
		FROM items i
		LEFT JOIN items_tags it ON i.id = it.item_id
		LEFT JOIN tags t ON it.tag_id = t.id
		WHERE i.isArchived = 0 
		GROUP BY i.id
		ORDER BY sortDate IS NULL, sortDate ASC, i.name
	`)
	_, err = db.Exec(`
		CREATE VIEW IF NOT EXISTS ActiveIdeas AS
		SELECT 
		i.id, 
		i.name, 
		i.status, 
		i.path, 
		i.createdDate, 
		i.priorityDate, 
		i.completedDate,
		i.fingerprint,
		GROUP_CONCAT(t.id || ':' || t.tag || ':' || t.color, ';') AS tags
		FROM ideas i
		LEFT JOIN ideas_tags it ON i.id = it.idea_id
		LEFT JOIN tags t ON it.tag_id = t.id
		WHERE i.isArchived = 0 
		GROUP BY i.id
		ORDER BY sortDate IS NULL, sortDate ASC, i.name
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
	if environment.GetEnvironment() == "dev" {
		os.WriteFile(filepath.Join(directory.GetCachePath(), "fingerprint-dev"), []byte(strconv.FormatInt(time.Now().UnixMilli(), 10)), 0666);
	} else {
		os.WriteFile(filepath.Join(directory.GetCachePath(), "fingerprint"), []byte(strconv.FormatInt(time.Now().UnixMilli(), 10)), 0666);
	}
}

func AddItem(Item *types.Item, isArchived bool, tx *sql.Tx){
	_, err := tx.Exec(
		`INSERT INTO items (id, type, name, status, endDate, startDate, priorityDate, completedDate, path, fingerprint, isArchived, verified) 
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?)
		ON CONFLICT(id) DO UPDATE SET
			type = excluded.type,
			name = excluded.name,
			status = excluded.status,
			endDate = excluded.endDate,
			startDate = excluded.startDate,
			priorityDate = excluded.priorityDate,
			completedDate = excluded.completedDate,
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
		Item.CompletedDate,
		Item.Path,
		Item.Fingerprint,
		isArchived,
		true);

	if err != nil {
		fmt.Println(err);
	}
}
func AddIdea(Item *types.Idea, isArchived bool, tx *sql.Tx){
	_, err := tx.Exec(
		`INSERT INTO ideas (id, name, status, createdDate, priorityDate, completedDate, path, fingerprint, isArchived, verified) 
		VALUES (?,?,?,?,?,?,?,?,?,?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			status = excluded.status,
			createdDate = excluded.createdDate,
			priorityDate = excluded.priorityDate,
			completedDate = excluded.completedDate,
			path = excluded.path,
			fingerprint = excluded.fingerprint,
			isArchived = excluded.isArchived,
			verified = excluded.verified;
		;`, 
		Item.Id, 
		Item.Name,
		Item.Status,
		Item.CreatedDate,
		Item.PriorityDate,
		Item.CompletedDate,
		Item.Path,
		Item.Fingerprint,
		isArchived,
		true);

	if err != nil {
		fmt.Println(err);
	}
}


func RemoveItemTag(id string, tx *sql.Tx){
	_, err := tx.Exec(
		`DELETE FROM items_tags WHERE item_id = ?`, id);
	if err != nil {
		fmt.Println(err);
	}
}
func RemoveIdeaTag(id string, tx *sql.Tx){
	_, err := tx.Exec(
		`DELETE FROM ideas_tags WHERE idea_id = ?`, id);
	if err != nil {
		fmt.Println(err);
	}
}
func AddItemTag(id string, tag_id string, tx *sql.Tx){
	fmt.Println("adding " + tag_id);
	_, err := tx.Exec(
		`INSERT OR IGNORE INTO items_tags (item_id, tag_id) values (?, ?)`, 
		id, tag_id);

	if err != nil {
		fmt.Println(err);
	}
}
func AddIdeaTag(id string, tag_id string, tx *sql.Tx){
	fmt.Println("adding " + tag_id);
	_, err := tx.Exec(
		`INSERT OR IGNORE INTO ideas_tags (idea_id, tag_id) values (?, ?)`, 
		id, tag_id);

	if err != nil {
		fmt.Println(err);
	}
}
func AddTags(id string, tag string, colors []string, tx *sql.Tx){
	color := strings.Join(colors, ",")
	_, err := tx.Exec(
		`
		INSERT INTO tags (id, tag, color, verified) values (?, ?, ?, 1) 
		ON CONFLICT (id) DO UPDATE SET
			tag = excluded.tag,
			color = excluded.color,
			verified = 1
		`,
		id, tag, color);
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
func QueryIdeaFingerprintByPath(path string) (int64, error){
	if db == nil {
		open();
	}
	var fingerprint int64;
	row := db.QueryRow(
		`SELECT fingerprint FROM ideas WHERE path=?;`, 
		path);
	
	err := row.Scan(&fingerprint);
	if err != nil {
		fmt.Println(err);
	}

	return fingerprint, err;
}


func RemoveItemByPath(path string, tx *sql.Tx){
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
func ResetActiveIdeas() {
	if db == nil {
		open();
	}
	_, err := db.Exec(`
		UPDATE ideas
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
func VerifyIdeas(paths []string){
	if db == nil {
		open();
	}
	tx, _ := db.Begin();
	stmt, _ := tx.Prepare(`
	UPDATE ideas
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
func DropUnverifiedIdeas(){
	if db == nil {
		open();
	}

	tx, _ := db.Begin();
	tx.Exec(`DELETE FROM ideas WHERE isArchived=0 AND verified=0`);
	tx.Commit();
}


func open() {
	var err error;
	if environment.GetEnvironment() == "dev" {
		db, err = sql.Open("sqlite", filepath.Join(directory.GetCachePath(), "dev.db"))
	} else {
		db, err = sql.Open("sqlite", filepath.Join(directory.GetCachePath(), "main.db"))
	}
	if err != nil {
		fmt.Println(err);
	}
}

package file

import (
	"fmt"
	"os"
	"strings"
	"bufio"

	"github.com/Stelare-Rose/Pyxis-Service/internal/app/types"
)

func ReadActiveFile(path string) types.Item {
	file, _ := os.Open(path);
	defer file.Close();

	scanner := bufio.NewScanner(file);
	var item types.Item;
	for scanner.Scan() {
		if scanner.Text() == "-----" {
			break;
		}
		splitted := strings.SplitN(scanner.Text(), ":", 2);
		switch (splitted[0]){
		case "[Type]":
			item.Type = splitted[1];		
		case "[ID]":
			item.Id = splitted[1];
		case "[Name]":
			item.Name = splitted[1];
		case "[Status]":
			item.Status = splitted[1];
		case "[Hard-Deadline]":
			item.HardDeadline = splitted[1];
		case "[Soft-Deadline]":
			item.SoftDeadline = splitted[1];
		case "[After-Task]":
			tasks := strings.Split(splitted[1], ",");
			for i := range tasks {
				tasks[i] = strings.TrimSpace(tasks[i]);
			}
			item.AfterTask = append(item.AfterTask, tasks...);
		case "[Tags]":
			tags := strings.Split(splitted[1], ",");
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i]);
			}
			item.Tags = append(item.Tags, tags...);
		default:
			fmt.Println(splitted[1]);
			}
	}
	return item;
}

func Fingerprint(path string) int64 {
	file, _ := os.Stat(path);
	return file.Size() + file.ModTime().Unix();
}

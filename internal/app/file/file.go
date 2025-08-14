package file

import (
	"fmt"
	"os"
	"strings"

	"github.com/Stelare-Rose/Pyxis-Service/internal/app/types"
)

func ReadActiveFile(path string) types.Item {
	rawData, _ := os.ReadFile(path);
	fmt.Println(string(rawData));
	rawData = []byte(strings.TrimSpace(string(rawData)));
	var item types.Item;
	for obj := range strings.SplitSeq(string(rawData), "\n") {
		splitted := strings.SplitN(obj, ":", 2);
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
		case "[After-Task]":
			tasks := strings.Split(splitted[1], ",");
			for i := range tasks {
				tasks[i] = strings.TrimSpace(tasks[i]);
			}
			item.AfterTasks = append(item.AfterTasks, tasks...);
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

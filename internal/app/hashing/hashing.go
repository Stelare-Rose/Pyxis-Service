package hashing

import (
	"fmt"
	"hash/fnv"
	"os"
	"strconv"
)

func MetadataHash(path string) uint64 {
	file, err := os.Stat(path);
	if err != nil {
		fmt.Println(err);
		return 0;
	}
	//fmt.Println(file.Name(), file.ModTime(), file.Size());
	h := getHash(file);
	return h;
}

func getHash(file os.FileInfo) uint64 {
	h := fnv.New64();
	h.Write([]byte(strconv.FormatInt(file.Size(), 10)));
	h.Write([]byte(strconv.FormatInt(file.ModTime().Unix(), 10)));
	if(h.Sum64() == 0) {
		return 1;
	}
	return h.Sum64();
}

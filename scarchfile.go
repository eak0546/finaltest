package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

type MyFile struct {
	Path    string
	Size    int64
	Name    string
	ModTime time.Time
}

func getDrives() (r []string) {
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		f, err := os.Open(string(drive) + ":\\")
		if err == nil {
			d := string(drive) + ":/"
			r = append(r, d)
			f.Close()
		}
	}
	return
}
func ProcessingExtension(dir string, f os.FileInfo, extension map[string]string, files *[]MyFile, wg *sync.WaitGroup) {
	defer wg.Done()
	filename := f.Name()
	index := strings.LastIndex(filename, ".")
	if index < 0 {
		return
	}
	index = index + 1
	size := len(filename)
	ext := filename[index:size]
	_, ok := extension[ext]
	if ok {
		var mf MyFile
		mf.Path = dir + "/" + f.Name()
		mf.Size = f.Size()
		mf.Name = f.Name()
		mf.ModTime = f.ModTime()
		*files = append(*files, mf)
	}
}
func FindFileFromExtension(extension map[string]string, dir string, files *[]MyFile, wg *sync.WaitGroup) {
	defer wg.Done()
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, f := range fs {
		var subWg sync.WaitGroup
		subWg.Add(1)
		if f.IsDir() {
			path := dir + "/" + f.Name()
			go FindFileFromExtension(extension, path, files, &subWg)

		} else {
			ProcessingExtension(dir, f, extension, files, &subWg)
		}
		subWg.Wait()
	}
}

func txtCreate() {
	file, err := os.Create("./TXT/1.txt")
	if err != nil {
		return
	}
	for x := 1; x <= 1000; x++ {
		defer file.Close()
		file.WriteString("Atiwan\nPhongam\n25\n")
	}
}

func main() {
	var wg sync.WaitGroup
	t1 := time.Now()
	myfiles := []MyFile{}
	extensionMap := make(map[string]string)
	extensionMap["jpg"] = ".jpg"
	extensionMap["txt"] = ".txt"
	drives := getDrives()
	wg.Add(len(drives))
	for _, drive := range drives {
		go FindFileFromExtension(extensionMap, drive, &myfiles, &wg)

	}
	wg.Wait()
	t2 := time.Now()
	diftime := t2.Sub(t1)
	fmt.Println("total files = ", len(myfiles))
	fmt.Println("total time = ", diftime)
}

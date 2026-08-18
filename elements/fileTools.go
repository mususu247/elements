package elements

// version 2026-05-21 s.sato
import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func FileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if err != nil {
		return false

	} else {
		if info.IsDir() {
			return false

		} else {
			return true
		}
	}
}

func FolderExists(folderName string) bool {
	info, err := os.Stat(folderName)
	if err != nil {
		return false

	} else {
		if info.IsDir() {
			return true

		} else {
			return false
		}
	}
}

func GetFullPath(fileName string) (string, error) {
	dir, fn := filepath.Split(fileName)
	dir = filepath.ToSlash(dir)
	if (dir == "./") || (dir == "") {
		exeFull, _ := os.Executable()
		exePath, _ := filepath.Split(exeFull)
		fileName = filepath.Join(exePath, fn)
	}
	fileName, err := filepath.Abs(fileName)
	if err != nil {
		return "", err
	}
	fileName = filepath.ToSlash(fileName)
	return fileName, nil
}

func GetFileName(filePath string) string {
	_, fn := filepath.Split(filePath)
	return fn
}

func GetPath(filePath string) string {
	dir, _ := filepath.Split(filePath)
	return dir
}

func PathSplit(filePath string) (path string, name string) {
	path, name = filepath.Split(filePath)
	return path, name
}

func GetFileModTime(fileName string) (time.Time, error) {
	var t time.Time
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return t, err

	} else {
		return fileInfo.ModTime(), nil
	}
}

func SetFileModTime(fileName string, stamp time.Time) {
	err := os.Chtimes(fileName, stamp, stamp)
	if err != nil {
		log.Printf("(Error) setFileModTime: %v (%v)", fileName, stamp)
	}
}

func TimeNow() string {
	const stamplayout = "2006-01-02 15:04:05.000"
	return time.Now().Format(stamplayout) + " "
}

func findSubFolder(basePath string, find string, deep int) []string {
	var list []string
	dirEntrys, err := os.ReadDir(basePath)
	if err != nil {
		return list
	}
	for i := range dirEntrys {
		dirEntry := dirEntrys[i]
		n := dirEntry.Name()
		fn := filepath.Join(basePath, n)
		if dirEntry.IsDir() {
			if deep == -1 {
				sub := findSubFolder(fn, find, deep)
				list = append(list, sub...)

			} else {
				deep--
				if deep >= 0 {
					sub := findSubFolder(fn, find, deep)
					list = append(list, sub...)
				}
			}
		} else {
			ok, _ := filepath.Match(find, n)
			if ok {
				list = append(list, fn)
			}
		}
	}
	return list
}

func FindFile(findPath string, deep int) ([]string, error) {
	var list []string
	if FileExists(findPath) {
		list = append(list, findPath)
		return list, nil
	}
	dir, find := PathSplit(findPath)
	log.Printf("dir:%v find:%v", dir, find)
	subList := findSubFolder(dir, find, deep)
	if len(subList) > 0 {
		list = append(list, subList...)
	}
	return list, nil
}

func findSubFolderEx(basePath string, find string, deep int, last time.Time) []string {
	var list []string
	dirEntrys, err := os.ReadDir(basePath)
	if err != nil {
		return list
	}
	for i := range dirEntrys {
		dirEntry := dirEntrys[i]
		n := dirEntry.Name()
		fn := filepath.Join(basePath, n)
		if dirEntry.IsDir() {
			if deep == -1 {
				sub := findSubFolderEx(fn, find, deep, last)
				list = append(list, sub...)

			} else {
				deep--
				if deep >= 0 {
					sub := findSubFolderEx(fn, find, deep, last)
					list = append(list, sub...)
				}
			}
		} else {
			ok, _ := filepath.Match(find, n)
			if ok {
				ft, _ := GetFileModTime(fn)
				if last.Before(ft) {
					//fmt.Printf("fn:%v,dt:%v\n",fn,ft)
					list = append(list, fn)
				}
			}
		}
	}
	return list
}

func FindFileEx(findPath string, deep int, last time.Time) ([]string, error) {
	var list []string
	if FileExists(findPath) {
		ft, _ := GetFileModTime(findPath)
		if last.Before(ft) {
			//fmt.Printf("fn:%v,dt:%v\n",findPath,ft)
			list = append(list, findPath)
		}
		return list, nil
	}
	dir, find := PathSplit(findPath)
	log.Printf("dir:%v find:%v", dir, find)
	subList := findSubFolderEx(dir, find, deep, last)
	if len(subList) > 0 {
		list = append(list, subList...)
	}
	return list, nil
}

func GetAbsolutePathName(filePath string) (string, error) {
	return filepath.Abs(filePath)
}

func GetExeName() string {
	path := os.Args[0]
	path = GetFileName(path)
	return path
}

func GetConfigName() string {
	defaultName := "config.json"
	exeName := GetExeName()
	exeName = strings.ReplaceAll(exeName, ".exe", ".json")
	if !FileExists(exeName) {
		exeName = defaultName
	}
	return "./" + exeName
}

func DeleteFile(fileName string) error {
	return os.Remove(fileName)
}

func GetExtensionName(fileName string) string {
	ext := filepath.Ext(fileName)
	ext = strings.ToLower(ext)
	return ext
}

func timeComp(dt1, dt2 time.Time) int {
	var x1 string
	var x2 string
	const layout = "2006-01-02 15:04:05"
	x1 = dt1.Format(layout)
	x2 = dt2.Format(layout)
	if x1 == x2 {
		return 0
	}
	if x1 > x2 {
		return -1
	}
	if x1 < x2 {
		return 1
	}
	return 0
}

func AddLocal(value time.Time) time.Time {
	var x []int
	for i := 0; i < 6; i++ {
		switch i {
		case 0:
			x = append(x, value.Year())
		case 1:
			x = append(x, int(value.Month()))
		case 2:
			x = append(x, value.Day())
		case 3:
			x = append(x, value.Hour())
		case 4:
			x = append(x, value.Minute())
		case 5:
			x = append(x, value.Second())
		}
	}
	loc, _ := time.LoadLocation("Local")
	return time.Date(x[0], time.Month(x[1]), x[2], x[3], x[4], x[5], 0, loc)
}

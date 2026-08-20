package elements

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func CheckZip(zipFile string) error {
	zr, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zr.Close()

	for _, zf := range zr.File {
		fmt.Printf("name:'%v' mod:#%v# \n", zf.Name, zf.Modified.Format("2006-01-02 15:00:00"))
	}
	return nil
}

func FromZip(srcFile string, dstPath string) error {
	_, err := os.Stat(dstPath)
	if err == nil {
		err = os.RemoveAll(dstPath)
		if err != nil {
			return err
		}
	}

	zr, err := zip.OpenReader(srcFile)
	if err != nil {
		log.Printf("(Error) %v", err)
		return err
	}
	defer zr.Close()

	for _, zf := range zr.File {
		fd := zf.Modified
		src, err := zf.Open()
		if err != nil {
			log.Printf("(Error) %v", err)
			return err
		}
		defer src.Close()

		unFile := filepath.Join(dstPath, zf.Name)
		dir := filepath.Dir(unFile)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Printf("(Error) %v", err)
			continue
		}

		dst, err := os.Create(unFile)
		if err != nil {
			log.Printf("(Error) %v", err)
			continue
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			log.Printf("(Error) %v", err)
			continue
		}
		SetFileModTime(unFile, AddLocal(fd))
	}
	return nil
}

func toZip(basePath string, srcPath string, zw *zip.Writer) error {
	list, err := os.ReadDir(srcPath)
	if err != nil {
		log.Printf("(Error) %v \n", err)
		return err
	}

	for i := range list {
		absPath := filepath.Join(srcPath, list[i].Name())
		relPath, _ := filepath.Rel(basePath, absPath)

		if list[i].IsDir() {
			toZip(basePath, absPath, zw)
		} else {
			fw, err := zw.Create(relPath)
			if err != nil {
				log.Printf("(Error) %v \n", err)
				continue
			}

			bytes, err := os.ReadFile(absPath)
			if err == nil {
				fw.Write(bytes)
			}
		}
	}
	return nil
}

func ToZip(srcPath string, dstFile string) error {
	_, err := os.Stat(dstFile)
	if err == nil {
		err = os.Remove(dstFile)
		if err != nil {
			return err
		}
	}

	zipFile, _ := os.Create(dstFile)
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	srcPath, err = filepath.Abs(srcPath)
	if err != nil {
		return err
	}

	basePath := srcPath + string(filepath.Separator)
	toZip(basePath, srcPath, zipWriter)
	return nil
}

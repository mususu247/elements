package elements

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
)

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
		var dirPath string
		if zf.FileInfo().IsDir() {
			dirPath = zf.Name
		} else {
			dirPath = path.Dir(zf.Name)
		}

		unzipPath := filepath.Join(dstPath, dirPath)
		if err := os.MkdirAll(unzipPath, 0755); err != nil {
			log.Printf("(Error) %v", err)
			return err
		}
		if !zf.FileInfo().IsDir() {
			fd := zf.Modified
			src, err := zf.Open()
			if err != nil {
				log.Printf("(Error) %v", err)
				return err
			}
			defer src.Close()

			unFile := filepath.Join(dstPath, zf.Name)
			dst, err := os.Create(unFile)
			if err != nil {
				log.Printf("(Error) %v", err)
				return err
			}
			defer dst.Close()

			if _, err := io.Copy(dst, src); err != nil {
				log.Printf("(Error) %v", err)
				return err
			}
			SetFileModTime(unFile, AddLocal(fd))
		}
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
			relPath = relPath + string(filepath.Separator)
			_, err = zw.Create(relPath)
			if err != nil {
				log.Printf("(Error) %v \n", err)
				continue
			}

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

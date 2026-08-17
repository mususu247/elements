package elements

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func FromZip(srcFile string, dstPath string) error {
	zr, err := zip.OpenReader(srcFile)
	if err != nil {
		fmt.Printf("(Error) %v", err)
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
			fmt.Printf("(Error) %v", err)
			return err
		}
		if !zf.FileInfo().IsDir() {
			fd := zf.Modified
			src, err := zf.Open()
			if err != nil {
				fmt.Printf("(Error) %v", err)
				return err
			}
			defer src.Close()

			unFile := filepath.Join(dstPath, zf.Name)
			dst, err := os.Create(unFile)
			if err != nil {
				fmt.Printf("(Error) %v", err)
				return err
			}
			defer dst.Close()
			if _, err := io.Copy(dst, src); err != nil {
				fmt.Printf("(Error) %v", err)
				return err
			}
			SetFileModTime(unFile, AddLocal(fd))
		}
	}
	return nil
}

func ToZip(srcPath string, dstFile string) error {
	findPath := filepath.Join(srcPath, "*.*")
	list, err := FindFile(findPath, -1)
	if err != nil {
		return err
	}
	// 新しいzipアーカイブファイルを作成
	zipFile, _ := os.Create(dstFile)
	defer zipFile.Close()

	// zipアーカイブのライターを作成
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for i := range list {
		// ファイルをzipアーカイブに追加
		srcFile := strings.ReplaceAll(list[i], srcPath, "")
		//fileInfo, err := os.Stat(srcFile)
		//if err != nil {
		//	continue
		//}

		//header, _ := zip.FileInfoHeader(fileInfo)
		//writer, _ := zipWriter.CreateHeader(header)
		writer, _ := zipWriter.Create(srcFile)
		file, err := os.Open(list[i])
		if err == nil {
			io.Copy(writer, file)
		}
	}

	return nil
}

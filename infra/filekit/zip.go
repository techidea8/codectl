package filekit

import (
	"archive/zip"
	"io"
	"os"
)

type ziper struct {
	zipfile *os.File
	files   []string
}

func NewZiper(zipname string) *ziper {
	file, _ := os.Create(zipname)
	return &ziper{
		zipfile: file,
		files:   make([]string, 0),
	}
}
func (s *ziper) Add(file ...string) *ziper {
	s.files = append(s.files, file...)
	return s
}

func (s *ziper) Zip() error {
	return ZipFiles(s.zipfile, s.files)
}

// ZipFiles compresses one or many files into a single zip archive file.
// Param 1: filename is the output zip file's name.
// Param 2: files is a list of files to add to the zip.
func ZipFiles(zipfile *os.File, files []string) error {
	defer zipfile.Close()

	zipWriter := zip.NewWriter(zipfile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err := AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

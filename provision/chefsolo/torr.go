package chefsolo

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Torr struct {
	Source string
	Base   string
	writer io.Writer
}

func NewTorr(sourcefile string) *Torr {
	return &Torr{Source: sourcefile}
}

func (t *Torr) untar() error {
	fmt.Fprintf(t.writer, "  untar (%s)\n", t.Source)

	file, err := os.Open(t.Source)
	if err != nil {
		return err
	}

	defer file.Close()
	var fileReader io.ReadCloser = file
	// just in case we are reading a tar.gz file, add a filter to handle gzipped file
	if strings.HasSuffix(t.Source, ".gz") {
		if fileReader, err = gzip.NewReader(file); err != nil {
			return err
		}
		defer fileReader.Close()
	}

	tarBallReader := tar.NewReader(fileReader)
	// Extracting tarred files
	for {
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		filename := filepath.Join(t.Base, header.Name[strings.Index(header.Name, "/")+1:])

		switch header.Typeflag {
		case tar.TypeDir:
			// handle directory
			fmt.Fprintf(t.writer, "  creating directory (%s)\n", filename)
			err = os.MkdirAll(filename, os.FileMode(header.Mode)) // or use 0755 if you prefer

			if err != nil {
				return err
			}

		case tar.TypeReg:
			// handle normal file
			fmt.Fprintf(t.writer, "  writing untarred (%s)\n", filename)
			writer, err := os.Create(filename)

			if err != nil {
				return err
			}

			io.Copy(writer, tarBallReader)

			err = os.Chmod(filename, os.FileMode(header.Mode))

			if err != nil {
				return err
			}

			writer.Close()
		default:
			fmt.Fprintf(t.writer, "  unable to untar type : %c in file (%s)\n", header.Typeflag, filename)
		}
	}
	fmt.Fprintf(t.writer, "  untar (%s) OK\n", t.Source)
	return nil
}

func (t *Torr) cleanup() error {
	fmt.Fprintf(t.writer, "  cleanup tar (%s) OK\n", t.Source)
	if _, err := os.Stat(t.Source); err == nil {
		if err = os.Remove(t.Source); err != nil {
			return err
		}
	}
	fmt.Fprintf(t.writer, "  cleanup tar (%s) OK\n", t.Source)
	return nil
}

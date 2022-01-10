package utils

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func UnTar(src io.Reader, dst string) error {
	gr, err := gzip.NewReader(src)
	if err != nil {
		return fmt.Errorf("error creating gzip reader: %w", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		header, err := tr.Next()

		switch {

		// If no more files are found return
		case err == io.EOF:
			return nil

		// Return any other error
		case err != nil:
			return err

		// If the header is nil, skip it
		case header == nil:
			continue

		// Skip the any files duplicated as hidden files
		case strings.HasPrefix(header.Name, "._"):
			continue
		}

		// The target location where the dir/file should be created
		segs := strings.Split(header.Name, string(filepath.Separator))
		segs = segs[1:]
		target := filepath.Join(dst, filepath.Join(segs...))

		fi := header.FileInfo()

		if fi.IsDir() {
			os.MkdirAll(target, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(target), os.ModePerm); err != nil {
			return err
		}

		fd, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fi.Mode())
		if err != nil {
			return err
		}

		// NOTE: We use looped CopyN() not Copy() to avoid gosec G110 (CWE-409):
		// Potential DoS vulnerability via decompression bomb.
		for {
			_, err := io.CopyN(fd, tr, 1024)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
		}

		fd.Close()
	}
}

// unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func UnZip(src string, dst string) ([]string, error) {
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		// The zip contains a folder, and inside that folder are the files we're
		// interested in. So while looping over the files (whose .Name field is the
		// full path including the containing folder) we strip out the first path
		// segment to ensure the files we need are extracted to the current directory.
		segs := strings.Split(f.Name, string(filepath.Separator))
		segs = segs[1:]
		fpath := filepath.Join(dst, filepath.Join(segs...))
		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		fd, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		// NOTE: We use looped CopyN() not Copy() to avoid gosec G110 (CWE-409):
		// Potential DoS vulnerability via decompression bomb.
		for {
			_, err := io.CopyN(fd, rc, 1024)
			if err != nil {
				if err == io.EOF {
					break
				}
				return filenames, err
			}
		}

		fd.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}

	return filenames, nil
}

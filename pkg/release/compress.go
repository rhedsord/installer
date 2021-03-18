package release

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

var tarfile *tar.Writer
var trimPath string

// compress the final bundle

func createArchive(bundle *BundleRoot) error {

	// Create bundle archive file
	file, err := os.Create(bundle.BaseDir + "/" + bundle.Version + ".tgz")
	if err != nil {
		return err
	}
	defer file.Close()
	gw := gzip.NewWriter(file)
	defer gw.Close()
	tarfile = tar.NewWriter(gw)
	defer tarfile.Close()

	logrus.Info("Creating bundle file: " + file.Name())

	trimPath = bundle.BaseDir + "/"
	err = filepath.Walk(bundle.BundleDir, addToTar)

	return err

}

func addToTar(path string, info os.FileInfo, err error) error {
	logrus.Debugln("Path: " + path)

	var link string

	// Check for and handle symbolic links
	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		link, _ = os.Readlink(path)
		logrus.Debugln("Symlink: ", link)
	}

	// Set the file header info, trim the base dir path so pathing is relative the
	// bundle directory
	fhdr, _ := tar.FileInfoHeader(info, link)
	fhdr.Name = strings.TrimPrefix(path, trimPath)
	tarfile.WriteHeader(fhdr)

	// Write the file to the archive if it is a regular file
	if info.Mode().IsRegular() {

		if file, err := os.Open(path); err == nil {

			defer file.Close()

			if ws, err := io.Copy(tarfile, file); err == nil {
				logrus.Debugln("Write: "+file.Name()+"; Size: ", ws)
			} else {
				logrus.Errorln(file.Name(), "; ", err)
				return err
			}
		} else {
			logrus.Errorln(err)
			return err
		}

	}

	return nil
}

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

func createArchive(bundle *BundleRoot) {

	file, _ := os.Create(bundle.BaseDir + "/" + bundle.Version + ".tgz")
	defer file.Close()
	gw := gzip.NewWriter(file)
	defer gw.Close()
	tarfile = tar.NewWriter(gw)
	defer tarfile.Close()

	logrus.Info("Creating bundle file: " + file.Name())

	trimPath = bundle.BaseDir + "/"
	filepath.Walk(bundle.BundleDir, addToTar)

}

func addToTar(path string, info os.FileInfo, err error) error {
	logrus.Debugln("Path: " + path)

	var link string

	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		link, _ = os.Readlink(path)
		logrus.Debugln("Symlink: ", link)
	}

	fhdr, _ := tar.FileInfoHeader(info, link)
	fhdr.Name = strings.TrimPrefix(path, trimPath)
	tarfile.WriteHeader(fhdr)

	if info.Mode().IsRegular() {

		if file, err := os.Open(path); err == nil {

			defer file.Close()

			if ws, err := io.Copy(tarfile, file); err == nil {
				logrus.Debugln("Write: "+file.Name()+"; Size: ", ws)
			} else {
				logrus.Errorln(file.Name(), "; ", err)
			}
		}

	}

	return nil
}

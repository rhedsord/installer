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

	fhdr, _ := tar.FileInfoHeader(info, "")
	fhdr.Name = strings.TrimPrefix(path, trimPath)
	tarfile.WriteHeader(fhdr)
	file, err := os.Open(path)
	defer file.Close()

	ws, err := io.Copy(tarfile, file)
	if err != nil {
		logrus.Errorln(err)
	} else {
		logrus.Infoln("Write Size: ", ws)
	}

	return nil
}

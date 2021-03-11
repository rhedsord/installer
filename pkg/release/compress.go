package release

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var tarfile *tar.Writer

// compress the final bundle

func createArchive(basedir string) {

	file, _ := os.Create("filename.tgz")
	defer file.Close()
	gw := gzip.NewWriter(file)
	defer gw.Close()
	tarfile = tar.NewWriter(gw)
	defer tarfile.Close()

	filepath.Walk(basedir, addToTar)

}

func addToTar(path string, info os.FileInfo, err error) error {
	logrus.Debugln("Path: " + path)

	fhdr, _ := tar.FileInfoHeader(info, info.Name())
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

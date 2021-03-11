package release

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/sirupsen/logrus"
)

//var BundleDir string

func (b *BundleRoot) createTree() {

	b.BaseDir = b.RootDir + "/bundle"
	b.BundleDir = b.BaseDir + "/" + b.Version
	logrus.Info("bundle dir: ", b.BundleDir)

	v := reflect.ValueOf(&b.SubTree).Elem()
	logrus.Info("children names ", v)

	// TODO: Iterate over the outer struct to dynamically
	// create the entire tree.
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		u := fmt.Sprintf("%v", f)
		n := b.BundleDir + "/" + u
		os.MkdirAll(n, 0777)

	}

}

func (b *BundleInfo) writeInfo(bundle *BundleRoot) {

	p, _ := json.MarshalIndent(b, "", "    ")
	logrus.Info(string(p))

	f := bundle.BundleDir + "/bundle-info.json"
	logrus.Info(f)

	err := ioutil.WriteFile(f, p, 0777)

	if err != nil {
		logrus.Error(err)

	}

}

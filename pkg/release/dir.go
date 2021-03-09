package release

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/sirupsen/logrus"
)

func (b BundleRoot) createTree(baseDir string) {

	bundleDir := baseDir + "/" + "bundle" + "/" + b.Version
	logrus.Info("bundle dir: ", bundleDir)

	v := reflect.ValueOf(&b.SubTree).Elem()
	logrus.Info("children names ", v)

	// TODO: Iterate over the outer struct to dynamically
	// create the entire tree.
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		u := fmt.Sprintf("%v", f)
		n := bundleDir + "/" + u
		os.MkdirAll(n, 0777)

	}

}

func (b BundleInfo) writeInfo(baseDir string) {

	p, _ := json.MarshalIndent(b, "", "    ")
	logrus.Info(p)

	f := baseDir + "bundle/" + b.Version + "/bundle-info.json"
	logrus.Info(f)

	err := ioutil.WriteFile(f, p, 0777)

	logrus.Info(err)

}

package release

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/openshift/installer/pkg/rhcos"
	"github.com/sirupsen/logrus"
)

func getRhcosVersion(v string) string {
	ctx, _ := context.WithTimeout(context.TODO(), 30*time.Second)
	meta, _ := rhcos.FetchRHCOSBuild(ctx, arch)
	logrus.Info("meta: ", meta)
	logrus.Info("rhcosver ", meta.OSTreeVersion)
	return meta.OSTreeVersion
}

func (b BundleInfo) bundleGetRhcosURL() (string, string) {
	// get Platform
	p := b.Platform

	// get arch
	a := b.Arch

	// build url

	logrus.Info(p, a)

	// return url
	return "url", "file"

}

func (b BundleInfo) downRhcos(baseDir string) error {

	// First Get RHCOS URL for platform
	// bundleGetRhcosURL(b)

	url := "URL"

	// Then get filepath to save

	filepath := baseDir + b.Version + "rhcos"

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil

}

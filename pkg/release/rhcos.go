package release

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/sirupsen/logrus"
)

func getRhcosVersion(v string) string {
	ctx, _ := context.WithTimeout(context.TODO(), 30*time.Second)
	meta, _ := rhcos.FetchRHCOSBuild(ctx, arch)
	logrus.Info("meta: ", meta)
	logrus.Info("rhcosver ", meta.OSTreeVersion)
	return meta.OSTreeVersion
}

// func (b BundleInfo) getRhcosName() (string, error) {
// 	ctx, _ := context.WithTimeout(context.TODO(), 30*time.Second)

// 	meta, err := rhcos.FetchRHCOSBuild(ctx, arch)
// 	if err != nil {
// 		return "", err //, errors.Wrap(err, "failed to fetch RHCOS metadata")
// 	}
// 	logrus.Info(meta)
// 	s := b.Platform

// 	var platform rhcos.s = rhcos.s(s)
// 	path := "meta.Images." + b.Platform + ".Path"

// 	rhcosNameURL, err := url.Parse(meta.Images.platform.Path)
// 	rhcosName := fmt.Sprintf("%v", rhcosNameURL)

// 	if err != nil {
// 		return "", err
// 	}

// 	return rhcosName, nil
// }

func (b BundleInfo) getRhcosURL() (string, error) {

	ctx, _ := context.WithTimeout(context.TODO(), 30*time.Second)

	// get Platform
	p := b.Platform

	// get arch
	var a types.Architecture = types.Architecture(b.Arch)

	var osimage string
	var err error
	switch p {
	case "aws":
		osimage, err = rhcos.AWS(ctx, a)
	case "gcp":
		osimage, err = rhcos.GCP(ctx, a)
	case "libvirt":
		osimage, err = rhcos.QEMU(ctx, a)
	case "openstack":
		//if oi := config.Platform.OpenStack.ClusterOSImage; oi != "" {
		//	osimage = oi
		//	break
		//}
		osimage, err = rhcos.OpenStack(ctx, arch)
	case "ovirt":
		osimage, err = rhcos.OpenStack(ctx, arch)
	case "kubevirt":
		osimage, err = rhcos.OpenStack(ctx, arch)
	case "azure":
		osimage, err = rhcos.VHD(ctx, arch)
	case "baremetal":
		// Check for RHCOS image URL override
		//if oi := config.Platform.BareMetal.ClusterOSImage; oi != "" {
		//	osimage = oi
		//	break
		//}
		// Note that baremetal IPI currently uses the OpenStack image
		// because this contains the necessary ironic config drive
		// ignition support, which isn't enabled in the UPI BM images
		osimage, err = rhcos.OpenStack(ctx, arch)
	case "vsphere":
		// Check for RHCOS image URL override
		//if config.Platform.VSphere.ClusterOSImage != "" {
		//	osimage = config.Platform.VSphere.ClusterOSImage
		//	break
		//}
		osimage, err = rhcos.VMware(ctx, arch)
	case none.Name:
	default:
		return "", errors.New("invalid Platform")
	}
	if err != nil {
		return "", err
	}
	return osimage, nil

}

func (b BundleInfo) getRhcosName() (string, error) {

	ctx, _ := context.WithTimeout(context.TODO(), 30*time.Second)
	meta, err := rhcos.FetchRHCOSBuild(ctx, arch)
	// if err != nil {
	// 	return "", errors.Wrap(err, "failed to fetch RHCOS metadata")
	// }

	// get Platform
	p := b.Platform

	// get arch
	// var a types.Architecture = types.Architecture(b.Arch)

	var filename string
	switch p {
	case "aws":
		filename = meta.Images.AWS.Path
	case "libvirt":
		filename = meta.Images.QEMU.Path
	case "openstack":
		//if oi := config.Platform.OpenStack.ClusterOSImage; oi != "" {
		//	osimage = oi
		//	break
		//}
		filename = meta.Images.OpenStack.Path
	case "ovirt":
		filename = meta.Images.OpenStack.Path
	case "kubevirt":
		filename = meta.Images.OpenStack.Path
	case "vsphere":
		// Check for RHCOS image URL override
		//if config.Platform.VSphere.ClusterOSImage != "" {
		//	osimage = config.Platform.VSphere.ClusterOSImage
		//	break
		//}
		filename = meta.Images.VMware.Path
	case none.Name:
	default:
		return "", errors.New("invalid Platform")
	}
	if err != nil {
		return "", err
	}
	return filename, nil

}

func (b BundleInfo) downRhcos(baseDir string) error {

	// First Get RHCOS URL for platform
	// bundleGetRhcosURL(b)

	url, err := b.getRhcosURL()
	logrus.Info("RHCOS URL ", url)
	logrus.Error(err)

	// Then get filepath to save
	filename, err := b.getRhcosName()
	logrus.Info(filename)

	filepath := baseDir + "/bundle/" + b.Version + "/rhcos/" + filename

	logrus.Info(filepath)

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

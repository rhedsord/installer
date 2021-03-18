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
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	meta, _ := rhcos.FetchRHCOSBuild(ctx, archAmd64)
	logrus.Debug("Meta: ", meta)
	logrus.Debug("RHCOS version: ", meta.OSTreeVersion)
	return meta.OSTreeVersion
}

func (bundle BundleInfo) getRhcosURL() (string, error) {

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)

	defer cancel()

	// get Platform
	platfrm := bundle.Platform

	// get arch
	var arch types.Architecture = types.Architecture(bundle.Arch)

	var osimage string
	var err error
	switch platfrm {
	case "aws":
		osimage, err = rhcos.AWS(ctx, arch)
	case "gcp":
		osimage, err = rhcos.GCP(ctx, arch)
	case "libvirt":
		osimage, err = rhcos.QEMU(ctx, arch)
	case "openstack":
		osimage, err = rhcos.OpenStack(ctx, archAmd64)
	case "ovirt":
		osimage, err = rhcos.OpenStack(ctx, archAmd64)
	case "kubevirt":
		osimage, err = rhcos.OpenStack(ctx, archAmd64)
	case "azure":
		osimage, err = rhcos.VHD(ctx, archAmd64)
	case "baremetal":
		// Note that baremetal IPI currently uses the OpenStack image
		// because this contains the necessary ironic config drive
		// ignition support, which isn't enabled in the UPI BM images
		osimage, err = rhcos.OpenStack(ctx, archAmd64)
	case "vsphere":
		osimage, err = rhcos.VMware(ctx, archAmd64)
	case none.Name:
	default:
		return "", errors.New("invalid Platform")
	}
	if err != nil {
		return "", err
	}
	return osimage, nil

}

func (bundle BundleInfo) getRhcosName() (string, error) {

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	meta, err := rhcos.FetchRHCOSBuild(ctx, archAmd64)
	if err != nil {
		return "", err
	}

	// get Platform
	platfrm := bundle.Platform

	// get arch
	// var a types.Architecture = types.Architecture(b.Arch)

	var filename string
	switch platfrm {
	case "aws":
		filename = meta.Images.AWS.Path
	case "libvirt":
		filename = meta.Images.QEMU.Path
	case "openstack":
		filename = meta.Images.OpenStack.Path
	case "ovirt":
		filename = meta.Images.OpenStack.Path
	case "kubevirt":
		filename = meta.Images.OpenStack.Path
	case "vsphere":
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

func (b BundleInfo) downRhcos(bundle *BundleRoot) error {

	// First Get RHCOS URL for platform
	url, err := b.getRhcosURL()
	if err == nil {
		logrus.Info("RHCOS URL ", url)
	} else {
		return err
	}

	// Then get filepath to save
	filename, err := b.getRhcosName()
	if err == nil {
		logrus.Info(filename)
	} else {
		return err
	}

	filepath := bundle.BundleDir + "/" + bundle.SubTree.Rhcos + "/" + filename

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

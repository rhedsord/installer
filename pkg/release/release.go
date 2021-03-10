package release

import (
	"context"
	"time"

	"github.com/openshift/installer/pkg/asset/installconfig"
	assetstore "github.com/openshift/installer/pkg/asset/store"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/version"
	"github.com/sirupsen/logrus"
)

const (
	arch = "amd64"
)

// BundleInfo is the metadata of the bundle
type BundleInfo struct {
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	RhcosVer string `json:"rhcosVer"`
}

// subtree is the content of a versioned bundle
type subTree struct {
	Client string `json:"client"`
	Images string `json:"images,omitempty"`
	Rhcos  string `json:"rhcos,omitempty"`
}

// BundleRoot is the directory/file structure of the bundle
type BundleRoot struct {
	Version string `json:"version,omitempty"`
	SubTree subTree
}

// CreateBundle retrieves all dependencies for OCP installation
// and places them into a directory or compressed file.
func CreateBundle(baseDir string) {

	if AssetStore, err := assetstore.NewStore(baseDir); err == nil {

		logrus.Info("AssestStore: ")
		if installConfig, err := AssetStore.Load(&installconfig.InstallConfig{}); err == nil && installConfig != nil {
			logrus.Debug("Install-config: ", installConfig.Name())
			logrus.Info("Platform: ", installConfig.(*installconfig.InstallConfig).Config.PullSecret)
		} else {
			logrus.Errorln("Install Config Load failed: ", err)
		}
	} else {
		logrus.Errorln("Assest Store failed: ", err)
	}
	// arch := config.ControlPlane.Architecture
	ctx, _ := context.WithTimeout(context.TODO(), 30*time.Second)
	logrus.Info("ctx: ", ctx)
	ocpVersion, _ := version.Version()
	logrus.Info("ocpVersion: ", ocpVersion)
	meta, _ := rhcos.FetchRHCOSBuild(ctx, arch)
	logrus.Info("meta: ", meta)
	rhcosVersion := meta.OSTreeVersion
	logrus.Info("rhcosver ", rhcosVersion)

	// logrus.Infoln("OCP Version: ", ocpVersion)
	// logrus.Infoln("RHCOS Version: ", rhcosVersion)
	// logrus.Debugln("Arch: ", arch)

	bundle := &BundleRoot{
		Version: ocpVersion,
		SubTree: subTree{
			Client: "client",
			Images: "images",
			Rhcos:  "rhcos",
		},
	}

	bundleInfo := &BundleInfo{

		Arch:     arch,
		Version:  ocpVersion,
		RhcosVer: rhcosVersion,
	}

	// }
	// logrus.Debugln(bundle)

	// create directory tree
	// if directory tree is not built or missing subdirectories
	bundle.createTree(baseDir)
	bundleInfo.writeInfo(baseDir)
	// else return nil
	//

	// Get release images
	// if all image digests are present/correct
	// return nil
	// else mirrorImages()
	mirrorImages(baseDir)

	// Get rhcos image
	// if rhcos image hash matches
	// return nil
	// else getRhcos()

	// Get binaries
	// if binary bundles hash match
	// return nil
	// else getBin()

	// Write bundle metadata
	// writeMeta()

	// Compress bundle
	// compressBundle()

}

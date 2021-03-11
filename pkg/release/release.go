package release

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	assetstore "github.com/openshift/installer/pkg/asset/store"
	"github.com/openshift/installer/pkg/version"
	"github.com/sirupsen/logrus"
)

const (
	arch = "amd64"
)

type bundleResource interface {
}

// BundleInfo is the metadata of the bundle
type BundleInfo struct {
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	RhcosVer string `json:"rhcosVer"`
	Platform string `json:"platform"`
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

func getInstallConfig(baseDir string) asset.Asset {
	assetStore, _ := assetstore.NewStore(baseDir)
	installConfig, _ := assetStore.Load(&installconfig.InstallConfig{})
	logrus.Info("InstallConfig ", installConfig)
	return installConfig
}

// CreateBundle retrieves all dependencies for OCP installation
// and places them into a directory or compressed file.
func CreateBundle(baseDir string) {

	// if AssetStore, err := assetstore.NewStore(baseDir); err == nil {

	// 	logrus.Info("AssestStore: ")
	// 	if InstallConfig, err := AssetStore.Load(&installconfig.InstallConfig{}); err == nil && installConfig != nil {
	// 		logrus.Debug("Install-config: ", InstallConfig.Name())
	// 		logrus.Info("Platform: ", installConfig.(*installconfig.InstallConfig).Config.PullSecret)
	// 	} else {
	// 		logrus.Errorln("Install Config Load failed: ", err)
	// 	}
	// } else {
	// 	logrus.Errorln("Assest Store failed: ", err)
	// }

	// ctx, _ := context.WithTimeout(context.TODO(), 30*time.Second)
	// logrus.Info("ctx: ", ctx)
	ocpVersion, _ := version.Version()
	// logrus.Info("ocpVersion: ", ocpVersion)
	// meta, _ := rhcos.FetchRHCOSBuild(ctx, arch)
	// logrus.Info("meta: ", meta)
	// rhcosVersion := meta.OSTreeVersion
	// logrus.Info("rhcosver ", rhcosVersion)

	rhcosVersion := getRhcosVersion(ocpVersion)

	platform := getInstallConfig(baseDir).(*installconfig.InstallConfig).Config.Platform.Name()
	logrus.Info("Platform ", platform)

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
		Platform: platform,
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

	bundleInfo.downRhcos(baseDir)

	// Get binaries
	// if binary bundles hash match
	// return nil
	// else getBin()

	// Write bundle metadata
	// writeMeta()

	// Compress bundle
	createArchive(baseDir)

}

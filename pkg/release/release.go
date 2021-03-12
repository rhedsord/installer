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
	Version   string `json:"version,omitempty"`
	RootDir   string `default:"."`
	BaseDir   string
	BundleDir string
	SubTree   subTree
}

func getInstallConfig(rootDir string) asset.Asset {
	assetStore, _ := assetstore.NewStore(rootDir)
	installConfig, _ := assetStore.Load(&installconfig.InstallConfig{})
	logrus.Info("InstallConfig ", installConfig)
	return installConfig
}

// CreateBundle retrieves all dependencies for OCP installation
// and places them into a directory or compressed file.
func CreateBundle(rootDir string) {

	ocpVersion, _ := version.Version()

	rhcosVersion := getRhcosVersion(ocpVersion)

	platform := getInstallConfig(rootDir).(*installconfig.InstallConfig).Config.Platform.Name()
	logrus.Info("Platform ", platform)

	logrus.Infoln("OCP Version: ", ocpVersion)
	logrus.Infoln("RHCOS Version: ", rhcosVersion)
	logrus.Debugln("Arch: ", arch)

	bundle := &BundleRoot{
		Version: ocpVersion,
		RootDir: rootDir,
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

	// logrus.Debugln(bundle)

	// create directory tree
	// if directory tree is not built or missing subdirectories
	bundle.createTree()
	bundleInfo.writeInfo(bundle)
	// else return nil
	//

	// Get release images
	// if all image digests are present/correct
	// return nil
	// else mirrorImages()
	mirrorImages(bundle)

	// Get rhcos image
	// if rhcos image hash matches
	// return nil
	// else getRhcos()

	// bundleInfo.downRhcos(baseDir)

	// Get binaries
	// if binary bundles hash match
	// return nil
	// else getBin()

	// Write bundle metadata
	// writeMeta()

	// Compress bundle
	createArchive(bundle)

}

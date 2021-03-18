package release

import (
	"strings"

	"github.com/openshift/installer/pkg/asset/installconfig"
	assetstore "github.com/openshift/installer/pkg/asset/store"
	"github.com/openshift/installer/pkg/version"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	archAmd64 = "amd64"
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

var installConfig *installconfig.InstallConfig

func getInstallConfig(rootDir string) (*installconfig.InstallConfig, error) {
	if installConfig == nil {
		if assetStore, err := assetstore.NewStore(rootDir); err == nil {
			if asset, err := assetStore.Load(&installconfig.InstallConfig{}); asset != nil {
				installConfig = asset.(*installconfig.InstallConfig)
			} else {
				return nil, errors.Wrap(err, "Failed to load install config assest")
			}
		} else {
			return nil, errors.Wrap(err, "Failed to create assetStore")

		}

		logrus.Info("InstallConfig ", installConfig)

	}
	return installConfig, nil
}

// CreateBundle retrieves all dependencies for OCP installation
// and places them into a directory or compressed file.
func CreateBundle(rootDir string) error {

	// Remove trailing '/' from directory path
	rootDir = strings.TrimRight(rootDir, "/")

	iconfig, err := getInstallConfig(rootDir)
	if iconfig == nil || err != nil {
		logrus.Errorln("Failed to load config file")
		logrus.Error(err)
		return err
	}

	ocpVersion, _ := version.Version()

	rhcosVersion := getRhcosVersion(ocpVersion)

	platform := iconfig.Config.Platform.Name()

	logrus.Infoln("OCP Version: ", ocpVersion)
	logrus.Infoln("Platform: ", platform)
	logrus.Infoln("RHCOS Version: ", rhcosVersion)
	logrus.Infoln("Arch: ", archAmd64)

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

		Arch:     archAmd64,
		Version:  ocpVersion,
		RhcosVer: rhcosVersion,
		Platform: platform,
	}

	// create directory tree
	bundle.createTree()
	bundleInfo.writeInfo(bundle)

	// Get release images and client binaries
	if err := downloadContent(bundle); err != nil {
		werr := errors.Wrap(err, "Failed to download release image content")
		return werr
	}

	// Get rhcos image
	if err := bundleInfo.downRhcos(bundle); err != nil {
		werr := errors.Wrap(err, "Failed to download RHCOS image")
		return werr
	}

	// Compress bundle
	if err := createArchive(bundle); err != nil {
		werr := errors.Wrap(err, "Faile to create archive bundle")
		return werr
	}

	logrus.Info("Bundle creation completed.") //TODO log next steps

	return nil
}

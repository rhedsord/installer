package release

import (
	"context"
	"time"

	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/version"
)

// BundleInfo is the metadata of the bundle
type BundleInfo struct {
	Arch     string
	Version  string
	RhcosVer string
}

// subtree is the
type subTree struct {
	Dir1 string
	Dir2 string
	Dir3 string
	Info BundleInfo
}

// BundleRoot is the directory/file structure of the bundle
type BundleRoot struct {
	Version string
	SubTree subTree
}

// CreateBundle retrieves all dependencies for OCP installation
// and places them into a directory or compressed file.
func CreateBundle() {

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)

	ocpVersion, err := version.Version()
	rhcosVersion, err := rhcos.VMware(ctx, arch)

	bundle := BundleRoot{
		Version: ocpVersion,
		SubTree: subTree{
			Dir1: "client",
			Dir2: "images",
			Dir3: "rhcos",
			Info: BundleInfo{
				Arch:     "x86_64",
				Version:  ocpVersion,
				RhcosVer: rhcosVersion,
			},
		},
	}

	// create directory tree
	// if directory tree is not built or missing subdirectories
	// createTree(version.Version())
	// else return nil

	//

	// Get release images
	// if all image digests are present/correct
	// return nil
	// else mirrorImages()

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

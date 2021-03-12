package rhcos

import (
	"context"
	"net/url"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/types"
)

// QEMU fetches the URL of the Red Hat Enterprise Linux CoreOS release.
func AWS(ctx context.Context, arch types.Architecture) (string, error) {
	meta, err := FetchRHCOSBuild(ctx, arch)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch RHCOS metadata")
	}

	base, err := url.Parse(meta.BaseURI)
	if err != nil {
		return "", err
	}

	relAWS, err := url.Parse(meta.Images.AWS.Path)
	if err != nil {
		return "", err
	}

	baseURL := base.ResolveReference(relAWS).String()

	// Attach sha256 checksum to the URL.  Always provide the
	// uncompressed SHA256; the cache will take care of
	// uncompressing before checksumming.
	baseURL += "?sha256=" + meta.Images.AWS.UncompressedSHA256

	// Check that we have generated a valid URL
	_, err = url.ParseRequestURI(baseURL)
	if err != nil {
		return "", err
	}

	return baseURL, nil
}

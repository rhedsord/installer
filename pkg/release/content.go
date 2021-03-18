package release

//
import (
	"io/ioutil"
	"os/exec"

	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset/releaseimage"
)

func downloadContent(bundle *BundleRoot) error {

	logrus.Infoln("Mirroring Images")

	psfile, err := createPullSecretFile(bundle)
	if err != nil {
		return err
	}

	if image, err := releaseimage.Default(); err == nil {
		logrus.Info("Release image: ", image)

		mirrorReleaseImagesToDir(image, bundle, psfile)

		extractCLITools(image, bundle, psfile)

	} else {
		return err
	}

	return nil
}

func extractCLITools(image string, bundle *BundleRoot, psfile string) {
	cmdExtract := exec.Command("oc", "adm", "release", "extract",
		"--tools", image, "--to", bundle.BundleDir+"/"+bundle.SubTree.Client,
		"--registry-config", psfile)
	if outex, errex := cmdExtract.CombinedOutput(); errex == nil {
		logrus.Debug(string(outex))
	} else {
		logrus.Debug(string(outex))
		logrus.Warn(errex)
	}
}

func mirrorReleaseImagesToDir(image string, bundle *BundleRoot, psfile string) {
	cmdMirror := exec.Command("oc", "adm", "release", "mirror",
		image, "--to-dir", bundle.BundleDir+"/"+bundle.SubTree.Images,
		"--registry-config", psfile)

	if outmir, errmir := cmdMirror.CombinedOutput(); errmir == nil {
		logrus.Debug(string(outmir))
	} else {
		// TODO: Why does it get an exit status 1 even when it runs successfully
		logrus.Debug(string(outmir))
		logrus.Warn(errmir)
	}
}

// Read pull secret from Install Config and store to a file
func createPullSecretFile(bundle *BundleRoot) (string, error) {
	var psfile string = bundle.BaseDir + "/.pull-secret"
	var pullsecret string = "{}"

	if iconfig, err := getInstallConfig(bundle.RootDir); err == nil {
		pullsecret = iconfig.Config.PullSecret
	} else {
		return psfile, err
	}

	ferr := ioutil.WriteFile(psfile, []byte(pullsecret), 0644)
	if ferr != nil {
		logrus.Error(ferr)
		return psfile, ferr
	}
	return psfile, nil
}

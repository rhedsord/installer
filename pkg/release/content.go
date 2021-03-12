package release

//
import (
	"io/ioutil"
	"os/exec"

	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	assetstore "github.com/openshift/installer/pkg/asset/store"
)

func downloadContent(bundle *BundleRoot) {
	logrus.Infoln("Mirroring Images")

	var pullsecret string = "{}"
	var psfile string = bundle.BundleDir + "/.pull-secret"
	// var installConfig installconfig.InstallConfig
	if assetStore, err := assetstore.NewStore(bundle.RootDir); err == nil {

		logrus.Info("AssestStore: ")
		if asset, err := assetStore.Load(&installconfig.InstallConfig{}); err == nil {
			logrus.Debug("Install-config: ", asset.Name())
			installConfig := asset.(*installconfig.InstallConfig)
			pullsecret = installConfig.Config.PullSecret
		} else {
			logrus.Errorln("Install Config Load failed: ", err)
		}
	} else {
		logrus.Errorln("Assest Store failed: ", err)
	}

	ferr := ioutil.WriteFile(psfile, []byte(pullsecret), 0644)
	if ferr != nil {
		logrus.Error(ferr)
	}

	// ocpVersion, _ := version.Version()
	// logrus.Info("OCP Version: ", ocpVersion)

	if image, err := releaseimage.Default(); err == nil {
		logrus.Info("Release image: ", image)
		cmdMirror := exec.Command("oc", "adm", "release", "mirror",
			image, "--to-dir", bundle.BundleDir+"/"+bundle.SubTree.Images,
			"--registry-config", psfile)

		if outmir, errmir := cmdMirror.Output(); errmir == nil {
			logrus.Info(string(outmir))
		} else {
			// TODO: Why does it get an exit status 1 even when it runs successfully
			logrus.Info(string(outmir))
			logrus.Error(errmir)
		}

		cmdExtract := exec.Command("oc", "adm", "release", "extract",
			"--tools", image, "--to", bundle.BundleDir+"/"+bundle.SubTree.Client,
			"--registry-config", psfile)
		if outex, errex := cmdExtract.Output(); errex == nil {
			logrus.Info(string(outex))
		} else {
			// TODO: Why does it get an exit status 1 even when it runs successfully
			logrus.Info(string(outex))
			logrus.Error(errex)
		}
	}

}

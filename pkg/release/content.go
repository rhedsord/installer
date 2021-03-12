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

func mirrorImages(bundle *BundleRoot) {
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
		logrus.Info("release image: ", image)
		cmd := exec.Command("oc", "adm", "release", "mirror",
			image, "--to-dir", bundle.BundleDir+"/"+bundle.SubTree.Images,
			"--registry-config", psfile)
		out, errout := cmd.Output()
		if errout == nil {
			logrus.Info(string(out))
		} else {
			// TODO: Why does it get an exit status 1 even when it runs successfully
			logrus.Info(string(out))
			logrus.Error(errout)
		}
	}

}

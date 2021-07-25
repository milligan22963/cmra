// Pakcage version is for all version related assets
package version

import "github.com/sirupsen/logrus"

// VersionString is what is assigned during build if at all
// to allow the user to display what they are currently running
var VersionString = "development"

type VersionInstance struct {
}

func (version *VersionInstance) Run() {
	logrus.Infof("version: \n\t%s", VersionString)
}

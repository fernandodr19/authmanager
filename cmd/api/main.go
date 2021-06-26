package main

import "github.com/sirupsen/logrus"

func main() {
	logger := logrus.New()
	logger.Infof("Build info: time[%s] git_hash[%s]", BuildTime, BuildGitCommit)
}

package expand

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	tcexe "github.com/Timothee-Cardoso/tc-exe"
	"github.com/google/shlex"
)

func findSh() (string, error) {
	shPath, err := tcexe.LookPath("sh")
	if err == nil {
		return shPath, nil
	}

	if runtime.GOOS == "windows" {
		winNotFoundErr := errors.New("unable to locate sh to execute the shell alias with. The sh.exe interpreter is typically distributed with Git for Windows.")
		// We can try and find a sh executable in a Git for Windows install
		gitPath, err := tcexe.LookPath("git")
		if err != nil {
			return "", winNotFoundErr
		}

		shPath = filepath.Join(filepath.Dir(gitPath), "..", "bin", "sh.exe")
		_, err = os.Stat(shPath)
		if err != nil {
			return "", winNotFoundErr
		}

		return shPath, nil
	}

	return "", errors.New("unable to locate sh to execute shell alias with")
}

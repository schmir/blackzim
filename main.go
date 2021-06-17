package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type Runner struct {
	BlackExecutable          string
	ZimportsExecutable       string
	ZimportsWorkingDirectory string
	zimportsCmd              *exec.Cmd
	blackCmd                 *exec.Cmd
}

func startpipe(cmds []*exec.Cmd) error {
	for i := 0; i < len(cmds)-1; i++ {
		r, w, err := os.Pipe()
		if err != nil {
			return err
		}
		cmds[i].Stdout = w
		cmds[i+1].Stdin = r

	}
	return nil
}
func (runner *Runner) Init() error {
	err := runner.findExecutables()
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	runner.ZimportsWorkingDirectory = findSetupCfg(cwd)

	return nil
}

func (runner *Runner) findExecutables() error {
	black, err := exec.LookPath("black")
	if err != nil {
		return err
	}

	zimports, err := exec.LookPath("zimports")
	if err != nil {
		return err
	}
	runner.BlackExecutable = black
	runner.ZimportsExecutable = zimports
	return nil
}

func (runner *Runner) Start() error {
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	defer r.Close()
	defer w.Close()

	runner.zimportsCmd = exec.Command(runner.ZimportsExecutable, "-")
	runner.zimportsCmd.Stdin = os.Stdin
	runner.zimportsCmd.Stdout = w
	runner.zimportsCmd.Stderr = os.Stderr
	runner.zimportsCmd.Dir = runner.ZimportsWorkingDirectory

	runner.blackCmd = exec.Command(runner.BlackExecutable, "-")
	runner.blackCmd.Stdin = r
	runner.blackCmd.Stdout = os.Stdout
	runner.blackCmd.Stderr = os.Stderr

	err = runner.zimportsCmd.Start()
	if err != nil {
		return err
	}

	err = runner.blackCmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func (runner *Runner) Wait() error {
	err := runner.zimportsCmd.Wait()
	if err != nil {
		return err
	}

	return runner.blackCmd.Wait()
}

func findSetupCfg(start string) string {
	dir := start
	for {
		if _, err := os.Stat(filepath.Join(dir, "setup.cfg")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func run() error {
	runner := Runner{}

	err := runner.Init()
	if err != nil {
		return err
	}
	err = runner.Start()
	if err != nil {
		return err
	}
	return runner.Wait()
}

func main() {
	err := run()
	if err != nil {
		log.Printf("Error: %s", err)
		os.Exit(1)
	}
	os.Exit(0)
}

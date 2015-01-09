package fsobserve

import (
	"bytes"
	"github.com/go-fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	shellKey     = "SHELL"
	defaultShell = "sh"
)

type Observer struct {
	config *Config
}

type Config struct {
	Command  string
	Dir      string
	Patterns []string
	Interval time.Duration
	Shell    string
}

func NewConfig(command, dir, patterns string, interval time.Duration) *Config {
	ps := []string{}
	for _, p := range strings.Split(patterns, " ") {
		rp := strings.Trim(p, " ")
		if len(rp) == 0 {
			continue
		}
		ps = append(ps, rp)
	}
	shell := os.Getenv(shellKey)
	if shell == "" {
		shell = defaultShell
	}
	return &Config{
		Command:  command,
		Dir:      dir,
		Patterns: ps,
		Interval: interval,
		Shell:    shell,
	}
}

func New(config *Config) *Observer {
	return &Observer{
		config: config,
	}
}

func (obs *Observer) Run() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	sh, opt := obs.getShell()
	command := []string{opt, obs.config.Command}
	done := make(chan bool)
	ticker := time.NewTicker(obs.config.Interval)
	hasPatterns := false
	if len(obs.config.Patterns) > 0 {
		hasPatterns = true
	}
	go func() {
		events := []*fsnotify.Event{}
		for {
			select {
			case event := <-watcher.Events:
				if event.Op == fsnotify.Write || event.Op == fsnotify.Create {
					if (!hasPatterns && !isHiddenPath(&event)) || obs.IsUnderWatch(&event) {
						events = append(events, &event)
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			case <-ticker.C:
				if len(events) > 0 {
					obs.callback(events, sh, command)
					events = []*fsnotify.Event{}
				}
			}
		}
	}()

	err = watcher.Add(obs.config.Dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
	return nil
}

func (obs *Observer) getShell() (string, string) {
	if runtime.GOOS == "windows" {
		return "cmd", "/c"
	}
	return obs.config.Shell, "-c"
}

func (obs *Observer) IsUnderWatch(ev *fsnotify.Event) bool {
	for _, p := range obs.config.Patterns {
		if Glob(p, ev.Name) {
			return true
		}
	}
	return false
}

func (obs *Observer) callback(events []*fsnotify.Event, sh string, command []string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(sh, command...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Start()
	err := cmd.Wait()
	if err != nil {
		log.Println(err)
	}
	out, err := ioutil.ReadAll(&stdout)
	if err != nil {
		panic(err)
	}

	errOut, err := ioutil.ReadAll(&stderr)
	if err != nil {
		panic(err)
	}

	if string(errOut) == "" {
		log.Printf("exec: %v\n%v", command[1], string(out))
	} else {
		log.Printf("error: \n%v", string(errOut))
	}
}

func isHiddenPath(ev *fsnotify.Event) bool {
	return strings.HasPrefix(ev.Name, ".")
}

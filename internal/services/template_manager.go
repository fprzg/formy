package services

import (
	"bytes"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/labstack/echo/v4"
)

type TemplateManager struct {
	sync.RWMutex
	templates *template.Template
	watcher   *fsnotify.Watcher
	e         *echo.Echo
}

const UserInterfaceDir = "../../ui"

func NewTemplateManager(watchChanges bool, e *echo.Echo) (*TemplateManager, error) {
	tm := &TemplateManager{
		e: e,
	}

	if err := tm.compileTemplates(); err != nil {
		return nil, err
	}

	if watchChanges {
		if err := tm.watchTemplateChanges(); err != nil {
			return nil, err
		}
	}

	return tm, nil
}

func (tm *TemplateManager) compileTemplates() error {
	tm.Lock()
	defer tm.Unlock()

	tmpl := template.New("")

	err := filepath.Walk(UserInterfaceDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(path, ".tmpl") || strings.HasSuffix(path, ".html")) {
			_, err := tmpl.ParseFiles(path)
			if err != nil {
				tm.e.Logger.Printf("Error parsing template %s: %v\n", path, err)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	tm.templates = tmpl
	tm.e.Logger.Printf("[template_manager] Templates compiled.\n")

	return nil
}

func (tm *TemplateManager) watchTemplateChanges() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	tm.watcher = watcher

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Create == fsnotify.Create ||
					event.Op&fsnotify.Remove == fsnotify.Remove {
					tm.e.Logger.Printf("[template_manager] Change detected: %s. Recompiling...\n", event.Name)
					_ = tm.compileTemplates()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				tm.e.Logger.Printf("[template_manager] Watcher error: %v\n", err)
			}
		}
	}()

	err = filepath.Walk(UserInterfaceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})

	return nil
}

func (tm *TemplateManager) ExecuteTemplate(name string, data interface{}) (string, error) {
	tm.RLock()
	defer tm.RUnlock()

	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(&buf, name, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (tm *TemplateManager) Close() error {
	if tm.watcher != nil {
		return tm.watcher.Close()
	}
	return nil
}

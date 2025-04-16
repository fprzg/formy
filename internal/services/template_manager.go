package services

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"formy.fprzg.net/internal/models"
	"github.com/fsnotify/fsnotify"

	//"github.com/justinas/nosurf"
	"github.com/labstack/echo/v4"
)

type TemplateManager struct {
	sync.RWMutex
	templates map[string]*template.Template
	watcher   *fsnotify.Watcher
	e         *echo.Echo
}

type TemplateData struct {
	Year  int
	Flash string
	//IsAuthenticated bool
	CSRFToken       string
	Dashboard       bool
	Toast           string
	FormsData       map[string]any
	SubmissionsData map[string]any
	UserData        models.User
}

const (
	UserInterfaceDir = "../../ui"
	BaseTemplatePath = "../../ui/base.tmpl.html"
	PagesDir         = "../../ui/pages"
)

func NewTemplateData(r *http.Request) *TemplateData {
	td := &TemplateData{
		Year:  time.Now().Year(),
		Toast: "",
		//FormsData:       make(map[string]any),
		//SubmissionsData: make(map[string]any),
		//UserData:        make(map[string]any),
		//IsAuthenticated: false,

		// NOTE(Farid): Se necesita el CSRFToken para recibir la respuesta de los form
		//CSRFToken: nosurf.Token(r),
	}

	if r != nil {
		//td.Flash = app.sessionManager.PopString(r.Context(), "flash")
		//td.IsAuthenticated = app.isAuthenticated(r)
	}

	return td
}

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

	templates := make(map[string]*template.Template)

	err := filepath.Walk(PagesDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(path, ".tmpl.html") || strings.HasSuffix(path, ".html")) {
			relPath, err := filepath.Rel(PagesDir, path)
			if err != nil {
				return err
			}

			tmplName := filepath.ToSlash(relPath)

			tmpl, err := template.New("base").ParseFiles(BaseTemplatePath, path)
			if err != nil {
				tm.e.Logger.Printf("Error parsing template %s: %v\n", path, err)
				return nil
			}

			templates[tmplName] = tmpl
		}

		return nil
	})

	if err != nil {
		return err
	}

	tm.templates = templates
	tm.e.Logger.Printf("[template_manager] Templates compiled: %d\n", len(templates))

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

	tmpl, ok := tm.templates[name]
	if !ok {
		return "", fmt.Errorf("template %s not found", name)
	}

	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "base", data)
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

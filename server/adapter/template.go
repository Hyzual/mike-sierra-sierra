/*
 *   Copyright (C) 2020-2021  Joris MASSON
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU Affero General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU Affero General Public License for more details.
 *
 *   You should have received a copy of the GNU Affero General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package adapter

import (
	"fmt"
	"html/template"
	"io"
	"path"
	"path/filepath"
)

// TemplateExecutor resolves the given relative template paths from the templates/ folder,
// parses the templates from the resolved files and executes them on the given writer.
// It returns an error if a template can't be found or if the execution fails.
type TemplateExecutor interface {
	Load(writer io.Writer, data interface{}, templatePaths ...string) error
}

// templateBaseExecutor implements TemplateExecutor for production code
type templateBaseExecutor struct {
	basePath string // absolute path to the /templates directory
}

// NewTemplateExecutor creates a new TemplateExecutor
func NewTemplateExecutor(basePath string) TemplateExecutor {
	return &templateBaseExecutor{basePath}
}

// Load resolves the given relative template paths from the templates/ folder,
// parses the templates from the resolved files and executes them on the given writer.
// It returns an error if a template can't be found or if the execution fails.
func (t *templateBaseExecutor) Load(writer io.Writer, data interface{}, templatePaths ...string) error {
	var cleanedPaths []string
	for _, templatePath := range templatePaths {
		cleanedPaths = append(cleanedPaths, path.Join(t.basePath, filepath.Clean(templatePath)))
	}

	tmpl, err := template.ParseFiles(cleanedPaths...)
	if err != nil {
		return fmt.Errorf("could not load the templates %v: %w", templatePaths, err)
	}
	err = tmpl.Execute(writer, data)
	if err != nil {
		return fmt.Errorf("could not execute the templates %v: %w", templatePaths, err)
	}
	return nil
}

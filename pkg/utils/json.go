// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package utils

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
)

func requirePointer(v any) error {
	if reflect.ValueOf(v).Kind() != reflect.Pointer {
		return Err("argument must be a pointer")
	}
	return nil
}

func WriteJSON(path string, data any) error {
	if path == "" {
		return Err("path cannot be empty")
	}

	if err := requirePointer(data); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return Wrapf(err, "failed to create directory for %s", path)
	}

	tmpPath := path + ".tmp"

	f, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return Wrapf(err, "failed to create temp file %s", tmpPath)
	}

	encErr := func() error {
		defer f.Close()

		bw := bufio.NewWriter(f)
		enc := json.NewEncoder(bw)
		enc.SetIndent("", "  ")

		if err := enc.Encode(data); err != nil {
			return Wrap("failed to encode JSON", err)
		}

		if err := bw.Flush(); err != nil {
			return Wrap("failed to flush buffer", err)
		}

		if err := f.Sync(); err != nil {
			return Wrap("fsync failed", err)
		}

		return nil
	}()

	if encErr != nil {
		_ = os.Remove(tmpPath)
		return encErr
	}

	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return Wrap("atomic rename failed", err)
	}

	return nil
}

func ReadJSON(path string, data any) error {
	if path == "" {
		return Err("path cannot be empty")
	}

	if err := requirePointer(data); err != nil {
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		return Wrap("failed to open JSON file", err)
	}
	defer func() { _ = f.Close() }()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()

	if err := dec.Decode(data); err != nil {
		return Wrap("failed to decode JSON file", err)
	}

	return nil
}

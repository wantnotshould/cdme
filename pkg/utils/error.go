// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package utils

import (
	"errors"
	"fmt"
	"strings"
)

func Err(parts ...string) error {
	switch len(parts) {
	case 0:
		return errors.New("unknown error")
	case 1:
		return errors.New(parts[0])
	default:
		return errors.New(strings.Join(parts, " "))
	}
}

func Wrap(context string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", context, err)
}

func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	context := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s: %w", context, err)
}

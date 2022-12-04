package service

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"text/template"
)

func buildURL(urlT string, address addressInfo) (string, error) {
	tpl, err := template.New("URL").Parse(urlT)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	if err := tpl.Execute(&b, address); err != nil {
		return "", err
	}
	return b.String(), nil
}

func buildCommand(ctx context.Context, commandT []string, address addressInfo) (*exec.Cmd, error) {
	if len(commandT) == 0 {
		return nil, errors.New("empty command")
	}
	var command []string
	for _, v := range commandT {
		tpl, err := template.New("command").Parse(v)
		if err != nil {
			return nil, err
		}
		var b bytes.Buffer
		if err := tpl.Execute(&b, address); err != nil {
			return nil, err
		}
		command = append(command, b.String())
	}
	return exec.CommandContext(ctx, command[0], command[1:]...), nil
}

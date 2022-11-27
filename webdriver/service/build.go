package service

import (
	"bytes"
	"errors"
	"os/exec"
	"text/template"
)

func buildURL(urlT string, address addressInfo) (string, error) {
	urlTemplate, err := template.New("URL").Parse(urlT)
	if err != nil {
		return "", err
	}
	urlBuffer := &bytes.Buffer{}
	if err := urlTemplate.Execute(urlBuffer, address); err != nil {
		return "", err
	}
	return urlBuffer.String(), nil
}

func buildCommand(commandT []string, address addressInfo) (*exec.Cmd, error) {
	if len(commandT) == 0 {
		return nil, errors.New("empty command")
	}

	var command []string
	for _, argument := range commandT {
		argTemplate, err := template.New("command").Parse(argument)
		if err != nil {
			return nil, err
		}
		var argBuffer bytes.Buffer
		if err := argTemplate.Execute(&argBuffer, address); err != nil {
			return nil, err
		}
		command = append(command, argBuffer.String())
	}
	return exec.Command(command[0], command[1:]...), nil
}

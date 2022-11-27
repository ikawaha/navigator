package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Service struct {
	mu       sync.Mutex
	URLT     string   // url template eg. "http://localhost:{{.Port}}"
	CommandT []string // command template eg. ["chromedriver", "--port={{.Port}}"]
	baseURL  string
	command  *exec.Cmd
}

func (s *Service) URL() string {
	return s.baseURL
}

func (s *Service) Start(debug bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.command != nil {
		return errors.New("already running")
	}

	address, err := getFreeAddress()
	if err != nil {
		return fmt.Errorf("failed to locate a free port: %s", err)
	}

	url, err := buildURL(s.URLT, address)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}
	s.baseURL = url

	command, err := buildCommand(s.CommandT, address)
	if err != nil {
		return fmt.Errorf("failed to parse command: %w", err)
	}
	if debug {
		os.Stderr.WriteString(command.String() + "\n")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
	}
	if err := command.Start(); err != nil {
		err = fmt.Errorf("failed to run command: %w", err)
		if debug {
			os.Stderr.WriteString("ERROR: " + err.Error() + "\n")
		}
		return err
	}
	s.command = command
	return nil
}

func (s *Service) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.command == nil {
		return errors.New("already stopped")
	}

	switch runtime.GOOS {
	case "windows":
		if err := s.command.Process.Kill(); err != nil {
			return fmt.Errorf("failed to stop command: %w", err)
		}
	default:
		if err := s.command.Process.Signal(syscall.SIGTERM); err != nil {
			return fmt.Errorf("failed to stop command: %w", err)
		}
	}
	s.command.Wait()
	s.command = nil
	s.baseURL = ""
	return nil
}

type addressInfo struct {
	Address string
	Host    string
	Port    string
}

func getFreeAddress() (addressInfo, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return addressInfo{}, err
	}
	defer listener.Close()

	address := listener.Addr().String()
	addressParts := strings.SplitN(address, ":", 2)
	return addressInfo{
		Address: address,
		Host:    addressParts[0],
		Port:    addressParts[1],
	}, nil
}

func (s *Service) WaitForBoot(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	wakeup := make(chan struct{})
	go func(ctx context.Context) {
		up := s.checkStatus()
		for !up {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(500 * time.Millisecond)
				up = s.checkStatus()
			}
		}
		wakeup <- struct{}{}
	}(ctx)
	select {
	case <-ctx.Done():
		return errors.New("failed to start before timeout")
	case <-wakeup:
		return nil
	}
}

func (s *Service) checkStatus() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	req, err := http.NewRequest(http.MethodGet, s.baseURL+"/status", nil)
	if err != nil {
		return false
	}
	client := &http.Client{
		Timeout: time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	return resp.StatusCode == 200
}

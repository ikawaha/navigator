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
	"sync"
	"syscall"
	"time"
)

// Service represents a web driver service.
type Service struct {
	mu       sync.Mutex
	urlT     string   // url template eg. "http://localhost:{{.Port}}"
	commandT []string // command template eg. ["chromedriver", "--port={{.Port}}"]
	baseURL  string
	command  *exec.Cmd
}

// New creates new web driver service.
func New(urlT string, commandT []string) *Service {
	return &Service{
		urlT:     urlT,
		commandT: commandT,
	}
}

// URL returns the base URL of the service.
func (s *Service) URL() string {
	return s.baseURL
}

// Start starts the service.
func (s *Service) Start(debug bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.command != nil {
		return errors.New("already running")
	}

	address, err := getFreeAddress()
	if err != nil {
		return fmt.Errorf("failed to locate a free port: %w", err)
	}

	url, err := buildURL(s.urlT, address)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}
	s.baseURL = url

	command, err := buildCommand(s.commandT, address)
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

// Stop stops the service.
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
	_ = s.command.Wait()
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
	var lc net.ListenConfig
	l, err := lc.Listen(context.TODO(), "tcp", "localhost:0")
	if err != nil {
		return addressInfo{}, err
	}
	defer l.Close()

	address := l.Addr().String()
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return addressInfo{}, err
	}
	return addressInfo{
		Address: address,
		Host:    host,
		Port:    port,
	}, nil
}

const bootWait = 500 * time.Millisecond

// WaitForBoot waits until the service starts.
func (s *Service) WaitForBoot(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	wakeup := make(chan struct{})
	go func(ctx context.Context) {
		up := s.checkStatus(ctx)
		for !up {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(bootWait)
				up = s.checkStatus(ctx)
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

func (s *Service) checkStatus(ctx context.Context) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/status", nil)
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
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	return resp.StatusCode == 200
}

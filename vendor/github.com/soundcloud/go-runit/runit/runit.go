package runit

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
	"time"
)

const (
	defaultServiceDir = "/etc/service"

	taiOffset = 4611686018427387914
	statusLen = 20

	posTimeStart = 0
	posTimeEnd   = 7
	posPidStart  = 12
	posPidEnd    = 15

	posWant  = 17
	posState = 19

	StateDown   = 0
	StateUp     = 1
	StateFinish = 2
)

var (
	ENoRunsv      = errors.New("runsv not running")
	StateToString = map[int]string{
		StateDown:   "down",
		StateUp:     "up",
		StateFinish: "finish",
	}
)

type SvStatus struct {
	Pid        int
	Duration   int
	Timestamp  time.Time
	State      int
	NormallyUp bool
	Want       int
}

type service struct {
	Name       string
	ServiceDir string
}

func GetServices(dir string) ([]*service, error) {
	if dir == "" {
		dir = defaultServiceDir
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	services := []*service{}
	for _, file := range files {
		if file.Mode()&os.ModeSymlink == os.ModeSymlink || file.IsDir() {
			services = append(services, GetService(file.Name(), dir))
		}
	}
	return services, nil
}

func GetService(name string, dir string) *service {
	if dir == "" {
		dir = defaultServiceDir
	}
	r := service{Name: name, ServiceDir: dir}
	return &r
}

func (s *service) file(file string) string {
	return fmt.Sprintf("%s/%s/supervise/%s", s.ServiceDir, s.Name, file)
}

func (s *service) runsvRunning() (bool, error) {
	file, err := os.OpenFile(s.file("ok"), os.O_WRONLY, 0)
	if err != nil {
		if err == syscall.ENXIO {
			return false, nil
		}
		return false, err
	}
	file.Close()
	return true, nil
}

func (s *service) status() ([]byte, error) {
	file, err := os.Open(s.file("status"))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	status := make([]byte, statusLen)
	_, err = file.Read(status)
	return status, err
}

func (s *service) NormallyUp() bool {
	_, err := os.Stat(s.file("down"))
	return err != nil
}

func (s *service) Status() (*SvStatus, error) {
	status, err := s.status()
	if err != nil {
		return nil, err
	}

	var pid int
	pid = int(status[posPidEnd])
	for i := posPidEnd - 1; i >= posPidStart; i-- {
		pid <<= 8
		pid += int(status[i])
	}

	tai := int64(status[posTimeStart])
	for i := posTimeStart + 1; i <= posTimeEnd; i++ {
		tai <<= 8
		tai += int64(status[i])
	}
	state := status[posState] // 0: down, 1: run, 2: finish

	tv := &syscall.Timeval{}
	if err := syscall.Gettimeofday(tv); err != nil {
		return nil, err
	}
	sS := SvStatus{
		Pid:        pid,
		Timestamp:  time.Unix(tai-taiOffset, 0), // FIXME: do we just select the wrong slice?
		Duration:   int(int64(tv.Sec) - (tai - taiOffset)),
		State:      int(state),
		NormallyUp: s.NormallyUp(),
	}

	switch status[posWant] {
	case 'u':
		sS.Want = StateUp
	case 'd':
		sS.Want = StateDown
	}

	return &sS, nil
}

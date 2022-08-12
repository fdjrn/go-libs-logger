package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *log.Logger

type LumberjackLogger struct {
	LogPath       string
	LastLogDate   string
	CompressLog   bool
	DailyRotate   bool
	LogToTerminal bool
	LumberjackLog *lumberjack.Logger
}

// SetLog :
func (p *LumberjackLogger) SetLog() error {
	if p.LogPath == "" {
		return fmt.Errorf("You need to specify log path in path flag")
	}

	// Get Log Directory
	splitLogDir := strings.Split(p.LogPath, "/")
	logDir := strings.Join(splitLogDir[:len(splitLogDir)-1], "/")

	// Create directory if not exists
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	e, err := os.OpenFile(p.LogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	Log = log.New(e, "", log.Ldate|log.Ltime)
	p.LumberjackLog = &lumberjack.Logger{
		Filename:  p.LogPath,
		Compress:  p.CompressLog,
		LocalTime: true,
	}

	if !p.LogToTerminal {
		Log.SetOutput(p.LumberjackLog)
	} else {
		mw := io.MultiWriter(os.Stdout, p.LumberjackLog)
		Log.SetOutput(mw)
	}

	p.LastLogDate = time.Now().Format("2006-01-02")

	if p.DailyRotate == true {
		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			p.LogRotator()
		}()

	}

	return nil
}

// LogRotator :
func (p *LumberjackLogger) LogRotator() {
	for {
		// If Lumberjacklog not set then skip
		if p.LumberjackLog == nil {
			continue
		}
		// Rotate Log If LastLog Date != Current Date
		if p.LastLogDate != time.Now().Format("2006-01-02") {
			// Set LastLogDate to Current Date
			p.LastLogDate = time.Now().Format("2006-01-02")
			p.LumberjackLog.Rotate()
			log.Println("| Log Rotated")
		}
		// Sleep every 5 seconds
		time.Sleep(time.Duration(5) * time.Second)
	}
}

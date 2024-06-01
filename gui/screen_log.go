package gui

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
	"golang.org/x/crypto/ssh"
)

const (
	SEKAI_LOG  = "Sekai"
	INTERX_LOG = "Interx"
	SHIDAI_LOG = "Shidai"
)

var cancelLogScreenDataBinding binding.DataListener

var (
	ansiEscape              = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	sekaiLogArtifactRemove  = regexp.MustCompile(`^[^\d]+(\d{2}:\d{2}:\d{2})[^\s]*\s`)
	interxLogArtifactRemove = regexp.MustCompile(`^.*?interx\[\d+\]:\s*`)
	shidaiLogArtifactRemove = regexp.MustCompile(`^.*?interx\[\d+\]:\s*|^.*?\[GIN\]\s*`)
)

func makeLogScreen(_ fyne.Window, g *Gui) fyne.CanvasObject {

	sekaiTab := container.NewTabItem("Sekai",
		makeLogTab(g, "http://localhost:8282/logs/sekai", SEKAI_LOG),
	)
	interxTab := container.NewTabItem("Interx",
		makeLogTab(g, "http://localhost:8282/logs/interx", INTERX_LOG),
	)
	shidaiTab := container.NewTabItem("Shidai",
		makeLogTab(g, "http://localhost:8282/logs/shidai", SHIDAI_LOG),
	)

	tabsMenu := container.NewAppTabs(sekaiTab, interxTab, shidaiTab)
	tabsMenu.OnUnselected = func(ti *container.TabItem) {
		log.Println("Unselected:", ti)
		g.LogCtxCancel()
		g.LogCtx, g.LogCtxCancel = context.WithCancel(context.Background())
	}

	return tabsMenu
}

func makeLogTab(g *Gui, addressToListen string, logType string) fyne.CanvasObject {
	g.LogCtx, g.LogCtxCancel = context.WithCancel(context.Background())

	cancelAndCreateFunc := func() {
		log.Printf("cancelAndCreateFunc")
		g.LogCtxCancel()

	}

	logTextLabel := widget.NewLabel("")
	logTextLabel.Wrapping = fyne.TextWrapBreak
	logTextLabel.Refresh()
	logScroll := container.NewVScroll(
		logTextLabel,
	)

	switchData := binding.NewBool()
	switchData.Set(true)
	startStopButton := widget.NewButton("Start", func() {})

	cancelLogScreenDataBinding = binding.NewDataListener(func() {

		log.Println("data changed")
		cancelAndCreateFunc()
	})
	startStopFunc := func() {
		resp, err := createTunnelForSSEConnection(g.sshClient, addressToListen)
		if err != nil {
			log.Println(err)
			return
		}
		state, _ := switchData.Get()
		log.Println("state:", state)
		if state {
			startStopButton.Text = "Stop"
			switchData.Set(false)
			g.LogCtx, g.LogCtxCancel = context.WithCancel(context.Background())
			go runLogScreenV2(resp, g.LogCtx, logTextLabel, logType)
		} else {
			switchData.Set(true)
			startStopButton.Text = "Start"
			cancelLogScreenDataBinding.DataChanged()
		}
		startStopButton.Refresh()

	}
	startStopButton.OnTapped = startStopFunc

	copyButton := widget.NewButtonWithIcon("Copy", theme.FileIcon(), func() {
		data := logTextLabel.Text

		err := clipboard.WriteAll(data)
		if err != nil {
			log.Println(err)
			return
		}
	})
	return container.NewBorder(copyButton, startStopButton, nil, nil, logScroll)

}

func runLogScreenV2(resp *http.Response, cancelCtx context.Context, infoLabel *widget.Label, logType string) {
	var writeData string
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)

	go func(ctx context.Context) {
		reader := bufio.NewReader(resp.Body)
		defer resp.Body.Close()
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Println("Context cancelled, exiting log reader")
				return
			default:
				// Set a short timeout to periodically check the context
				lineCh := make(chan string)
				errCh := make(chan error)
				go func() {
					line, err := reader.ReadString('\n')
					if err != nil {
						errCh <- err
					} else {
						lineCh <- line
					}
				}()
				select {
				case <-ctx.Done():
					log.Println("Context cancelled, exiting log reader")
					return
				case line := <-lineCh:
					line = stripANSI(line)
					switch logType {
					case SEKAI_LOG:
						// line = sekaiRemoveLogPrefix(line)
					case INTERX_LOG:
						line = interxRemoveLogPrefix(line)
					case SHIDAI_LOG:
						line = shidaiRemoveLogPrefix(line)
					}
					mu.Lock()
					writeData += line
					writeData = truncateString(writeData, 12000)
					mu.Unlock()

				case err := <-errCh:
					if err == io.EOF {
						return
					}
					log.Printf("Failed to read response body: %v", err)
					return
				}
			}
		}
	}(cancelCtx)

	go func(ctx context.Context) {
		defer wg.Done()
		var previousData string
		for {
			select {
			case <-ctx.Done():
				log.Println("Context cancelled, exiting log writer")
				return
			default:
				mu.Lock()
				data := writeData
				mu.Unlock()
				if previousData != writeData {
					previousData = writeData
					infoLabel.SetText(data)
					log.Println("infoLabel updated")
				}
				time.Sleep(time.Second * 1)
				// time.Sleep(time.Second * 1)
			}
		}
	}(cancelCtx)

	wg.Wait()
}

func sekaiRemoveLogPrefix(s string) string {
	return sekaiLogArtifactRemove.ReplaceAllString(s, "")
}
func interxRemoveLogPrefix(s string) string {
	return interxLogArtifactRemove.ReplaceAllString(s, "")
}
func shidaiRemoveLogPrefix(s string) string {
	return shidaiLogArtifactRemove.ReplaceAllString(s, "")
}

func stripANSI(s string) string {
	return ansiEscape.ReplaceAllString(s, "")
}

func truncateString(s string, maxSize int) string {
	if len(s) > maxSize {
		excessLength := len(s) - maxSize
		truncatePos := strings.Index(s[excessLength:], "\n")
		if truncatePos != -1 {
			return s[excessLength+truncatePos+1:]
		}
		return ""
	}
	return s
}

func createTunnelForSSEConnection(sshClient *ssh.Client, address string) (*http.Response, error) {
	dialer := func(network, addr string) (net.Conn, error) {
		return sshClient.Dial(network, addr)
	}

	httpTransport := &http.Transport{
		Dial: dialer,
	}

	httpClient := &http.Client{
		Transport: httpTransport,
	}

	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		log.Printf("Failed to create HTTP request: %v", err)
		return nil, err
	}
	req.Header.Set("Accept", "text/event-stream")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Failed to send HTTP request: %v", err)
		return nil, err

	}

	return resp, nil
}

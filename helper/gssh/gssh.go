package gssh

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"golang.org/x/crypto/ssh"
)

type ResultV2 struct {
	Err error
}

func MakeSHH_ClientWithPassword(ipAndPort, user, psswrd string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(psswrd),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the SSH server
	client, err := ssh.Dial("tcp", ipAndPort, config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func MakeSSH_ClientWithPrivKey(ipAndPort, user string, key []byte) (*ssh.Client, error) {
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", ipAndPort, config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func MakeSSH_ClientWithPrivKeyAndPassphrase(ipAndPort, user string, key, passphrase []byte) (*ssh.Client, error) {
	signer, err := ssh.ParsePrivateKeyWithPassphrase(key, passphrase)
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", ipAndPort, config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func CheckIfPassphraseNeeded(privateKeyBytes []byte) (bool, error) {
	_, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		if _, ok := err.(*ssh.PassphraseMissingError); ok {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func ExecuteSSHCommandV2(client *ssh.Client, command string, outputChan chan<- string, resultChan chan<- ResultV2) {
	log.Printf("RUNNING CMD:\n%s", command)
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		resultChan <- ResultV2{Err: err}

	}()
	defer close(outputChan)
	defer cancel() // This will be called if an early return occurs
	session, err := client.NewSession()
	if err != nil {
		if err == io.EOF {
			err = fmt.Errorf("ssh EOF, probably ssh server was down, please restart Kensho: %w", err)
		}
		log.Println("Error when creating new session: ", err.Error())
		// cancel()
		// close(outputChan)
		// resultChan <- ResultV2{Err: err}
		return
	}
	defer session.Close()

	// Setting up stdout and stderr pipes
	stdoutPipe, err := session.StdoutPipe()
	if err != nil {
		log.Println("Error when creating stdoutPipe: ", err.Error())
		// cancel()
		// close(outputChan)
		// resultChan <- ResultV2{Err: err}
		return
	}
	stderrPipe, err := session.StderrPipe()
	if err != nil {
		log.Println("Error when creating stderrPipe: ", err.Error())
		// cancel()
		// close(outputChan)
		// resultChan <- ResultV2{Err: err}
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)
	// Start the command
	err = session.Start(command)
	if err != nil {
		log.Println("Error when starting new session: ", err.Error())
		// cancel()
		// close(outputChan) // Close the channel on error
		// resultChan <- ResultV2{Err: err}
		return
	}

	// Read from stdout and stderr concurrently
	go streamOutput(ctx, stdoutPipe, outputChan, &wg)
	go streamOutput(ctx, stderrPipe, outputChan, &wg)

	// go monitorConnection(ctx, client.Conn, outputChan, &wg)
	err = session.Wait()
	if err != nil {
		// cancel()
		// close(outputChan)
		// resultChan <- ResultV2{Err: err}
		return
	}
	wg.Wait()
	// cancel()
	// close(outputChan) // Close the channel when done

}

func streamOutput(ctx context.Context, reader io.Reader, outputChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(reader)
	for {
		select {
		case <-ctx.Done():
			return // Exit if context is cancelled
		default:
			if scanner.Scan() {
				select {
				case outputChan <- scanner.Text():
					log.Println(scanner.Text())
				case <-ctx.Done():
					return
				}
			} else {
				return
			}
		}
	}
}

func ExecuteSSHCommandV3(client *ssh.Client, command string, outputChan chan<- string, resultChan chan<- ResultV2) {
	session, err := client.NewSession()
	if err != nil {
		if err == io.EOF {
			err = fmt.Errorf("ssh EOF, probably ssh server was down, please restart Kensho: %w", err)
		}
		log.Println("Error when creating new session: ", err.Error())

		close(outputChan)
		resultChan <- ResultV2{Err: err}
		return
	}
	defer session.Close()

	o, err := session.CombinedOutput(command)

	// outputChan <- "Command was executed successfully"
	outputChan <- string(o)
	resultChan <- ResultV2{Err: err}
	close(outputChan)
}

func MakeSSHsessionForTerminal(client *ssh.Client) (*ssh.Session, error) {
	// Create a session
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, err
	}

	// Request a pty (pseudo-terminal)
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // Enable echoing
		ssh.TTY_OP_ISPEED: 14400, // Input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // Output speed = 14.4kbaud
	}

	if err := session.RequestPty("ansi", 80, 40, modes); err != nil {
		session.Close()
		client.Close()
		return nil, err
	}

	return session, nil
}

// for testing
func MakeSSHsessionForTerminalV2(client *ssh.Client) (*ssh.Session, error) {
	// Create a session
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, err
	}

	// Request a pty (pseudo-terminal)
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // Enable echoing
		ssh.TTY_OP_ISPEED: 14400, // Input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // Output speed = 14.4kbaud
	}

	if err := session.RequestPty("ansi", 80, 40, modes); err != nil {
		session.Close()
		client.Close()
		return nil, err
	}

	return session, nil
}

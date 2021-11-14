package ssh

import (
	"Cube-back/log"
	"Cube-back/models/common/configure"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"path"
	"time"
)

var SftpClient *sftp.Client

var config = map[string]interface{}{}

type Ssh struct {
	FileServerIp       string
	FileServerPort     int
	FileServerUser     string
	FileServerPassword string
}

func connect(user, password, host string, port int) (*ssh.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //ssh.FixedHostKey(hostKey),
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	return sshClient, nil
}

func sftpClientCreate(sshClient *ssh.Client) (*sftp.Client, error) {
	// create sftp client
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}
	return sftpClient, nil
}

func sessionCreate(sshClient *ssh.Client) (*ssh.Session, error) {
	var session *ssh.Session
	session, err := sshClient.NewSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func CommandExecute(command string) {
	session, err := createSession()
	if err != nil {
		log.Error(err)
	}
	defer session.Close()
	session.Run(command)
}

func UploadFile(filename string, remotePath string, data []uint8) bool {
	sftpClient, err := createSftpClient()
	if err != nil {
		log.Error(err)
		return false
	}
	defer sftpClient.Close()
	err = sftpClient.MkdirAll(remotePath)
	if err != nil {
		log.Error(err)
		return false
	}
	dstFile, err := sftpClient.Create(path.Join(remotePath, filename))
	if err != nil {
		log.Error(err)
		return false
	}
	_, error := dstFile.Write(data)
	if error != nil {
		log.Error(error)
		return false
	}
	dstFile.Close()
	return true
}

func RemoveFile(path string) {
	sftpClient, err := createSftpClient()
	if err != nil {
		log.Error(err)
	}
	defer sftpClient.Close()
	err = sftpClient.Remove(path)
	if err != nil {
		log.Error(err)
	}
}

func RemoveDirectory(remotePath string) *sftp.Client {
	sftpClient, err := createSftpClient()
	if err != nil {
		log.Error(err)
	}
	defer sftpClient.Close()
	remoteFiles, err := sftpClient.ReadDir(remotePath)
	if err != nil {
		return sftpClient
	}
	defer sftpClient.RemoveDirectory(remotePath)
	for _, backupDir := range remoteFiles {
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		if backupDir.IsDir() {
			RemoveDirectory(remoteFilePath)
		} else {
			sftpClient.Remove(path.Join(remoteFilePath))
		}
	}
	return sftpClient
}

func createSftpClient() (*sftp.Client, error) {
	sshClient, err := connect(
		config["FileServerUser"].(string),
		config["FileServerPassword"].(string),
		config["FileServerIp"].(string),
		config["FileServerPort"].(int))
	if sshClient != nil {
		return sftpClientCreate(sshClient)
	} else {
		return nil, err
	}
}

func createSession() (*ssh.Session, error) {
	sshClient, err := connect(
		config["FileServerUser"].(string),
		config["FileServerPassword"].(string),
		config["FileServerIp"].(string),
		config["FileServerPort"].(int))
	if sshClient != nil {
		return sessionCreate(sshClient)
	} else {
		return nil, err
	}
}

func init() {
	conf := new(Ssh)
	configure.Get(&conf)
	config["FileServerUser"] = conf.FileServerUser
	config["FileServerPassword"] = conf.FileServerPassword
	config["FileServerIp"] = conf.FileServerIp
	config["FileServerPort"] = conf.FileServerPort
}

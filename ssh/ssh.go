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

type Ssh struct {
	FileServerIp       string
	FileServerPort     int
	FileServerUser     string
	FileServerPassword string
}

func connect(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
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

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}
	return sftpClient, nil
}

func UploadFile(filename string, remotePath string, data []uint8) bool {
	sftpClient, err := createInstance()
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
	sftpClient, err := createInstance()
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
	sftpClient, err := createInstance()
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

func createInstance() (*sftp.Client, error) {
	conf := new(Ssh)
	configure.Get(&conf)
	sftpClient, err := connect(conf.FileServerUser, conf.FileServerPassword, conf.FileServerIp, conf.FileServerPort)
	return sftpClient, err
}

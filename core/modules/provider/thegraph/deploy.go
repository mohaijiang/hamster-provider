package thegraph

import (
	"github.com/ThomasRooney/gexpect"
	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

// InstallDocker install docker
func InstallDocker() error {
	pwd := "12345678"
	//Determine whether to save and install docker
	cmd := exec.Command("docker", "-v")
	_, err := cmd.CombinedOutput()
	if err != nil {
		//install docker
		curlCmd := exec.Command("curl", "-fsSL", "https://get.docker.com", "-o", "get-docker.sh")
		curlErr := curlCmd.Run()
		if curlErr != nil {
			log.Printf("cmd.Run() failed with %s\n", curlErr)
			return curlErr
		}
		installCmd := "sh 'get-docker.sh'"
		//exec install docker
		child, installErr := gexpect.Spawn(installCmd)
		if installErr != nil {
			log.Printf("cmd.Run() failed with %s\n", installErr)
			return installErr
		}
		//input sudo password
		sendErr := child.SendLine(pwd)
		if sendErr != nil {
			log.Printf("password error")
			return sendErr
		}
		waitErr := child.Wait()
		if waitErr != nil {
			log.Printf("cmd.Run() failed with %s\n", waitErr)
			return waitErr
		}
		//start docker
		startCmd := "sudo service docker start"
		//startCmd := exec.Command("systemctl","start","docker")
		//startErr := startCmd.Run()
		childStart, startErr := gexpect.Spawn(startCmd)
		if startErr != nil {
			log.Printf("cmd.Run() failed with %s\n", startErr)
			return startErr
		}
		//input password
		startSendErr := childStart.SendLine(pwd)
		if startSendErr != nil {
			log.Printf("start pws error")
			return startSendErr
		}
		startWaitErr := childStart.Wait()
		if startWaitErr != nil {
			log.Printf("start wait error")
			return startWaitErr
		}
		return nil
	}
	//start docker
	startCmd := "sudo service docker start"
	//startCmd := exec.Command("systemctl","start","docker")
	//startErr := startCmd.Run()
	childStart, startErr := gexpect.Spawn(startCmd)
	if startErr != nil {
		log.Printf("cmd.Run() installed start failed with %s\n", startErr)
		return startErr
	}
	startSendErr := childStart.SendLine(pwd)
	if startSendErr != nil {
		log.Printf("installed start password error")
		return startSendErr
	}
	return nil
}

// TemplateInstance Docker compose file instantiation
func templateInstance(data DeployParams) error {

	t, err := template.ParseFiles("/home/mohaijiang/IdeaProject/hamster-provider/templates/graph-docker-compose.text")
	if err != nil {
		log.Printf("template failed with %s\n", err)
		return err
	}
	//create file in .hamster-provider
	url := strings.Join([]string{config.DefaultConfigDir(), "docker-compose.yml"}, string(os.PathSeparator))
	file, createErr := os.Create(url)
	if createErr != nil {
		log.Printf("create file failed %s\n", err)
		return createErr
	}
	writeErr := t.Execute(file, data)
	if writeErr != nil {
		log.Printf("template write file failed %s\n", err)
		return writeErr
	}
	return nil
}

// StartDockerCompose exec docker-compose
func startDockerCompose() error {
	cmd := exec.Command("docker", "compose", "up", "-d")
	println(config.DefaultConfigDir())
	cmd.Dir = config.DefaultConfigDir()
	return cmd.Run()
}

// StopDockerCompose  停止docker compose 服务
func stopDockerCompose() error {
	cmd := exec.Command("docker", "compose", "down", "-v")
	println(config.DefaultConfigDir())
	cmd.Dir = config.DefaultConfigDir()
	return cmd.Run()
}

type ComposeStatus int

const (
	STOP = 1
)

func composeStatus() ComposeStatus {

	return STOP
}

package myssh

import (
	"fmt"
	"os"
	"testing"
)

func TestClient_Run(t *testing.T) {
	c, err := NewClient(os.Getenv("HOST"), os.Getenv("USER"), os.Getenv("PASSWORD"), 22)
	if err != nil {
		t.Error("new client err:", err)
		return
	}
	defer c.Close()

	//if output, err := c.Run("ls"); err != nil {
	//	t.Error("run ls err:", err)
	//	return
	//} else {
	//	t.Log("ls :", string(output))
	//}

	c.Runs(func(cmd string, output []byte, err error) bool {
		if err != nil {
			fmt.Printf("%s: err \n %s", cmd, err.Error())
			return false
		}

		fmt.Printf("%s:\n %s", cmd, string(output))
		return true
	}, "ls", "systemctl list-units | grep node", "docker ps -a")
}

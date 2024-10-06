package main

import (
	"fmt"
	"net"
	"os/exec"
	"time"
)

func log(message string) {
	fmt.Printf("%s - %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
}

func checkInternet() bool {
	_, err := net.DialTimeout("tcp", "1.1.1.1:53", time.Second*5)
	return err == nil
}

func restartWiFi() error {
	cmds := []struct {
		name string
		args []string
	}{

		{"networksetup", []string{"-setairportpower", "en0", "off"}},
		{"sleep", []string{"5"}},
		{"networksetup", []string{"-setairportpower", "en0", "on"}},
	}

	for _, cmd := range cmds {
		output, err := exec.Command(cmd.name, cmd.args...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s実行エラー: %v, 出力: %s", cmd.name, err, output)
		}
	}
	return nil

}

func main() {
	log("Wi-Fi監視スクリプトを開始しました")

	for {
		if !checkInternet() {
			log("インターネット接続が切断されました、Wi-Fiを再起動します")

			err := restartWiFi()
			if err != nil {
				log(fmt.Sprintf("Wi-Fi再起動エラー: %v", err))
			} else {
				log("Wi-Fiを再起動しました、再接続を確認中...")
			}

			for i := 1; i < 5; i++ {
				time.Sleep(10 * time.Second)
				if checkInternet() {
					log("インターネット接続が復旧しました")
					break
				} else if i == 5 {
					log("インターネット接続の復旧に失敗しました、次の確認まで待機します")
				} else {
					log(fmt.Sprintf("再接続を確認中... (試行%d/5)", i))
				}
			}
		} else {
			log("インターネット接続は正常です")
		}
		time.Sleep(15 * time.Second)
	}
}

/**
 * @Author: Resynz
 * @Date: 2022/1/6 19:01
 */
package util

import (
	"fmt"
	"os"
)

func FormatDir(df string) error {
	f,err := os.Stat(df)
	if err!=nil{
		if os.IsNotExist(err) {
			return os.MkdirAll(df, 0755)
		}
		return err
	}
	if !f.IsDir() {
		return fmt.Errorf("invalid path:%s. It should be a dir",df)
	}
	return nil
}

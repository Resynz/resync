/**
 * @Author: Resynz
 * @Date: 2022/1/7 15:13
 */
package util

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"sync"
)

func RunCommand(ctx context.Context,writer io.Writer, cmdPath,command string, args ...string) error {
	var cmd *exec.Cmd
	if ctx == nil {
		cmd=exec.Command(command,args...)
	}else {
		cmd=exec.CommandContext(ctx,command,args...)
	}
	cmd.Dir = cmdPath
	stderr,err:=cmd.StderrPipe()
	if err!=nil{
		return err
	}
	stdout,err:=cmd.StdoutPipe()
	if err!=nil{
		return err
	}
	if err = cmd.Start();err!=nil{
		return err
	}
	wg:=&sync.WaitGroup{}
	wg.Add(2)
	go readLog(wg,writer,stderr)
	go readLog(wg,writer,stdout)
	if err = cmd.Wait();err!=nil{
		return err
	}
	wg.Wait()
	return nil
}

func readLog(wg *sync.WaitGroup, writer io.Writer, reader io.Reader) {
	defer wg.Done()
	r:=bufio.NewReader(reader)
	for {
		line,_,err:=r.ReadLine()
		if err == io.EOF || err!=nil{
			return
		}
		if len(line) > 0 {
			writer.Write(line)
			writer.Write([]byte("\n"))
		}
	}
}

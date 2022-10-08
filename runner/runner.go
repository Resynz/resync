/**
 * @Author: Resynz
 * @Date: 2022/1/6 17:09
 */
package runner

import (
	"context"
	"fmt"
	"github.com/rosbit/dbx"
	"net/url"
	"os"
	"resync/config"
	"resync/db"
	"resync/db/model"
	"resync/util"
	"strings"
	"time"
)

type Runner struct {
	Id           int64
	ExecutorId   int64
	ExecutorName string
	Task         *model.Task
	Detail       *model.TaskDetail
	ActionList   []*model.Action
	LogFile      *os.File
	CodeDir      string

	Context context.Context
	Cancel  context.CancelFunc
	isQuit  bool

	PercentTotal int64
	Percent      int64
}

func (s *Runner) initLogFile() error {
	df := fmt.Sprintf("./data/logs/runner/%d", s.Id)
	err := util.FormatDir(df)
	if err != nil {
		return err
	}
	name := "log.log"
	s.LogFile, err = os.OpenFile(fmt.Sprintf("%s/%s", df, name), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *Runner) writePreLog() error {
	_, err := s.LogFile.WriteString(fmt.Sprintf("Started at [%s] by user [%s]\n", time.Now().Format("2006-01-02 15:04:05"), s.ExecutorName))
	return err
}

func (s *Runner) dealSourceCode() error {
	if s.Detail.SourceCodeType == model.SourceCodeTypeNone {
		return nil
	}
	if s.Detail.SourceCodeType == model.SourceCodeTypeGit {
		return s.dealSourceCodeByGit()
	}
	return nil
}

func (s *Runner) dealSourceCodeByGit() error {
	s.LogFile.WriteString("Clone source code ...\n")
	cmd := "clone"
	if s.Detail.Branch != "" {
		cmd = fmt.Sprintf("%s -b %s", cmd, s.Detail.Branch)
	}
	_tmpRepository := s.Detail.RepositoryUrl
	if s.Detail.CodeAuthId > 0 {
		var codeAuth model.CodeAuth
		has, err := db.Handler.XStmt(codeAuth.GetTableName()).Where(dbx.Eq("id", s.Detail.CodeAuthId)).Get(&codeAuth)
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("git code auth not found")
		}
		hs := strings.Split(s.Detail.RepositoryUrl, "//")
		s.Detail.RepositoryUrl = fmt.Sprintf("%s//%s:%s@%s", hs[0], url.QueryEscape(codeAuth.UserName), url.QueryEscape(codeAuth.Password), hs[1])
		_tmpRepository = fmt.Sprintf("%s//%s:***@%s", hs[0], codeAuth.UserName, hs[1])
	}
	codeDir := fmt.Sprintf("./data/task/%s", s.Task.Name)
	s.CodeDir = codeDir
	_ = os.RemoveAll(codeDir)
	err := util.FormatDir(codeDir)
	if err != nil {
		return err
	}
	s.LogFile.WriteString("git " + fmt.Sprintf("%s %s %s", cmd, _tmpRepository, codeDir) + "\n")

	if s.isQuit {
		return fmt.Errorf("runner quit")
	}

	args := strings.Split(fmt.Sprintf("%s %s %s", cmd, s.Detail.RepositoryUrl, codeDir), " ")
	ctx, _ := context.WithCancel(s.Context)
	if err = util.RunCommand(ctx, s.LogFile, config.Conf.Pwd, "git", args...); err != nil {
		return err
	}
	if s.isQuit {
		return fmt.Errorf("runner quit")
	}
	args = strings.Split("git show --stat -q", " ")
	ctx, _ = context.WithCancel(s.Context)
	if err = util.RunCommand(ctx, s.LogFile, s.CodeDir, args[0], args[1:]...); err != nil {
		return err
	}
	s.LogFile.WriteString("Clone source code done.\n")
	return nil
}

func (s *Runner) dealActions() error {
	if len(s.ActionList) == 0 {
		return nil
	}
	s.LogFile.WriteString("Running actions ...\n")
	for _, v := range s.ActionList {
		if s.isQuit {
			return fmt.Errorf("runner quit")
		}
		if v.Type == model.ActionTypeShell {
			lines := strings.Split(v.Content, "\n")
			for _, l := range lines {
				if s.isQuit {
					return fmt.Errorf("runner quit")
				}
				s.LogFile.WriteString(fmt.Sprintf("%s\n", l))
				cds := strings.Split(l, " ")
				ctx, _ := context.WithCancel(s.Context)
				if err := util.RunCommand(ctx, s.LogFile, s.CodeDir, cds[0], cds[1:]...); err != nil {
					return err
				}
			}
		}
		s.Percent += 1
	}
	s.LogFile.WriteString("Run actions done.\n")
	return nil
}

func (s *Runner) Run() error {
	var err error
	go func() {
		<-s.Context.Done()
		if err == nil {
			s.LogFile.WriteString("User canceled.\n")
		}
		s.isQuit = true
	}()
	defer func() {
		if !s.isQuit {
			s.Cancel()
		}
		if err != nil {
			s.LogFile.WriteString(fmt.Sprintf("[FAILED] %s\n", err.Error()))
		}
		s.LogFile.Close()
	}()
	if s.isQuit {
		err = fmt.Errorf("runner quit")
		return err
	}
	if err = s.initLogFile(); err != nil {
		return err
	}
	if s.isQuit {
		err = fmt.Errorf("runner quit")
		return err
	}
	// step 1
	if err = s.writePreLog(); err != nil {
		return err
	}
	s.Percent += 1
	if s.isQuit {
		err = fmt.Errorf("runner quit")
		return err
	}
	// step 2
	if err = s.dealSourceCode(); err != nil {
		return err
	}
	if s.isQuit {
		err = fmt.Errorf("runner quit")
		return err
	}
	s.Percent += 1
	// step 3
	if err = s.dealActions(); err != nil {
		return err
	}
	s.Percent += 1
	return nil
}

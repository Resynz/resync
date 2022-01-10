/**
 * @Author: Resynz
 * @Date: 2022/1/5 11:02
 */
package model

type SourceCodeType uint8
const (
	SourceCodeTypeNone SourceCodeType = iota
	SourceCodeTypeGit
)

type AuthType uint8
const (
	AuthTypeUnknown AuthType = iota
	AuthTypeGit
)

type ActionType uint8
const (
	ActionTypeUnknown ActionType = iota
	ActionTypeShell
)

type AccountStatus uint8
const (
	AccountStatusUnknown AccountStatus = iota
	AccountStatusEnable
	AccountStatusDisable
)

type TaskStatus uint8
const (
	TaskStatusUnknown TaskStatus = iota
	TaskStatusPending
	TaskStatusProcess
	TaskStatusSuccess
	TaskStatusFailed
	TaskStatusCancel
)
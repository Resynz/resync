/**
 * @Author: Resynz
 * @Date: 2022/1/5 13:51
 */
package controller

import (
	"resync/code"
	"resync/common"
)

func Ping(ctx *common.Context) {
	common.HandleResponse(ctx,code.SuccessCode,nil)
}

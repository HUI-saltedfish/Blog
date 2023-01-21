package controllers

import "github.com/gin-gonic/gin"

func BackupPost(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJSON(c, res)
	err = BackupMysql()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func RestorePost(c *gin.Context) {

}

func BackupMysql() (err error) {
	// TODO
	return nil
}

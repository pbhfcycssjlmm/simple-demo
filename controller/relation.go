package controller

import (
	"net/http"
	"strconv"

	"github.com/RaymondCode/simple-demo/respository"
	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	respository.Response
	UserList []respository.User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
//关注功能
func RelationAction(c *gin.Context) {
	to_user_id := c.Query("to_user_id")
	touserid, _ := strconv.ParseInt(to_user_id, 10, 64)
	token := c.Query("token")
	follower := respository.UsersLoginInfo[token]
	actiontype := c.Query("action_type")
	var follow_follower respository.FollowFollower
	follow, err := respository.NewUserDaoInstance().QueryUserById(touserid)
	find := respository.Db.Table("follow_followers").Where("follow_id = ?", follow.Id).Where("follower_id = ?", follower.Id).Find(&follow_follower)
	if err != nil {
		c.JSON(http.StatusOK, respository.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
	if touserid == follower.Id {
		//作者不能关注自己
		c.JSON(http.StatusOK, respository.Response{StatusCode: 1, StatusMsg: "User is Author"})
	} else if _, exist := respository.UsersLoginInfo[token]; exist {
		if actiontype == "1" {
			follow.FollowCount++
			follower.FollowerCount++
			respository.NewUserDaoInstance().SaveUser(follower)
			respository.NewUserDaoInstance().SaveUser(*follow)
			if find != nil {
				follow_follower.FollowId = follow.Id
				follow_follower.FollowerId = follower.Id
				follow_follower.IsFavorite = true
				respository.Db.Save(&follow_follower)
			} else {
				respository.Db.Create(&follow_follower)
			}
		}
		if actiontype == "2" {
			follow.FollowCount--
			follower.FollowerCount--
			respository.NewUserDaoInstance().SaveUser(follower)
			respository.NewUserDaoInstance().SaveUser(*follow)
			follow_follower.FollowId = follow.Id
			follow_follower.FollowerId = follower.Id
			follow_follower.IsFavorite = false
			respository.Db.Save(&follow_follower)
		}
		c.JSON(http.StatusOK, respository.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, respository.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	userid := c.Query("user_id")
	uid, _ := strconv.ParseInt(userid, 10, 64)
	FollowList := respository.QueryFollowListByUserId(uid)
	c.JSON(http.StatusOK, UserListResponse{
		Response: respository.Response{
			StatusCode: 0,
		},
		UserList: FollowList,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	userid := c.Query("user_id")
	uid, _ := strconv.ParseInt(userid, 10, 64)
	FollowerList := respository.QueryFollowerListByUserId(uid)
	c.JSON(http.StatusOK, UserListResponse{
		Response: respository.Response{
			StatusCode: 0,
		},
		UserList: FollowerList,
	})
}

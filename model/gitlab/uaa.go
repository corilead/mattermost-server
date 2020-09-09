package oauthgitlab

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/mattermost/mattermost-server/v5/einterfaces"
	"github.com/mattermost/mattermost-server/v5/model"
)

type UaaProvider struct {
}

type UaaUser struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Name     string `json:"full_name"`
}

func init() {
	provider := &UaaProvider{}
	einterfaces.RegisterOauthProvider(model.USER_AUTH_SERVICE_UAA, provider)
}

func userFromUaaUser(glu *UaaUser) *model.User {
	user := &model.User{}
	username := glu.Username
	if username == "" {
		username = glu.Login
	}
	user.Username = model.CleanUsername(username)
	splitName := strings.Split(glu.Name, " ")
	if len(splitName) == 2 {
		user.FirstName = splitName[0]
		user.LastName = splitName[1]
	} else if len(splitName) >= 2 {
		user.FirstName = splitName[0]
		user.LastName = strings.Join(splitName[1:], " ")
	} else {
		user.FirstName = glu.Name
	}
	user.Email = glu.Email
	user.Email = strings.ToLower(user.Email)
	userId := glu.getAuthData()
	user.AuthData = &userId
	user.AuthService = model.USER_AUTH_SERVICE_UAA

	return user
}

func uaaUserFromJson(data io.Reader) *UaaUser {
	decoder := json.NewDecoder(data)
	var glu UaaUser
	err := decoder.Decode(&glu)
	fmt.Println("err ", err)
	if err == nil {
		return &glu
	} else {
		return nil
	}
}

func (glu *UaaUser) ToJson() string {
	b, err := json.Marshal(glu)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func (glu *UaaUser) IsValid() bool {
	if glu.Id == "" {
		return false
	}

	//if len(glu.Email) == 0 {
	//	return false
	//}

	return true
}

func (glu *UaaUser) getAuthData() string {
	return glu.Id
}

func (m *UaaProvider) GetUserFromJson(data io.Reader) *model.User {
	glu := uaaUserFromJson(data)
	fmt.Println(glu)
	if glu.IsValid() {
		return userFromUaaUser(glu)
	}

	return &model.User{}
}

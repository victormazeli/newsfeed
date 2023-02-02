package handlers

import (
	"encoding/json"
	"errors"
	"gatewayservice/graph/model"
	"gatewayservice/httpclients"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"net/http"
)

var httpClient = resty.New()

type AuthHandler struct{}

func (a AuthHandler) UserLogin(credential model.LoginCredential) (*model.AuthResponse, error) {
	resp, err := httpClient.R().SetBody(credential).Post(httpclients.UserService + "/auth/token")

	if err != nil {
		logrus.Errorf("Error login user, %s", err)
		return nil, err
	}

	logrus.Infof("Response code: %d", resp.StatusCode())
	if resp.StatusCode() == http.StatusOK {
		authResponse := model.AuthResponse{}
		err := json.Unmarshal(resp.Body(), &authResponse)
		if err != nil {
			logrus.Errorf("Error parsing response, %s", err)
			return nil, err
		}
		return &authResponse, nil
	} else {
		return nil, errors.New("response code not matching, got " + string(rune(resp.StatusCode())))
	}

}

func (a AuthHandler) UserRegistration(detail model.UserSignUpDetail) (*model.User, error) {

	resp, err := httpClient.R().SetBody(detail).Post(httpclients.UserService + "/auth/register")

	if err != nil {
		logrus.Errorf("Error login user, %s", err)
		return nil, err
	}

	logrus.Infof("Response code: %d", resp.StatusCode())
	if resp.StatusCode() == http.StatusCreated {
		response := model.User{}
		err := json.Unmarshal(resp.Body(), &response)
		if err != nil {
			logrus.Errorf("Error parsing response, %s", err)
			return nil, err
		}
		return &response, nil
	} else {
		return nil, errors.New("response code not matching, got " + string(rune(resp.StatusCode())))
	}

}

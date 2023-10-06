package httpcontroller

import (
	"github.com/gin-gonic/gin"

	"github.com/taraslis453/solid-software-test/internal/entity"
	"github.com/taraslis453/solid-software-test/internal/service"
	"github.com/taraslis453/solid-software-test/pkg/errs"
)

type userRoutes struct {
	routerContext
}

func newUserRoutes(options RouterOptions) {
	r := &userRoutes{
		routerContext{
			services: options.Services,
			logger:   options.Logger.Named("userRoutes"),
			cfg:      options.Config,
		},
	}

	p := options.Handler.Group("/users")
	{
		p.POST("register/", errorHandler(options, r.registerUser))
		p.GET("/login", errorHandler(options, r.loginUser))
		p.POST("/refresh-token", errorHandler(options, r.refreshToken))
		p.GET("/:id", newAuthMiddleware(options), errorHandler(options, r.getUserResponse))
		p.PUT("", newAuthMiddleware(options), errorHandler(options, r.updateUser))
	}
}

type registerUserRequestBody struct {
	Name         string `json:"name" binding:"required"`
	Surname      string `json:"surname" binding:"required"`
	EmailAddress string `json:"email" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Phone        string `json:"phone"`
}

type registerUserResponse struct {
}

func (r *userRoutes) registerUser(c *gin.Context) (interface{}, *httpErr) {
	logger := r.logger.Named("registerUser").WithContext(c)

	var body registerUserRequestBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		logger.Info("failed to parse body", "err", err)
		return nil, &httpErr{Type: httpErrTypeClient, Message: "invalid request body", Details: err}
	}
	logger = logger.With("body", body)
	logger.Debug("parsed request body")

	err = r.services.User.RegisterUser(c, service.RegisterUserOptions{
		Name:         body.Name,
		Surname:      body.Surname,
		Phone:        body.Phone,
		EmailAddress: body.EmailAddress,
		Password:     body.Password,
	})
	if err != nil {
		if errs.IsExpected(err) {
			logger.Info(err.Error())
			return nil, &httpErr{Type: httpErrTypeClient, Message: err.Error(), Code: errs.GetCode(err)}
		}

		logger.Error("failed to register user", "err", err)
		return nil, &httpErr{Type: httpErrTypeServer, Message: "failed to register user", Details: err}
	}

	logger.Info("successfully registered user")
	return registerUserResponse{}, nil
}

type loginUserRequestQuery struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type loginUserResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	UserID       string `json:"userId"`
}

func (r *userRoutes) loginUser(c *gin.Context) (interface{}, *httpErr) {
	logger := r.logger.Named("loginUser").WithContext(c)

	var query loginUserRequestQuery
	err := c.ShouldBindQuery(&query)
	if err != nil {
		logger.Info("failed to parse query", "err", err)
		return nil, &httpErr{Type: httpErrTypeClient, Message: "invalid request query", Details: err}
	}
	logger = logger.With("query", query)
	logger.Debug("parsed query")

	output, err := r.services.User.LoginUser(c, service.LoginUserOptions{
		EmailAddress: query.Email,
		Password:     query.Password,
	})
	if err != nil {
		if errs.IsExpected(err) {
			logger.Info(err.Error())
			return nil, &httpErr{Type: httpErrTypeClient, Message: err.Error(), Code: errs.GetCode(err)}
		}

		logger.Error("failed to login user", "err", err)
		return nil, &httpErr{Type: httpErrTypeServer, Message: "failed to login user", Details: err}
	}

	logger.Info("successfully logged in user")
	return loginUserResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
		UserID:       output.UserID,
	}, nil
}

type refreshTokenResponseBody struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (r *userRoutes) refreshToken(c *gin.Context) (interface{}, *httpErr) {
	logger := r.logger.Named("refreshToken").WithContext(c)

	token, err := getAuthToken(c.GetHeader("Authorization"))
	if err != nil {
		logger.Info(err.Error())
		return nil, &httpErr{Type: httpErrTypeClient, Message: err.Error()}
	}

	refreshedToken, err := r.services.User.RefreshUserToken(c, token)
	if err != nil {
		if errs.IsExpected(err) {
			logger.Info(err.Error())
			return nil, &httpErr{Type: httpErrTypeClient, Message: err.Error(), Code: errs.GetCode(err)}
		}
		logger.Error("failed to refresh token", "err", err)
		return nil, &httpErr{Type: httpErrTypeServer, Message: "failed to refresh token", Details: err}
	}

	logger.Info("token successfully refreshed")
	return refreshTokenResponseBody{
		AccessToken:  refreshedToken.AccessToken,
		RefreshToken: refreshedToken.RefreshToken,
	}, nil
}

type getUserPathParams struct {
	ID string `uri:"id" binding:"required"`
}

type getUserResponse struct {
	User *entity.User `json:"user"`
}

func (r *userRoutes) getUserResponse(c *gin.Context) (interface{}, *httpErr) {
	logger := r.logger.Named("getUserResponse").WithContext(c)

	var pathParams getUserPathParams
	err := c.ShouldBindUri(&pathParams)
	if err != nil {
		logger.Info("failed to parse path params", "err", err)
		return nil, &httpErr{Type: httpErrTypeClient, Message: "invalid path params", Details: err}
	}

	user, err := r.services.User.GetUser(c, service.GetUserOptions{
		ID: pathParams.ID,
	})
	if err != nil {
		if errs.IsExpected(err) {
			logger.Info(err.Error())
			return nil, &httpErr{Type: httpErrTypeClient, Message: err.Error(), Code: errs.GetCode(err)}
		}

		logger.Error("failed to get user", "err", err)
		return nil, &httpErr{Type: httpErrTypeServer, Message: "failed to get user", Details: err}
	}

	logger.Info("successfully got me")
	return getUserResponse{
		User: user,
	}, nil
}

type updateUserRequestBody struct {
	User *entity.User `json:"user" binding:"required"`
}

type updateUserResponse struct {
	User *entity.User `json:"user" binding:"required"`
}

func (r *userRoutes) updateUser(c *gin.Context) (interface{}, *httpErr) {
	logger := r.logger.Named("updateUser").WithContext(c)

	var body updateUserRequestBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		logger.Info("failed to parse body", "err", err)
		return nil, &httpErr{Type: httpErrTypeClient, Message: "invalid request body", Details: err}
	}
	logger = logger.With("body", body)
	logger.Debug("parsed request body")

	updatedUser, err := r.services.User.UpdateUser(c, body.User)
	if err != nil {
		if errs.IsExpected(err) {
			logger.Info(err.Error())
			return nil, &httpErr{Type: httpErrTypeClient, Message: err.Error(), Code: errs.GetCode(err)}
		}

		logger.Error("failed to update user", "err", err)
		return nil, &httpErr{Type: httpErrTypeServer, Message: "failed to update user", Details: err}
	}

	logger.Info("successfully updated user")
	return updateUserResponse{
		User: updatedUser,
	}, nil
}

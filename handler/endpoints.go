package handler

import (
	"net/http"

	"InterviewBackendSawitProGolang/generated"
	"InterviewBackendSawitProGolang/pkg/jwt"
	"InterviewBackendSawitProGolang/pkg/validator"
	"InterviewBackendSawitProGolang/repository"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var err error

func (s *Server) Hello(ctx echo.Context) error {
	fmt.Println(ctx.Request().Context().Value("user_id"))
	return ctx.JSON(http.StatusOK, "hello world!")
}
func (s *Server) GetProfile(ctx echo.Context) error {
	var resp generated.GetProfileResponse
	userUUID, err := s.convertUserIDtoUUID(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to convert uuid")
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "internal server error",
		})
	}

	var input repository.GetUserByIDInput = repository.GetUserByIDInput{
		ID: userUUID,
	}

	output, err := s.Repository.GetUserByID(ctx.Request().Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Failed get profile")
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "internal server error",
		})
	}

	resp.FullName = &output.FullName
	resp.PhoneNumber = &output.PhoneNumber

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) Register(ctx echo.Context) error {
	var req generated.RegisterJSONBody
	var resp generated.RegisterResponse
	errors := map[string]interface{}{}
	if err := ctx.Bind(&req); err != nil {
		log.Error().Err(err).Msg("Unable to bind request")
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "internal server error",
		})
	}
	passwordSalt := randstr.String(10)
	var input repository.User = repository.User{
		UserInfo: repository.UserInfo{
			PhoneNumber: *req.PhoneNumber,
			FullName:    *req.FullName,
		},
		UserSecret: repository.UserSecret{
			Password:     *req.Password,
			PasswordSalt: passwordSalt,
		},
	}

	err = validator.Validate(input)
	if err != nil {
		json.Unmarshal([]byte(err.Error()), &errors)
		return ctx.JSON(http.StatusBadRequest, generated.ErrorBadRequestResponse{
			Errors: &errors,
		})
	}

	hashPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(*req.Password+passwordSalt), bcrypt.MinCost)
	if err != nil {
		log.Error().Err(err).Msg("Unable to hash password")
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "internal server error",
		})
	}

	input.Password = string(hashPasswordBytes)

	user, err := s.getProfileByPhoneNumber(ctx, input.PhoneNumber)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "internal server error",
		})
	}

	if user.PhoneNumber != "" {
		errors["phoneNumber"] = []string{
			"Phonenumber already exists",
		}

		return ctx.JSON(http.StatusBadRequest, generated.ErrorBadRequestResponse{
			Errors: &errors,
		})
	}

	output, err := s.Repository.InsertUser(ctx.Request().Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert user")
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "internal server error",
		})
	}
	resp.Id = &output.ID

	return ctx.JSON(http.StatusCreated, resp)
}

func (s *Server) Login(ctx echo.Context) error {
	curr := time.Now()
	var resp generated.LoginResponse
	var req generated.LoginJSONBody

	if err := ctx.Bind(&req); err != nil {
		log.Error().Err(err).Msg("Unable to bind request")
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "internal server error",
		})
	}

	user, err := s.getProfileByPhoneNumber(ctx, *req.PhoneNumber)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "internal server error",
		})
	}

	if user.ID.String() == "" {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "phonenumber or password is wrong",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(*req.Password+user.PasswordSalt))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "phonenumber or password is wrong",
		})
	}

	resp.Id = &user.ID
	token, err := jwt.CreateJWSWithClaims(map[string]interface{}{"id": user.ID.String()})
	if err != nil {
		log.Error().Err(err).Msg("Unable to create JWT Token")
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "internal server error",
		})
	}

	sToken := string(token)
	resp.Token = &sToken

	if err := s.Repository.UpdateLastLoginAndSuccessfullyLogin(ctx.Request().Context(), repository.UpdateLastLoginAndSuccessfullyLoginInput{
		ID:        user.ID,
		LastLogin: &curr,
	}); err != nil {
		log.Error().Err(err).Msg("Failed update successfully login and last login")
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "internal server error",
		})
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) UpdateProfile(ctx echo.Context) error {
	var resp generated.UpdateProfileResponse
	var req generated.UpdateProfileJSONBody
	errors := map[string]interface{}{}

	if err := ctx.Bind(&req); err != nil {
		log.Error().Err(err).Msg("Unable to bind request")
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "internal server error",
		})
	}

	userUUID, err := s.convertUserIDtoUUID(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to convert uuid")
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "internal server error",
		})
	}

	var input repository.UpdateUserInput = repository.UpdateUserInput{
		ID:          userUUID,
		PhoneNumber: *req.PhoneNumber,
		FullName:    *req.FullName,
	}

	err = validator.Validate(input)
	if err != nil {
		json.Unmarshal([]byte(err.Error()), &errors)
		return ctx.JSON(http.StatusBadRequest, generated.ErrorBadRequestResponse{
			Errors: &errors,
		})
	}

	user, err := s.getProfileByPhoneNumber(ctx, input.PhoneNumber)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "internal server error",
		})
	}

	if user.PhoneNumber != "" {
		if user.ID != userUUID {
			return ctx.JSON(http.StatusConflict, generated.ErrorResponse{
				Message: "Phonenumber already exists",
			})
		}
	}

	if err := s.Repository.UpdateUser(ctx.Request().Context(), input); err != nil {
		log.Error().Err(err).Msg("Failed to Update Profile")
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "internal server error",
		})
	}

	resp.Id = &userUUID
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) convertUserIDtoUUID(ctx echo.Context) (uuid.UUID, error) {
	var userID string = ctx.Get("user_id").(string)
	return uuid.Parse(userID)
}

func (s *Server) getProfileByPhoneNumber(ctx echo.Context, phoneNumber string) (repository.User, error) {
	var getUserByPhoneInput repository.GetUserByPhoneNumberInput = repository.GetUserByPhoneNumberInput{
		PhoneNumber: phoneNumber,
	}

	outputUser, err := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), getUserByPhoneInput)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error().Err(err).Msg("Failed to get profile")
			return repository.User{}, err
		}
	}

	return outputUser, nil
}

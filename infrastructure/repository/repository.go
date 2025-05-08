package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Nerzal/gocloak/v8"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/mitchellh/mapstructure"
	"github.com/patricksferraz/auth-service/application/grpc/pb"
	"github.com/patricksferraz/auth-service/domain/entity"
	"github.com/patricksferraz/auth-service/infrastructure/external"
	"github.com/patricksferraz/auth-service/logger"
	"github.com/patricksferraz/auth-service/utils"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
)

type Repository struct {
	K              *external.Keycloak
	Kafka          *external.Kafka
	EmployeeClient *external.EmployeeClient
}

func NewRepository(keycloak *external.Keycloak, kafka *external.Kafka, employeeClient *external.EmployeeClient) *Repository {
	return &Repository{
		K:              keycloak,
		Kafka:          kafka,
		EmployeeClient: employeeClient,
	}
}

func (r *Repository) Login(ctx context.Context, auth *entity.Auth) (*entity.JWT, error) {
	span, ctx := apm.StartSpan(ctx, "Login", "repository")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	jwt, err := r.K.Client.Login(ctx, r.K.ClientID, r.K.ClientSecret, r.K.Realm, auth.Username, auth.Password)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	return &entity.JWT{
		AccessToken:      jwt.AccessToken,
		IDToken:          jwt.IDToken,
		ExpiresIn:        jwt.ExpiresIn,
		RefreshExpiresIn: jwt.RefreshExpiresIn,
		RefreshToken:     jwt.RefreshToken,
		TokenType:        jwt.TokenType,
		NotBeforePolicy:  jwt.NotBeforePolicy,
		SessionState:     jwt.SessionState,
		Scope:            jwt.Scope,
	}, nil
}

func (r *Repository) RefreshToken(ctx context.Context, refreshToken string) (*entity.JWT, error) {
	span, ctx := apm.StartSpan(ctx, "RefreshToken", "repository")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	jwt, err := r.K.Client.RefreshToken(ctx, refreshToken, r.K.ClientID, r.K.ClientSecret, r.K.Realm)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	return &entity.JWT{
		AccessToken:      jwt.AccessToken,
		IDToken:          jwt.IDToken,
		ExpiresIn:        jwt.ExpiresIn,
		RefreshExpiresIn: jwt.RefreshExpiresIn,
		RefreshToken:     jwt.RefreshToken,
		TokenType:        jwt.TokenType,
		NotBeforePolicy:  jwt.NotBeforePolicy,
		SessionState:     jwt.SessionState,
		Scope:            jwt.Scope,
	}, nil
}

func (r *Repository) FindClaimsByToken(ctx context.Context, accessToken string) (*entity.Claims, error) {
	span, ctx := apm.StartSpan(ctx, "FindClaimsByToken", "repository")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	jwt, _, err := r.K.Client.DecodeAccessToken(ctx, accessToken, r.K.Realm, r.K.Audience)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	fmt.Println(jwt.Claims)

	claims := new(entity.Claims)
	mapstructure.Decode(jwt.Claims, claims)
	log.WithField("claims", claims).Info("claims mapstructure")

	type ResourceAccess struct {
		ResourceAccess map[string]map[string][]string `mapstructure:"resource_access"`
	}

	ra := new(ResourceAccess)
	mapstructure.Decode(jwt.Claims, ra)

	roles := ra.ResourceAccess[r.K.ClientID]["roles"]
	claims.Roles = roles
	log.WithField("claims", claims).Info("claims with roles")

	return claims, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *entity.User, accessToken string) error {
	gUser := gocloak.User{
		Username: &user.Username,
		Enabled:  &user.Enabled,
	}
	gUser.Attributes = utils.StructToAttr(user)

	userID, err := r.K.Client.CreateUser(ctx, accessToken, r.K.Realm, gUser)
	if err != nil {
		return err
	}

	user.ID = userID
	return nil
}

func (r *Repository) FindUser(ctx context.Context, id string, accessToken string) (*entity.User, error) {
	e, err := r.K.Client.GetUserByID(ctx, accessToken, r.K.Realm, id)
	if err != nil {
		return nil, err
	}

	var employeeID string
	if e.Attributes != nil {
		employeeID = (*e.Attributes)["employee_id"][0]
	}

	user := &entity.User{
		Username:   *e.Username,
		Enabled:    *e.Enabled,
		EmployeeID: employeeID,
	}
	user.ID = *e.ID
	user.CreatedAt = time.Unix(0, *e.CreatedTimestamp*int64(time.Millisecond))

	return user, nil
}

func (r *Repository) SearchUsers(ctx context.Context, filter *entity.Filter, accessToken string) ([]*entity.User, error) {
	first := *filter.Page * *filter.PageSize
	gUsers, err := r.K.Client.GetUsers(
		ctx,
		accessToken,
		r.K.Realm,
		gocloak.GetUsersParams{
			Username: filter.Username,
			Enabled:  filter.Enabled,
			First:    &first,
			Max:      filter.PageSize,
		},
	)
	if err != nil {
		return nil, err
	}

	var users []*entity.User
	for _, u := range gUsers {
		var employeeID string
		if u.Attributes != nil {
			employeeID = (*u.Attributes)["employee_id"][0]
		}
		user := &entity.User{
			Username:   *u.Username,
			Enabled:    *u.Enabled,
			EmployeeID: employeeID,
		}
		user.ID = *u.ID
		user.CreatedAt = time.Unix(0, *u.CreatedTimestamp*int64(time.Millisecond))
		users = append(users, user)
	}

	return users, nil
}

func (r *Repository) SetPassword(ctx context.Context, pass *entity.PasswordInfo, accessToken string) error {
	err := r.K.Client.SetPassword(ctx, accessToken, pass.UserID, r.K.Realm, pass.Password, pass.Temporary)
	return err
}

func (r *Repository) PublishEvent(ctx context.Context, msg, topic, key string) error {
	message := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &topic, Partition: ckafka.PartitionAny},
		Value:          []byte(msg),
		Key:            []byte(key),
	}
	err := r.Kafka.Producer.Produce(message, r.Kafka.DeliveryChan)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) FindEmployee(ctx context.Context, employeeID string) error {
	req := &pb.FindEmployeeRequest{
		Id: employeeID,
	}
	_, err := r.EmployeeClient.C.FindEmployee(ctx, req)
	return err
}

func (r *Repository) Logout(ctx context.Context, refreshToken string) error {
	err := r.K.Client.Logout(ctx, r.K.ClientID, r.K.ClientSecret, r.K.Realm, refreshToken)
	return err
}

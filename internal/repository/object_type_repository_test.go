package repository

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/goccy/go-json"
	"github.com/krobus00/storage-service/internal/infrastructure"
	"github.com/krobus00/storage-service/internal/model"
	"github.com/krobus00/storage-service/internal/utils"
	"github.com/spf13/viper"
)

func newObjectTypeRepoMock(t *testing.T) (model.ObjectTypeRepository, sqlmock.Sqlmock, *miniredis.Miniredis) {
	dbConn, dbMock := utils.NewDBMock()
	miniRedis := miniredis.RunT(t)
	viper.Set("redis.cache_host", fmt.Sprintf("redis://%s", miniRedis.Addr()))
	redisClient, err := infrastructure.NewRedisClient()
	utils.ContinueOrFatal(err)
	objectTypeRepo := NewObjectTypeRepository()
	err = objectTypeRepo.InjectDB(dbConn)
	utils.ContinueOrFatal(err)
	err = objectTypeRepo.InjectRedisClient(redisClient)
	utils.ContinueOrFatal(err)

	return objectTypeRepo, dbMock, miniRedis
}

func Test_objectTypeRepository_Create(t *testing.T) {
	var (
		objectTypeID = utils.GenerateUUID()
	)
	type args struct {
		objectType *model.ObjectType
	}
	tests := []struct {
		name    string
		args    args
		mockErr error
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				objectType: &model.ObjectType{
					ID:   objectTypeID,
					Name: "image",
				},
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "error create",
			args: args{
				objectType: &model.ObjectType{
					ID:   objectTypeID,
					Name: "image",
				},
			},
			mockErr: errors.New("db error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()

			r, dbMock, _ := newObjectTypeRepoMock(t)

			dbMock.ExpectBegin()
			dbMock.ExpectExec("INSERT INTO \"object_types\"").
				WithArgs(tt.args.objectType.ID, tt.args.objectType.Name).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.mockErr)

			if tt.wantErr {
				dbMock.ExpectRollback()
			} else {
				dbMock.ExpectCommit()
			}

			if err := r.Create(ctx, tt.args.objectType); (err != nil) != tt.wantErr {
				t.Errorf("objectTypeRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_objectTypeRepository_FindByID(t *testing.T) {
	var (
		objectTypeID = utils.GenerateUUID()
	)
	type mockSelect struct {
		objectType *model.ObjectType
		err        error
	}
	type mockCache struct {
		objectType *model.ObjectType
	}
	type args struct {
		id string
	}
	tests := []struct {
		name       string
		args       args
		mockSelect *mockSelect
		mockCache  *mockCache
		want       *model.ObjectType
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				id: objectTypeID,
			},
			mockSelect: &mockSelect{
				objectType: &model.ObjectType{
					ID:   objectTypeID,
					Name: "image",
				},
			},
			mockCache: nil,
			want: &model.ObjectType{
				ID:   objectTypeID,
				Name: "image",
			},
			wantErr: false,
		},
		{
			name: "success found in cache",
			args: args{
				id: objectTypeID,
			},
			mockCache: &mockCache{
				objectType: &model.ObjectType{
					ID:   objectTypeID,
					Name: "image",
				},
			},
			want: &model.ObjectType{
				ID:   objectTypeID,
				Name: "image",
			},
			wantErr: false,
		},
		{
			name: "error object type not found",
			args: args{
				id: objectTypeID,
			},
			mockSelect: &mockSelect{
				objectType: nil,
				err:        nil,
			},
			mockCache: nil,
			want:      nil,
			wantErr:   false,
		},
		{
			name: "error find object type",
			args: args{
				id: objectTypeID,
			},
			mockSelect: &mockSelect{
				objectType: nil,
				err:        errors.New("db error"),
			},
			mockCache: nil,
			want:      nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			r, dbMock, redisMock := newObjectTypeRepoMock(t)

			cacheKey := model.NewObjectTypeCacheKeyByID(tt.args.id)
			if tt.mockSelect != nil {
				row := sqlmock.NewRows([]string{"id", "name"})
				if tt.mockSelect.objectType != nil {
					objectType := tt.mockSelect.objectType
					row.AddRow(
						objectType.ID,
						objectType.Name,
					)
				}

				dbMock.ExpectQuery("^SELECT .+ FROM \"object_types\"").
					WithArgs(tt.args.id).
					WillReturnRows(row).
					WillReturnError(tt.mockSelect.err)
			}
			if tt.mockCache != nil {
				cacheData, err := json.Marshal(tt.mockCache.objectType)
				if err != nil {
					utils.ContinueOrFatal(err)
				}
				_ = redisMock.Set(cacheKey, string(cacheData))
			}

			got, err := r.FindByID(ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("objectTypeRepository.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("objectTypeRepository.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_objectTypeRepository_FindByName(t *testing.T) {
	var (
		objectTypeID = utils.GenerateUUID()
	)
	type mockSelect struct {
		objectType *model.ObjectType
		err        error
	}
	type mockCache struct {
		objectType *model.ObjectType
	}
	type args struct {
		name string
	}
	tests := []struct {
		name       string
		args       args
		mockSelect *mockSelect
		mockCache  *mockCache
		want       *model.ObjectType
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				name: "image",
			},
			mockSelect: &mockSelect{
				objectType: &model.ObjectType{
					ID:   objectTypeID,
					Name: "image",
				},
			},
			mockCache: nil,
			want: &model.ObjectType{
				ID:   objectTypeID,
				Name: "image",
			},
			wantErr: false,
		},
		{
			name: "success found in cache",
			args: args{
				name: "image",
			},
			mockCache: &mockCache{
				objectType: &model.ObjectType{
					ID:   objectTypeID,
					Name: "image",
				},
			},
			want: &model.ObjectType{
				ID:   objectTypeID,
				Name: "image",
			},
			wantErr: false,
		},
		{
			name: "error object type not found",
			args: args{
				name: "image",
			},
			mockSelect: &mockSelect{
				objectType: nil,
				err:        nil,
			},
			mockCache: nil,
			want:      nil,
			wantErr:   false,
		},
		{
			name: "error find object type",
			args: args{
				name: "imae",
			},
			mockSelect: &mockSelect{
				objectType: nil,
				err:        errors.New("db error"),
			},
			mockCache: nil,
			want:      nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			r, dbMock, redisMock := newObjectTypeRepoMock(t)

			cacheKey := model.NewObjectTypeCacheKeyByName(tt.args.name)
			if tt.mockSelect != nil {
				row := sqlmock.NewRows([]string{"id", "name"})
				if tt.mockSelect.objectType != nil {
					objectType := tt.mockSelect.objectType
					row.AddRow(
						objectType.ID,
						objectType.Name,
					)
				}

				dbMock.ExpectQuery("^SELECT .+ FROM \"object_types\"").
					WithArgs(tt.args.name).
					WillReturnRows(row).
					WillReturnError(tt.mockSelect.err)
			}
			if tt.mockCache != nil {
				cacheData, err := json.Marshal(tt.mockCache.objectType)
				if err != nil {
					utils.ContinueOrFatal(err)
				}
				_ = redisMock.Set(cacheKey, string(cacheData))
			}

			got, err := r.FindByName(ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("objectTypeRepository.FindByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("objectTypeRepository.FindByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_objectTypeRepository_DeleteByID(t *testing.T) {
	var (
		objectTypeID = utils.GenerateUUID()
	)
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		mockErr error
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				id: objectTypeID,
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "error delete object type",
			args: args{
				id: objectTypeID,
			},
			mockErr: errors.New("db error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, dbMock, _ := newObjectTypeRepoMock(t)

			dbMock.ExpectBegin()
			row := sqlmock.NewRows([]string{"id", "name"})

			row.AddRow(tt.args.id, "image")

			dbMock.ExpectQuery("DELETE FROM \"object_types\"").
				WithArgs(tt.args.id).
				WillReturnRows(row).
				WillReturnError(tt.mockErr)

			if tt.wantErr {
				dbMock.ExpectRollback()
			} else {
				dbMock.ExpectCommit()
			}
			if err := r.DeleteByID(context.TODO(), tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("objectTypeRepository.DeleteByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

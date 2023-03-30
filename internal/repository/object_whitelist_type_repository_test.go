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
	"gorm.io/gorm"
)

func newObjecteWhitelistTypeRepoMock(t *testing.T) (model.ObjectWhitelistTypeRepository, sqlmock.Sqlmock, *miniredis.Miniredis) {
	dbConn, dbMock := utils.NewDBMock()
	miniRedis := miniredis.RunT(t)
	viper.Set("redis.cache_host", fmt.Sprintf("redis://%s", miniRedis.Addr()))
	redisClient, err := infrastructure.NewRedisClient()
	utils.ContinueOrFatal(err)
	objectWhitelistTypeRepo := NewObjectWhitelistTypeRepository()
	err = objectWhitelistTypeRepo.InjectDB(dbConn)
	utils.ContinueOrFatal(err)
	err = objectWhitelistTypeRepo.InjectRedisClient(redisClient)
	utils.ContinueOrFatal(err)

	return objectWhitelistTypeRepo, dbMock, miniRedis
}

func Test_objectWhitelistTypeRepository_Create(t *testing.T) {
	var (
		typeID = utils.GenerateUUID()
	)
	type args struct {
		objectWhitelistType *model.ObjectWhitelistType
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
				objectWhitelistType: &model.ObjectWhitelistType{
					TypeID:    typeID,
					Extension: ".jpg",
				},
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "error create object whitelist type",
			args: args{
				objectWhitelistType: &model.ObjectWhitelistType{
					TypeID:    typeID,
					Extension: ".jpg",
				},
			},
			mockErr: errors.New("db error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()

			r, dbMock, _ := newObjecteWhitelistTypeRepoMock(t)

			dbMock.ExpectBegin()
			dbMock.ExpectExec("INSERT INTO \"object_whitelist_types\"").
				WithArgs(tt.args.objectWhitelistType.TypeID, tt.args.objectWhitelistType.Extension).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.mockErr)

			if tt.wantErr {
				dbMock.ExpectRollback()
			} else {
				dbMock.ExpectCommit()
			}
			if err := r.Create(ctx, tt.args.objectWhitelistType); (err != nil) != tt.wantErr {
				t.Errorf("objectWhitelistTypeRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_objectWhitelistTypeRepository_FindByTypeIDAndExt(t *testing.T) {
	var (
		typeID = utils.GenerateUUID()
	)
	type mockSelect struct {
		res *model.ObjectWhitelistType
		err error
	}
	type mockCache struct {
		res *model.ObjectWhitelistType
	}
	type args struct {
		typeID string
		ext    string
	}
	tests := []struct {
		name       string
		args       args
		mockSelect *mockSelect
		mockCache  *mockCache
		want       *model.ObjectWhitelistType
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				typeID: typeID,
				ext:    ".jpg",
			},
			mockSelect: &mockSelect{
				res: &model.ObjectWhitelistType{
					TypeID:    typeID,
					Extension: ".jpg",
				},
				err: nil,
			},
			want: &model.ObjectWhitelistType{
				TypeID:    typeID,
				Extension: ".jpg",
			},
			wantErr: false,
		},
		{
			name: "success found in cache",
			args: args{
				typeID: typeID,
				ext:    ".jpg",
			},
			mockCache: &mockCache{
				res: &model.ObjectWhitelistType{
					TypeID:    typeID,
					Extension: ".jpg",
				},
			},
			want: &model.ObjectWhitelistType{
				TypeID:    typeID,
				Extension: ".jpg",
			},
			wantErr: false,
		},
		{
			name: "error object whitelist type not found",
			args: args{
				typeID: typeID,
				ext:    ".jpg",
			},
			mockSelect: &mockSelect{
				res: nil,
				err: gorm.ErrRecordNotFound,
			},
			wantErr: false,
		},
		{
			name: "error find object whitelist type",
			args: args{
				typeID: typeID,
				ext:    ".jpg",
			},
			mockSelect: &mockSelect{
				res: nil,
				err: errors.New("db error"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			r, dbMock, redisMock := newObjecteWhitelistTypeRepoMock(t)

			cacheBucketKey := utils.NewBucketKey(model.NewObjectWhitelistTypeCacheKey(tt.args.typeID), tt.args.ext)

			if tt.mockSelect != nil {
				row := sqlmock.NewRows([]string{"type_id", "extension"})
				if tt.mockSelect.res != nil {
					objectWhitelistType := tt.mockSelect.res
					row.AddRow(
						objectWhitelistType.TypeID,
						objectWhitelistType.Extension,
					)
				}

				dbMock.ExpectQuery("^SELECT .+ FROM \"object_whitelist_types\"").
					WithArgs(tt.args.typeID, tt.args.ext).
					WillReturnRows(row).
					WillReturnError(tt.mockSelect.err)
			}
			if tt.mockCache != nil {
				cacheData, err := json.Marshal(tt.mockCache.res)
				if err != nil {
					utils.ContinueOrFatal(err)
				}
				redisMock.HSet(cacheBucketKey, tt.args.ext, string(cacheData))
			}

			got, err := r.FindByTypeIDAndExt(ctx, tt.args.typeID, tt.args.ext)
			if (err != nil) != tt.wantErr {
				t.Errorf("objectWhitelistTypeRepository.FindByTypeIDAndExt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("objectWhitelistTypeRepository.FindByTypeIDAndExt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_objectWhitelistTypeRepository_DeleteByTypeIDAndExt(t *testing.T) {
	var (
		typeID = utils.GenerateUUID()
	)
	type args struct {
		typeID string
		ext    string
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
				typeID: typeID,
				ext:    ".jpg",
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "error delete object whitelist type",
			args: args{
				typeID: typeID,
				ext:    ".jpg",
			},
			mockErr: errors.New("db error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, dbMock, _ := newObjecteWhitelistTypeRepoMock(t)

			dbMock.ExpectBegin()
			row := sqlmock.NewRows([]string{"type_id", "extension"})

			row.AddRow(tt.args.typeID, tt.args.ext)

			dbMock.ExpectQuery("DELETE FROM \"object_whitelist_types\"").
				WithArgs(tt.args.typeID, tt.args.ext).
				WillReturnRows(row).
				WillReturnError(tt.mockErr)

			if tt.wantErr {
				dbMock.ExpectRollback()
			} else {
				dbMock.ExpectCommit()
			}
			if err := r.DeleteByTypeIDAndExt(context.TODO(), tt.args.typeID, tt.args.ext); (err != nil) != tt.wantErr {
				t.Errorf("objectWhitelistTypeRepository.DeleteByTypeIDAndExt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

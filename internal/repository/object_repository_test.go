package repository

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/goccy/go-json"
	"github.com/golang/mock/gomock"
	"github.com/krobus00/storage-service/internal/infrastructure"
	"github.com/krobus00/storage-service/internal/model"
	"github.com/krobus00/storage-service/internal/model/mock"
	"github.com/krobus00/storage-service/internal/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func newObjectRepoMock(t *testing.T) (model.ObjectRepository, sqlmock.Sqlmock, *miniredis.Miniredis) {
	dbConn, dbMock := utils.NewDBMock()
	miniRedis := miniredis.RunT(t)
	viper.Set("redis.cache_host", fmt.Sprintf("redis://%s", miniRedis.Addr()))
	redisClient, err := infrastructure.NewRedisClient()
	utils.ContinueOrFatal(err)
	objectRepo := NewObjectRepository()
	err = objectRepo.InjectDB(dbConn)
	utils.ContinueOrFatal(err)
	err = objectRepo.InjectRedisClient(redisClient)
	utils.ContinueOrFatal(err)

	return objectRepo, dbMock, miniRedis
}

func Test_objectRepository_Create(t *testing.T) {
	var (
		objectID = utils.GenerateUUID()
		userID   = utils.GenerateUUID()
	)
	type mockPutObject struct {
		res *s3.PutObjectOutput
		err error
	}
	type args struct {
		userID string
		data   *model.ObjectPayload
	}
	tests := []struct {
		name          string
		args          args
		mockPutObject *mockPutObject
		mockErr       error
		wantErr       bool
	}{
		{
			name: "success",
			args: args{
				userID: userID,
				data: &model.ObjectPayload{
					Object: &model.Object{
						ID:         objectID,
						FileName:   "test",
						Key:        "/object/test",
						UploadedBy: userID,
						IsPublic:   false,
						Type:       "image",
					},
				},
			},
			mockPutObject: &mockPutObject{
				res: &s3.PutObjectOutput{},
				err: nil,
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "error create object",
			args: args{
				userID: userID,
				data: &model.ObjectPayload{
					Object: &model.Object{
						ID:         objectID,
						FileName:   "test",
						Key:        "/object/test",
						UploadedBy: userID,
						IsPublic:   false,
						Type:       "image",
					},
				},
			},
			mockPutObject: &mockPutObject{
				res: &s3.PutObjectOutput{},
				err: nil,
			},
			mockErr: errors.New("db error"),
			wantErr: true,
		},
		{
			name: "error put object",
			args: args{
				userID: userID,
				data: &model.ObjectPayload{
					Object: &model.Object{
						ID:         objectID,
						FileName:   "test",
						Key:        "/object/test",
						UploadedBy: userID,
						IsPublic:   false,
						Type:       "image",
					},
				},
			},
			mockPutObject: &mockPutObject{
				res: &s3.PutObjectOutput{},
				err: errors.New("s3 error"),
			},
			mockErr: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.TODO()

			r, dbMock, _ := newObjectRepoMock(t)
			s3Client := mock.NewMockS3Client(ctrl)
			err := r.InjectS3Client(s3Client)
			utils.ContinueOrFatal(err)

			file, err := os.Open("../../tests/sample-image.png")
			utils.ContinueOrFatal(err)
			defer file.Close()

			buf := bytes.NewBuffer(nil)
			_, err = io.Copy(buf, file)
			utils.ContinueOrFatal(err)

			tt.args.data.Src = buf.Bytes()

			if tt.mockPutObject != nil {
				s3Client.EXPECT().
					PutObject(gomock.Any(), gomock.Any()).
					Times(1).Return(tt.mockPutObject.res, tt.mockPutObject.err)
			}

			object := tt.args.data.Object

			dbMock.ExpectBegin()
			dbMock.ExpectExec("INSERT INTO \"objects\"").
				WithArgs(object.ID, fmt.Sprintf("%s.png", object.FileName), fmt.Sprintf("%s.png", object.Key), object.UploadedBy, object.IsPublic, object.TypeID, sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.mockErr)

			if tt.wantErr {
				dbMock.ExpectRollback()
			} else {
				dbMock.ExpectCommit()
			}

			if err := r.Create(ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("objectRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_objectRepository_FindByID(t *testing.T) {
	var (
		objectID = utils.GenerateUUID()
		userID   = utils.GenerateUUID()
		typeID   = utils.GenerateUUID()
	)
	type args struct {
		id string
	}
	type mockSelect struct {
		object *model.Object
		err    error
	}
	type mockCache struct {
		object *model.Object
	}
	tests := []struct {
		name       string
		args       args
		mockSelect *mockSelect
		mockCache  *mockCache
		want       *model.Object
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				id: objectID,
			},
			mockSelect: &mockSelect{
				object: &model.Object{
					ID:         objectID,
					FileName:   "test.jpg",
					Key:        "/object/test.jpg",
					UploadedBy: userID,
					IsPublic:   false,
					TypeID:     typeID,
				},
				err: nil,
			},
			mockCache: nil,
			want: &model.Object{
				ID:         objectID,
				FileName:   "test.jpg",
				Key:        "/object/test.jpg",
				UploadedBy: userID,
				IsPublic:   false,
				TypeID:     typeID,
			},
			wantErr: false,
		},
		{
			name: "success found in cache",
			args: args{
				id: objectID,
			},
			mockSelect: nil,
			mockCache: &mockCache{
				object: &model.Object{
					ID:         objectID,
					FileName:   "test.jpg",
					Key:        "/object/test.jpg",
					UploadedBy: userID,
					IsPublic:   false,
					TypeID:     typeID,
				},
			},
			want: &model.Object{
				ID:         objectID,
				FileName:   "test.jpg",
				Key:        "/object/test.jpg",
				UploadedBy: userID,
				IsPublic:   false,
				TypeID:     typeID,
			},
			wantErr: false,
		},
		{
			name: "error object not found",
			args: args{
				id: objectID,
			},
			mockSelect: &mockSelect{
				object: nil,
				err:    gorm.ErrRecordNotFound,
			},
			mockCache: nil,
			want:      nil,
			wantErr:   false,
		},
		{
			name: "error find object",
			args: args{
				id: objectID,
			},
			mockSelect: &mockSelect{
				object: nil,
				err:    errors.New("db error"),
			},
			mockCache: nil,
			want:      nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			r, dbMock, redisMock := newObjectRepoMock(t)

			cacheKey := model.NewObjectCacheKey(tt.args.id)
			if tt.mockSelect != nil {
				row := sqlmock.NewRows([]string{"id", "file_name", "key", "uploaded_by", "is_public", "type_id", "created_at"})
				if tt.mockSelect.object != nil {
					object := tt.mockSelect.object
					row.AddRow(
						object.ID,
						object.FileName,
						object.Key,
						object.UploadedBy,
						object.IsPublic,
						object.TypeID,
						object.CreatedAt,
					)
				}

				dbMock.ExpectQuery("^SELECT .+ FROM \"objects\"").
					WithArgs(tt.args.id).
					WillReturnRows(row).
					WillReturnError(tt.mockSelect.err)
			}
			if tt.mockCache != nil {
				cacheData, err := json.Marshal(tt.mockCache.object)
				if err != nil {
					utils.ContinueOrFatal(err)
				}
				_ = redisMock.Set(cacheKey, string(cacheData))
			}

			got, err := r.FindByID(ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("objectRepository.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("objectRepository.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_objectRepository_GeneratePresignedURL(t *testing.T) {
	var (
		objectID = utils.GenerateUUID()
		userID   = utils.GenerateUUID()
		typeID   = utils.GenerateUUID()
	)
	type mockPresignGetObject struct {
		res *v4.PresignedHTTPRequest
		err error
	}
	type mockCache struct {
		res *model.GetPresignedURLResponse
	}
	type args struct {
		object *model.Object
	}
	tests := []struct {
		name                 string
		args                 args
		mockPresignGetObject *mockPresignGetObject
		mockCache            *mockCache
		want                 *model.GetPresignedURLResponse
		wantErr              bool
	}{
		{
			name: "success",
			args: args{
				object: &model.Object{
					ID:         objectID,
					FileName:   "test.jpg",
					Key:        "/object/test.jpg",
					UploadedBy: userID,
					IsPublic:   false,
					TypeID:     typeID,
				},
			},
			mockPresignGetObject: &mockPresignGetObject{
				res: &v4.PresignedHTTPRequest{
					URL: "https://s3.bucket/test.jpg",
				},
				err: nil,
			},
			mockCache: nil,
			want: &model.GetPresignedURLResponse{
				ID:         objectID,
				Filename:   "test.jpg",
				Type:       "image",
				URL:        "https://s3.bucket/test.jpg",
				IsPublic:   false,
				UploadedBy: userID,
			},
			wantErr: false,
		},
		{
			name: "success found in cache",
			args: args{
				object: &model.Object{
					ID:         objectID,
					FileName:   "test.jpg",
					Key:        "/object/test.jpg",
					UploadedBy: userID,
					IsPublic:   false,
					TypeID:     typeID,
				},
			},
			mockPresignGetObject: nil,
			mockCache: &mockCache{
				res: &model.GetPresignedURLResponse{
					ID:         objectID,
					Filename:   "test.jpg",
					Type:       "image",
					URL:        "https://s3.bucket/test.jpg",
					IsPublic:   false,
					UploadedBy: userID,
				},
			},
			want: &model.GetPresignedURLResponse{
				ID:         objectID,
				Filename:   "test.jpg",
				Type:       "image",
				URL:        "https://s3.bucket/test.jpg",
				IsPublic:   false,
				UploadedBy: userID,
			},
			wantErr: false,
		},
		{
			name: "error find object",
			args: args{
				object: &model.Object{
					ID:         objectID,
					FileName:   "test.jpg",
					Key:        "/object/test.jpg",
					UploadedBy: userID,
					IsPublic:   false,
					TypeID:     typeID,
				},
			},
			mockPresignGetObject: &mockPresignGetObject{
				res: nil,
				err: errors.New("s3 error"),
			},
			mockCache: nil,
			want:      nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.TODO()
			cacheKey := model.NewObjectPresignedURLCacheKey(tt.args.object.ID)

			r, _, redisMock := newObjectRepoMock(t)
			s3Client := mock.NewMockS3Client(ctrl)
			err := r.InjectS3Client(s3Client)
			utils.ContinueOrFatal(err)

			if tt.mockPresignGetObject != nil {
				s3Client.EXPECT().
					PresignGetObject(gomock.Any(), gomock.Any()).
					Times(1).
					Return(tt.mockPresignGetObject.res, tt.mockPresignGetObject.err)
			}

			if tt.mockCache != nil {
				cacheData, err := json.Marshal(tt.mockCache.res)
				if err != nil {
					utils.ContinueOrFatal(err)
				}
				_ = redisMock.Set(cacheKey, string(cacheData))
			}

			got, err := r.GeneratePresignedURL(ctx, tt.args.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("objectRepository.GeneratePresignedURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !assert.Equal(t, got.URL, tt.want.URL) {
					t.Errorf("objectRepository.GeneratePresignedURL() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

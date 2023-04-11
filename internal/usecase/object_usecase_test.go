package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	authMock "github.com/krobus00/auth-service/pb/auth/mock"
	"github.com/krobus00/storage-service/internal/constant"
	"github.com/krobus00/storage-service/internal/model"
	"github.com/krobus00/storage-service/internal/model/mock"
	"github.com/krobus00/storage-service/internal/utils"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func Test_objectUsecase_Upload(t *testing.T) {
	var (
		userID   = utils.GenerateUUID()
		objectID = utils.GenerateUUID()
		typeID   = utils.GenerateUUID()
	)
	type args struct {
		userID  string
		payload *model.ObjectPayload
	}
	type mockHasAccess struct {
		err       error
		hasAccess *wrapperspb.BoolValue
	}
	type mockFindObjectType struct {
		res *model.ObjectType
		err error
	}
	type mockFindByTypeIDAndExt struct {
		res *model.ObjectWhitelistType
		err error
	}
	type mockCreate struct {
		res *model.Object
		err error
	}
	tests := []struct {
		name                   string
		args                   args
		mockHasAccess          *mockHasAccess
		mockFindObjectType     *mockFindObjectType
		mockFindByTypeIDAndExt *mockFindByTypeIDAndExt
		mockCreate             *mockCreate
		want                   *model.Object
		wantErr                bool
	}{
		{
			name: "success",
			args: args{
				userID: userID,
				payload: &model.ObjectPayload{
					Object: &model.Object{
						Type:     "image",
						FileName: "test",
						IsPublic: false,
					},
				},
			},
			mockHasAccess: &mockHasAccess{
				hasAccess: wrapperspb.Bool(true),
				err:       nil,
			},
			mockFindObjectType: &mockFindObjectType{
				res: &model.ObjectType{
					ID:   typeID,
					Name: "image",
				},
				err: nil,
			},
			mockFindByTypeIDAndExt: &mockFindByTypeIDAndExt{
				res: &model.ObjectWhitelistType{
					TypeID:    typeID,
					Extension: ".png",
				},
				err: nil,
			},
			mockCreate: &mockCreate{
				res: &model.Object{
					ID:         objectID,
					FileName:   "test.png",
					Key:        "/object/test.png",
					UploadedBy: userID,
					IsPublic:   false,
					TypeID:     typeID,
					Type:       "image",
				},
				err: nil,
			},
			want: &model.Object{
				ID:         objectID,
				FileName:   "test.png",
				Key:        "/object/test.png",
				UploadedBy: userID,
				IsPublic:   false,
				TypeID:     typeID,
				Type:       "image",
			},
			wantErr: false,
		},
		{
			name: "error unauthorized access",
			args: args{
				userID: userID,
				payload: &model.ObjectPayload{
					Object: &model.Object{
						Type:     "image",
						FileName: "test",
						IsPublic: false,
					},
				},
			},
			mockHasAccess: &mockHasAccess{
				hasAccess: wrapperspb.Bool(false),
				err:       model.ErrUnauthorizeAccess,
			},
			wantErr: true,
		},
		{
			name: "error object type not found",
			args: args{
				userID: userID,
				payload: &model.ObjectPayload{
					Object: &model.Object{
						Type:     "image",
						FileName: "test",
						IsPublic: false,
					},
				},
			},
			mockHasAccess: &mockHasAccess{
				hasAccess: wrapperspb.Bool(true),
				err:       nil,
			},
			mockFindObjectType: &mockFindObjectType{
				res: nil,
				err: nil,
			},
			wantErr: true,
		},
		{
			name: "error find object type",
			args: args{
				userID: userID,
				payload: &model.ObjectPayload{
					Object: &model.Object{
						Type:     "image",
						FileName: "test",
						IsPublic: false,
					},
				},
			},
			mockHasAccess: &mockHasAccess{
				hasAccess: wrapperspb.Bool(true),
				err:       nil,
			},
			mockFindObjectType: &mockFindObjectType{
				res: nil,
				err: errors.New("db error"),
			},
			wantErr: true,
		},
		{
			name: "error file extension not allowed",
			args: args{
				userID: userID,
				payload: &model.ObjectPayload{
					Object: &model.Object{
						Type:     "image",
						FileName: "test",
						IsPublic: false,
					},
				},
			},
			mockHasAccess: &mockHasAccess{
				hasAccess: wrapperspb.Bool(true),
				err:       nil,
			},
			mockFindObjectType: &mockFindObjectType{
				res: &model.ObjectType{
					ID:   typeID,
					Name: "image",
				},
				err: nil,
			},
			mockFindByTypeIDAndExt: &mockFindByTypeIDAndExt{
				res: nil,
				err: nil,
			},
			wantErr: true,
		},
		{
			name: "error find extension whitelis",
			args: args{
				userID: userID,
				payload: &model.ObjectPayload{
					Object: &model.Object{
						Type:     "image",
						FileName: "test",
						IsPublic: false,
					},
				},
			},
			mockHasAccess: &mockHasAccess{
				hasAccess: wrapperspb.Bool(true),
				err:       nil,
			},
			mockFindObjectType: &mockFindObjectType{
				res: &model.ObjectType{
					ID:   typeID,
					Name: "image",
				},
				err: nil,
			},
			mockFindByTypeIDAndExt: &mockFindByTypeIDAndExt{
				res: nil,
				err: errors.New("db error"),
			},
			wantErr: true,
		},
		{
			name: "error create object",
			args: args{
				userID: userID,
				payload: &model.ObjectPayload{
					Object: &model.Object{
						Type:     "image",
						FileName: "test",
						IsPublic: false,
					},
				},
			},
			mockHasAccess: &mockHasAccess{
				hasAccess: wrapperspb.Bool(true),
				err:       nil,
			},
			mockFindObjectType: &mockFindObjectType{
				res: &model.ObjectType{
					ID:   typeID,
					Name: "image",
				},
				err: nil,
			},
			mockFindByTypeIDAndExt: &mockFindByTypeIDAndExt{
				res: &model.ObjectWhitelistType{
					TypeID:    typeID,
					Extension: ".png",
				},
				err: nil,
			},
			mockCreate: &mockCreate{
				res: nil,
				err: errors.New("db error"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.TODO()
			ctx = context.WithValue(ctx, constant.KeyUserIDCtx, tt.args.userID)

			objectRepo := mock.NewMockObjectRepository(ctrl)
			objectTypeRepo := mock.NewMockObjectTypeRepository(ctrl)
			objectWhitelistTypeRepo := mock.NewMockObjectWhitelistTypeRepository(ctrl)
			authClientMock := authMock.NewMockAuthServiceClient(ctrl)

			if tt.mockHasAccess != nil {
				authClientMock.EXPECT().
					HasAccess(gomock.Any(), gomock.Any()).
					Times(1).
					Return(tt.mockHasAccess.hasAccess, tt.mockHasAccess.err)
			}

			if tt.mockFindObjectType != nil {
				objectTypeRepo.EXPECT().
					FindByName(gomock.Any(), tt.args.payload.Object.Type).
					Times(1).
					Return(tt.mockFindObjectType.res, tt.mockFindObjectType.err)
			}

			if tt.mockFindByTypeIDAndExt != nil {
				objectWhitelistTypeRepo.EXPECT().
					FindByTypeIDAndExt(gomock.Any(), tt.mockFindObjectType.res.ID, gomock.Any()).
					Times(1).
					Return(tt.mockFindByTypeIDAndExt.res, tt.mockFindByTypeIDAndExt.err)
			}

			if tt.mockCreate != nil {
				objectRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(1).
					DoAndReturn(func(ctx context.Context, data *model.ObjectPayload) error {
						data.Object = tt.mockCreate.res
						return tt.mockCreate.err
					})
			}

			uc := NewObjectUsecase()
			err := uc.InjectObjectRepo(objectRepo)
			utils.ContinueOrFatal(err)
			err = uc.InjectObjectTypeRepo(objectTypeRepo)
			utils.ContinueOrFatal(err)
			err = uc.InjectObjectWhitelistTypeRepo(objectWhitelistTypeRepo)
			utils.ContinueOrFatal(err)
			err = uc.InjectAuthClient(authClientMock)
			utils.ContinueOrFatal(err)

			got, err := uc.Upload(ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("objectUsecase.Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("objectUsecase.Upload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_objectUsecase_GeneratePresignedURL(t *testing.T) {
	var (
		userID   = utils.GenerateUUID()
		objectID = utils.GenerateUUID()
		typeID   = utils.GenerateUUID()
	)
	type mockHasAccess struct {
		err       error
		hasAccess *wrapperspb.BoolValue
	}
	type mockFindObjectByID struct {
		res *model.Object
		err error
	}
	type mockFindObjectType struct {
		res *model.ObjectType
		err error
	}
	type mockGeneratePresignedURL struct {
		res *model.GetPresignedURLResponse
		err error
	}
	type args struct {
		userID  string
		payload *model.GetPresignedURLPayload
	}
	tests := []struct {
		name                     string
		args                     args
		mockHasAccess            *mockHasAccess
		mockFindObjectByID       *mockFindObjectByID
		mockFindObjectType       *mockFindObjectType
		mockGeneratePresignedURL *mockGeneratePresignedURL
		want                     *model.GetPresignedURLResponse
		wantErr                  bool
	}{
		{
			name: "success",
			args: args{
				userID: userID,
				payload: &model.GetPresignedURLPayload{
					ObjectID: objectID,
				},
			},
			mockHasAccess: &mockHasAccess{
				hasAccess: wrapperspb.Bool(true),
				err:       nil,
			},
			mockFindObjectByID: &mockFindObjectByID{
				res: &model.Object{
					ID:         objectID,
					UploadedBy: userID,
					Type:       typeID,
				},
				err: nil,
			},
			mockFindObjectType: &mockFindObjectType{
				res: &model.ObjectType{
					ID:   typeID,
					Name: "image",
				},
			},
			mockGeneratePresignedURL: &mockGeneratePresignedURL{
				res: &model.GetPresignedURLResponse{
					ID:         objectID,
					Filename:   "test.png",
					Type:       "image",
					URL:        "https://s3.bucket/test.jpg",
					IsPublic:   false,
					UploadedBy: userID,
				},
			},
			want: &model.GetPresignedURLResponse{
				ID:         objectID,
				Filename:   "test.png",
				Type:       "image",
				URL:        "https://s3.bucket/test.jpg",
				IsPublic:   false,
				UploadedBy: userID,
			},
			wantErr: false,
		},
		{
			name: "success get other user private object",
			args: args{
				userID: userID,
				payload: &model.GetPresignedURLPayload{
					ObjectID: objectID,
				},
			},
			mockHasAccess: &mockHasAccess{
				hasAccess: wrapperspb.Bool(true),
				err:       nil,
			},
			mockFindObjectByID: &mockFindObjectByID{
				res: &model.Object{
					ID:         objectID,
					UploadedBy: "other-user",
					Type:       typeID,
					IsPublic:   false,
				},
				err: nil,
			},
			mockFindObjectType: &mockFindObjectType{
				res: &model.ObjectType{
					ID:   typeID,
					Name: "image",
				},
			},
			mockGeneratePresignedURL: &mockGeneratePresignedURL{
				res: &model.GetPresignedURLResponse{
					ID:         objectID,
					Filename:   "test.png",
					Type:       "image",
					URL:        "https://s3.bucket/test.jpg",
					IsPublic:   false,
					UploadedBy: userID,
				},
			},
			want: &model.GetPresignedURLResponse{
				ID:         objectID,
				Filename:   "test.png",
				Type:       "image",
				URL:        "https://s3.bucket/test.jpg",
				IsPublic:   false,
				UploadedBy: userID,
			},
			wantErr: false,
		},
		{
			name: "success get other user public object",
			args: args{
				userID: userID,
				payload: &model.GetPresignedURLPayload{
					ObjectID: objectID,
				},
			},
			mockHasAccess: &mockHasAccess{
				hasAccess: wrapperspb.Bool(true),
				err:       nil,
			},
			mockFindObjectByID: &mockFindObjectByID{
				res: &model.Object{
					ID:         objectID,
					UploadedBy: "other-user",
					Type:       typeID,
					IsPublic:   true,
				},
				err: nil,
			},
			mockFindObjectType: &mockFindObjectType{
				res: &model.ObjectType{
					ID:   typeID,
					Name: "image",
				},
			},
			mockGeneratePresignedURL: &mockGeneratePresignedURL{
				res: &model.GetPresignedURLResponse{
					ID:         objectID,
					Filename:   "test.png",
					Type:       "image",
					URL:        "https://s3.bucket/test.jpg",
					IsPublic:   false,
					UploadedBy: userID,
				},
			},
			want: &model.GetPresignedURLResponse{
				ID:         objectID,
				Filename:   "test.png",
				Type:       "image",
				URL:        "https://s3.bucket/test.jpg",
				IsPublic:   false,
				UploadedBy: userID,
			},
			wantErr: false,
		},
		{
			name: "error object not found",
			args: args{
				userID: userID,
				payload: &model.GetPresignedURLPayload{
					ObjectID: objectID,
				},
			},
			mockFindObjectByID: &mockFindObjectByID{
				res: nil,
				err: nil,
			},
			wantErr: true,
		},
		{
			name: "error find object",
			args: args{
				userID: userID,
				payload: &model.GetPresignedURLPayload{
					ObjectID: objectID,
				},
			},
			mockFindObjectByID: &mockFindObjectByID{
				res: nil,
				err: errors.New("db error"),
			},
			wantErr: true,
		},
		{
			name: "error unauthorized object access",
			args: args{
				userID: userID,
				payload: &model.GetPresignedURLPayload{
					ObjectID: objectID,
				},
			},
			mockHasAccess: &mockHasAccess{
				hasAccess: wrapperspb.Bool(false),
				err:       nil,
			},
			mockFindObjectByID: &mockFindObjectByID{
				res: &model.Object{
					ID:         objectID,
					UploadedBy: "other-user",
					Type:       typeID,
					IsPublic:   false,
				},
				err: nil,
			},
			wantErr: true,
		},
		{
			name: "error object type not found",
			args: args{
				userID: userID,
				payload: &model.GetPresignedURLPayload{
					ObjectID: objectID,
				},
			},
			mockHasAccess: &mockHasAccess{
				hasAccess: wrapperspb.Bool(true),
				err:       nil,
			},
			mockFindObjectByID: &mockFindObjectByID{
				res: &model.Object{
					ID:         objectID,
					UploadedBy: userID,
					Type:       typeID,
				},
				err: nil,
			},
			mockFindObjectType: &mockFindObjectType{
				res: nil,
				err: nil,
			},
			wantErr: true,
		},
		{
			name: "error find object type",
			args: args{
				userID: userID,
				payload: &model.GetPresignedURLPayload{
					ObjectID: objectID,
				},
			},
			mockHasAccess: &mockHasAccess{
				hasAccess: wrapperspb.Bool(true),
				err:       nil,
			},
			mockFindObjectByID: &mockFindObjectByID{
				res: &model.Object{
					ID:         objectID,
					UploadedBy: userID,
					Type:       typeID,
				},
				err: nil,
			},
			mockFindObjectType: &mockFindObjectType{
				res: nil,
				err: errors.New("db error"),
			},
			wantErr: true,
		},
		{
			name: "error sign object",
			args: args{
				userID: userID,
				payload: &model.GetPresignedURLPayload{
					ObjectID: objectID,
				},
			},
			mockHasAccess: &mockHasAccess{
				hasAccess: wrapperspb.Bool(true),
				err:       nil,
			},
			mockFindObjectByID: &mockFindObjectByID{
				res: &model.Object{
					ID:         objectID,
					UploadedBy: userID,
					Type:       typeID,
				},
				err: nil,
			},
			mockFindObjectType: &mockFindObjectType{
				res: &model.ObjectType{
					ID:   typeID,
					Name: "image",
				},
			},
			mockGeneratePresignedURL: &mockGeneratePresignedURL{
				res: nil,
				err: errors.New("db error"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.TODO()
			ctx = context.WithValue(ctx, constant.KeyUserIDCtx, tt.args.userID)

			objectRepo := mock.NewMockObjectRepository(ctrl)
			objectTypeRepo := mock.NewMockObjectTypeRepository(ctrl)
			authClientMock := authMock.NewMockAuthServiceClient(ctrl)

			if tt.mockFindObjectByID != nil {
				objectRepo.EXPECT().
					FindByID(gomock.Any(), tt.args.payload.ObjectID).
					Times(1).
					Return(tt.mockFindObjectByID.res, tt.mockFindObjectByID.err)

				object := tt.mockFindObjectByID.res
				if object != nil {
					if tt.mockHasAccess != nil && !object.IsPublic && object.UploadedBy != tt.args.userID {
						authClientMock.EXPECT().
							HasAccess(gomock.Any(), gomock.Any()).
							Times(1).
							Return(tt.mockHasAccess.hasAccess, tt.mockHasAccess.err)
					}
				}
			}

			if tt.mockFindObjectType != nil {
				objectTypeRepo.EXPECT().
					FindByID(gomock.Any(), tt.mockFindObjectByID.res.TypeID).
					Times(1).
					Return(tt.mockFindObjectType.res, tt.mockFindObjectType.err)
			}

			if tt.mockGeneratePresignedURL != nil {
				objectRepo.EXPECT().
					GeneratePresignedURL(gomock.Any(), tt.mockFindObjectByID.res).
					Times(1).
					Return(tt.mockGeneratePresignedURL.res, tt.mockGeneratePresignedURL.err)
			}

			uc := NewObjectUsecase()
			err := uc.InjectObjectRepo(objectRepo)
			utils.ContinueOrFatal(err)
			err = uc.InjectObjectTypeRepo(objectTypeRepo)
			utils.ContinueOrFatal(err)
			err = uc.InjectAuthClient(authClientMock)
			utils.ContinueOrFatal(err)

			got, err := uc.GeneratePresignedURL(ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("objectUsecase.GeneratePresignedURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("objectUsecase.GeneratePresignedURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

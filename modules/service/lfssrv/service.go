package lfssrv

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"regexp"
	"time"
	"zgit/modules/model/lfsmd"
	"zgit/modules/model/repomd"
	"zgit/modules/model/usermd"
	"zgit/pkg/lfs"
	"zgit/setting"
	"zgit/util"
)

var (
	oidPattern = regexp.MustCompile(`^[a-f\d]{64}$`)
)

type LockReqDTO struct {
	Repo     repomd.RepoInfo
	Operator usermd.UserInfo
}

type ListLockReqDTO struct {
	Repo     repomd.RepoInfo
	Operator usermd.UserInfo
	Path     string
	Cursor   string
	Limit    int
	RefName  string
}

type ListLockRespDTO struct {
	LockList []lfsmd.LfsLock
	Next     string
}

type UnlockReqDTO struct {
	Repo     repomd.RepoInfo
	LockId   int64
	Force    bool
	Operator usermd.UserInfo
}

type PointerDTO struct {
	Oid  string
	Size int64
}

type VerifyReqDTO struct {
	PointerDTO
	Repo     repomd.RepoInfo
	Operator usermd.UserInfo
}

type DownloadReqDTO struct {
	Oid      string
	Repo     repomd.RepoInfo
	Operator usermd.UserInfo
	FromByte int64
	ToByte   int64
}

type UploadReqDTO struct {
	Oid      string
	Size     int64
	Repo     repomd.RepoInfo
	Operator usermd.UserInfo
	Body     io.Reader
}

type DownloadRespDTO struct {
	io.ReadCloser
	FromByte int64
	ToByte   int64
	Length   int64
}

type BatchReqDTO struct {
	Repo     repomd.RepoInfo
	Operator usermd.UserInfo
	Objects  []PointerDTO
	IsUpload bool
	RefName  string
}

type LinkDTO struct {
	Href      string
	Header    map[string]string
	ExpiresAt *time.Time
}

// ObjectErrDTO defines the JSON structure returned to the client in case of an error.
type ObjectErrDTO struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ObjectDTO struct {
	PointerDTO
	Err error
}

type BatchRespDTO struct {
	ObjectList []ObjectDTO
}

func Lock(ctx context.Context, req LockReqDTO) (lfsmd.LfsLock, error) {
	// todo 获取仓库
	// todo 转换jwt token
	// todo 检查权限
	// todo 添加锁
	return lfsmd.LfsLock{
		Id:      1,
		RepoId:  "aa",
		OwnerId: "1",
		Path:    "xx",
		Created: time.Now(),
	}, nil
}

func ListLock(ctx context.Context, req ListLockReqDTO) (ListLockRespDTO, error) {
	// todo 获取仓库
	// todo 转换jwt token
	// todo 检查权限
	// todo 添加锁
	return ListLockRespDTO{
		LockList: []lfsmd.LfsLock{},
		Next:     "",
	}, nil
}

func Unlock(ctx context.Context, req UnlockReqDTO) (lfsmd.LfsLock, error) {
	return lfsmd.LfsLock{
		Id:      1,
		RepoId:  "aa",
		OwnerId: "1",
		Path:    "xx",
		Created: time.Now(),
	}, nil
}

func Verify(ctx context.Context, req VerifyReqDTO) error {
	object, err := lfs.StorageImpl.Stat(ctx, convertPointerPath(req.Oid))
	if err != nil {
		return err
	}
	if object.Size() != req.Size {
		return errors.New("invalid size")
	}
	return nil
}

func Download(ctx context.Context, reqDTO DownloadReqDTO) (DownloadRespDTO, error) {
	object, err := lfs.StorageImpl.Open(ctx, reqDTO.Oid)
	if err != nil {
		return DownloadRespDTO{}, err
	}
	if reqDTO.FromByte < 0 {
		reqDTO.FromByte = 0
	}
	if reqDTO.ToByte < reqDTO.FromByte {
		reqDTO.ToByte = reqDTO.FromByte
	}
	stat, err := object.Stat()
	if err != nil {
		return DownloadRespDTO{}, err
	}
	if reqDTO.ToByte > stat.Size() {
		reqDTO.ToByte = stat.Size()
	}
	if reqDTO.FromByte > 0 {
		_, err = object.Seek(reqDTO.FromByte, io.SeekStart)
		if err != nil {
			return DownloadRespDTO{}, err
		}
	}
	return DownloadRespDTO{
		ReadCloser: object,
		FromByte:   reqDTO.FromByte,
		ToByte:     reqDTO.ToByte,
		Length:     reqDTO.ToByte + 1 - reqDTO.FromByte,
	}, nil
}

func Upload(ctx context.Context, reqDTO UploadReqDTO) error {
	_, err := lfs.StorageImpl.Save(ctx, convertPointerPath(reqDTO.Oid), reqDTO.Body)
	return err
}

func convertPointerPath(oid string) string {
	if len(oid) < 5 {
		return oid
	}
	return path.Join(oid[0:2], oid[2:4], oid[4:])
}

func IsValidPointer(oid string, size int64) error {
	if len(oid) != 64 {
		return errors.New("oid length should equals 64")
	}
	if !oidPattern.MatchString(oid) {
		return errors.New("oid format error")
	}
	if size < 0 {
		return errors.New("size should greater than 0")
	}
	return nil
}

func Batch(ctx context.Context, req BatchReqDTO) (BatchRespDTO, error) {
	ret := make([]ObjectDTO, 0, len(req.Objects))
	for _, object := range req.Objects {
		if err := IsValidPointer(object.Oid, object.Size); err != nil {
			return BatchRespDTO{}, fmt.Errorf("%s format err: %v", object.Oid, err)
		}
		meta, b, err := lfsmd.GetMetaObjectByOid(ctx, object.Oid)
		if err != nil {
			return BatchRespDTO{},
				fmt.Errorf("find %s err: %v",
					object.Oid,
					util.VolumeReadable(setting.MaxLfsFileSize()),
				)
		}
		if b && meta.Size != object.Size {
			// 大小不一致
			return BatchRespDTO{}, fmt.Errorf("%s size err", object.Oid)
		}
		// 文件存在 但没有落库
		exists, _ := lfs.StorageImpl.Exists(ctx, convertPointerPath(object.Oid))
		if req.IsUpload {
			if !exists && object.Size > setting.MaxLfsFileSize() {
				return BatchRespDTO{},
					fmt.Errorf("%s file size exceeded: %v",
						object.Oid,
						util.VolumeReadable(setting.MaxLfsFileSize()),
					)
			}
			if exists && !b {
				if err = lfsmd.InsertMetaObject(lfsmd.MetaObject{
					RepoId: req.Repo.Id,
					Oid:    object.Oid,
					Size:   object.Size,
				}); err != nil {
					return BatchRespDTO{}, fmt.Errorf("insert %s err", object.Oid)
				}
			}
			ret = append(ret, ObjectDTO{
				PointerDTO: object,
			})
		} else {
			if !exists || !b {
				ret = append(ret, ObjectDTO{
					Err: errors.New("not found"),
				})
			} else {
				ret = append(ret, ObjectDTO{
					PointerDTO: object,
				})
			}
		}
	}
	return BatchRespDTO{
		ObjectList: ret,
	}, nil
}

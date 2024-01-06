package lfssrv

import (
	"context"
	"errors"
	"fmt"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"io"
	"path/filepath"
	"zgit/pkg/git/lfs"
	"zgit/pkg/i18n"
	"zgit/pkg/perm"
	"zgit/standalone/modules/model/lfsmd"
	"zgit/standalone/modules/model/projectmd"
	"zgit/standalone/modules/model/repomd"
	"zgit/standalone/modules/model/usermd"
	"zgit/util"
)

func Lock(ctx context.Context, reqDTO LockReqDTO) (lfsmd.LfsLock, error) {
	if err := reqDTO.IsValid(); err != nil {
		return lfsmd.LfsLock{}, err
	}
	// 检查仓库访问权限
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	p, err := getPerm(ctx, reqDTO.Repo, reqDTO.Operator)
	if err != nil {
		return lfsmd.LfsLock{}, err
	}
	if !p.GetRepoPerm(reqDTO.Repo.RepoId).CanPush {
		return lfsmd.LfsLock{}, util.UnauthorizedError()
	}
	lock, err := lfsmd.InsertLock(ctx, lfsmd.InsertLockReqDTO{
		RepoId: reqDTO.Repo.RepoId,
		Owner:  reqDTO.Operator.Account,
		Path:   reqDTO.Path,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return lfsmd.LfsLock{}, util.InternalError()
	}
	// 添加锁
	return lock, nil
}

func ListLock(ctx context.Context, reqDTO ListLockReqDTO) (ListLockRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return ListLockRespDTO{}, err
	}
	// 检查仓库访问权限
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	p, err := getPerm(ctx, reqDTO.Repo, reqDTO.Operator)
	if err != nil {
		return ListLockRespDTO{}, err
	}
	if !p.GetRepoPerm(reqDTO.Repo.RepoId).CanAccess {
		return ListLockRespDTO{}, util.UnauthorizedError()
	}
	// 查询lock
	return ListLockRespDTO{
		LockList: []lfsmd.LfsLock{},
		Next:     "",
	}, nil
}

func Unlock(ctx context.Context, reqDTO UnlockReqDTO) (lfsmd.LfsLock, error) {
	if err := reqDTO.IsValid(); err != nil {
		return lfsmd.LfsLock{}, err
	}
	// 检查仓库访问权限
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	p, err := getPerm(ctx, reqDTO.Repo, reqDTO.Operator)
	if err != nil {
		return lfsmd.LfsLock{}, err
	}
	if !p.GetRepoPerm(reqDTO.Repo.RepoId).CanPush {
		return lfsmd.LfsLock{}, util.UnauthorizedError()
	}
	// 查找lock是否存在
	lock, b, err := lfsmd.GetLockById(ctx, reqDTO.LockId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return lfsmd.LfsLock{}, util.InternalError()
	}
	if !b {
		return lfsmd.LfsLock{}, util.InvalidArgsError()
	}
	_, err = lfsmd.DeleteLock(ctx, lock.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return lfsmd.LfsLock{}, util.InternalError()
	}
	return lock, nil
}

func Verify(ctx context.Context, reqDTO VerifyReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	// 检查仓库访问权限
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	p, err := getPerm(ctx, reqDTO.Repo, reqDTO.Operator)
	if err != nil {
		return err
	}
	if !p.GetRepoPerm(reqDTO.Repo.RepoId).CanAccess {
		return util.UnauthorizedError()
	}
	object, err := lfs.StorageImpl.Stat(ctx, convertPointerPath(reqDTO.Repo.Path, reqDTO.Oid))
	if err != nil {
		return err
	}
	if object.Size() != reqDTO.Size {
		return errors.New("invalid size")
	}
	return nil
}

func Download(ctx context.Context, reqDTO DownloadReqDTO) (DownloadRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return DownloadRespDTO{}, err
	}
	// 检查仓库访问权限
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	p, err := getPerm(ctx, reqDTO.Repo, reqDTO.Operator)
	if err != nil {
		return DownloadRespDTO{}, err
	}
	if !p.GetRepoPerm(reqDTO.Repo.RepoId).CanAccess {
		return DownloadRespDTO{}, util.UnauthorizedError()
	}
	object, err := lfs.StorageImpl.Open(ctx, convertPointerPath(reqDTO.Repo.Path, reqDTO.Oid))
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
	// 检查仓库访问权限
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	p, err := getPerm(ctx, reqDTO.Repo, reqDTO.Operator)
	if err != nil {
		return err
	}
	if !p.GetRepoPerm(reqDTO.Repo.RepoId).CanPush {
		return util.UnauthorizedError()
	}
	_, err = lfs.StorageImpl.Save(ctx, convertPointerPath(reqDTO.Repo.Path, reqDTO.Oid), reqDTO.Body)
	return err
}

func convertPointerPath(repoPath, oid string) string {
	if len(oid) < 5 {
		return filepath.Join(repoPath, oid)
	}
	return filepath.Join(repoPath, oid[0:2], oid[2:4], oid[4:])
}

func Batch(ctx context.Context, reqDTO BatchReqDTO) (BatchRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return BatchRespDTO{}, err
	}
	ret := make([]ObjectDTO, 0, len(reqDTO.Objects))
	for _, object := range reqDTO.Objects {
		meta, b, err := lfsmd.GetMetaObjectByOid(ctx, object.Oid)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return BatchRespDTO{}, util.InternalError()
		}
		if b && meta.Size != object.Size {
			logger.Logger.WithContext(ctx).Errorf("obj: %s size is not equal %d != %d", object.Oid, meta.Size, object.Size)
			// 大小不一致
			return BatchRespDTO{}, util.InvalidArgsError()
		}
		// 文件存在 但没有落库
		exists, _ := lfs.StorageImpl.Exists(ctx, convertPointerPath(reqDTO.Repo.Path, object.Oid))
		if reqDTO.IsUpload {
			// 检查是否超过单个lfs文件配置大小
			if !exists && reqDTO.Repo.Cfg.SingleLfsFileLimitSize > 0 && object.Size > reqDTO.Repo.Cfg.SingleLfsFileLimitSize {
				return BatchRespDTO{},
					fmt.Errorf(i18n.GetByKey(i18n.LfsExceedSingleFileLimitSize),
						object.Oid,
						util.VolumeReadable(object.Size),
						util.VolumeReadable(reqDTO.Repo.Cfg.SingleLfsFileLimitSize),
					)
			}
			if exists && !b {
				if err = lfsmd.InsertMetaObject(lfsmd.MetaObject{
					RepoId: reqDTO.Repo.Path,
					Oid:    object.Oid,
					Size:   object.Size,
				}); err != nil {
					logger.Logger.WithContext(ctx).Error(err)
					return BatchRespDTO{}, fmt.Errorf(i18n.GetByKey(i18n.SystemInternalError))
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

func getPerm(ctx context.Context, repo repomd.RepoInfo, operator usermd.UserInfo) (perm.Detail, error) {
	p, b, err := projectmd.GetProjectUserPermDetail(ctx, repo.ProjectId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return perm.Detail{}, util.InternalError()
	}
	if !b {
		return perm.Detail{}, util.UnauthorizedError()
	}
	return p.PermDetail, nil
}

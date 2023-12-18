package lfsapi

import (
	"encoding/base64"
	"fmt"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"zgit/modules/model/lfsmd"
	"zgit/modules/model/repomd"
	"zgit/modules/model/usermd"
	"zgit/modules/service/lfssrv"
	"zgit/modules/service/reposrv"
	"zgit/modules/service/usersrv"
	"zgit/pkg/lfs"
	"zgit/setting"
)

const (
	// MediaType contains the media type for LFS server requests
	MediaType = "application/vnd.git-lfs+json"
)

var (
	rangeHeaderRegexp = regexp.MustCompile(`bytes=(\d+)\-(\d*).*`)
)

type BatchReqVO struct {
	Operation string `json:"operation"`
	// 没使用到
	Transfers []string    `json:"transfers,omitempty"`
	Ref       ReferenceVO `json:"ref,omitempty"`
	Objects   []PointerVO `json:"objects"`
	HashAlgo  string      `json:"hash_algo"`
}

// ReferenceVO contains a git reference.
// https://github.com/git-lfs/git-lfs/blob/main/docs/api/batch.md#ref-property
type ReferenceVO struct {
	Name string `json:"name"`
}

// PointerVO contains LFS pointer data
type PointerVO struct {
	Oid  string `json:"oid"`
	Size int64  `json:"size"`
}

// BatchRespVO contains multiple object metadata Representation structures
// for use with the batch API.
// https://github.com/git-lfs/git-lfs/blob/main/docs/api/batch.md#successful-responses
type BatchRespVO struct {
	Transfer string         `json:"transfer,omitempty"`
	Objects  []ObjectRespVO `json:"objects"`
}

// ObjectRespVO is object metadata as seen by clients of the LFS server.
type ObjectRespVO struct {
	PointerVO
	Actions map[string]LinkVO `json:"actions,omitempty"`
	Error   *ObjectErrVO      `json:"error,omitempty"`
}

// LinkVO provides a structure with information about how to access a object.
type LinkVO struct {
	Href      string            `json:"href"`
	Header    map[string]string `json:"header,omitempty"`
	ExpiresAt *time.Time        `json:"expires_at,omitempty"`
}

// ObjectErrVO defines the JSON structure returned to the client in case of an error.
type ObjectErrVO struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// LockVO represent a lock
// for use with the locks API.
type LockVO struct {
	Id       string       `json:"id"`
	Path     string       `json:"path"`
	LockedAt time.Time    `json:"locked_at"`
	Owner    *LockOwnerVO `json:"owner"`
}

// LockOwnerVO represent a lock owner
// for use with the locks API.
type LockOwnerVO struct {
	Name string `json:"name"`
}

type ErrVO struct {
	Message       string `json:"message"`
	Documentation string `json:"documentation_url,omitempty"`
	RequestID     string `json:"request_id,omitempty"`
}

type PostLockReqVO struct {
	Path string      `json:"path"`
	Ref  ReferenceVO `json:"ref"`
}

type PostLockRespVO struct {
	Lock LockVO `json:"lock"`
}

type ListLockReqVO struct {
	Path    string `json:"path" form:"path"`
	Id      string `json:"id" form:"id"`
	Cursor  string `json:"cursor" form:"cursor"`
	Limit   int    `json:"limit" form:"limit"`
	RefSpec string `json:"refspec" form:"refspec"`
}

type ListLockVerifyReqVO struct {
	Cursor string      `json:"cursor"`
	Limit  int         `json:"limit"`
	Ref    ReferenceVO `json:"ref"`
}

type ListLockRespVO struct {
	Locks []LockVO `json:"locks"`
	Next  string   `json:"next_cursor,omitempty"`
}

type UnlockReqVO struct {
	Force bool `json:"force"`
}

type UnlockRespVO struct {
	Lock LockVO `json:"lock"`
}

type ListLockVerifyRespVO struct {
	Ours   []LockVO `json:"ours"`
	Theirs []LockVO `json:"theirs"`
	Next   string   `json:"next_cursor,omitempty"`
}

func InitLfsHttpApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		infoLfs := e.Group(":companyId/:clusterId/:repoName/info/lfs", packRepoPath)
		{
			infoLfs.POST("/objects/batch", checkMediaType, batch)
			infoLfs.PUT("/objects/:oid/:size", upload)
			infoLfs.GET("/objects/:oid/:filename", download)
			infoLfs.GET("/objects/:oid", download)
			infoLfs.POST("/verify", checkMediaType, verify)
			locks := infoLfs.Group("/locks", checkMediaType)
			{
				locks.GET("/", listLock)
				locks.POST("/", lock)
				locks.POST("/verify", listLockVerify)
				locks.POST("/:id/unlock", unlock)
			}
		}
	})
}

func packRepoPath(c *gin.Context) {
	companyId := c.Param("companyId")
	clusterId := c.Param("clusterId")
	repoName := c.Param("repoName")
	repoPath := filepath.Join(companyId, clusterId, repoName)
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		c.JSON(http.StatusBadRequest, ErrVO{
			Message: "auth not found",
		})
		c.Abort()
		return
	}
	token, err := jwt.ParseWithClaims(authorization, new(lfs.Claims), func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return setting.LfsJwtSecretBytes(), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrVO{
			Message: err.Error(),
		})
		c.Abort()
		return
	}
	claims, ok := token.Claims.(*lfs.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrVO{
			Message: "invalid token",
		})
		c.Abort()
		return
	}
	ctx := c.Request.Context()
	repo, b, err := reposrv.GetRepoInfoByRelativePath(ctx, repoPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrVO{
			Message: "internal error",
		})
		c.Abort()
		return
	}
	if !b {
		c.JSON(http.StatusUnauthorized, ErrVO{
			Message: "unknown repo",
		})
		c.Abort()
		return
	}
	userInfo, b, err := usersrv.GetUserInfoByUserId(ctx, claims.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrVO{
			Message: "internal error",
		})
		c.Abort()
		return
	}
	if !b {
		c.JSON(http.StatusUnauthorized, ErrVO{
			Message: "invalid token",
		})
		c.Abort()
		return
	}
	c.Set("operator", userInfo)
	c.Set("claims", claims)
	c.Set("authorization", authorization)
	c.Set("repo", repo)
	c.Next()
}

func checkMediaType(c *gin.Context) {
	header := c.GetHeader("Accept")
	accepts := strings.Split(header, ";")
	if len(accepts) == 0 || accepts[0] != MediaType {
		c.JSON(http.StatusUnsupportedMediaType, ErrVO{
			Message: "unsupported media type",
		})
		c.Abort()
		return
	} else {
		c.Next()
	}
}

func batch(c *gin.Context) {
	var req BatchReqVO
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrVO{
			Message: "bad request",
		})
		return
	}
	var isUpload bool
	if req.Operation == "upload" {
		isUpload = true
	} else if req.Operation == "download" {
		isUpload = false
	} else {
		c.JSON(http.StatusBadRequest, ErrVO{
			Message: "bad request",
		})
		return
	}
	ctx := c.Request.Context()
	reqDTO := lfssrv.BatchReqDTO{
		Repo:     getRepo(c),
		Operator: getOperator(c),
		IsUpload: isUpload,
		RefName:  req.Ref.Name,
	}
	reqDTO.Objects, _ = listutil.Map(req.Objects, func(t PointerVO) (lfssrv.PointerDTO, error) {
		return lfssrv.PointerDTO{
			Oid:  t.Oid,
			Size: t.Size,
		}, nil
	})
	respDTO, err := lfssrv.Batch(ctx, reqDTO)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrVO{
			Message: err.Error(),
		})
		return
	}
	authorization := c.MustGet("authorization").(string)
	header := map[string]string{
		"Authorization": authorization,
	}
	verifyHeader := map[string]string{
		"Accept":        MediaType,
		"Authorization": authorization,
	}
	var resp BatchRespVO
	repoPath := getRepo(c).RelativePath
	resp.Objects, _ = listutil.Map(respDTO.ObjectList, func(t lfssrv.ObjectDTO) (ObjectRespVO, error) {
		if t.Err == nil {
			var actions map[string]LinkVO
			if isUpload {
				actions = map[string]LinkVO{
					"upload": {
						Href:   fmt.Sprintf("%s/%s/info/lfs/objects/%s/%d", setting.AppUrl(), repoPath, t.Oid, t.Size),
						Header: header,
					},
					"verify": {
						Href:   fmt.Sprintf("%s/%s/info/lfs/verify", setting.AppUrl(), repoPath),
						Header: verifyHeader,
					},
				}
			} else {
				actions = map[string]LinkVO{
					"download": {
						Href:   fmt.Sprintf("%s/%s/info/lfs/objects/%s", setting.AppUrl(), repoPath, t.Oid),
						Header: header,
					},
				}
			}
			return ObjectRespVO{
				PointerVO: PointerVO{
					Oid:  t.Oid,
					Size: t.Size,
				},
				Actions: actions,
			}, nil
		} else {
			return ObjectRespVO{
				Error: &ObjectErrVO{
					Code:    http.StatusUnprocessableEntity,
					Message: err.Error(),
				},
			}, nil
		}
	})
	c.JSON(http.StatusOK, resp)
}

func getOperator(c *gin.Context) usermd.UserInfo {
	return c.MustGet("operator").(usermd.UserInfo)
}

func getRepo(c *gin.Context) repomd.RepoInfo {
	return c.MustGet("repo").(repomd.RepoInfo)
}

func lock(c *gin.Context) {
	var req PostLockReqVO
	err := c.ShouldBindJSON(&req)
	if err != nil {
		writeRespMessage(c, http.StatusBadRequest, "bad request")
		return
	}
	ctx := c.Request.Context()
	operator := getOperator(c)
	singleLock, err := lfssrv.Lock(ctx, lfssrv.LockReqDTO{
		Repo:     getRepo(c),
		Operator: operator,
	})
	if err != nil {
		c.JSON(http.StatusOK, ErrVO{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, PostLockRespVO{
		Lock: model2LockVO(singleLock, operator),
	})
}

func listLock(c *gin.Context) {
	var req ListLockReqVO
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrVO{
			Message: err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	operator := getOperator(c)
	listResp, err := lfssrv.ListLock(ctx, lfssrv.ListLockReqDTO{
		Repo:     getRepo(c),
		Operator: operator,
		Cursor:   req.Cursor,
		Limit:    req.Limit,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrVO{
			Message: err.Error(),
		})
		return
	}
	listVO, _ := listutil.Map(listResp.LockList, func(lock lfsmd.LfsLock) (LockVO, error) {
		return model2LockVO(lock, operator), nil
	})
	c.JSON(http.StatusOK, ListLockRespVO{
		Locks: listVO,
		Next:  listResp.Next,
	})
}

func unlock(c *gin.Context) {
	var req UnlockReqVO
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrVO{
			Message: err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	operator := getOperator(c)
	singleLock, err := lfssrv.Unlock(ctx, lfssrv.UnlockReqDTO{
		Repo:     getRepo(c),
		LockId:   id,
		Force:    req.Force,
		Operator: operator,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrVO{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, UnlockRespVO{
		Lock: model2LockVO(singleLock, operator),
	})
}

func listLockVerify(c *gin.Context) {
	var req ListLockVerifyReqVO
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrVO{
			Message: err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	operator := getOperator(c)
	listResp, err := lfssrv.ListLock(ctx, lfssrv.ListLockReqDTO{
		Repo:     getRepo(c),
		Operator: operator,
		Cursor:   req.Cursor,
		Limit:    req.Limit,
		RefName:  req.Ref.Name,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrVO{
			Message: err.Error(),
		})
		return
	}
	voList := listResp.LockList
	ours, _ := listutil.Filter(voList, func(lock lfsmd.LfsLock) (bool, error) {
		return lock.OwnerId == operator.Id, nil
	})
	oursRet, _ := listutil.Map(ours, func(lock lfsmd.LfsLock) (LockVO, error) {
		return model2LockVO(lock, operator), nil
	})
	theirs, _ := listutil.Filter(voList, func(lock lfsmd.LfsLock) (bool, error) {
		return lock.OwnerId != operator.Id, nil
	})
	theirsRet, _ := listutil.Map(theirs, func(lock lfsmd.LfsLock) (LockVO, error) {
		return model2LockVO(lock, operator), nil
	})
	respVO := ListLockVerifyRespVO{
		Ours:   oursRet,
		Theirs: theirsRet,
		Next:   listResp.Next,
	}
	c.JSON(http.StatusOK, respVO)
}

func verify(c *gin.Context) {
	var req PointerVO
	err := c.ShouldBindJSON(&req)
	if err != nil {
		writeRespMessage(c, http.StatusBadRequest, "bad request")
		return
	}
	ctx := c.Request.Context()
	err = lfssrv.Verify(ctx, lfssrv.VerifyReqDTO{
		PointerDTO: lfssrv.PointerDTO{
			Oid:  req.Oid,
			Size: req.Size,
		},
		Repo:     getRepo(c),
		Operator: getOperator(c),
	})
	if err != nil {
		writeRespMessage(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeRespMessage(c, http.StatusOK, "")
}

func download(c *gin.Context) {
	oid := c.Param("oid")
	var (
		fromByte int64 = -1
		toByte   int64 = -1
	)
	rangeStr := c.GetHeader("Range")
	if rangeStr != "" {
		match := rangeHeaderRegexp.FindStringSubmatch(rangeStr)
		if len(match) > 1 {
			fromByte, _ = strconv.ParseInt(match[1], 10, 32)
			if match[2] != "" {
				toByte, _ = strconv.ParseInt(match[2], 10, 32)
			}
		}
	}
	ctx := c.Request.Context()
	downloadResp, err := lfssrv.Download(ctx, lfssrv.DownloadReqDTO{
		Oid:      oid,
		Repo:     getRepo(c),
		Operator: getOperator(c),
		FromByte: fromByte,
		ToByte:   toByte,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrVO{
			Message: err.Error(),
		})
		return
	}
	defer downloadResp.Close()
	extraHeader := make(map[string]string)
	if downloadResp.FromByte > 0 {
		extraHeader["Content-Range"] = fmt.Sprintf("bytes %d-%d/%d", downloadResp.FromByte, downloadResp.ToByte, downloadResp.Length)
		extraHeader["Access-Control-Expose-Headers"] = "Content-Range"
	}
	filename := c.Param("filename")
	if filename != "" {
		decodedFilename, err := base64.RawURLEncoding.DecodeString(filename)
		if err == nil {
			extraHeader["Content-Disposition"] = fmt.Sprintf("attachment; filename=\"%s\"", string(decodedFilename))
			extraHeader["Access-Control-Expose-Headers"] = "Content-Disposition"
		}
	}
	c.DataFromReader(http.StatusOK, downloadResp.Length, "application/octet-stream", downloadResp, extraHeader)
}

func upload(c *gin.Context) {
	size, err := strconv.ParseInt(c.Param("size"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrVO{
			Message: "wrong size",
		})
		return
	}
	oid := c.Param("oid")
	body := c.Request.Body
	ctx := c.Request.Context()
	defer body.Close()
	err = lfssrv.Upload(ctx, lfssrv.UploadReqDTO{
		Oid:      oid,
		Size:     size,
		Repo:     getRepo(c),
		Operator: getOperator(c),
		Body:     body,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrVO{
			Message: err.Error(),
		})
		return
	}
	writeRespMessage(c, http.StatusOK, "")
}

func model2LockVO(lock lfsmd.LfsLock, locker usermd.UserInfo) LockVO {
	return LockVO{
		Id:       strconv.FormatInt(lock.Id, 10),
		Path:     lock.Path,
		LockedAt: lock.Created.Round(time.Second),
		Owner: &LockOwnerVO{
			Name: locker.Name,
		},
	}
}

func writeRespMessage(c *gin.Context, code int, message string) {
	c.Data(code, MediaType, []byte(message))
}

package e

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",

	INVALID_PARSE_FORM:              "解析绑定表单错误",
	INVALID_PARAMS_VERIFY:           "参数校验错误",
	INVALID_MAIL_VERIFY:             "邮箱校验错误",
	ERROR_USERNAME_PASSWORD:         "用户名密码不正确",
	ERROR_AUTH_CHECK_TOKEN_FAIL:     "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT:  "Token已超时,请重新登录",
	ERROR_AUTH_TOKEN:                "Token生成失败",
	ERROR_AUTH:                      "Token错误",
	ERROR_UPLOAD_SAVE_IMAGE_FAIL:    "保存图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:   "检查图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT: "校验图片错误，图片格式或大小有问题",

	ERROR_UPLOAD_SAVE_FILE_FAIL:    "保存文件失败",
	ERROR_UPLOAD_CHECK_FILE_FAIL:   "检查文件失败",
	ERROR_UPLOAD_CHECK_FILE_FORMAT: "校验文件错误，文件格式或大小有问题",

	ERROR_USERNAME_EXIST:  "用户名已存在",
	ERROR_ADDUSER_FAIL:    "新增用户失败",
	ERROR_LOGIN_MAIL_FAIL: "登录失败, 用户邮箱不存在",
	ERROR_GET_USER_FAIL:   "获取用户信息错误",

	ERROR_LOGIN_PWD_FAIL: "登录失败, 用户密码错误",
	ERROR_USERMAIL_EXIST: "用户邮箱已被注册",

	ERROR_UPLOAD_NDFILE_FAIL:   "上传文件失败",
	ERROR_GET_NDFILELIST_FAIL:  "获取文件列表失败",
	ERROR_MOVE_TO_TRASH_FAIL:   "移动到回收站失败",
	ERROR_DELETE_NDFILE_FAIL:   "删除文件失败",
	ERROR_GET_DIR_LIST_FAIL:    "获取文件夹列表失败",
	ERROR_ADD_DIRL_FAIL:        "新建文件夹失败",
	ERROR_UPDATE_DIR_FAIL:      "修改文件夹失败",
	ERROR_DELETE_DIR_FAIL:      "删除文件夹失败",
	ERROR_DELETE_DIR_IS_PARENT: "删除失败,文件夹含子文件夹",
	ERROR_DELETE_DIR_HAS_FILE:  "删除失败,文件夹非空",

	ERROR_GET_USERBYMOBILE_FAIL: "用户手机号不正确",
	ERROR_GET_USERMAIL_FAIL:     "用户邮箱不正确",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}

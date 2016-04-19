package v1

import "git.zxq.co/ripple/rippleapi/common"

type setAllowedData struct {
	UserID  int `json:"user_id"`
	Allowed int `json:"allowed"`
}

// UserManageSetAllowedPOST allows to set the allowed status of an user.
func UserManageSetAllowedPOST(md common.MethodData) common.CodeMessager {
	data := setAllowedData{}
	if err := md.RequestData.Unmarshal(&data); err != nil {
		return ErrBadJSON
	}
	if data.Allowed < 0 || data.Allowed > 2 {
		return common.SimpleResponse(400, "Allowed status must be between 0 and 2")
	}
	_, err := md.DB.Exec("UPDATE users SET allowed = ? WHERE id = ?", data.Allowed, data.UserID)
	if err != nil {
		md.Err(err)
		return Err500
	}
	query := `
SELECT users.id, users.username, register_datetime, rank,
	latest_activity, users_stats.username_aka,
	users_stats.country, users_stats.show_country
FROM users
LEFT JOIN users_stats
ON users.id=users_stats.id
WHERE users.id=?
LIMIT 1`
	return userPuts(md, md.DB.QueryRow(query, data.UserID))
}

package utils

import (
	"github.com/qiniu/api/rs"
)

func GetQiniuUptoken(bucketName string) string {
	putPolicy := rs.PutPolicy{
		Scope: bucketName,
	}
	return putPolicy.Token(nil)
}
func GetQiniuPrivateDownloadUrl(domain, key string) string {
	baseUrl := rs.MakeBaseUrl(domain, key)
	policy := rs.GetPolicy{}
	return policy.MakeRequest(baseUrl, nil)
}

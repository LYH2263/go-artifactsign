// Package artifactsign 提供制品摘要、HMAC 签名/验签、信任根与吊销。
//
// 能力概览：
//   - Digest / Sign / Verify：对制品字节做 SHA256 摘要与 HMAC 签名
//   - 信任根注册与 ExportTrust；密钥轮换
//   - 吊销指纹；Verify 时检查
//   - Manifest 多文件摘要聚合与临时清单落盘
//   - Close：先刷持久化再释放密钥；之后 Sign 返回 ErrClosed
//
// 演示用 HMAC/SHA256，不引入第三方密码学库。
package artifactsign

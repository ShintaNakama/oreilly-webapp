# My Example Go
## 1.chat-app
- アバター(アバター画像)がインターフェース化されている
```go
type Avatar interface {
  GetAvatarURL(ChatUser) (string, error)
}
```
- 下記のアバター画像を取得する方法が抽象化されている
  - 複数のアバター画像の取得
  - localの画像ファイルからの取得
  - 認証サービスからの取得
  - Gravatarからの取得

- チャットユーザーもインターフェース化されている
```go
import	gomniauthcommon "github.com/stretchr/gomniauth/common"
type ChatUser interface {
	UniqueID() string
	AvatarURL() string
}

type chatUser struct {
	gomniauthcommon.User
	uniqueID string
}

func (u chatUser) UniqueID() string {
	return u.uniqueID
}
```
- チャットユーザーの構造体自体は、privateにし、インターフェースとして必要な機能を提供している
- GomniauthのUserインターフェースでも、AvatarURLメソッドが定義されているため、chatuserは、UniqueIDメソッドの実装だけで、インターフェースを満たしているということになる
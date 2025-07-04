# OGP 検証サービス API ドキュメント

このディレクトリには、OGP 検証サービス API の OpenAPI 仕様が含まれています。

## API 仕様

- **OpenAPI バージョン**: 3.0.3
- **仕様ファイル**: `openapi.yaml`

## ドキュメントの表示

### 方法 1: Swagger UI (推奨)

Swagger UI を使用してインタラクティブな API ドキュメントを表示するには：

```bash
# backend ディレクトリから
go run cmd/swagger/main.go
```

その後、ブラウザで http://localhost:8081 を開きます。

### 方法 2: オンライン Swagger エディター

1. https://editor.swagger.io/ にアクセス
2. `openapi.yaml` の内容をコピー
3. エディターに貼り付け

### 方法 3: VS Code 拡張機能

VS Code で "OpenAPI (Swagger) Editor" 拡張機能をインストールして、シンタックスハイライトと検証機能付きで仕様を表示・編集できます。

## API エンドポイント

### コアエンドポイント

- **POST /api/v1/ogp/verify** - 指定された URL の OGP メタデータを検証
- **GET /health** - ヘルスチェックエンドポイント

### レート制限

API はレート制限を実装しています：
- **制限**: IP アドレスあたり 10 リクエスト/分
- **レスポンス**: 制限を超えた場合は HTTP 429

### CORS サポート

API は以下のヘッダーで CORS をサポートしています：
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: POST, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type`

## プラットフォーム固有の制限

API はプラットフォーム固有の制限に対してコンテンツを検証します：

### Twitter/X
- タイトル: 70 文字
- 説明: 200 文字
- 推奨画像: 1200x630px

### Facebook
- タイトル: 100 文字
- 説明: 300 文字
- 推奨画像: 1200x630px

### Discord
- タイトル: 256 文字
- 説明: 2048 文字
- 画像サイズ: 柔軟

## リクエスト例

```bash
curl -X POST http://localhost:8080/api/v1/ogp/verify \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com"}'
```

## レスポンス例

```json
{
  "url": "https://github.com",
  "ogp_data": {
    "title": "GitHub · Build and ship software on a single, collaborative platform",
    "description": "Join the world's most widely adopted...",
    "image": "https://github.githubassets.com/assets/home24-5939032587c9.jpg",
    "url": "https://github.com/",
    "type": "object",
    "site_name": "GitHub"
  },
  "validation": {
    "is_valid": true,
    "warnings": [],
    "errors": [],
    "checks": {
      "has_title": true,
      "has_description": true,
      "has_image": true,
      "image_valid": true,
      "url_valid": true
    }
  },
  "previews": {
    "twitter": {
      "platform": "twitter",
      "title": "GitHub · Build and ship software...",
      "description": "Join the world's most widely adopted...",
      "image": "https://github.githubassets.com/assets/home24-5939032587c9.jpg",
      "is_valid": true,
      "warnings": [],
      "title_length": 69,
      "desc_length": 186,
      "max_title_len": 70,
      "max_desc_len": 200
    },
    "facebook": { ... },
    "discord": { ... }
  },
  "timestamp": "2025-07-03T18:00:00Z"
}
```

## セキュリティに関する考慮事項

1. **プライベート IP ブロック**: サービスはプライベート IP アドレスへのリクエストをブロックします
2. **レート制限**: 悪用を防ぐため IP あたり 10 リクエスト/分
3. **リクエストタイムアウト**: 外部 URL の取得は 10 秒でタイムアウト
4. **CORS**: クロスオリジンリクエストを許可するよう設定

## 関連ドキュメント

- [Open Graph Protocol](https://ogp.me/)
- [Twitter Cards](https://developer.twitter.com/en/docs/twitter-for-websites/cards/overview/abouts-cards)
- [Facebook Sharing](https://developers.facebook.com/docs/sharing/webmasters/)
- [Discord Embeds](https://discord.com/developers/docs/resources/channel#embed-object)
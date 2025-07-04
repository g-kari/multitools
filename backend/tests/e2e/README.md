# OGP 検証サービスの E2E テスト

このディレクトリには、OGP 検証サービス API のエンドツーエンドテストが含まれています。

## 概要

E2E テストは、実行中の API インスタンスに対して実際の HTTP リクエストを行い、レスポンスを検証することで、サービスの完全な機能性を確認します。

## テスト対象

### コア機能
- ✅ ヘルスエンドポイントの検証
- ✅ 完全なデータでの OGP 検証（GitHub）
- ✅ 良好なデータでの OGP 検証（Wikipedia）
- ✅ 最小限のデータでの OGP 検証（Example.com）

### エラーハンドリング
- ✅ 空の URL の検証
- ✅ 無効な JSON の処理
- ✅ プライベート IP のブロック

### セキュリティとパフォーマンス
- ✅ レート制限（IP あたり 10 リクエスト/分）
- ✅ CORS ヘッダーの検証
- ✅ 許可されていないメソッドのレスポンス
- ✅ API レイテンシーテスト（< 10 秒）
- ✅ 同時リクエストの処理

## テストの実行

### 前提条件
- Docker と Docker Compose
- テスト用ウェブサイトへのネットワークアクセス（GitHub、Wikipedia、Example.com）

### すべての E2E テストを実行

```bash
cd backend/tests/e2e
./run-e2e.sh
```

### 手動テスト

1. テスト環境を起動：
```bash
docker-compose -f docker-compose.e2e.yml up --build
```

2. 特定のテストを実行：
```bash
docker-compose -f docker-compose.e2e.yml exec e2e-tests go test -v -run TestHealthEndpoint
```

3. クリーンアップ：
```bash
docker-compose -f docker-compose.e2e.yml down --volumes --remove-orphans
```

## テスト構造

### テストファイル
- `main_test.go` - コア E2E テスト実装
- `docker-compose.e2e.yml` - テスト環境設定
- `Dockerfile.e2e` - テストランナーコンテナ
- `run-e2e.sh` - テスト実行スクリプト

### テスト環境
- **バックエンドサービス**: ポート 8080 でヘルスチェック付きで実行
- **テストランナー**: curl と jq ツールを含む Go テスト環境
- **ネットワーク**: テスト用の分離された Docker ネットワーク

## テストシナリオ

### 1. ヘルスチェック（`TestHealthEndpoint`）
`/health` エンドポイントが 200 OK を返すことを確認します。

### 2. GitHub OGP テスト（`TestOGPVerifyEndpoint_GitHub`）
完全な OGP メタデータを持つウェブサイトでテスト：
- 必要なすべての OGP フィールドが存在することを検証
- プラットフォーム固有のプレビューをチェック
- 検証結果を確認

### 3. Wikipedia OGP テスト（`TestOGPVerifyEndpoint_Wikipedia`）
良好な OGP メタデータを持つウェブサイトでテスト：
- タイトルと説明の存在を検証
- 検証フラグをチェック

### 4. Example.com テスト（`TestOGPVerifyEndpoint_ExampleCom`）
最小限の OGP 実装でテスト：
- 警告は出るが有効なレスポンスを期待
- 欠落データのエラーハンドリングを検証

### 5. エラーケース（`TestOGPVerifyEndpoint_ErrorCases`）
様々なエラー条件をテスト：
- 空の URL（400 Bad Request）
- 無効な JSON（400 Bad Request）
- プライベート IP（500 Internal Server Error）

### 6. レート制限（`TestRateLimiting`）
10 リクエスト/分のレート制限を検証：
- 10 回の成功リクエストを実行
- 11 回目のリクエストが 429 を返すことを確認

### 7. CORS ヘッダー（`TestCORSHeaders`）
CORS プリフライトの処理を検証：
- OPTIONS メソッドをテスト
- 必要な CORS ヘッダーをチェック

### 8. メソッド検証（`TestMethodNotAllowed`）
サポートされていない HTTP メソッドが 405 を返すことをテスト。

### 9. パフォーマンス（`TestAPILatency`）
API レスポンス時間が 10 秒未満であることを検証。

### 10. 同時実行（`TestConcurrentRequests`）
複数の同時リクエストの処理をテスト。

## 期待される結果

すべてのテストが以下の条件で成功するはずです：
- ✅ 10/10 テストが成功
- テストの失敗やエラーがない
- 許容範囲内のパフォーマンス
- 適切なエラーハンドリング

## トラブルシューティング

### よくある問題

1. **ネットワーク接続の問題**
   - Docker がインターネットにアクセスできることを確認
   - テスト用ウェブサイトにアクセスできるかチェック

2. **開発中のレート制限**
   - テスト実行間隔を 1 分空ける
   - またはバックエンドサービスを再起動

3. **ヘルスチェックの失敗**
   - docker-compose.e2e.yml でヘルスチェックタイムアウトを増加
   - バックエンドログで起動問題をチェック

### デバッグコマンド

```bash
# バックエンドログを表示
docker-compose -f docker-compose.e2e.yml logs backend-e2e

# テストログを表示
docker-compose -f docker-compose.e2e.yml logs e2e-tests

# 特定のエンドポイントを手動でテスト
docker-compose -f docker-compose.e2e.yml exec e2e-tests curl -X POST http://backend-e2e:8080/api/v1/ogp/verify -H "Content-Type: application/json" -d '{"url":"https://github.com"}'
```

## CI/CD との統合

これらの E2E テストは GitHub Actions やその他の CI/CD パイプラインに統合できます：

```yaml
- name: Run E2E Tests
  run: |
    cd backend/tests/e2e
    ./run-e2e.sh
```
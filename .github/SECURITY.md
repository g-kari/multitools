# セキュリティポリシー

## サポート対象バージョン

以下のバージョンでセキュリティアップデートを積極的にサポートしています：

| バージョン | サポート状況        |
| --------- | ------------------ |
| 1.0.x     | :white_check_mark: |
| < 1.0     | :x:                |

## 脆弱性の報告

セキュリティ脆弱性を深刻に受け止めています。セキュリティ脆弱性を発見した場合は、以下の手順に従ってください：

### 報告方法

1. セキュリティ脆弱性について公開の GitHub issue を作成**しないでください**
2. 脆弱性の詳細についてリポジトリオーナーにメールを送信してください
3. 以下の情報を含めてください：
   - 脆弱性の説明
   - 問題を再現する手順
   - 潜在的な影響
   - 修正案（ある場合）

### 対応について

- **受領確認**: 脆弱性レポートの受領を48時間以内に確認します
- **評価**: 脆弱性を評価し、5営業日以内に深刻度を判定します
- **進捗報告**: 5営業日ごとに進捗状況を定期的に更新します
- **解決**: 重要な脆弱性については30日以内の解決を目指します

### セキュリティベストプラクティス

このプロジェクトは以下のセキュリティベストプラクティスに従っています：

- Dependabot による定期的な依存関係更新
- CI/CD パイプラインでのセキュリティスキャン
- 入力検証とサニタイゼーション
- レート制限と CORS 保護
- プライベート IP アドレスのブロック
- セキュアな Docker 設定

### 責任ある開示

以下についてご協力をお願いします：

- 公開開示前に脆弱性を修正するための合理的な時間を与えてください
- 研究中にデータへのアクセス、変更、削除を避けてください
- ユーザーのプライバシーとシステムの整合性を尊重してください

プロジェクトのセキュリティ確保にご協力いただき、ありがとうございます！
# OGP Verification Service API Documentation

This directory contains the OpenAPI specification for the OGP Verification Service API.

## API Specification

- **OpenAPI Version**: 3.0.3
- **Specification File**: `openapi.yaml`

## Viewing the Documentation

### Option 1: Swagger UI (Recommended)

To view the interactive API documentation using Swagger UI:

```bash
# From the backend directory
go run cmd/swagger/main.go
```

Then open http://localhost:8081 in your browser.

### Option 2: Online Swagger Editor

1. Go to https://editor.swagger.io/
2. Copy the contents of `openapi.yaml`
3. Paste into the editor

### Option 3: VS Code Extension

Install the "OpenAPI (Swagger) Editor" extension in VS Code to view and edit the specification with syntax highlighting and validation.

## API Endpoints

### Core Endpoints

- **POST /api/v1/ogp/verify** - Verify OGP metadata for a given URL
- **GET /health** - Health check endpoint

### Rate Limiting

The API implements rate limiting:
- **Limit**: 10 requests per minute per IP address
- **Response**: HTTP 429 when limit exceeded

### CORS Support

The API supports CORS with the following headers:
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: POST, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type`

## Platform-Specific Limits

The API validates content against platform-specific limits:

### Twitter/X
- Title: 70 characters
- Description: 200 characters
- Recommended image: 1200x630px

### Facebook
- Title: 100 characters
- Description: 300 characters
- Recommended image: 1200x630px

### Discord
- Title: 256 characters
- Description: 2048 characters
- Flexible image dimensions

## Example Request

```bash
curl -X POST http://localhost:8080/api/v1/ogp/verify \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com"}'
```

## Example Response

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

## Security Considerations

1. **Private IP Blocking**: The service blocks requests to private IP addresses
2. **Rate Limiting**: 10 requests per minute per IP to prevent abuse
3. **Request Timeout**: 10-second timeout for fetching external URLs
4. **CORS**: Configured to allow cross-origin requests

## Further Documentation

- [Open Graph Protocol](https://ogp.me/)
- [Twitter Cards](https://developer.twitter.com/en/docs/twitter-for-websites/cards/overview/abouts-cards)
- [Facebook Sharing](https://developers.facebook.com/docs/sharing/webmasters/)
- [Discord Embeds](https://discord.com/developers/docs/resources/channel#embed-object)
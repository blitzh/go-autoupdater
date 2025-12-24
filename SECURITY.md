# Security Policy

## Supported Versions
This project is currently in early stage. Only the latest release is supported.

## Reporting a Vulnerability
If you believe you have found a security vulnerability, please do **not** open a public issue.

Instead, report it privately:
- Email: (add your email here)  
- Or: open a GitHub Security Advisory (recommended)

## Important Notes (SHA256-only)
Current releases verify downloads using **SHA256** integrity checks only. This protects against corrupted/partial downloads, but **does not** guarantee authenticity if an update server (or DNS) is compromised.

For production/higher-risk environments, use:
- HTTPS-only update endpoints
- Controlled hosting (least privilege)
- (Planned) signature verification (e.g., Ed25519) for manifest/artifacts
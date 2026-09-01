# Security Policy

## Reporting Vulnerabilities

If you discover a security vulnerability in this project, please **do not** open a public issue. Instead, please report it responsibly by sending an email to the project maintainer.

When reporting a vulnerability, please provide:
- Description of the vulnerability
- Steps to reproduce (if applicable)
- Potential impact
- Suggested fix (if you have one)

We take security seriously and will investigate and respond to vulnerability reports promptly.

## Security Features

This URL Shortener implements multiple security hardening measures:

### Input Validation
- **URL Format Validation**: All URLs are validated for proper HTTP/HTTPS scheme
- **Length Limits**: URLs have maximum length constraints (2048 characters)
- **SSRF Prevention**: Blocks requests to localhost, private IPs (10.x, 172.16-31.x, 192.168.x), and loopback addresses
- **Character Encoding**: UTF-8 validation for all input

### Network Security
- **Security Headers**: Implements CSP, X-Frame-Options, X-Content-Type-Options
- **CORS Policy**: Strict CORS configuration with whitelisted origins
- **HSTS**: HTTP Strict Transport Security enabled for HTTPS
- **Rate Limiting**: Ready for rate limiting implementation

### Configuration Management
- **Environment Variables**: All sensitive configuration via env vars (never hardcoded)
- **Configuration Validation**: All config values validated on startup
- **No Secrets in Code**: Zero credentials in source code

### Container Security
- **Non-root User**: Container runs as unprivileged user (UID 1000)
- **Minimal Base Image**: Alpine Linux for smallest attack surface
- **Security Scanning**: Dockerfile and dependencies scanned in CI/CD
- **Health Checks**: Built-in health check endpoint
- **Version Pinning**: Specific version pins (not :latest)

### Application Security
- **Error Handling**: Structured error handling without exposing internals
- **Panic Recovery**: Graceful panic recovery middleware
- **Concurrency Safety**: Thread-safe storage with mutex locks
- **Race Condition Testing**: All code tested with `-race` flag

## Dependency Management

All dependencies are tracked and regularly audited for known vulnerabilities:

```bash
# Check for vulnerable dependencies
go list -json -m all | nancy sleuth

# Security scanning
gosec ./...
```

## Best Practices for Users

When using this URL shortener:

1. **Only shorten public URLs**: Never shorten URLs containing sensitive data
2. **Validate Origin**: Always verify the short URL source before accessing
3. **Use HTTPS**: Always access the service over HTTPS in production
4. **Environment Variables**: Use secure secret management for configuration
5. **Access Control**: Implement authentication if exposing to untrusted networks

## Known Security Considerations

### In-Memory Storage
- URLs are stored in application memory only
- Data is not persisted to disk
- All data is lost on application restart
- **Recommendation**: For production use, implement persistent storage with database encryption

### No Authentication
- The API does not include authentication by default
- **Recommendation**: Implement API keys, OAuth, or other authentication mechanisms in production

### Single Instance
- Designed for single-instance deployment
- **Recommendation**: For high-availability, implement distributed storage (Redis, etc.)

## Security Hardening Checklist

Before deploying to production:

- [ ] Use HTTPS/TLS only (redirect HTTP to HTTPS)
- [ ] Implement rate limiting or DDoS protection
- [ ] Set up Web Application Firewall (WAF) rules
- [ ] Enable security monitoring and alerting
- [ ] Implement authentication for API access
- [ ] Use persistent encrypted storage instead of in-memory
- [ ] Regular security audits and penetration testing
- [ ] Keep Go and dependencies up to date
- [ ] Monitor container image for vulnerabilities
- [ ] Set up logging and monitoring
- [ ] Implement backup and disaster recovery

## Security Updates

- Critical security updates will be released as soon as possible
- Follow the project repository for security announcements
- Subscribe to Go security mailing list for framework updates

## Contact

For security questions or concerns, please contact the project maintainers.

---

**Last Updated**: September 2026  
**Version**: 1.0.0  
**Review Schedule**: Quarterly

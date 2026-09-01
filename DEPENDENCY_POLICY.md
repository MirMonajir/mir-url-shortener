# Dependency Management Policy

## Overview

This document outlines the dependency management strategy for the URL Shortener project, including:
- Dependency selection criteria
- Vulnerability management
- License compliance
- Update procedures

## Approved Dependencies

### Core Framework
- **github.com/gin-gonic/gin**: v1.11.0+ - Web framework (MIT License)
  - Mature, widely-used HTTP framework for Go
  - Active community support
  - Regular security updates

### Testing
- **github.com/stretchr/testify**: v1.11.1+ - Testing utilities (MIT License)
  - De facto standard for Go testing
  - Assertion library for clear test cases

### Utilities
- **golang.org/x/net**: v0.44.0+ - Network utilities (BSD License)
  - Official Go project (trusted)
  - Used for public suffix handling in domain extraction
  - Regular updates from Go team

- **github.com/joho/godotenv**: v1.5.1+ - .env file loading (MIT License)
  - Standard library for environment variable management
  - Minimal, well-maintained

### Domain/Infrastructure
- **golang.org/x/crypto**: - Cryptographic functions (BSD License)
  - Official Go project
  - Used for secure password handling (if implemented)
  - Included via transitive dependencies

## Vulnerability Management

### Scanning
All dependencies are scanned for known vulnerabilities:

1. **Go Mod Audit** (Built-in)
   ```bash
   go list -json -m all | nancy sleuth
   ```

2. **Gosec** (Go Security Scanner)
   ```bash
   gosec ./...
   ```

3. **Trivy** (Container Image Scanning)
   ```bash
   trivy image mir-url-shortener:latest
   ```

4. **GitHub Dependabot** (Automated alerts)
   - Configured in `.github/dependabot.yml`
   - Alerts on known vulnerabilities

### Vulnerability Response
- **Critical**: Fixed immediately, security patch released
- **High**: Fixed in next release cycle
- **Medium**: Planned for next update
- **Low**: Monitored, updated in regular maintenance

## Update Procedures

### Regular Updates
1. Run `go get -u ./...` monthly
2. Review changelog for breaking changes
3. Run full test suite
4. Run security scans
5. Commit and test in CI/CD

### Security Updates
1. Immediate action on critical vulnerabilities
2. Test thoroughly before deployment
3. Release security patch immediately
4. Notify users of security fixes

### Breaking Changes
1. Evaluate impact on application
2. Plan migration if needed
3. Update code before upgrading dependency
4. Run full test suite
5. Document breaking changes in changelog

## License Compliance

### License Allowlist
The project uses only dependencies with permissive open-source licenses:
- MIT License ✅
- BSD License ✅
- Apache 2.0 License ✅
- ISC License ✅
- Unlicense ✅

### Prohibited Licenses
- GPL v3 or later ❌
- AGPL ❌
- Any license requiring source disclosure for binary distribution ❌

### License Verification
Before adding a dependency:
1. Check license file
2. Verify it's in allowlist
3. Document license in this file
4. Update CI/CD scanning if needed

## Dependency Audit Process

### Weekly
- Review GitHub Dependabot alerts
- Assess vulnerability severity
- Plan fixes for critical issues

### Monthly
- Run full dependency audit
- Check for outdated packages
- Review new versions
- Update non-breaking dependencies

### Quarterly
- Full security review
- License audit
- Evaluate dependency removal opportunities
- Plan major version upgrades

## Adding New Dependencies

Before adding a new dependency:

1. **Necessity Check**
   - Is this functionality not available in stdlib?
   - Can we implement it ourselves?
   - Is the dependency actively maintained?

2. **Security Check**
   - Review security history
   - Check for known vulnerabilities
   - Verify integrity (checksums)

3. **License Check**
   - Verify license is in allowlist
   - No GPL/AGPL licenses
   - Document license

4. **Size Check**
   - Check dependency size
   - Review transitive dependencies
   - Ensure no unnecessary bloat

5. **Community Check**
   - Verify active maintenance
   - Check GitHub stars/activity
   - Look for active issue resolution

## Current Dependency Tree

```
github.com/MirMonajir/mir-url-shortener
├── github.com/gin-gonic/gin v1.11.0
│   └── [other dependencies...]
├── github.com/joho/godotenv v1.5.1
├── github.com/stretchr/testify v1.11.1
│   ├── github.com/davecgh/go-spew v1.1.1
│   ├── github.com/pmezard/go-difflib v1.0.0
│   └── gopkg.in/yaml.v3 v3.0.1
└── golang.org/x/net v0.44.0
    └── golang.org/x/text v0.29.0
```

## Monitoring

### Automated Checks
- ✅ GitHub Dependabot enabled
- ✅ Go module vulnerability scanning in CI/CD
- ✅ Container image scanning in CI/CD
- ✅ gosec security scanning in CI/CD

### Manual Review
- Monthly dependency audit
- Quarterly full security review
- Annual license compliance audit

## References

- [Go Module Best Practices](https://golang.org/doc/modules/managing-dependencies)
- [OWASP Dependency-Check](https://owasp.org/www-project-dependency-check/)
- [Open Source License Guide](https://choosealicense.com/)

---

**Last Updated**: September 2026  
**Maintained By**: Project Team  
**Review Schedule**: Quarterly

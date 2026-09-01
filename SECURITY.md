# Security Policy

## Supported Versions

Only the latest released version is supported. Security fixes are not backported to older tags.

| Version  | Supported          |
| -------- | ------------------ |
| latest   | :white_check_mark: |
| < latest | :x:                |

## Reporting a Vulnerability

Please follow these steps if you discover a security vulnerability in this project:

### Do Not

- **Do not** open a public GitHub issue for security vulnerabilities
- **Do not** disclose the vulnerability publicly until it has been addressed

### Do

1. **Report privately** via [GitHub Security Advisories](https://github.com/miikkak/rcon-cli/security/advisories/new) <!-- markdownlint-disable-line MD013 -->
2. **Include in your report:**
   - Description of the vulnerability
   - Steps to reproduce the issue
   - Potential impact
   - Suggested fix (if you have one)

3. **Response timeline:**
   - You should receive an acknowledgment within 48 hours
   - We'll provide a detailed response within 7 days
   - We'll work with you to understand and fix the issue
   - We'll release a fix as soon as possible

## Scope

This is an RCON client - it connects outbound to an RCON server using a host, port, and password
supplied via flags, environment variables, or a config file, and sends whatever command it's
given. It doesn't listen on any network port itself, doesn't execute anything derived from the
RCON server's response beyond parsing and displaying it, and doesn't persist the password
anywhere beyond the process's own memory and whatever config file you chose to put it in. If you
find a way this tool could be used to leak the RCON password, execute unintended commands, or
otherwise misbehave beyond its documented behavior, that's exactly the kind of thing to report.

## Security Best Practices

When using this tool:

- Always use a specific released version, not a locally built development binary, in production
- Treat the RCON password like any other credential - prefer the environment variable or config
  file over passing `--password` on the command line, where it can leak into shell history or
  process listings
- Restrict the config file's permissions if it contains a password (the file itself is your
  responsibility - this tool doesn't create or manage it)
- Keep the tool updated - check releases periodically or watch the repository

## Security Scanning

This project uses automated security scanning:

- **Trivy** (filesystem scan against `go.sum`) for dependency vulnerability scanning, on a
  weekly schedule and on demand
- **golangci-lint** (including security-relevant linters) on every PR
- **Renovate** for automated dependency updates

## Other Automated Review

Every pull request also gets an AI code review. This is a general correctness/quality review,
not a vulnerability scanner - don't rely on it as a substitute for the security scanning above.

## Disclosure Policy

- Security issues are fixed in private before public disclosure
- After a fix is released, we publish a security advisory
- We credit reporters in the advisory (unless they prefer anonymity)

## Past Security Advisories

No security advisories have been published yet.

## Contact

For security-related questions or concerns, please use the reporting method above rather than
public channels.

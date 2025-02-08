# go-regal

**RegAL - Registry Access Lock: Your keys to the kingdom.**

[![GoDoc](https://godoc.org/github.com/yourusername/go-regal?status.svg)](https://godoc.org/github.com/yourusername/go-regal) [![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE) **A simple, straightforward, and flexible authentication and authorization library for Go, designed to secure your domain.**

---

## What is RegAL?

(Project still in development, and unreleased)

RegAL is a Go library that provides an easy-to-use yet robust foundation for authentication and authorization in your Go applications.  Built with simplicity, flexibility, and security in mind, RegAL is designed to:

*   **Simplify Access Control:**  Provides clear and intuitive APIs for defining and enforcing access policies.
*   **Offer Domain-Oriented Authentication:**  Supports diverse authentication methods, from traditional username/password to modern approaches like email login, SSO, and future credential types (with extensibility for physical access control).
*   **Enhance Security:**  Prioritizes security best practices, including hashing secrets (even short-lived ones) whenever possible.
*   **Promote Clean and Extensible Code:**  Designed for easy integration and extension to fit your specific needs.

RegAL leverages a core architecture based on:

*   **Registry:**  Manages user identities and authentication methods.
*   **Access:**  Handles authentication logic and credential verification.
*   **Door Lock:**  Enforces authorization policies and access decisions.

**In essence, RegAL gives you the "keys to the kingdom" – providing you with the tools to confidently control and secure access to your Go applications and services.**

---

## Features

*   **Modular Architecture:**  Built with independent components (Registry, Access, Door Lock) for flexibility and maintainability.
*   **Extensible Authentication:**  Supports multiple authentication methods out-of-the-box and is designed for easy addition of new methods.
*   **Focus on Ease of Use:**  Provides a clean and straightforward API for developers.
*   **Security First:**  Implements security best practices to protect sensitive data.
*   **Go Native:**  Built in Go for performance and seamless integration with Go ecosystems.
*   **Domain-Oriented Design:**  Adaptable to various application domains, from web applications to potentially physical access control systems.

---

## Getting Started
(To be completed)

# Next steps
## For the MVP
### MUST
#### User registration
- [ ] add user registry persistence implementation
- [ ] email authentication method persistence implementation
- [ ] add user registry persistence abstraction
- [x] login Access using default authentication method
- [ ] login DoorLock 
- [x] add default authentication method

### SHOULD
- [ ] admin access to list users

### COULD
- [ ] add password authentication method (needs password reset functionality)
- [ ] add sms/whatsapp authentication method
- [ ] add SSO authentication method
- [ ] associate user with authentication method on the access component
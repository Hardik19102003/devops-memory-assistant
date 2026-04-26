# DevOps Memory Assistant 🚀

A tool to store and recall DevOps issues, so you never debug the same problem twice.

---

## 💡 Why this project?

As a DevOps engineer, I often encountered recurring issues like:
- CrashLoopBackOff
- Misconfigurations
- Pod failures

Debugging the same problem again and again is frustrating.

This tool helps you:
- Save issues
- Recall solutions instantly

## 🔍 Features

- Save DevOps issues (error, cause, fix)
- Search past issues instantly
- PostgreSQL-backed storage

---

## 🔌 API

### Save Issue

POST /issue

### Search Issues

GET /search?error=CrashLoopBackOff
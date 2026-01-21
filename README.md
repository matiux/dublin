Dublin
======

Dublin is a maintained and modernized fork of the original **Broadway** project.

It provides infrastructure and testing helpers for building **CQRS** and
**Event Sourced** applications in PHP, with a strong focus on explicit
models, clear boundaries and minimal framework intrusion.

Dublin continues the original vision of Broadway while evolving it to support
modern PHP versions and current development practices.

![build status](https://github.com/matiux/dublin/actions/workflows/ci.yml/badge.svg)

---

## Fork origin and intent

Dublin is a fork of the **Broadway** project  
(https://github.com/broadway/broadway), originally developed by the Broadway
contributors and the team at Qandidate.

The original project is no longer actively maintained.  
The goal of Dublin is to:

- keep the project **alive and maintained**
- support **modern PHP versions (8.4+)**
- modernize tooling, CI and testing infrastructure
- evolve the codebase while preserving the original design philosophy

Where possible, Dublin aims to remain **API-compatible** with Broadway.

---

## About

Like Broadway, Dublin provides a set of loosely coupled components that can be
used independently or together to build CQRS and Event Sourced systems.

Dublin intentionally stays out of your way: no mandatory frameworks, no hidden
magic, no forced architectural choices.

---

## Architecture documentation (Structurizr)

This project includes Structurizr documentation. After building the project
and starting the containers, the docs are available at:

http://localhost:8080/

Build and run:

```bash
make build-php ARG="--no-cache"
make upd
```

---

## Installation

```bash
composer require matiux/dublin
```

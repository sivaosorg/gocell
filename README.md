# Gocell

![GitHub contributors](https://img.shields.io/github/contributors/sivaosorg/gocell)
![GitHub followers](https://img.shields.io/github/followers/sivaosorg)
![GitHub User's stars](https://img.shields.io/github/stars/pnguyen215)
![GitHub language count](https://img.shields.io/github/languages/count/sivaosorg/gocell)
![GitHub top language](https://img.shields.io/github/languages/top/sivaosorg/gocell)
![GitHub issues](https://img.shields.io/github/issues/sivaosorg/gocell)
![GitHub closed issues](https://img.shields.io/github/issues-closed/sivaosorg/gocell)
![GitHub repo size](https://img.shields.io/github/repo-size/sivaosorg/gocell)

GoCell is a robust and scalable base platform for building web applications using the Go programming language.

## Table of Contents

- [Gocell](#gocell)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Prerequisites](#prerequisites)
  - [Getting Started](#getting-started)
    - [Running the REST API Server](#running-the-rest-api-server)
    - [Building the REST API Server](#building-the-rest-api-server)
    - [Running the Job Server](#running-the-job-server)
    - [Building the Job Server](#building-the-job-server)
  - [Modules Support](#modules-support)
    - [Running Tests](#running-tests)
    - [Tidying up Modules](#tidying-up-modules)
    - [Upgrading Dependencies](#upgrading-dependencies)
    - [Cleaning Dependency Cache](#cleaning-dependency-cache)
  - [Tools Commands](#tools-commands)
    - [Swagger Documentation](#swagger-documentation)
  - [Component](#component)

## Introduction

GoCell is a robust and scalable base platform for building web applications using the Go programming language. It provides a solid foundation with essential features and a modular architecture that allows developers to quickly kickstart their projects and focus on building core functionality.

## Prerequisites

Golang version v1.20

## Getting Started

Explain how users can get started with Gocell project.

### Running the REST API Server

To run the REST API server, use the following command:

```bash
make run
```

This will start the server using the provided configuration files.

### Building the REST API Server

To build the REST API server, use the following command:

```bash
make build
```

This will compile the server executable.

### Running the Job Server

To run the Job server, use the following command:

```bash
make run-job
```

This will start the Job server using the specified configuration file.

### Building the Job Server

To build the Job server, use the following command:

```bash
make build-job
```

This will compile the Job server executable.

## Modules Support

Explain how users can interact with the various modules in Gocell project.

### Running Tests

To run tests for all modules, use the following command:

```bash
make test
```

### Tidying up Modules

To tidy up the project's Go modules, use the following command:

```bash
make tidy
```

### Upgrading Dependencies

To upgrade project dependencies, use the following command:

```bash
make deps-upgrade
```

### Cleaning Dependency Cache

To clean the Go module cache, use the following command:

```bash
make deps-clean-cache
```

## Tools Commands

### Swagger Documentation

To generate Swagger documentation, use the following command:

```bash
make swaggo
```

This will initialize Swagger and generate the API documentation.

## Component

- [x] Integrated Postgres
- [x] Integrated MySQL
- [x] Integrated RabbitMQ
- [x] Integrated Redis
- [x] Integrated Swagger
- [x] Integrated Websocket
- [x] Integrated base middlewares (CORS, ERROR, ...)
- [x] Integrated Telegram / Slack bot
- [x] Integrated Logger
- [ ] Integrating Kafka
- [ ] Add Dockerfile
- [ ] Add cronjob

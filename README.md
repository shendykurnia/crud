# Senior Software Engineer Test

The goal of this test is to assert (to some degree) your coding and architectural skills. You're given a simple problem so you can focus on showcasing development techniques. It's up to you to strike the right balance between rapidly meeting our test criteria and showing off what you can do.

You're **allowed and encouraged** to use third party libraries, as long as you put them together yourself **without relying on a framework or microframework** to do it for you. An effective developer knows what to build and what to reuse, but also how his/her tools work. Be prepared to answer some questions about those libraries, like why you chose them and what other alternatives you're familiar with.

As this is a code review process, please avoid adding generated code to the project. This makes our jobs as reviewers more difficult, as we can't review code you didn't write.

## Prerequsites

We use [Docker](https://www.docker.com/products/docker) to administer this test. This ensures that we get an identical result to you when we test your application out. If you don't have it already, you'll need Docker installed on your machine. **The application MUST run in the Docker containers** - if it doesn't we cannot accept your submission. You **MAY** edit the containers or add additional ones if you like, but this **MUST** be clearly documented.

We have provided some containers to help build your application in Go, with a variety of persistence layers available to use.

### Technology

- Valid Go 1.7/1.8 code.
- Persist data to either Postgres, Redis, or MongoDB (in the provided containers).
    - Postgres connection details:
        - host: `postgres`
        - port: `5432`
        - dbname: `tokopedia`
        - username: `tokopedia`
        - password: `tokopedia`
    - Redis connection details:
        - host: `redis`
        - port: `6379`
    - NSQ connection details:
        - host: `nsqd`
        - port:
            - tcp: `4150`
            - http: `4151`
            
- Use the provided `docker-compose.yml` file in the root of this repository. You are free to add more containers to this if you like.

## Instructions

1. Fork this repository.
- Run `docker-compose up -d` to start the development environment.
- Visit `http://localhost:9000` to see the contents of the web container and develop your application.
- Create a pull request from your `fork` to the `tokopedia-interview:master` branch. 
- Pull Request should contain setup instructions for your application and a breakdown of the packages you chose to use and design decisions you made.

### API

Your application **MUST** have an interface for interaction to another application. You are free to choose any kind of interface you like.

## Requirements

We'd like you to build a simple Order Management Service. The API **MUST** provide the following functionality:

- Get List order
- Create order
- Process order
- Cancel order
- Finalize order
- Search order

### Schema

- **Order**
    - Order ID
    - Shop ID
    - Customer ID
    - Order Status (finish)
    - Products [multiple products]

## Submission
Just simple create a Pull Request, and we will review it as soon as possible.

---
Good Luck !

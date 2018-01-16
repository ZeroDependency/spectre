Spectre - A Golang Gin Gonic Middleware Framework for Distributed Simulation Testing
============

Description
------------

Inspired by https://medium.com/netflix-techblog/https-medium-com-netflix-techblog-simone-a-distributed-simulation-service-b2c85131ca1b I wanted to write a similar system to use in the numerous projects that I'm working on. 

As a small team it is important to be able to include as much automated testing as possible to ensure that the code I'm working does the job I intend. In larger teams when you have dozens or even hundreds of developers working on a project, it becomes impossible for a QA team to easily test everything that is being worked on so, again, automated testing comes in handy.

There are a number of API Testing frameworks around already such as Jasmine that allow you to define your API Specification and run tests against them as part of your CI process which can allow you make sure you don't break any API contracts, but I have personally found these dauting to fully understand or write from scratch and maintaining a Test codebase alongside your feature development can result in a lot of time being spent writing these tests (which is not necessarily a bad thing!).

The idea behind Spectre is to allow developers to write extremely small configuration files that determine when a test should trigger, and what response to provide. Spectre then sits as a piece of Middleware as part of a Gin Gonic API Server and can be deployed into any environment (even production!) with very little impact on end users.

Whilst many testing frameworks require you to create seed data this can get tedious when you have large objects that you need to configure before the tests can even run. Or perhaps you are more interested in testing your mobile app than the API itself.

Take the rather trivial example of a new user sign up, and wanting to ensure that the mobile app prompts the user to verify their email address if it hasn't been done so. Using a normal testing framework you would need to either enter data directly into a database or hit an endpoint to generate a new user. If you're working with an existing database, you need to ensure that the email address and user details you're using in your test are unique. You're then also creating a new user every time the test runs.

Spectre allows you to do away with that and instead provide a simple JSON configuration file to the Spectre Server that says "When x conditions are met, return y with an HTTP status code of z". You're not needing to insert data for the lifetime of the test, you're not testing the business logic of the API itself (which is not what you're concerned about) instead you are getting a perfectly valid response under conditions you control, and ensure that your application performs how it should. With no cleanup required aferwards.

How To Configure
------------

Example Test Configurations
------------

Performance Impacts
------------

Future Improvements
------------

* Store Test Configuration in Database
* Provide Spectre UI webpage for creating/editing tests
* Provide Logging to see Requests that triggered a test
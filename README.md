Spectre - Golang Gin Gonic Middleware Framework for Distributed Simulation Testing
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

Becase Spectre is designed to allow you to run Simulation Testing against Production environments, it specifies an invocationCount in the configuration, once this number of invocations have been reached, the test will be no longer be run. This helps ensure that any bad tests do not cause any runaway issues. 

Specifying items in the Trigger also means that ALL those conditions have to be met, allowing you to be as granular as necessary to reduce impact on normal users.

How To Configure
------------

~~~
go get -u github.com/ZeroDependency/spectre
~~~

Simply include the SpectreTest middleware at a Global, Route or Handler level (best used as a global middleware)

~~~
r := gin.Default()

// The value passed here should match the 'service' parameter from the test definition
r.Use(middleware.SpectreTest("passgen"))

r.POST("/password/generate/:id", httpPostPassword)
r.Run(fmt.Sprintf("%v:%v", ipAddress, port))
~~~

Running Spectre Server

~~~
make
./bin/spectre
~~~

This will have the Spectre Server running on port 18080 and searching a folder call definitions for test definitions. Ensure your servers with the Spectre Middleware have an appropriate SPECTRE_SERVER environment variable set up (for eg "http://localhost:18080")

Test Definition Configuration
------------

Test Definition Configuration files should be stored in a folder called definitions in the Spectre server working directory. They are JSON files with one test definition per file. The filename is not imporant. Below is an example configuration that triggers when a member of the Body data structure called "length" has a value of 12.

    {
        "id": "password-generation-bad-request",
        "name": "Password Generation Bad Request Test",
        "service": "passgen",
        "url": "/password/generate/1",
        "invocationCount": 100,
        "response": null,
        "responseCode": 400,
        "trigger": {
            "headers": null,
            "parameters": null,
            "query": null,
            "body": {
                "length": 12
            }
        }
    }

| Parameter | Use | Type |
| --- | --- | --- |
| id | Unique identifier of the test definition | string |
| name | Friendly name to desribe the test | string |
| service | The service identifier this is related to. You specify a service in the Middleware configuration of the service under simulation | string |
| url | The URL to match the test to | string |
| invocationCount | Once this number of invocations has been reached, the test is removed from the definition cache and is no longer executed | number |
| response | The JSON body to respond with when a call matches the trigger requirements | object |
| responseCode | The HTTP Status code to respond with when a call matches the trigger requirements | number |
| trigger | This structure describes the trigger conditions. If all sections are null, the test is ignored. It is essentially a number of key/value lists | object |

Performance Impacts
------------

TBC

Future Improvements
------------

* Store Test Configuration in Database
* Provide Spectre UI webpage for creating/editing tests
* Provide Logging to see Requests that triggered a test
* Remove coupling with Gin Gonic (move away from gin.Context to http.Request)
* Performance Improvements
* Unit/Integration Testing and Coverage Reports
* Memcache-style Cache Layer to prevent querying Spectre Server
* User definable Cache Timeout for Middleware Client (currently hard-coded to 60 seconds)

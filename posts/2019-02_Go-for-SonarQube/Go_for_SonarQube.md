
# Go for SonarQube

Static code analysis is a great and easy way to discover bugs, race conditions, code smells or to check whether code matches the coding conventions. I will motivate why it’s useful to use SonarQube for static analysis of Go code and show you how it’s done.

## Go Vet

Go already ships with the tool *vet* which does static code analysis of Go code. To use *vet* just run `go vet source/directory/*.go`. Vet will perform various checks like detecting shadowed variables, nil function comparison and many others (see `go doc cmd/vet` for a full list) and list all violations. Vet is great because it performs many basic checks very fast.

## Go Linters

Yet if we want to detect style mistakes, duplicated code or even check for security problems go vet is not enough. All these checks can be done by linters like [Go Meta Linter](https://github.com/alecthomas/gometalinter) or [GolangCI-Lint](https://github.com/golangci/golangci-lint). These are [classic linters](https://en.wikipedia.org/wiki/Lint_(software)) that perform static code analysis on go code and report their results in a standardized format. These linters can be integrated in editors like [VisualStudio Code](https://code.visualstudio.com/), [Vim](https://www.vim.org/) or [GoLand](https://www.jetbrains.com/go/). So, these linters perform in-depth analysis of Go code and can detect many kinds of violations. Depending on the linter configuration the linting can take a significant time in your build pipeline.

Static code analysis is also useful to calculate key metrics of your source code like lines of code or cyclomatic complexity. Combined with code coverage or unit test count these key metrics can provide a good first impression on the state of your source code. They can be displayed on a dashboard like [atlasboard](https://bitbucket.org/atlassian/atlasboard) in your project for all to see.

But where do we get these key metrics from? Linters like Go Meta Linter can calculate most of the key metrics like[ cyclometric complexity](https://en.wikipedia.org/wiki/Cyclomatic_complexity). Other key metrics like code coverage of unit test count are calculated from the tests. That way we can calculate the key metrics during the continuous integration pipeline.

## Are Linters not enough? Why SonarQube?

What if we want to know how the key metrics of our source code evolved over time? Have we been able to improve code coverage since the last release? How many lines of code do we add each sprint? Which components in our projects are the biggest?

That’s where [SonarQube](https://www.sonarqube.org) comes into play. SonarQube analyses source code using static code analysis, code coverage and unit tests over time. That way SonarQube is able to answer all the questions above. You see that version 1.0.0 had 1.562 lines of code with a coverage of 85%, whereas now you are at version 2.1.0 with 3.842 lines of code. If you look closer, you see that the package *auth* has grown the most. In SonarQube you see how your key metrics and your source code have evolved over time and over specific versions like shown below.

![SonarQube results over time](sonarqube_over_time.png?raw=true)*SonarQube results over time*

If you have a project or a company developing services in different programming languages, you have another benefit of using SonarQube for key source code metrics. SonarQube is able to calculate these key metrics for many different programming languages like Go, Java, C#, JavaScript and many others. That way you can calculate combined key metrics even for projects or companies which use many different programming languages. E.g. you your project overall has 12.346 lines of code with an overall code coverage of 76%. The JavaScript frontend has 3.704 lines of code which is round about one third of your Java backend.

Let’s see how we can analyse our Go code using SonarQube.

## SonarQube for Go

To analyse Go code with SonarQube you need a running SonarQube. You can run SonarQube locally with docker run --name sonarqube -p 9000:9000 sonarqube. The Go plugin comes with SonarQube so no need to install it manually. Now you are ready to analyse your source code. For that we use the [sonar-scanner tool](https://docs.sonarqube.org/display/SCAN/Analyzing+with+SonarQube+Scanner) provided by SonarQube, which needs a Java Runtime Environment. Next we need a configuration file for SonarQube named *sonar-project.properties* in the root of your project. The configuration file tells SonarQube which sources to analyse, the name and version of the project and where to find the test coverage. See an example *sonar-project.properties* below:

    sonar.projectKey=de.red6:service_sonar
    sonar.projectName=service_sonar
    sonar.projectVersion=1.0.0
    sonar.host.url=http://localhost:9000
    sonar.login=**SECRET**
    
    sonar.sources=.
    sonar.exclusions=**/*_test.go,**/vendor/**,**/testdata/*
     
    sonar.tests=.
    sonar.test.inclusions=**/*_test.go
    sonar.test.exclusions=**/vendor/**
    sonar.go.coverage.reportPaths=bin/cov.out

Before you can analyse your code with SonarQube you need to run the tests and record the code coverage (e.g. with ``go test -short -coverprofile=bin/cov.out `go list./..|grep -v vendor/```). SonarQube will use the recorded code coverage.

Now you can run your first SonarQube analysis using *sonar-scanner*. Your code will be analysed and uploaded to SonarQube including the code coverage.

## Example

You find a sample project [service_sonar](https://github.com/remast/service_sonar) with the full setup for SonarQube. In the sample project you can run the SonarQube analysis with *make sonar. *If you don’t have the sonar-scanner tool installed on your machine you can also use a dockerized version of the sonar scanner (see [sample project](https://github.com/remast/service_sonar/blob/master/README.md)). The dockerized version is a great fit for CI environments like [GitlabCI](https://about.gitlab.com/product/continuous-integration/   ).

All settings in the configuration file *sonar-project.properties* can also be provided using the command line. This makes sense especially for the SonarQube host and login. So in you CI pipeline you can run `sonar-scanner-Dsonar.host.url="http://sonar.mycompany.com" -Dsonar.login="**MY_TOKEN**"`.

The sample project also includes a configuration for Gitlab CI in *.gitlab-ci.yml*. The Gitlab CI pipeline performs a SonarQube analysis in the stage sonar. It expects the host and login for SonarQube in the environment variables `SONAR_HOST` and `SONAR_API_TOKEN`.

After running the SonarQube analysis you see the results in SonarQube:

![Result of first SonarQube analysis](sonarqube_first_analysis.png?raw=true)*Result of first SonarQube analysis*

![Overview of SonarQube results](sonarqube_overview_results.png?raw=true)*Overview of SonarQube results*

## Summary

We have learned how to statically analyse Go code using the build in tool *vet*, linters like Go Meta Lint or the SonarQube plugin for Go. All of these tools are capable of calculating key metrics of your source code.

SonarQube provides two benefits:

* First it can track changes of your key metrics over time and by that provide useful insights on how your code quality is evolving.

* Second SonarQube can analyse code in many different programming languages and so you can calculate overall key metrics.

So you should definitely give [SonarQube](https://www.sonarqube.org) for Go a try.

*This post originally appeared on [medium](https://medium.com/red6-es/go-for-sonarqube-ffff5b74f33a).*
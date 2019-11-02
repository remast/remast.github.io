# Go for SonarCloud with Github Actions

Metrics like lines of code or test coverage are great to track and improve the quality of your source code. SonarQube can calculate these metrics for your project and track how they evolve over time. Since SonarQube natively supports Go it's a great fit to calculate metrics fo your Go project. 

Learn the basics of analyzing a Go project with SonarQube in my post [Go for SonarQube](https://medium.com/red6-es/go-for-sonarqube-ffff5b74f33a). In this post I'll show you how to use [Github Actions]() to analyze your Go project with [SonarCloud](https://sonarcloud.io). SonarCloud offers SonarQube as a service in the cloud.

## Set up a Build with Github Actions

We start by setting up a basic build pipeline for our Go project using Github Actions. We use the Github's official Go action [setup-go](https://github.com/actions/setup-go) for our build. So we create our pipeline for Github Actions in the file `.github/workflows/build.yml` with the following content:

```yaml
on: push
name: Main Workflow
jobs:
  build:
    name: Build, Test and Analyze
    runs-on: ubuntu-latest
    steps:
      - name: Setup Go
        uses: actions/setup-go@v1
        with:
          go-version: '1.13'
      - name: Clone Repository
        uses: actions/checkout@master
      - name: Build and Test
        run: make test
```

Now Github runs our pipeline on every push to the git repository. The pipeline set's up an Go environment, clones the repository and finally builds and tests with `make test`. The tests also calculate the code coverage of the tests and save it in `bin/cov.out`.

## Run the SonarCloud Analysis as Step
Now we can add the SonarCloud analysis as step in our pipeline using the [sonarcloud-github-action](https://github.com/SonarSource/sonarcloud-github-action). But before that you need to register for SonarCloud. You can use your Github account to sign up. 

The step to add the SonarCloud analysis is:

```yaml
      - name: Analyze with SonarCloud
        uses: sonarsource/sonarcloud-github-action@master
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          SONAR_TOKEN: ${{ secrets.SONAR_TOKEN }}
```
The step analyzes our Go code using the [sonar-scanner tool](https://docs.sonarqube.org/latest/analysis/scan/sonarscanner/).

The SonarCloud Action needs two environment variables. The first one is `GITHUB_TOKEN` which is already provided by Github (see [Virtual environments for GitHub Actions](https://help.github.com/en/github/automating-your-workflow-with-github-actions/virtual-environments-for-github-actions)). The second one is the `SONAR_TOKEN` to authenticate the Github Action with SonarCloud. 

To generate the access token `SONAR_TOKEN`  log into SonarCloud. Now click on your profile and then go to 'My Account' and then 'Security'. Or go directly to [account/security](https://sonarcloud.io/account/security/). Generate your access token for SonarCloud here. The access token is provided to the build pipeline as a secret environment variable. Go to your repository settings in Github. Then to 'Secrets' and add a new secret with name `SONAR_TOKEN` and use the generated SonarCloud access token as value, as shown below.

![Github Repostory Secrets](https://github.com/remast/remast.github.io/raw/master/articles/2019-11_Go-for-SonarCloud-with-Github-Actions/github-action-secrets.png)

Our new pipeline to build, test and analyze the source code with SonarCloud is shown below.

```yaml
on: push
name: Main Workflow
jobs:
  build:
    name: Build, Test and Analyze
    runs-on: ubuntu-latest
    steps:
      - name: Setup Go
        uses: actions/setup-go@v1
        with:
          go-version: '1.13'
      - name: Clone Repository
        uses: actions/checkout@master
      - name: Build and Test
        run: make test
      - name: Analyze with SonarCloud
        uses: sonarsource/sonarcloud-github-action@master
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          SONAR_TOKEN: ${{ secrets.SONAR_TOKEN }}
```

## Add Code Coverage to SonarCloud
SonarCloud is able to analyze our source code. Yet it's not able to pick up the code coverage of our tests which is stored in the file `bin/cov.out`. And this is a bit tricky. The sonar-scanner looks for the code coverage in file `/home/runner/work/service_sonar/service_sonar/bin/cov.out`. That's where it is after the tests. The problem is the sonar-scanner tool is executed with docker and mounts the current directory as a volume like `docker run --workdir /github/workspace [...] -v "/home/runner/work/service_sonar/service_sonar":"/github/workspace"`. Within the docker run there is no file `/home/runner/work/service_sonar/service_sonar/bin/cov.out` but the path of that file within docker is actually `/github/workspace/bin/cov.out`. So we need to use that path as `sonar.go.coverage.reportPaths` in the SonarQube settings file `sonar-project.properties`. The full `sonar-project.properties` for Github Actions are shown below: 

```
# Github organization linked to sonarcloud
sonar.organization=remast

# Project key from sonarcloud dashboard for Github Action, otherwise pick a project key you like
sonar.projectKey=de.red6:service_sonar

sonar.projectName=service_sonar
sonar.projectVersion=1.0.0

sonar.sources=.
sonar.exclusions=**/*_test.go,**/vendor/**,**/testdata/*
 
sonar.tests=.
sonar.test.inclusions=**/*_test.go
sonar.test.exclusions=**/vendor/**
sonar.go.coverage.reportPaths=/github/workspace/bin/cov.out
```

A live example of the pipeline can be found at [service_sonar](https://github.com/remast/service_sonar/tree/feature/gh-pipeline-simplify) (branch `feature/gh-pipeline-simplify`).

## Our project in the SonarCloud Dashboard
After the first SonarCloud analysis you can see the results on the SonarCloud project dashboard shown below.

![SonarCloud Dashboard](https://github.com/remast/remast.github.io/raw/master/articles/2019-11_Go-for-SonarCloud-with-Github-Actions/go-sonarcloud-dashboard.png)

Check out dashboard of the sample project at https://sonarcloud.io/dashboard?id=de.red6%3Aservice_sonar.

## Run the SonarCloud Analysis as Job
Currently our pipeline has two main steps. The first step is build and test and the second step is the SonarCloud analysis. In big projects the SonarClould analysis can take up to several minutes. Since steps are executed synchronously one after the other the pipeline can only continue once SonarCloud is done analyzing. But what about our integration tests? Say we want to run some integration tests after the build. We want to give feedback about the integration tests as soon as possible. This means we don't want to wait for the SonarCloud analysis but proceed to the integration tests right after build and test. This is easily possible if we separate out the SonarCloud analysis as a separate job.

Github Actions runs jobs in a workflow in parallel. So we can split our pipeline in three jobs. The first job is build and test. Then we have a second job to for the SonarCloud analysis and a third job that runs the integration tests. The job to build and test has to run first. After that both the job for the SonarCloud analysis and the integration tests can run in parallel. 

### First pipeline draft
The first draft of our improved pipeline shown below. It has the three jobs `build`, `sonarCloudTrigger` and `integrationTest`. The jobs `sonarCloudTrigger` and `integrationTest` both depend on the job build which is declared with `needs: build`.

```yaml
on: push
name: Main Workflow
jobs:
  build:
    name: Compile and Test
    runs-on: ubuntu-latest
    steps:
      - name: Clone Repository
        uses: actions/checkout@master
      - name: Setup go
        uses: actions/setup-go@v1
        with:
          go-version: '1.13'
      - run: make test

  sonarCloudTrigger:
    needs: build
    name: SonarCloud Trigger
    runs-on: ubuntu-latest
    steps:
      - name: Clone Repository
        uses: actions/checkout@master
      - name: Analyze with SonarCloud
        uses: sonarsource/sonarcloud-github-action@master
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          SONAR_TOKEN: ${{ secrets.SONAR_TOKEN }}

  integrationTest:
    needs: build
    name: Integration Test
    runs-on: ubuntu-latest
    steps:
      - run: echo Should run integration tests.
```
The pipeline is executed the way we wanted. The job `build` runs first and after that `sonarCloudTrigger` and `integrationTest` run in parallel. It has only one flaw though. Do you see it? If you check the logs you will notice that SonarCloud is not able to read the code coverage of the unit tests. SonarCloud looks for the code coverage in the file `bin/cov.out`. This file is only present in the `build` job but not in the job `sonarCloudTrigger`. So we need to transfer the file from one job to the other.

## Revised pipeline with Code Coverage
To transfer the file `bin/cov.out` from the `build` job to the `sonarCloudTriggerJob` we'll use the actions [upload-artifact](https://github.com/actions/upload-artifact) and [download-artifact](https://github.com/actions/download-artifact). We add the action `upload-artifact` to the `build` job and `download-artifact` to the `sonarCloudTriggerJob` job like below:

```yaml
on: push
name: Main Workflow
jobs:
  build:
    name: Compile and Test
    runs-on: ubuntu-latest
    steps:
      - name: Clone Repository
        uses: actions/checkout@master
      - name: Setup go
        uses: actions/setup-go@v1
        with:
          go-version: '1.13'
      - run: make test
      - name: Archive code coverage results
        uses: actions/upload-artifact@v1
        with:
          name: code-coverage-report
          path: bin

  sonarCloudTrigger:
    needs: build
    name: SonarCloud Trigger
    runs-on: ubuntu-latest
    steps:
      - name: Clone Repository
        uses: actions/checkout@master
      - name: Download code coverage results
        uses: actions/download-artifact@v1
        with:
          name: code-coverage-report
          path: bin
      - name: Analyze with SonarCloud
        uses: sonarsource/sonarcloud-github-action@master
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          SONAR_TOKEN: ${{ secrets.SONAR_TOKEN }}

  integrationTest:
    needs: build
    name: Integration Test
    runs-on: ubuntu-latest
    steps:
      - run: echo Should run integration tests.
```
A live example of the pipeline can be found at [service_sonar](https://github.com/remast/service_sonar/tree/master) (branch `master`).

## Add a SonarCloud Badge
We can add a SonarCloud badge to our project to quickly show the SonarCloud status of our project from within the readme. To create SonarCloud badge go to your SonarCloud project dashboard and click 'Get project badges'. You can choose between three badges as shown below.

![SonarCloud Badges](https://github.com/remast/remast.github.io/raw/master/articles/2019-11_Go-for-SonarCloud-with-Github-Actions/sonarcloud-badges.png)

## Wrap Up
Now you have learned how to analyse Go code using SonarCloud from a Github Action pipeline. By using a separate job for the SonarCloud analysis we are able to run the  time-consuming analysis in parallel to other important tasks like integration tests or deployment.

All code is provided as a running example in the Github project [service_sonar](https://github.com/remast/service_sonar).
#Using SonarCloud with Github Actions and Maven

In this post you will will learn how to analyse your Java Maven project with SonarCloud using Github Actions. 

Starting point is a simple Java project with a Maven build. First we'll use SonarCloud to analyze our source code from our local dev environment. Then we'll use Github Actions to run the Maven build. So finally we have a fully functional ci pipeline which builds and analyzes our code using Github Actions.

## Set up SonarCloud

### Step 1: Create a Project
In order to use [SonarCloud](https://sonarcloud.io/) you need to create an account and set up a project. So first create an account and log in. Now you can create a new project [here](https://sonarcloud.io/projects/create) or using the '+' button. A project in SonarCloud must belong to an organization. SonarCloud automatically imports your Github organizations. So you can use any of your Github organizations or use the default organization by your Github user name. 

![SonarCloud Create Project](https://github.com/remast/remast.github.io/raw/master/posts/2020-03_SonarCloud_with_GithubAction_and_Maven/SonarCloud_CreateProject.png)

After you've created your project, your project has an organization key and a project key. You'll need both to run an SonarCloud analysis. You can always look up organization key and project key from the dashboard
of your project like shown below.

![SonarCloud Project Key](https://github.com/remast/remast.github.io/raw/master/posts/2020-03_SonarCloud_with_GithubAction_and_Maven/SonarCloud_ProjectKey.png)

### Step 2: Generate a SonarCloud Token
Now we'll set up a secure token as authentication for SonarCloud. Generate a new token from the tab 'security' in your account settings (which is [here](https://sonarcloud.io/account/security/)). Make sure to store the token since you'll only see it right after you've greated it. Now you're all set up to run a first analysis.

![SonarCloud Generate Token](https://github.com/remast/remast.github.io/raw/master/posts/2020-03_SonarCloud_with_GithubAction_and_Maven/SonarCloud_GenerateToken.png)

## Run SonarCloud analysis locally using Maven
You can run the SonarCloud analysis using maven. The organization key, project key and the generated token must be passed to the [Sonar Maven Plugin](https://docs.sonarqube.org/latest/analysis/scan/sonarscanner-for-maven/) as well as the url for SonarCloud. Replace `<GENERATED_TOKEN>` with the SonarCloud token you generated in the previous step. So the command is:

```bash
mvn sonar:sonar \
   -Dsonar.projectKey=baralga \
   -Dsonar.organization=baralga \
   -Dsonar.host.url=https://sonarcloud.io \
   -Dsonar.login=<GENERATED_TOKEN>
```

After you ran the analysis the results will shortly be online in the SonarCloud dashboard at `https://sonarcloud.io/dashboard?id=<projectKey>`. The dashboard for our sample project `baralga` is available at [https://sonarcloud.io/dashboard?id=baralga](https://sonarcloud.io/dashboard?id=baralga).

## Run SonarCloud analysis using Github Actions
Now we will use Github Actions to run the SonarCloud analysis from our ci pipeline. We'll use Maven for that like we did before. We set up Github Action that runs the SonarClound analysis using Maven. Like before we pass organization key and project key as parameters. Additionally we need to provide the SonarCloud token and the [Github Token](https://help.github.com/en/actions/configuring-and-managing-workflows/authenticating-with-the-github_token).

The token for SonarCloud is stored as a encrypted secret as described [here](https://help.github.com/en/actions/configuring-and-managing-workflows/creating-and-storing-encrypted-secrets). We can access it in our Github Action with `${{ secrets.SONAR_TOKEN }}`. The [Github Token](https://help.github.com/en/actions/configuring-and-managing-workflows/authenticating-with-the-github_token) is already provided by Github Actions itself and we can access it with `${{ secrets.GITHUB_TOKEN }}`.


And here's the complete Github Action workflow. Save it in the file `.github\workflows\sonar.yml` and off you go.
```yaml
name: SonarCloud
on:
  push:
    branches:
      - master
jobs:
  build:
    runs-on: ubuntu-16.04
    steps:
    - uses: actions/checkout@v1
    - name: Set up JDK
      uses: actions/setup-java@v1
      with:
        java-version: '11'
    - name: Analyze with SonarCloud
      run: ./mvnw -B verify sonar:sonar -Dsonar.projectKey=baralga -Dsonar.organization=baralga -Dsonar.host.url=https://sonarcloud.io -Dsonar.login=$SONAR_TOKEN
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        SONAR_TOKEN: ${{ secrets.SONAR_TOKEN }}
```

Now every build of your Github Actions pipeline analyzes the code using SonarCloud. That looks like below or see it live in action in the [Baralga Actions](https://github.com/Baralga/baralga/actions).

![Baralga SonarCloud Build](https://github.com/remast/remast.github.io/raw/master/posts/2020-03_SonarCloud_with_GithubAction_and_Maven/SonarCloud_GithubActionsBuild.png)

## Topping it off with Code Coverage
As last step we add the code coverage of our unit tests to SonarCloud. 

### Step 1: Calculate Test Coverage with Jacoco
We use [Jacoco](https://www.eclemma.org/jacoco/) to calculate the code coverage of our tests. For that we add the [jacoco-maven-plugin](https://www.eclemma.org/jacoco/trunk/doc/maven.html) to our `pom.xml`:

```xml
<plugins>
  ...
  <plugin>
    <groupId>org.jacoco</groupId>
    <artifactId>jacoco-maven-plugin</artifactId>
    <version>0.8.5</version>
    <executions>
      <execution>
        <id>prepare-agent</id>
        <goals>
          <goal>prepare-agent</goal>
        </goals>
      </execution>
      <execution>
        <id>report</id>
          <phase>prepare-package</phase>
          <goals>
            <goal>report</goal>
        </goals>
      </execution>
    </executions>
  </plugin>
<plugins>
```

We tell SonarCloud where to find the calculated code coverage using the parameter `-Dsonar.coverage.jacoco.xmlReportPaths=${project.build.directory}/site/jacoco/jacoco.xml` for our Maven build. After the next build the code coverage will show up in SonarCloud.

## Summary

Step by step we introduced SonarCloud to analyze our code within our ci pipeline using Github Actions. Whenever the ci pipeline runs, the code is analyzed using SonarCloud and the results and metrics are available
in the SonarCloud dashboard. You can find a working example at [baralga](https://github.com/Baralga/baralga).

{% github github.com/Baralga/baralga no-readme %}
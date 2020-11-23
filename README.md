# camas
A zero-dependency tool for finding secrets in source code directories.

## Building
Maven is required for building.
```bash
# Building
CGO_ENABLED=0 GOOS=darwin go build -a -installsuffix cgo . 

# Running
./camas -configFile my-camas-config.json -inputDirectory /path/to/my/directory -outputFormat json -outputFile myAppSecrets.json
```

## Arguments
Arg | Description | Required
------ | ------ | ------
inputDirectory | Directory containing source code to analyze (Required) | true
configFile | Configuration File (Required) | true
outputFile | Output File | false
outputFormat | Output Format [txt, json] (defaults to txt) | false
noiseLevel | minimum noise level to report on (defaults to 0/all) | false
numWorkers | number of go workers to execute (default to #cpus) | false
cpuProfile | write cpu profile to file (only used for debugging) | false

## Report

#### The output txt is as follows:

```
[/path/to/my/app/connect.pp:5 (Username and password in URI) "  options => 'http-proxy="http://proxyuser:proxypass@example.org:3128"',"]
```

#### The output txt is as follows:

```Javascript
  {
    "rule-name": "Username and password in URI",
    "absolute-file-path": "/path/to/my/app/connect.pp",
    "line-number": 5,
    "content": "  options => 'http-proxy=\"http://proxyuser:proxypass@example.org:3128\"',",
    "noise-level": 0
  }
```

## Configuration File
To generate a config file, you can use the rule_generation/load_data.sh script if you like. It pulls in data from several other open-source projects. Camas does not attempt to generate our own rules, rather leverages the rules from other projects. The generation script can be run by using:

```bash
./load_data.sh | jq . > camas_config.json
```

## Docker Support
As camas is a zero-dependency golang app, the docker container is built from scratch, which makes it small and secure. 

```bash
docker build -t <tag_of_your_choice> .
```

Your docker container will need at least one mount point for the directory containing your app. Here is an example.
```bash
docker run --read-only -v /source/path/to/app:/path/to/app/in/container -it <tag_built_with> -inputDirectory /path/to/app/in/container -configFile camas-config.json
```
The --mount variant of mounting a volume can also be used if desired. If you want to write the output to a location outside of your container, then you will have to set a second mount point or re-use the existing one. If memory issues are encounterd, try running container with increased memory using the -m argument.
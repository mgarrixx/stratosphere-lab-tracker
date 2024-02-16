# Stratosphere-lab-Tracker

A simple tools for collecting packet capture files (as well as malware samples) from [Stratosphere Research Laboratory](https://www.stratosphereips.org/).

![](https://media.licdn.com/dms/image/C4E16AQH93jzMATzgGQ/profile-displaybackgroundimage-shrink_200_800/0/1517698799721?e=2147483647&v=beta&t=-7mcvU9luGmhFt4321oamQMOrwMJe_ILpO90vY0vi58)

## Usage
Clone this repo
```
$ git clone github.com/mgarrixx/stratosphere-lab-tracker
$ cd stratosphere-lab-tracker
```

And then, just `go`
```
# [Data Collection Type]: Choose one among `android`|`malware`|`normal`|`iot`
# [Directory PATH]: Directory where the downloaded resource will be placed in

$ go run main.go -source=[Data Collection Type] -save_path=[Directory PATH]
```
# Radish
*v0.1.0*

Radish is a simple and powerful Redis GUI administration panel

## Features:
As for now Radish supports only managing (adding, deleting, updating) different types of keys and values

### Features soming soon (or later...):
* Tree view for keys
* Slowlog
* Monitoring instances load
* Keyboard shortcuts
* Authorization

## Running Radish
The easiest way to run Radish is using Docker image:

```
$ docker pull sad0vnikov/radish
$ mv config.json.example config.json
```
Edit config.json and add your Redis hosts

```
$ docker run docker run -p 8080:8080 -d -v /full/path/to/config.json:/config.json sad0vnikov/radish
```

Your Radish instance we'll be accessible on localhost:8080

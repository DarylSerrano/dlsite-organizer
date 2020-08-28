# DLsite Organizer
Simple cli tool for organizing DLsite works on local folder.

## About
* This tool creates a sqlite db file called `data.db` on working directory or provided directory. 
* For each file containing an `RJxxxxxx`, scraps the information of that work and saves them into the db. 
* Filters works saved on the db by VA, Circle, Tags and SFW/NSFW.

## Instalation
```bash
$ go get -u github.com/DarylSerrano/dlsite-organizer
$ go install github.com/DarylSerrano/dlsite-organizer
```
* For Windows OS you must have the `gcc` compiler installed. You can use [TDM-GCC](https://jmeubank.github.io/tdm-gcc/)

## Usage
```bash
$ dlsite-organizer [command]
$ dlsite-organizer refresh
$ dlsite-organizer filter sfw # Filter by sfw
$ dlsite-organizer filter tags --db "./test/testdata" "./test/test" # Filter by tags
```

## Contribute
This is a personal project but feel free to contribute (PR-s and issues are welcomed).

## License
[Apache license 2.0](./LICENSE)
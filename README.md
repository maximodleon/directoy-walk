# Directory walk

This a simple tool to list, archive and delete files in a given folder

## Options
```
  -root
   Root directory to start
  -logFile
   Log deletes to this file
  -list
   List files only
  -del
   Delete file
  -archive
   Archive directory
  -ext
   File extension to filter out
  -size
   Minimum file size
```

## Examples

List files

```bash
go run . -root <DIR> -ext <FILE EXTENSION> -log <LOG FILE PATH>
```

Delete files

```bash
go run . -root <DIR> -ext <FILE EXTENSION> -log <LOG FILE PATH> -d
```

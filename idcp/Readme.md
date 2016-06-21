# idcp

Idempotant copy of a file between locations.

## Usage

Short form:

`idcp source_file target_file`

Long form (recommended for scripts):

```
Usage of idcp:
  -dest string
    	destination that the file should be copied to
  -group string
    	group of the destination file
  -mode string
    	file mode of the destination file (default "0644")
  -owner string
    	owner of the destination file
  -source string
    	source of the file to be copied
```

## Details

If the source file and target file are identical, no copy is made. If they are
not the same, then the source file is copied to the target file location with a
'.tmp' extension, and atomically renamed to the target file name. This will
copy the file without affecting any other files which were hardlinked to the
original target file.

After any copy operation, the ownership and mode of the target file are
changed. If not specified, the current user and a mode of 0644 are used. Note
that if the group is not explicitly specified but the owner is, the primary
group of the provided owner is.

## Bugs

The application of user, group, and mode to the target file are not idempotant
right now, as the Go library does not provide a method for obtaining the
current values. This is high on the list of priorities to change.

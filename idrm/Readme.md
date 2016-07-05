# idrm

Idempotant removal of a file

## Usage

Short Form:

`idrm file_path`

Long Form (recommended for scripts):

```
Usage of ./idrm:
  -filepath string
        path of the file to be removed
```

## Details

If the file path exists, it is unlinked. Unlinking occurs in the following
manner:

1) Create a hardlink to filepath + ".rmtmp"
2) Unlink the original filepath
3) Unlink the temporary filepath

Unlinking occurs in this manner since the act of unlinking large files on some
file systems creates a long lock on anything which is watching or depends upon
that file. This method ensures that programs which look for the original file
can resume as quickly as possible, regardless of how long the actual file
removal takes.

Does not support removing directories.

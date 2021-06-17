# blackzim
blackzim is an apheleia formatter, which runs zimports and black.

## Usage
Install the blackzim executable with
```
go install github.com/schmir/blackzim@latest
```

Declare the blackzim formatter inside your emacs init file:

```elisp
(setf (alist-get 'blackzim apheleia-formatters)
      '("blackzim"))
```
Please consult apheleia's documentation on how to use that formatter.

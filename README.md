WIP
---
This project is a *work in progress*. The implementation is *incomplete* and
subject to change. The documentation can be inaccurate.

Vanilj
======
Vanilj is a mandelbrot image generator. It's name is a reference to "vanilla
dreams" which is an awesome Swedish cookie.

Installation
------------

`go get github.com/karlek/vanilj/cmd/vanilj`

Generate an image
-----------------
```shell
$ vanilj -z 500 -cr 1 -width 500 -height 300 -o "1.png"
```

API documentation
-----------------

* [canvas][]
* [fractal][]

	- [mandel][]

[canvas]: http://godoc.org/github.com/karlek/vanilj/canvas
[fractal]: http://godoc.org/github.com/karlek/vanilj/fractal
[mandel]: http://godoc.org/github.com/karlek/vanilj/fractal/mandel

Public domain
-------------
I hereby release this code into the [public domain](https://creativecommons.org/publicdomain/zero/1.0/).

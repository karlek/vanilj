WIP
---
This project is a *work in progress*. The implementation is *incomplete* and
subject to change. The documentation can be inaccurate.

Vanilj
======
Vanilj is a mandelbrot explorer written in go. It's name is a reference to "vanilla dreams" which is an awesome Swedish cookie.

Inspiration
-----------

[sub.blue](http://sub.blue/fractal-lab)

[sub.blue old blog](http://2008.sub.blue/blog.html)

[orbit traps](http://www.fractaldomains.com/tutorial/use-orbit-traps/)

Installation
------------

`go get github.com/karlek/vanilj/cmd/vanilj`

Generate an image
-----------------

```shell
$ vanilj -z 500 -cr 1.5018 -ci 1 -i 200 -t random -width 500 -height 300 -o "1.png"
```

Flags:
------

* __cr:__
	Real value offset from origio to zoom in on.
* __ci:__
	Imaginary value offset from origio to zoom in on.
* __i:__
	Number of iterations performed.
* __t:__
	Which color scheme to use, valid options are ["random", "pretty", "smooth", "pedagogical"].
* __z:__
	Zoom value. How many magnifications to make on the center point.
* __o:__
	Output filename.
* __width:__
	Width of created image.
* __height:__
	Height of created image.

Examples
--------

Pedagogical uses only five colors and cycles between them to represent the rate of divergence.

![Pedagogical representation of the Mandelbrot set](https://github.com/karlek/vanilj/blob/master/cmd/vanilj/pedagogical.png?raw=true)

Pretty uses a gradient of colors proportional to the number of iterations. This will make the image softer and more aesthetically pleasing.

![Pretty representation of the Mandelbrot set](https://github.com/karlek/vanilj/blob/master/cmd/vanilj/pretty.png?raw=true)

Random creates a unique gradient each time, with a number of colors proportional to the number of iterations.

![Random representation of the Mandelbrot set](https://github.com/karlek/vanilj/blob/master/cmd/vanilj/random.png?raw=true)

Smooth uses the "Normalized Iteration Count Algorithm" which removes the bands of colors with a smooth gradient. The bands can be observed in the Pedagogical representation.

![Smooth representation of the Mandelbrot set](https://github.com/karlek/vanilj/blob/master/cmd/vanilj/smooth.png?raw=true)

All images were created with:
```shell
$ vanilj -cr -1.5018 -z 1000 -i 800
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

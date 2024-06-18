<div align="center">
    <img alt="pixv" width="285" src="./assets/logo.svg">
    <p><i>Pixel art vectorization tool</i></p>
</div>

# What?

A CLI tool for converting raster pictures into vector format. The app can convert any PNG/JPEG images but for use with pixel art. Using photos or any other images with a high pixel density, a wide range of colors, or small chunks of adjacent pixels might not produce the desired result.

# Why?

Large amount of software processes raster images with interpolation. That mean that software will try to make soft color transition between adjacent pixels. Usually, such behavior is undesirable with small sized images and especially pixel art pictures as they become to blurry. To solve that problem you can use vectorized version of your image.

> Although some applications allow to disable interpolation, that almost always requires tweaking default rendering behavior.

# How does it work?

Currently tool can convert images only to svg using one of the following methods

### Square

Simply takes every pixel from image and draws square for each. Produces large sized files but with perfect Width x Hight pixel matrix. Will probably work better with some legacy software.

![pixv_square](./assets/pixv_square.svg)

```
original (png): 268 bytes
vectorized:     4495 bytes
```

### Rectangle

Similar to the previous one, but instead of drawing each pixel separately, it combines adjacent pixels of the same color in rectangular chunks. Usually generated files are much smaller.

![pixv_rectangle](./assets/pixv_rectangle.svg)

```
original (png): 268 bytes
vectorized:     2035 bytes
```

### Path

>  In development

There is also third the most efficient method that will draw path around chunks of the same color instead of breaking that chunk in a bunch of rectangles like Greedy algorithm do.

# How to use?

Following command will create svg variant of the given image in the current directory.

```
pixv image.png
```

You can use flags to customize the result

- `--method [method]`, `-m [method]` - Choose vectorization method. Accepts `square` or `rectangle`. Default method is `rectangle`

- `--scale [multiplier]`, `-s [multiplier]` - Change the scale of the pixels. Currently accepts only `integers`. Default multiplier is `1`

# How to build?

Make sure you have installed

- `Go` version 1.22.4 or higher

Use following commands to build the project

```
git clone https://github.com/axseem/pixv.git
cd pixv
go build
```


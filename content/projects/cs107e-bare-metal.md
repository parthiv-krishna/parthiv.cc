+++
title = "CS 107e"
weight = 50
extra = { featured_image = "/img/cs107e/larson-scanner.jpg", featured_image_alt = "Larson scanner written in ARM assembly." }
+++

As a student in CS 107e: Computer Systems from the Ground Up, I implemented various aspects of a computer system on a Raspberry Pi. The code was all written in bare metal C with no operating system. The aspects that I directly worked on were:

- A Larson Scanner, purely in ARM assembly!
- A clock, and a complete GPIO interfaced with the appropriate RasPi registers to enable digital reads/writes.
- A string library (`strcmp`, `strlcat`, `printf`, etc).
- A memory allocation library (`malloc`, `free`, `realloc`, etc).
- A PS/2 keyboard library and console with simple commands like `poke`, `peek`, `echo`, etc.
- A graphics library with text and shape (rectangles, anti-aliased lines, anti-aliased triangles, etc) drawing functionality.
- An interrupt handling library, to make the PS/2 keyboard more responsive when combined with costly graphics operations.

By the end of the quarter, I had created a terminal. After plugging in an HDMI monitor and PS/2 keyboard to the Raspberry Pi, I could type on the keyboard, see the characters show up on the monitor, and execute simple commands, like changing the values in a GPIO register to turn on an LED. Additionally, my version of the console terminal had a "history" capability where you could click the up and down arrow keys to access recent commands.

{{ img(path="static/img/cs107e/clock.jpg", alt="Clock breadboard.") }}

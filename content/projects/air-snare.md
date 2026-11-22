+++
title = "Air Snare"
weight = 40
extra = { featured_image = "/img/air-snare/airsnare.jpg", featured_image_alt = "Air Snare" }
+++

As part of the course, I worked on an open ended final project called *Air Snare*, which is a drum kit you can play without drums. The hardware was two "drum sticks" (pencils) with LSM6DS33 IMUs and push buttons, the Raspberry Pi, and an HDMI monitor for sound output. One stick played a kick drum or snare drum (selected with the push button), and the other played a hi-hat or crash symbol. From a software standpoint, the project involved:

- Development of I2C driver for two LSM6DS33 IMUs to read raw data.
- Filtering and gesture extraction from the IMU data.
- Audio sequencing to queue sounds to be played based on the detected gestures.
- Delivery of audio samples to the Raspberry Pi GPU to output over HDMI to the attached screen.

All of this was implemented in bare-metal C with no operating system and no `stdlib`.

{{ img(path="static/img/air-snare/drumstick.jpg", alt="Close-up of Air Snare drumstick.") }}

{{ img(path="static/img/air-snare/block-diagram.png", alt="Block diagram of Air Snare system.") }}

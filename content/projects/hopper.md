+++
title = "Hopper"
weight = 70
extra = { featured_image = "/img/hopper/cadmodel.png", featured_image_alt = "A render of the Hopper CAD model." }
+++

Additionally, as part of the course CS 229: Machine Learning, I worked on implementing reinforcement learning-based control for Hopper. Specifically, I developed a physics simulator using PyBullet, then trained two RL algorithms to control the robot in various scenarios, including stable balancing, velocity control, and jumping onto a step.

- DDPG: Deep Deterministic Policy Gradient. DDPG consists of a critic network and actor network. A copy of both networks are employed to calculate target values, and their weights are updated to slowly track those of the learned networks in order to increase stability.
- SAC: Soft Actor-Critic. SAC also uses an actor network and a critic network, as the name suggests. In contrast to the deterministic actions selected by DDPG, SAC selects actions probabilistically. The goal of this algorithm is to maximize not only the reward but also the entropy of the policy in order to encourage exploration.

We compared the results of the RL-based controllers to an LQR-based controller that my project partner developed, and found pretty positive results.

{{ img(path="static/img/hopper/velocity-control.png", alt="Results of velocity control for various algorithms.") }}

{{ img(path="static/img/hopper/jump-before.jpg", alt="Hopper simulator, before jump.") }}

{{ img(path="static/img/hopper/jump-after.jpg", alt="Hopper simulator, after jump.") }}

[Paper PDF](/pdf/hopper-rl-control.pdf)

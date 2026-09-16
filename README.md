# Pyxis
> This project is a work in progress and immature. Although it is usable, it is currently undergoing an internal rewrite and is feature incomplete.

Pyxis is a part of the Constellation Project, a set of personal applications designed around my personal workflow. The focus is on **capture** and **scheduling** actions backed by a fast and reliable local-first data layer.

The name "Pyxis" originates from the constellation Pyxis, The Mariner's Compass. The goal was to help me keep direction across long stretches of time.

# Purpose
Nearly all task managers ask me to adapt to their system. Pyxis was built as a very opinionated personal system that adapts to my own system. The core loop is inspired by Getting Things Done, and the file format being open plaintext files was inspired by Obsidian.

# Architecture
> This is the Service version of the repo. For the Front-end, see [Pyxis](https://github.com/Stelare-Rose/Pyxis)

Pyxis is split into two parts:
 - Front-End  
 The front-end is written in Nuxt and Tauri, as a standalone user-facing application. Currently, only the Linux version works and is tested.
 - Background Service  
 The background service is written in Go, and handles generation of the cache layer from the given data.
    
The data flow is intentionally a unidirectional cycle:
```
Data -> Indexer -> Cache -> Front-End App -> Data
```

# Data Contracts and Intentions

The application is designed with the following ideas in mind
1. The user should be able to edit their data freely (Through the front-end or otherwise), given it stays within the schema and file format.
2. If the app automatically modifies your data, it must produce the exact same file on every device
3. No other party modifies your data

# Platform Support
Currently the project only supports Linux, and the service can be built by running the build.sh file. It's not currently recommended to use this application though as it's still immature.

# Roadmap
- [ ] Rewrite in Rust and Svelte
- [ ] Ideas Feature
- [ ] Projects Feature
- [ ] Nix Flake Distribution
- [ ] Windows and Android Support

# License
This project is under the MIT License - see LICENSE file

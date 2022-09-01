# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

### [0.1.6](https://github.com/everdrone/grab/compare/v0.1.5...v0.1.6) (2022-09-01)


### Features

* better logging, add zerolog ([f48d6dd](https://github.com/everdrone/grab/commit/f48d6dda453f431594d6c80285e5043aca601723))
* checks for updates using github api ([f2a417e](https://github.com/everdrone/grab/commit/f2a417ebc40c518be17bfc797d8022859a2bec50))


### Bug Fixes

* make latest version check concurrent ([4a13fd8](https://github.com/everdrone/grab/commit/4a13fd8f2cbbc339234c3c6aa01f8746d99e46f0))
* make update notification parallel to download ([64384d9](https://github.com/everdrone/grab/commit/64384d98271f813ee956373b697aed6185dd1116))
* replace ioutil with io (SA1019) ([b5895df](https://github.com/everdrone/grab/commit/b5895dfcd0fdc9feb3d092b76e0e89e4ac5c19cc))

### [0.1.5](https://github.com/everdrone/grab/compare/v0.1.4...v0.1.5) (2022-08-21)


### Bug Fixes

* remove build info from root command ([bcc861f](https://github.com/everdrone/grab/commit/bcc861f8d29eb5b0594bbc7dd226fe29f5d98c0e))

### [0.1.4](https://github.com/everdrone/grab/compare/v0.1.3...v0.1.4) (2022-08-20)


### Bug Fixes

* fix getting the go version from go.mod ([d2a3faf](https://github.com/everdrone/grab/commit/d2a3fafe5d0578ef9ecf1898a313f615d1dcb174))
* fix the default configuration to work out of the box ([a27ebc3](https://github.com/everdrone/grab/commit/a27ebc3565fa8bd7ace20a1fa22a9cc902a7b267))
* fix workflow build script ([d3a96c1](https://github.com/everdrone/grab/commit/d3a96c1411e9688f0faa0f790d1e1aef71f06090))
* use the latest version of actions/setup-go@v3 during build ([ebbbbf4](https://github.com/everdrone/grab/commit/ebbbbf4f8fccd6c366e82cbb48ac07ef8ea4a1ef))

### [0.1.3](https://github.com/everdrone/grab/compare/v0.1.2...v0.1.3) (2022-08-19)

- Fixes a bug when passing files as arguments on windows
- Test coverage increased from 52% to 82%

### 0.1.2 (2022-08-18)

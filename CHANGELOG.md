# Changelog

## [0.4.0](https://github.com/instill-ai/usage-client/compare/v0.3.0-alpha...v0.4.0) (2025-07-04)


### Bug Fixes

* **mod:** update golang.org/x/net module to fix vulnerability issue ([e29b1f0](https://github.com/instill-ai/usage-client/commit/e29b1f01a52e779c8aac5bb74e629e5585deaf39))


### Miscellaneous

* **deps:** bump golang.org/x/net from 0.17.0 to 0.23.0 ([#26](https://github.com/instill-ai/usage-client/issues/26)) ([08e8ef7](https://github.com/instill-ai/usage-client/commit/08e8ef774f38fcb3c03de8a77f5450f55db8dd8a))
* **deps:** bump golang.org/x/net from 0.33.0 to 0.36.0 ([#28](https://github.com/instill-ai/usage-client/issues/28)) ([49504d9](https://github.com/instill-ai/usage-client/commit/49504d982f187756394fffebd7e77c855b394939))
* **deps:** bump golang.org/x/net from 0.36.0 to 0.38.0 ([#29](https://github.com/instill-ai/usage-client/issues/29)) ([1989dcf](https://github.com/instill-ai/usage-client/commit/1989dcf7cfb55857d7525e8edba9179a5ebcb69a))
* **deps:** bump google.golang.org/grpc from 1.56.0 to 1.56.3 ([#19](https://github.com/instill-ai/usage-client/issues/19)) ([128cfc5](https://github.com/instill-ai/usage-client/commit/128cfc5ddca339700a58d73a4979be92630b7efc))
* **deps:** bump google.golang.org/protobuf from 1.30.0 to 1.33.0 ([#25](https://github.com/instill-ai/usage-client/issues/25)) ([4a3a39f](https://github.com/instill-ai/usage-client/commit/4a3a39f2fd614f44e4556afded4f40ab925b7369))
* **protobuf:** update service type ([#31](https://github.com/instill-ai/usage-client/issues/31)) ([89526ea](https://github.com/instill-ai/usage-client/commit/89526ea95ca7f4fa7ec60ac3999a8ad4e40ed0c4))
* release v0.4.0 ([dd3d0a7](https://github.com/instill-ai/usage-client/commit/dd3d0a77c6aa6bc5d9e0c2919cdd4d4ce83cda5d))
* update config.json and formatting ([#30](https://github.com/instill-ai/usage-client/issues/30)) ([0571825](https://github.com/instill-ai/usage-client/commit/057182528ad69416fc3f5b2d5b47b8b3a5308566))

## [0.3.0-alpha](https://github.com/instill-ai/usage-client/compare/v0.2.4-alpha...v0.3.0-alpha) (2024-02-16)


### Features

* adopt latest protobuf package ([#23](https://github.com/instill-ai/usage-client/issues/23)) ([62ace2c](https://github.com/instill-ai/usage-client/commit/62ace2cbc15518c9eb0e95215bb7ff6ba701a1c8))

## [0.2.4-alpha](https://github.com/instill-ai/usage-client/compare/v0.2.3-alpha...v0.2.4-alpha) (2023-06-22)


### Miscellaneous Chores

* release v0.2.4-alpha ([9fc955e](https://github.com/instill-ai/usage-client/commit/9fc955e1cd4dabcc4f2e3af40572d8c6fd7dd4f4))

## [0.2.3-alpha](https://github.com/instill-ai/usage-client/compare/v0.2.2-alpha...v0.2.3-alpha) (2023-04-07)


### Miscellaneous Chores

* release v0.2.3-alpha ([77feabb](https://github.com/instill-ai/usage-client/commit/77feabbb897a22de3030adc2b1d347ad0bc17b06))

## [0.2.2-alpha](https://github.com/instill-ai/usage-client/compare/v0.2.1-alpha...v0.2.2-alpha) (2023-02-20)


### Miscellaneous Chores

* release v0.2.2-alpha ([58be415](https://github.com/instill-ai/usage-client/commit/58be415ba729573e875764b298e81d5ffbd48a80))

## [0.2.1-alpha](https://github.com/instill-ai/usage-client/compare/v0.1.2-alpha...v0.2.1-alpha) (2022-09-13)


### Miscellaneous Chores

* release v0.2.1-alpha ([05d5f63](https://github.com/instill-ai/usage-client/commit/05d5f638f2eb632e18435260bc40ca16cc4e5ace))

## [0.1.2-alpha](https://github.com/instill-ai/usage-client/compare/v0.1.1-alpha...v0.1.2-alpha) (2022-08-21)


### Bug Fixes

* separate init and start reporter ([b863502](https://github.com/instill-ai/usage-client/commit/b8635029d06ce812feaf3e32c2dc3439b4d59540))

## [0.1.1-alpha](https://github.com/instill-ai/usage-client/compare/v0.1.0-alpha...v0.1.1-alpha) (2022-06-26)


### Bug Fixes

* remove normalizing session ([90bd718](https://github.com/instill-ai/usage-client/commit/90bd71834fddb9ad4f3b122da5b12b4f69db8382))

## [0.1.0-alpha](https://github.com/instill-ai/usage-client/compare/v0.0.0-alpha...v0.1.0-alpha) (2022-06-22)


### Features

* add usage package ([5b233d2](https://github.com/instill-ai/usage-client/commit/5b233d2747eb7981167e7f18728052f01eae4055))


### Bug Fixes

* add edition checking ([3130efc](https://github.com/instill-ai/usage-client/commit/3130efc502aba0167eee59b9b27fd2c2c649429b))
* change the report freq to 1h ([4013a37](https://github.com/instill-ai/usage-client/commit/4013a3777d9f3d0e8db1929addde907c90adefb8))
* fix callback function call ([c17424a](https://github.com/instill-ai/usage-client/commit/c17424a565e463c882c4a64a390564840bd863ca))
* refactor package name ([17dcc6c](https://github.com/instill-ai/usage-client/commit/17dcc6ca283e27904f70c6fd7230551c6eb0bdde))
* replace db with backend repository interface ([064008a](https://github.com/instill-ai/usage-client/commit/064008af9dc5c1bbda3df5b0c37efad8f192fd60))

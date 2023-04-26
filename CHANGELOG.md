# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

### [0.0.9](https://github.com/hidromatologia-v2/models/compare/v0.0.8...v0.0.9) (2023-04-26)


### Features

* added default controller creation ([957607e](https://github.com/hidromatologia-v2/models/commit/957607efc954584f9af5097c213d65d065c83148))
* Config environment variables struct ([5fdbf09](https://github.com/hidromatologia-v2/models/commit/5fdbf09d671442f5997f941be886cc23b95c40c7))


### Bug Fixes

* added memphis container to docker setup ([d07d423](https://github.com/hidromatologia-v2/models/commit/d07d423374c23d481541047b12fb95e14f625780))

### 0.0.8 (2023-04-25)

### 0.0.7 (2023-04-24)

### 0.0.6 (2023-04-24)

### 0.0.5 (2023-04-19)

### [0.0.4](https://github.com/hidromatologia-v2/models/compare/v0.0.3...v0.0.4) (2023-04-19)


### Features

* Admin table support ([29c7a4c](https://github.com/hidromatologia-v2/models/commit/29c7a4cc35e75260e5a20b1b4ed52b02c697342c))
* Alert create operation ([16d120a](https://github.com/hidromatologia-v2/models/commit/16d120a26986775b0f90b1426347bb5a49979a88))
* Alert QueryOne controller ([debf62c](https://github.com/hidromatologia-v2/models/commit/debf62c56a5b1bcbb456029f1b4a225c700c5cfe))
* Alerts query many ([c4a544c](https://github.com/hidromatologia-v2/models/commit/c4a544cc123dcce97832e93a3cc28c3835b73c46))
* Authorize API keys of stations ([5746fba](https://github.com/hidromatologia-v2/models/commit/5746fba86fc23f1f66bd240c1d2f2187053d2100))
* authorize user ([6b3b586](https://github.com/hidromatologia-v2/models/commit/6b3b58675b687e67976c1938ef8aecd93c2b2773))
* Create station and Add/Delete sensors ([fbc482f](https://github.com/hidromatologia-v2/models/commit/fbc482f607438dc3a847406714a81dad7230589a))
* created Controller ([a89de99](https://github.com/hidromatologia-v2/models/commit/a89de9927c9ced7e3a8f40f65bdda8d416e30e05))
* delete alerts operation ([7775594](https://github.com/hidromatologia-v2/models/commit/7775594bd598fd1643e18c8f9c85355417539dee))
* delete station controller ([5c29d1c](https://github.com/hidromatologia-v2/models/commit/5c29d1c4fb713d249046505d207b96921c49fbbf))
* numeric random operations ([f89f7f8](https://github.com/hidromatologia-v2/models/commit/f89f7f8ef91043b887d02ae14c9af4150cc64711))
* Query Account details ([07b2c67](https://github.com/hidromatologia-v2/models/commit/07b2c673247ffa0430b24c12efc63c8438c6fb09))
* Query many stations controller operation ([f633eed](https://github.com/hidromatologia-v2/models/commit/f633eedfdfc86fd7632044557ebc117760f4ecca))
* query one station controller operation ([375335c](https://github.com/hidromatologia-v2/models/commit/375335c9e9f6cf60571e413892fdbd635b846dc0))
* query sensor historical data by UUID ([2811350](https://github.com/hidromatologia-v2/models/commit/28113506a48218fa1da7c19e98eecc3d3b8bf48f))
* registration controller ([62314ed](https://github.com/hidromatologia-v2/models/commit/62314ed7b0fd06e4832247f30f0752f6be415f29))
* Update Alert controller operation ([9d24b25](https://github.com/hidromatologia-v2/models/commit/9d24b254320c9bbadbc4d76c78b252c70f63cd16))
* Update alerts system ([ecf764f](https://github.com/hidromatologia-v2/models/commit/ecf764f54aba956600ae1c6cb907f61a1bec6908))
* Update station controller operation ([f389d9c](https://github.com/hidromatologia-v2/models/commit/f389d9c69e6ca9f251d7759e78ebdf6da3dca4ab))
* user authentication ([1990a6d](https://github.com/hidromatologia-v2/models/commit/1990a6d04281efe7c725df2d2e46fa34b7f04d04))
* User details update ([8a875e8](https://github.com/hidromatologia-v2/models/commit/8a875e8de5a74fca5e2a5a4d33d734f5ff426dbb))


### Bug Fixes

* added APIKey column to stations ([c7e61f2](https://github.com/hidromatologia-v2/models/commit/c7e61f2da1e0c48edc450db9cb128d622c174664))
* added enabled to alert model ([17f5112](https://github.com/hidromatologia-v2/models/commit/17f5112a5bbacd14a47c43c3d13e01d5c10be768))
* added missing Controller close ([379fa15](https://github.com/hidromatologia-v2/models/commit/379fa152eb85f4e748f30ea0123ad8945a21e3a4))
* added on cascade on relation arrays ([ecf6d94](https://github.com/hidromatologia-v2/models/commit/ecf6d94e2a0a1f64d4d9f89bb32d7c1d6bd1760d))
* added support for checking if the account was at the end able to query a specific alert ([7468a4e](https://github.com/hidromatologia-v2/models/commit/7468a4e87c8ef8b1055799ec50fc6f5bc0790e66))
* added support for max int64 ([179ba7a](https://github.com/hidromatologia-v2/models/commit/179ba7a5e78a41bffa345e8b28714ee357e2acb1))
* after updating phone or email account should be unconfirmed again ([d96e289](https://github.com/hidromatologia-v2/models/commit/d96e2890bd11affd40900255bcdeded2e539adb4))
* CI/CD with my correct email address ([7d7d902](https://github.com/hidromatologia-v2/models/commit/7d7d90284ea38f72c012e6752b29fc9f4d85a8d8))
* comparing pointed value instead of address ([84ac972](https://github.com/hidromatologia-v2/models/commit/84ac9728982594616141d5213f228bca06f7d30c))
* ensure at least one station is deleted ([e6e2809](https://github.com/hidromatologia-v2/models/commit/e6e280941c720c6c245767184a00e35db185ae4f))
* fixed typo in users struct ([d96b1bb](https://github.com/hidromatologia-v2/models/commit/d96b1bb47603f33e2c401e68eacdeabfb08cacef))
* organized it accord to the repo docs/Database.md ([608971a](https://github.com/hidromatologia-v2/models/commit/608971a1be8d370592fdd8d2dd0f53757e7e720c))
* removed duplicate condition field in alerts struct ([51d5117](https://github.com/hidromatologia-v2/models/commit/51d51178e89ffa7d437b3dafc7faeab8f6c51792))
* updated tables to be compatible with gorm zero value UPDATE option ([947c0d1](https://github.com/hidromatologia-v2/models/commit/947c0d1631fb217e1b8076569ec39c05ac31ed95))

### 0.0.3 (2023-04-16)


### Bug Fixes

* added coverage TOKEN ([033a5bc](https://github.com/hidromatologia-v2/models/commit/033a5bc86018d1c9739b72fac32a3d2ebfa9c41a))

### 0.0.2 (2023-04-16)


### Bug Fixes

* removed yaml from docker compose up command ([2d1cea7](https://github.com/hidromatologia-v2/models/commit/2d1cea7e07c274cd31e42641b5e775d4ebf4463c))

### 0.0.1 (2023-04-12)


### Bug Fixes

* removed DEV_OPS token ([aed78e6](https://github.com/hidromatologia-v2/models/commit/aed78e65924676a50506cb307f37d86c328932ff))
* Updated CI/CD pipelines ([4b90c04](https://github.com/hidromatologia-v2/models/commit/4b90c04e91cc686635c8a7d5d6b86dbfdea17f40))
* Updated go.mod with the information of this repository ([adf2c4f](https://github.com/hidromatologia-v2/models/commit/adf2c4fcd3a645413b0e5290f7ac1caa41eb1d45))

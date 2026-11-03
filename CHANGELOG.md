# Changelog

## [0.11.0](https://github.com/andersmmg/bytewing/compare/v0.10.0...v0.11.0) (2026-11-03)


### ⚠ BREAKING CHANGES

* PUID and PGID are removed. Remove them from your compose — they are ignored. Ownership is now adopted from the mounted data directory; if it is owned by root, chown it to a non-root user.

### Features

* account menu in the sidebar footer ([73f693b](https://github.com/andersmmg/bytewing/commit/73f693b31fa36c1f067606c22484be05753f45d8))
* account settings — change password and manage sessions ([5d40ab9](https://github.com/andersmmg/bytewing/commit/5d40ab96f4ee2dd9bba07f23b7914ead71306a50))
* add passkeys support and bump Go to 1.27 ([f10e619](https://github.com/andersmmg/bytewing/commit/f10e6195d18a56c15a4cf6d8219fc392b19299c5))
* adopt the data directory's owner instead of PUID/PGID ([3ef70d1](https://github.com/andersmmg/bytewing/commit/3ef70d1310763a1db398b99baaecf257ea876873))
* API tokens for programmatic access ([892ebdd](https://github.com/andersmmg/bytewing/commit/892ebdd5eb3fffa51c0cde4b7e6a4746de573048))
* audit log of who did what ([242d5a3](https://github.com/andersmmg/bytewing/commit/242d5a3bb7784d671afc6f4f3493a807b7fbd1ca))
* auto adjust thumbnail crops ([4fa0784](https://github.com/andersmmg/bytewing/commit/4fa07848bc23d9c1891b524a729ea16c8bc6d19a))
* back up and restore the database from the CLI ([d4a9797](https://github.com/andersmmg/bytewing/commit/d4a979768cd3240e50a305e68a3ef44a35961a7a))
* background job queue, copying folders and large files ([f40133a](https://github.com/andersmmg/bytewing/commit/f40133a0029a3a63d5e32b70f04b17785ab083b9))
* box selection ([f10e0be](https://github.com/andersmmg/bytewing/commit/f10e0be8e76771048bb4113bffe68588184184ce))
* cached file metadata for loads ([817dc00](https://github.com/andersmmg/bytewing/commit/817dc008bff3f28a07e037dbb0df39f4f6e81d67))
* clear visual for cut files ([01bf6e4](https://github.com/andersmmg/bytewing/commit/01bf6e4b6fe38132ac61c0dc62abfd2bd2136536))
* colors per file/folder entry per-user ([6b9274f](https://github.com/andersmmg/bytewing/commit/6b9274ff3751cc886b501e8dd488c27329a83ff0))
* copy and move in the selection action bar ([eebbeb5](https://github.com/andersmmg/bytewing/commit/eebbeb54623f1a226d2ec04095c5d274428cd778))
* copy, cut and paste files and folders ([a2b33a3](https://github.com/andersmmg/bytewing/commit/a2b33a36cbae9e2b5153e6713be135be79e9af8d))
* download a selection or a folder as a zip ([63f4ca5](https://github.com/andersmmg/bytewing/commit/63f4ca57fffc8a5249b5a7b8a9314a6ede76de00))
* download a URL into storage, guarded against SSRF ([c02c773](https://github.com/andersmmg/bytewing/commit/c02c77318886233f0427603539d754b6c88ada2a))
* drag a file onto a folder to move it ([6cbf180](https://github.com/andersmmg/bytewing/commit/6cbf180bab6f1e33ccef4ff959ebabe6c312ed2d))
* drop a folder onto the page to upload its whole tree ([2ab708c](https://github.com/andersmmg/bytewing/commit/2ab708c4d83f12ef4321eb1100aa228bbc3cc7f9))
* edit permissions inline in the manage-access dialog ([f16dc97](https://github.com/andersmmg/bytewing/commit/f16dc97ba985e7319d2617341892a16399adb35f))
* entry references/links ([1e716cd](https://github.com/andersmmg/bytewing/commit/1e716cdb3ef25ad2cdb33a3d61be0bf80ab814ad))
* extract ZIP archives as a background job ([7db77d2](https://github.com/andersmmg/bytewing/commit/7db77d2a75fb34df35c4cee73b6b386a055d2adb))
* faster uploads, concurrent ([2a3a6ca](https://github.com/andersmmg/bytewing/commit/2a3a6cad6f7231a8b6e76f497ecebe7da0c9107f))
* file info dialog ([ca33af0](https://github.com/andersmmg/bytewing/commit/ca33af0dbbfaf4fbafd916b3cf7aa1c8081d0d2c))
* folder contents counts ([35cb71f](https://github.com/andersmmg/bytewing/commit/35cb71f11fe7bc4d6c09fe5250599f3e45efc4b5))
* grant access and enforce the share permission ([aec2218](https://github.com/andersmmg/bytewing/commit/aec221813634131a6c73f6efb9e10286f6464bf2))
* history navigation, keyboard shortcuts, drag to breadcrumbs ([e89b626](https://github.com/andersmmg/bytewing/commit/e89b626ff84d72803215008fbf34146073dfb409))
* inline rename, auto refresh ([83c617a](https://github.com/andersmmg/bytewing/commit/83c617a82e5bb1ef57bae6ab6f5ef11def49399f))
* interface themes ([748a19c](https://github.com/andersmmg/bytewing/commit/748a19c01efc6fa80ef09b840f200a457de98472))
* invite people to create accounts ([e2e7243](https://github.com/andersmmg/bytewing/commit/e2e7243d8c57fa3ed3411df5f7007b93d9840010))
* load thumbnails from the endpoint, show the upload queue ([54e419b](https://github.com/andersmmg/bytewing/commit/54e419bcb5bd6241b1f360aee9366d787a7afa47))
* make Share the only external link, expiring by default ([015fa27](https://github.com/andersmmg/bytewing/commit/015fa274105c068fa7c6964a4946e04c9d333e3c))
* manage existing accounts from the Members screen ([c1a6c17](https://github.com/andersmmg/bytewing/commit/c1a6c17513ad02b47c86c186da8ced86334c377b))
* mark shared files with a link badge in the grid ([39d5313](https://github.com/andersmmg/bytewing/commit/39d5313afbc73418131019ed1b5962fbc00b9e56))
* more keyboard navigation of files ([1b482e1](https://github.com/andersmmg/bytewing/commit/1b482e1650eeae7ccae84c4206dd41db2d5711a5))
* password option in the share dialog ([986c3cd](https://github.com/andersmmg/bytewing/commit/986c3cd7dc68eef839da9a94167a66544a36bdd2))
* password-protected share links ([ff484b8](https://github.com/andersmmg/bytewing/commit/ff484b8e1f48286321e78ba27efb7b84ad6441ed))
* paste button inside the download URL field ([822c07d](https://github.com/andersmmg/bytewing/commit/822c07d5783dd3aa369b3fefbe2fab4c2fca06a5))
* pause, resume, cancel and show the speed of a download ([c8052cd](https://github.com/andersmmg/bytewing/commit/c8052cd7fe9a83e29ed5be86d5eb37499c499c8c))
* pdf thumbnails ([591f8b5](https://github.com/andersmmg/bytewing/commit/591f8b5740b6066bd0bc6235bad8f20b15c4ce43))
* preview cache ([a4aeef2](https://github.com/andersmmg/bytewing/commit/a4aeef288f3fa9069e702e739d84c21d1c661b6c))
* preview video, audio and text as well as images and PDF ([a04b034](https://github.com/andersmmg/bytewing/commit/a04b034431d02acb220f784eca220b8eeebd369e))
* public folder shares with confined browsing ([57ab3a9](https://github.com/andersmmg/bytewing/commit/57ab3a9361f8b92c551ba0dd3395979273625071))
* public share links for files ([e705232](https://github.com/andersmmg/bytewing/commit/e705232f1c9003e9c46e8298c7fa7d280f9787d8))
* quotas and improved userdirs ([2f6054a](https://github.com/andersmmg/bytewing/commit/2f6054a12ea7f67d84106fe0d47e3258c775d43d))
* recursive file search from the toolbar ([93b0007](https://github.com/andersmmg/bytewing/commit/93b000717a05c3fb76b7415adea2d9537252f306))
* reset a forgotten password with recovery codes ([8fb663d](https://github.com/andersmmg/bytewing/commit/8fb663deb9c4206555d1c34a4d5d6bc558c5b00d))
* retention, checksums, and a uniform settings UI ([42d3ae0](https://github.com/andersmmg/bytewing/commit/42d3ae03127d7c6c40956884516f64390c9ebe6c))
* s3 checksums support ([c8fc9a7](https://github.com/andersmmg/bytewing/commit/c8fc9a76bfd580b24dc88f8b8a261c9aa7b32ac2))
* **s3:** basic s3 backend support ([b6208cb](https://github.com/andersmmg/bytewing/commit/b6208cb1763635679640c4308310f204646f2bed))
* **s3:** working previews on s3 ([1fd8e4e](https://github.com/andersmmg/bytewing/commit/1fd8e4e11fc94a028293d23d125e7bbbdea21a87))
* serve compressed image thumbnails from the server ([de52537](https://github.com/andersmmg/bytewing/commit/de525374cb7ea99afed6471c477592b933396a89))
* server downloads with s3, also off by default ([8588944](https://github.com/andersmmg/bytewing/commit/8588944d3e638a46fc4a9b68fb311d6034ca5ec7))
* share dialog and a Shared section ([2a3d04b](https://github.com/andersmmg/bytewing/commit/2a3d04b85701966b781e62a091ccbe6bdbccdae5))
* share folders from the context menu ([e4e6e81](https://github.com/andersmmg/bytewing/commit/e4e6e81955aa33562d4f61b4192141546792337a))
* show a place, not an IP, for each session ([f293f7c](https://github.com/andersmmg/bytewing/commit/f293f7c2974b2551d3a50ac068448a30db76b718))
* show only the actions the caller is allowed to perform ([abef97e](https://github.com/andersmmg/bytewing/commit/abef97e31cc81785c54b78c999db7c48253d4180))
* show the logo in the sign-in, sidebar and breadcrumb ([d61d62d](https://github.com/andersmmg/bytewing/commit/d61d62dbce3ea303a8cb592143377e6d68bce91d))
* show the share badge in list view too ([1a174d6](https://github.com/andersmmg/bytewing/commit/1a174d6a3cfb574ac968d65ab1e4a569a8cfd6e9))
* snapshot before migrations and report divergences on restore ([e1afd62](https://github.com/andersmmg/bytewing/commit/e1afd62a3538f72927e63b185d8852bab45263dc))
* thumbnail generation ([9c60fea](https://github.com/andersmmg/bytewing/commit/9c60fea0654a3fb6a92f81b14b416667b05f2690))
* two-factor authentication (TOTP) ([eaa9f88](https://github.com/andersmmg/bytewing/commit/eaa9f8861292dbfc7122d4174d3563f6bd0aa36b))
* unify creating with a New menu and empty-area right-click ([c06d396](https://github.com/andersmmg/bytewing/commit/c06d39669f5e30e27c307ed09bb497c6444c5ec5))
* upload cancel button ([f966435](https://github.com/andersmmg/bytewing/commit/f9664351cfcd22186faee5fe35816d5cce695462))
* user groups for granting access to a team at once ([2ec1e81](https://github.com/andersmmg/bytewing/commit/2ec1e812545ad77893a5d02d9cf4fc146be64c86))


### Bug Fixes

* allow previewing pdf when using s3 ([5ef1194](https://github.com/andersmmg/bytewing/commit/5ef11944df9017e3e6e748ad9b49a83fe8152df9))
* box select improvements ([c3213fc](https://github.com/andersmmg/bytewing/commit/c3213fcd42a49284a9c9c0ddbecbb1ed8b37ae59))
* clarify move to trash action ([1c1f1a0](https://github.com/andersmmg/bytewing/commit/1c1f1a08a7bbe18ed6e90f949fc7700197961ffa))
* clear multi-file drags ([8bb567e](https://github.com/andersmmg/bytewing/commit/8bb567efe57baf2e52d51ea0fe52344b1991a121))
* don't truncate list rows with the function ([9da6fbe](https://github.com/andersmmg/bytewing/commit/9da6fbed48b96b247c29408e14843521f068b784))
* Download saves the file instead of opening it inline ([0e0a2d6](https://github.com/andersmmg/bytewing/commit/0e0a2d6fa03f35a63c2a9f2b3dbdedb3d4972190))
* gosec checks now pass correctly ([#2](https://github.com/andersmmg/bytewing/issues/2)) ([5b9ec1d](https://github.com/andersmmg/bytewing/commit/5b9ec1dc39b455ec562fe5cf1709e1f2b1395eea))
* handle ctrl+a and right clicks ([51976b7](https://github.com/andersmmg/bytewing/commit/51976b74216b8dc1bd6b1df77247c4784de31171))
* larger file icon sizes ([1be2077](https://github.com/andersmmg/bytewing/commit/1be2077e9699cf5e3e8e7be772a38aead6784a72))
* let a granted user reach files through folders above them ([be04664](https://github.com/andersmmg/bytewing/commit/be04664979afd6f90e60a0ba6d7119fcffb5d322))
* let the PDF preview render by allowing scripts in its iframe ([75b2a14](https://github.com/andersmmg/bytewing/commit/75b2a144bca582095c667c9f973727888138b99d))
* pasting clears the clipboard so Paste no longer lingers ([22cf970](https://github.com/andersmmg/bytewing/commit/22cf9706c60c8f1f24f79da7c0fe64582e338f73))
* record the real client IP behind a reverse proxy ([0822320](https://github.com/andersmmg/bytewing/commit/0822320576922eda75b156507d6d2124cd5b98aa))
* renaming a folder carries its permissions and ownership ([8b5a599](https://github.com/andersmmg/bytewing/commit/8b5a599fe6cba9889c1f3f436534c8f2bdc2de1b))
* show a favicon in the browser tab ([94b39a6](https://github.com/andersmmg/bytewing/commit/94b39a61144d9fd82b6e8730f7e4fd501aeee78b))
* show admins full instance total storage ([53de0fa](https://github.com/andersmmg/bytewing/commit/53de0fa1b5bfd69f425e3a07e3d46e429f0d5c91))
* stamp the version from version.txt when none is passed to the build ([119558c](https://github.com/andersmmg/bytewing/commit/119558ccbbeca3dbc8ff0f639f2e247ba855c22b))
* the context-menu Download actually downloads ([41cd63a](https://github.com/andersmmg/bytewing/commit/41cd63a74efade6be019361eddd3f927302f21b0))
* **theme:** deeper background blacks, clearer colors ([0080ab6](https://github.com/andersmmg/bytewing/commit/0080ab620ea24034da9ca7315bcdc2ab21800755))
* thumbnails never loaded because the image was display:none ([74d6ea8](https://github.com/andersmmg/bytewing/commit/74d6ea8fd1d1249ef24214b4f5866281b4a4585d))

## [0.10.0](https://github.com/andersmmg/bytewing/compare/v0.9.0...v0.10.0) (2026-08-10)


### ⚠ BREAKING CHANGES

* PUID and PGID are removed. Remove them from your compose — they are ignored. Ownership is now adopted from the mounted data directory; if it is owned by root, chown it to a non-root user.

### Features

* adopt the data directory's owner instead of PUID/PGID ([3ef70d1](https://github.com/andersmmg/bytewing/commit/3ef70d1310763a1db398b99baaecf257ea876873))


### Bug Fixes

* let the PDF preview render by allowing scripts in its iframe ([75b2a14](https://github.com/andersmmg/bytewing/commit/75b2a144bca582095c667c9f973727888138b99d))
* show a favicon in the browser tab ([94b39a6](https://github.com/andersmmg/bytewing/commit/94b39a61144d9fd82b6e8730f7e4fd501aeee78b))

## [0.9.0](https://github.com/andersmmg/bytewing/compare/v0.8.0...v0.9.0) (2026-08-06)


### Features

* download a URL into storage, guarded against SSRF ([c02c773](https://github.com/andersmmg/bytewing/commit/c02c77318886233f0427603539d754b6c88ada2a))
* paste button inside the download URL field ([822c07d](https://github.com/andersmmg/bytewing/commit/822c07d5783dd3aa369b3fefbe2fab4c2fca06a5))
* pause, resume, cancel and show the speed of a download ([c8052cd](https://github.com/andersmmg/bytewing/commit/c8052cd7fe9a83e29ed5be86d5eb37499c499c8c))

## [0.8.0](https://github.com/andersmmg/bytewing/compare/v0.7.0...v0.8.0) (2026-08-05)


### Features

* extract ZIP archives as a background job ([7db77d2](https://github.com/andersmmg/bytewing/commit/7db77d2a75fb34df35c4cee73b6b386a055d2adb))

## [0.7.0](https://github.com/andersmmg/bytewing/compare/v0.6.0...v0.7.0) (2026-08-04)


### Features

* back up and restore the database from the CLI ([d4a9797](https://github.com/andersmmg/bytewing/commit/d4a979768cd3240e50a305e68a3ef44a35961a7a))
* reset a forgotten password with recovery codes ([8fb663d](https://github.com/andersmmg/bytewing/commit/8fb663deb9c4206555d1c34a4d5d6bc558c5b00d))
* retention, checksums, and a uniform settings UI ([42d3ae0](https://github.com/andersmmg/bytewing/commit/42d3ae03127d7c6c40956884516f64390c9ebe6c))
* show a place, not an IP, for each session ([f293f7c](https://github.com/andersmmg/bytewing/commit/f293f7c2974b2551d3a50ac068448a30db76b718))
* show the logo in the sign-in, sidebar and breadcrumb ([d61d62d](https://github.com/andersmmg/bytewing/commit/d61d62dbce3ea303a8cb592143377e6d68bce91d))
* snapshot before migrations and report divergences on restore ([e1afd62](https://github.com/andersmmg/bytewing/commit/e1afd62a3538f72927e63b185d8852bab45263dc))
* two-factor authentication (TOTP) ([eaa9f88](https://github.com/andersmmg/bytewing/commit/eaa9f8861292dbfc7122d4174d3563f6bd0aa36b))


### Bug Fixes

* record the real client IP behind a reverse proxy ([0822320](https://github.com/andersmmg/bytewing/commit/0822320576922eda75b156507d6d2124cd5b98aa))

## [0.6.0](https://github.com/andersmmg/bytewing/compare/v0.5.0...v0.6.0) (2026-08-02)


### Features

* account menu in the sidebar footer ([73f693b](https://github.com/andersmmg/bytewing/commit/73f693b31fa36c1f067606c22484be05753f45d8))
* account settings — change password and manage sessions ([5d40ab9](https://github.com/andersmmg/bytewing/commit/5d40ab96f4ee2dd9bba07f23b7914ead71306a50))
* API tokens for programmatic access ([892ebdd](https://github.com/andersmmg/bytewing/commit/892ebdd5eb3fffa51c0cde4b7e6a4746de573048))
* audit log of who did what ([242d5a3](https://github.com/andersmmg/bytewing/commit/242d5a3bb7784d671afc6f4f3493a807b7fbd1ca))
* download a selection or a folder as a zip ([63f4ca5](https://github.com/andersmmg/bytewing/commit/63f4ca57fffc8a5249b5a7b8a9314a6ede76de00))
* preview video, audio and text as well as images and PDF ([a04b034](https://github.com/andersmmg/bytewing/commit/a04b034431d02acb220f784eca220b8eeebd369e))
* user groups for granting access to a team at once ([2ec1e81](https://github.com/andersmmg/bytewing/commit/2ec1e812545ad77893a5d02d9cf4fc146be64c86))


### Bug Fixes

* Download saves the file instead of opening it inline ([0e0a2d6](https://github.com/andersmmg/bytewing/commit/0e0a2d6fa03f35a63c2a9f2b3dbdedb3d4972190))
* stamp the version from version.txt when none is passed to the build ([119558c](https://github.com/andersmmg/bytewing/commit/119558ccbbeca3dbc8ff0f639f2e247ba855c22b))
* the context-menu Download actually downloads ([41cd63a](https://github.com/andersmmg/bytewing/commit/41cd63a74efade6be019361eddd3f927302f21b0))

## [0.5.0](https://github.com/andersmmg/bytewing/compare/v0.4.0...v0.5.0) (2026-07-31)


### Features

* copy and move in the selection action bar ([eebbeb5](https://github.com/andersmmg/bytewing/commit/eebbeb54623f1a226d2ec04095c5d274428cd778))
* edit permissions inline in the manage-access dialog ([f16dc97](https://github.com/andersmmg/bytewing/commit/f16dc97ba985e7319d2617341892a16399adb35f))
* grant access and enforce the share permission ([aec2218](https://github.com/andersmmg/bytewing/commit/aec221813634131a6c73f6efb9e10286f6464bf2))
* invite people to create accounts ([e2e7243](https://github.com/andersmmg/bytewing/commit/e2e7243d8c57fa3ed3411df5f7007b93d9840010))
* manage existing accounts from the Members screen ([c1a6c17](https://github.com/andersmmg/bytewing/commit/c1a6c17513ad02b47c86c186da8ced86334c377b))
* show only the actions the caller is allowed to perform ([abef97e](https://github.com/andersmmg/bytewing/commit/abef97e31cc81785c54b78c999db7c48253d4180))


### Bug Fixes

* let a granted user reach files through folders above them ([be04664](https://github.com/andersmmg/bytewing/commit/be04664979afd6f90e60a0ba6d7119fcffb5d322))
* pasting clears the clipboard so Paste no longer lingers ([22cf970](https://github.com/andersmmg/bytewing/commit/22cf9706c60c8f1f24f79da7c0fe64582e338f73))
* renaming a folder carries its permissions and ownership ([8b5a599](https://github.com/andersmmg/bytewing/commit/8b5a599fe6cba9889c1f3f436534c8f2bdc2de1b))

## [0.4.0](https://github.com/andersmmg/bytewing/compare/v0.3.0...v0.4.0) (2026-07-31)


### Features

* background job queue, copying folders and large files ([f40133a](https://github.com/andersmmg/bytewing/commit/f40133a0029a3a63d5e32b70f04b17785ab083b9))
* copy, cut and paste files and folders ([a2b33a3](https://github.com/andersmmg/bytewing/commit/a2b33a36cbae9e2b5153e6713be135be79e9af8d))
* drag a file onto a folder to move it ([6cbf180](https://github.com/andersmmg/bytewing/commit/6cbf180bab6f1e33ccef4ff959ebabe6c312ed2d))
* drop a folder onto the page to upload its whole tree ([2ab708c](https://github.com/andersmmg/bytewing/commit/2ab708c4d83f12ef4321eb1100aa228bbc3cc7f9))
* make Share the only external link, expiring by default ([015fa27](https://github.com/andersmmg/bytewing/commit/015fa274105c068fa7c6964a4946e04c9d333e3c))
* mark shared files with a link badge in the grid ([39d5313](https://github.com/andersmmg/bytewing/commit/39d5313afbc73418131019ed1b5962fbc00b9e56))
* password option in the share dialog ([986c3cd](https://github.com/andersmmg/bytewing/commit/986c3cd7dc68eef839da9a94167a66544a36bdd2))
* password-protected share links ([ff484b8](https://github.com/andersmmg/bytewing/commit/ff484b8e1f48286321e78ba27efb7b84ad6441ed))
* public folder shares with confined browsing ([57ab3a9](https://github.com/andersmmg/bytewing/commit/57ab3a9361f8b92c551ba0dd3395979273625071))
* public share links for files ([e705232](https://github.com/andersmmg/bytewing/commit/e705232f1c9003e9c46e8298c7fa7d280f9787d8))
* recursive file search from the toolbar ([93b0007](https://github.com/andersmmg/bytewing/commit/93b000717a05c3fb76b7415adea2d9537252f306))
* share dialog and a Shared section ([2a3d04b](https://github.com/andersmmg/bytewing/commit/2a3d04b85701966b781e62a091ccbe6bdbccdae5))
* share folders from the context menu ([e4e6e81](https://github.com/andersmmg/bytewing/commit/e4e6e81955aa33562d4f61b4192141546792337a))
* show the share badge in list view too ([1a174d6](https://github.com/andersmmg/bytewing/commit/1a174d6a3cfb574ac968d65ab1e4a569a8cfd6e9))
* unify creating with a New menu and empty-area right-click ([c06d396](https://github.com/andersmmg/bytewing/commit/c06d39669f5e30e27c307ed09bb497c6444c5ec5))

## [0.3.0](https://github.com/andersmmg/bytewing/compare/v0.2.0...v0.3.0) (2026-07-30)


### Features

* load thumbnails from the endpoint, show the upload queue ([54e419b](https://github.com/andersmmg/bytewing/commit/54e419bcb5bd6241b1f360aee9366d787a7afa47))
* serve compressed image thumbnails from the server ([de52537](https://github.com/andersmmg/bytewing/commit/de525374cb7ea99afed6471c477592b933396a89))


### Bug Fixes

* thumbnails never loaded because the image was display:none ([74d6ea8](https://github.com/andersmmg/bytewing/commit/74d6ea8fd1d1249ef24214b4f5866281b4a4585d))

# todo add check to verify sbot is installed

release.patch:
	sbot release version -m patch
	sbot push version
	sbot get version

release.minor:
	sbot release version -m minor
	sbot push version
	sbot get version

release.major:
	sbot release version -m major
	sbot push version
	sbot get version
# Node exporter Ubuntu Snap Generator

- You need to setup a Snapcraft build environment to run this script. (Currently a Ubuntu 16.04 build environment)
- See: https://docs.snapcraft.io/build-snaps/get-started-snapcraft

## Login to Snapcraft first
```
snapcraft login
```
## Run the build script, set VERSION as required
```
export VERSION="0.16.0"
./build-snap.sh ${VERSION}
```

## Upload the new Snap to the Ubuntu Snap Store
```
snapcraft push prometheus-node-exporter${VERSION}_amd64.snap
```

## Make the Snap available for the general public, Note: <snap_version> is initially 1
```
snapcraft release prometheus-node-exporter <snap_version> stable
```

## Snap licensing processes have not been fully formalized yet

- See: https://forum.snapcraft.io/t/snap-license-metadata/856/54


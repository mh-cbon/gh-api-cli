PREBUMP=
  666 git fetch --tags origin master
  666 git pull origin master

PREVERSION=
  philea -s "666 go vet %s" "666 go-fmt-fail %s"
  666 go run main.go -v
  666 changelog finalize --version !newversion!
  666 commit -q -m "changelog: !newversion!" -f change.log
  666 go install --ldflags "-X main.VERSION=!newversion!"
  emd gen -out README.md
  666 commit -q -m "README: !newversion!" -f README.md
  666 changelog md -o CHANGELOG.md --vars='{"name":"gh-api-cli"}'
  666 commit -q -m "changelog: !newversion!" -f CHANGELOG.md
  666 build-them-all build main.go -o "build/&os-&arch/&pkg" --os darwin --ldflags "-X main.VERSION=!newversion!"

POSTVERSION=
  666 git push
  666 git push --tags
  666 gh-api-cli create-release -n release -o mh-cbon -r gh-api-cli \
    --ver !newversion! -c "changelog ghrelease --version !newversion!" \
    --draft !isprerelease!
  666 go install --ldflags "-X main.VERSION=!newversion!"
  philea -s -S -p "build/*/**" "666 archive create -f -o=assets/%dname.tar.gz -C=build/%dname/ ."
  666 gh-api-cli create-release -n release --guess \
    --ver !newversion! -c "changelog ghrelease --version !newversion!" \
    --draft !isprerelease!
  666 gh-api-cli upload-release-asset -n release --glob "assets/*" --guess --ver !newversion!
  666 rm-glob -r build/
  666 rm-glob -r assets/

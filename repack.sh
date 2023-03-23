#!/bin/bash

repack() {
    # replace
    rm -rf ./.repack
    dpkg-deb -x $1 ./.repack
    dpkg-deb -e $1 ./.repack/DEBIAN
    chmod -R 755 ./.repack
    if [ -f ./.repack/DEBIAN/control ]; then
        sed -i 's/^Package: .*\.wxhbts/Package: me.otokaze/' ./.repack/DEBIAN/control
        sed -i "/^SileoDepiction.*$/Id" ./.repack/DEBIAN/control
        sed -i "/^Depiction.*$/Id" ./.repack/DEBIAN/control
        pkg=$(grep -i "^Package: " ./.repack/DEBIAN/control | sed -e "s/Package: //i")
        #ver=$(grep -i "^Version: " ./.repack/DEBIAN/control | sed -e "s/Version: //i")
        #ach=$(grep -i "^Architecture: " ./.repack/DEBIAN/control | sed -e "s/Architecture: //i")
        echo "Depiction: https://apt.otokaze.me/web/detail.html?pkg=$pkg" >> ./.repack/DEBIAN/control
        echo "SileoDepiction: https://apt.otokaze.me/debs/$pkg.deb.json" >> ./.repack/DEBIAN/control
    fi
    # backup
    mv $1 $1.old
    # repack
    find ./.repack -name '.DS_Store' -type f -delete
    dpkg-deb -b ./.repack $(dirname $1)/$pkg.deb
}

path=$1
if [ "$path" = "" ]; then
    path=./
elif [ -f $1 ]; then
    repack $1
    exit $?
fi

for file in $(find $path -name "*.deb")
do
    repack $file
done

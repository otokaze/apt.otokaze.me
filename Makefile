.PHONY: clean

update:
	@dpkg-scanpackages -m ./debs > Packages
	@cp Packages Packages.1
	@if [ -f Packages.bz2 ]; \
	then \
	    rm -rf Packages.bz2; \
	fi
	@bzip2 Packages
	@mv Packages.1 Packages
	@gpg --armor --detach-sign --sign -o Release.gpg Release
	@gpg --clearsign -o InRelease Release

clean:
	@find ./debs -name "*.old" -delete
	@rm -rf ./.repack
BINS = archive2disk cexec grub2disk kexec oci2disk qemuimg2disk rootio slurp syslinux writefile

include rules.mk

all: $(BINS)

.PHONY: $(BINS)
$(BINS):
	docker buildx build --platform linux/amd64 --load -t  $@:latest -f ./$@/Dockerfile .

formatters: $(toolBins)
	git ls-files '*.go' | xargs -I% sh -c 'sed -i "/^import (/,/^)/ { /^\s*$$/ d }" % && bin/gofumpt -w %'
	git ls-files '*.go' | xargs -I% bin/goimports -w %

tidy-all:
	for d in $(BINS); do (cd $$d; go mod tidy); done

include lint.mk

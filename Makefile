# Don't touch this variable
export GO111MODULE := on

files := main.go
os := $(shell go env GOOS)
arch := $(shell go env GOARCH)

# Output directory for binaries
output_dir := build

# Don't touch this definition
define newline


endef

all: format lint compile

format:
	gofmt -d -w $(files)

lint:
	$(foreach file,$(files),golangci-lint run --fix $(file)$(newline))

compile: normal

clean:
	rm -rf $(output_dir)

normal: ext = $(if $(findstring windows,$(os)),.exe)
normal:
	$(foreach file,$(files),\
		$(eval out := $(output_dir)/$(file:.go=)_$(os)-$(arch)$(ext))\
			GOOS=$(os) GOARCH=$(arch) go build -o $(out) $(file)$(newline))
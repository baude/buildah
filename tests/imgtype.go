package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	is "github.com/containers/image/storage"
	"github.com/containers/image/types"
	"github.com/containers/storage"
	"github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/projectatomic/buildah"
	"github.com/projectatomic/buildah/docker"
)

func main() {
	expectedManifestType := ""
	expectedConfigType := ""

	storeOptions := storage.DefaultStoreOptions
	root := flag.String("root", storeOptions.GraphRoot, "storage root directory")
	runroot := flag.String("runroot", storeOptions.RunRoot, "storage runtime directory")
	driver := flag.String("storage-driver", storeOptions.GraphDriverName, "storage driver")
	opts := flag.String("storage-opts", "", "storage option list (comma separated)")
	policy := flag.String("signature-policy", "", "signature policy file")
	mtype := flag.String("expected-manifest-type", buildah.OCIv1ImageManifest, "expected manifest type")
	showm := flag.Bool("show-manifest", false, "output the manifest JSON")
	showc := flag.Bool("show-config", false, "output the configuration JSON")
	flag.Parse()
	switch *mtype {
	case buildah.OCIv1ImageManifest:
		expectedManifestType = *mtype
		expectedConfigType = v1.MediaTypeImageConfig
	case buildah.Dockerv2ImageManifest:
		expectedManifestType = *mtype
		expectedConfigType = docker.V2S2MediaTypeImageConfig
	case "*":
		expectedManifestType = ""
		expectedConfigType = ""
	default:
		logrus.Fatalf("unknown -expected-manifest-type value, expected either %q or %q or %q",
			buildah.OCIv1ImageManifest, buildah.Dockerv2ImageManifest, "*")
	}
	if root != nil {
		storeOptions.GraphRoot = *root
	}
	if runroot != nil {
		storeOptions.RunRoot = *runroot
	}
	if driver != nil {
		storeOptions.GraphDriverName = *driver
	}
	if opts != nil && *opts != "" {
		storeOptions.GraphDriverOptions = strings.Split(*opts, ",")
	}
	systemContext := &types.SystemContext{
		SignaturePolicyPath: *policy,
	}
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return
	}
	store, err := storage.GetStore(storeOptions)
	if err != nil {
		logrus.Fatalf("error opening storage: %v", err)
	}
	defer store.Shutdown(false)

	for _, image := range args {
		oImage := v1.Image{}
		dImage := docker.V2Image{}
		oManifest := v1.Manifest{}
		dManifest := docker.V2S2Manifest{}
		manifestType := ""
		configType := ""

		ref, err := is.Transport.ParseStoreReference(store, image)
		if err != nil {
			logrus.Fatalf("error parsing reference %q: %v", image, err)
		}

		src, err := ref.NewImageSource(systemContext, []string{expectedManifestType})
		if err != nil {
			logrus.Fatalf("error opening source image %q: %v", image, err)
		}
		defer src.Close()

		manifest, manifestType, err := src.GetManifest()
		if err != nil {
			logrus.Fatalf("error reading manifest from %q: %v", image, err)
		}

		img, err := ref.NewImage(systemContext)
		if err != nil {
			logrus.Fatalf("error opening image %q: %v", image, err)
		}
		defer img.Close()

		config, err := img.ConfigBlob()
		if err != nil {
			logrus.Fatalf("error reading configuration from %q: %v", image, err)
		}

		switch expectedManifestType {
		case buildah.OCIv1ImageManifest:
			err = json.Unmarshal(manifest, &oManifest)
			if err != nil {
				logrus.Fatalf("error parsing manifest from %q: %v", image, err)
			}
			err = json.Unmarshal(config, &oImage)
			if err != nil {
				logrus.Fatalf("error parsing config from %q: %v", image, err)
			}
			manifestType = v1.MediaTypeImageManifest
			configType = oManifest.Config.MediaType
		case buildah.Dockerv2ImageManifest:
			err = json.Unmarshal(manifest, &dManifest)
			if err != nil {
				logrus.Fatalf("error parsing manifest from %q: %v", image, err)
			}
			err = json.Unmarshal(config, &dImage)
			if err != nil {
				logrus.Fatalf("error parsing config from %q: %v", image, err)
			}
			manifestType = dManifest.MediaType
			configType = dManifest.Config.MediaType
		}
		if expectedManifestType != "" && manifestType != expectedManifestType {
			logrus.Fatalf("expected manifest type %q in %q, got %q", expectedManifestType, image, manifestType)
		}
		if expectedConfigType != "" && configType != expectedConfigType {
			logrus.Fatalf("expected config type %q in %q, got %q", expectedConfigType, image, configType)
		}
		if showm != nil && *showm {
			fmt.Println(string(manifest))
		}
		if showc != nil && *showc {
			fmt.Println(string(config))
		}
	}
}

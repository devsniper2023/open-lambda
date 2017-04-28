package cache

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/open-lambda/open-lambda/worker/config"
	"github.com/open-lambda/open-lambda/worker/dockerutil"

	docker "github.com/fsouza/go-dockerclient"
	sb "github.com/open-lambda/open-lambda/worker/sandbox"
)

func InitCacheFactory(opts *config.Config, cluster string) (cf *BufferedCacheFactory, root *sb.DockerSandbox, rootDir, rootCID string, err error) {
	cf, root, rootDir, rootCID, err = NewBufferedCacheFactory(opts, cluster)
	if err != nil {
		return nil, nil, "", "", err
	}

	return cf, root, rootDir, rootCID, nil
}

// CacheFactory is a SandboxFactory that creates docker sandboxes for the cache.
type CacheFactory struct {
	client  *docker.Client
	cmd     []string
	caps    []string
	labels  map[string]string
	pkgsDir string
}

// emptySBInfo wraps sandbox information necessary for the buffer.
type emptySBInfo struct {
	sandbox    *sb.DockerSandbox
	sandboxDir string
}

// BufferedCacheFactory maintains a buffer of sandboxes created by another factory.
type BufferedCacheFactory struct {
	delegate *CacheFactory
	buffer   chan *emptySBInfo
	errors   chan error
	dir      string
}

// NewCacheFactory creates a CacheFactory.
func NewCacheFactory(cluster, pkgsDir string) (*CacheFactory, error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, err
	}

	cmd := []string{"/init"}

	caps := []string{"SYS_ADMIN"}

	labels := map[string]string{
		dockerutil.DOCKER_LABEL_CLUSTER: cluster,
		dockerutil.DOCKER_LABEL_TYPE:    dockerutil.POOL,
	}

	cf := &CacheFactory{client, cmd, caps, labels, pkgsDir}
	return cf, nil
}

// Create creates a docker sandbox from the pool directory.
func (cf *CacheFactory) Create(sandboxDir string, cmd []string) (*sb.DockerSandbox, string, error) {
	volumes := []string{
		fmt.Sprintf("%s:%s", sandboxDir, "/host"),
		fmt.Sprintf("%s:%s:ro", cf.pkgsDir, "/packages"),
	}

	container, err := cf.client.CreateContainer(
		docker.CreateContainerOptions{
			Config: &docker.Config{
				Image:  dockerutil.CACHE_IMAGE,
				Labels: cf.labels,
				Cmd:    cmd,
			},
			HostConfig: &docker.HostConfig{
				Binds:      volumes,
				PidMode:    "host",
				CapAdd:     cf.caps,
				AutoRemove: true,
			},
		},
	)
	if err != nil {
		return nil, "", err
	}

	sandbox := sb.NewDockerSandbox(sandboxDir, "", "", container, cf.client)
	return sandbox, container.ID, nil
}

// NewBufferedCacheFactory creates a BufferedCacheFactory and starts a go routine to
// fill the sandbox buffer.
func NewBufferedCacheFactory(opts *config.Config, cluster string) (*BufferedCacheFactory, *sb.DockerSandbox, string, string, error) {
	cacheDir := opts.Import_cache_dir
	pkgsDir := opts.Pkgs_dir
	buffer := opts.Import_cache_buffer
	indexHost := opts.Index_host
	indexPort := opts.Index_port

	rootCmd := []string{"python", "server.py"}
	if indexHost != "" && indexPort != "" {
		rootCmd = append(rootCmd, indexHost, indexPort)
	}

	delegate, err := NewCacheFactory(cluster, pkgsDir)
	if err != nil {
		return nil, nil, "", "", err
	}

	bf := &BufferedCacheFactory{
		delegate: delegate,
		buffer:   make(chan *emptySBInfo, buffer),
		errors:   make(chan error, buffer),
		dir:      cacheDir,
	}

	if err := os.MkdirAll(cacheDir, os.ModeDir); err != nil {
		return nil, nil, "", "", fmt.Errorf("failed to create pool directory at %s: %v", cacheDir, err)
	}

	// create the root container
	rootDir := filepath.Join(bf.dir, "root")
	if err := os.MkdirAll(rootDir, os.ModeDir); err != nil {
		return nil, nil, "", "", fmt.Errorf("failed to create cache entry directory at %s: %v", cacheDir, err)
	}

	root, rootCID, err := bf.delegate.Create(rootDir, rootCmd)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("failed to create cache entry sandbox: %v", err)
	} else if err := root.Start(); err != nil {
		return nil, nil, "", "", fmt.Errorf("failed to start cache entry sandbox: %v", err)
	}

	// fill the sandbox buffer
	var shared_idx int64 = 0
	for i := 0; i < 5; i++ {
		go func(idxptr *int64) {
			for {
				sandboxDir := filepath.Join(bf.dir, fmt.Sprintf("%d", atomic.AddInt64(idxptr, 1)))
				if err := os.MkdirAll(sandboxDir, os.ModeDir); err != nil {
					bf.buffer <- nil
					bf.errors <- err
				} else if sandbox, _, err := bf.delegate.Create(sandboxDir, []string{"/init"}); err != nil {
					bf.buffer <- nil
					bf.errors <- err
				} else if err := sandbox.Start(); err != nil {
					bf.buffer <- nil
					bf.errors <- err
				} else if err := sandbox.Pause(); err != nil {
					bf.buffer <- nil
					bf.errors <- err
				} else {
					bf.buffer <- &emptySBInfo{sandbox, sandboxDir}
					bf.errors <- nil
				}
			}
		}(&shared_idx)
	}

	log.Printf("filling cache buffer")
	for len(bf.buffer) < cap(bf.buffer) {
		time.Sleep(20 * time.Millisecond)
	}
	log.Printf("cache buffer full")

	return bf, root, rootDir, rootCID, nil
}

// Returns a sandbox ready for a cache interpreter
func (bf *BufferedCacheFactory) Create() (*sb.DockerSandbox, string, error) {
	info, err := <-bf.buffer, <-bf.errors
	if err != nil {
		return nil, "", err
	}

	if err := info.sandbox.Unpause(); err != nil {
		return nil, "", err
	}

	return info.sandbox, info.sandboxDir, nil
}

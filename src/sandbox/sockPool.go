package sandbox

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/open-lambda/open-lambda/ol/common"
)

// the first program is executed on the host, which sets up the
// container, running the second program inside the container
const SOCK_HOST_INIT = "/usr/local/bin/sock-init"
const SOCK_GUEST_INIT = "/ol-init"

var nextId int64 = 0

// SOCKPool is a ContainerFactory that creats docker containeres.
type SOCKPool struct {
	name          string
	rootDirs      *common.DirMaker
	cgPool        *CgroupPool
	mem           *MemPool
	eventHandlers []SandboxEventFunc
	debugger
}

// NewSOCKPool creates a SOCKPool.
func NewSOCKPool(name string, mem *MemPool) (cf *SOCKPool, err error) {
	cgPool, err := NewCgroupPool(name)
	if err != nil {
		return nil, err
	}

	rootDirs, err := common.NewDirMaker("root-"+name, true)
	if err != nil {
		return nil, err
	}

	pool := &SOCKPool{
		name:          name,
		mem:           mem,
		cgPool:        cgPool,
		rootDirs:      rootDirs,
		eventHandlers: []SandboxEventFunc{},
	}

	pool.debugger = newDebugger(pool)

	return pool, nil
}

func sbStr(sb Sandbox) string {
	if sb == nil {
		return "<nil>"
	}
	return fmt.Sprintf("<SB %s>", sb.ID())
}

func (pool *SOCKPool) Create(parent Sandbox, isLeaf bool, codeDir, scratchDir string, meta *SandboxMeta) (sb Sandbox, err error) {
	id := fmt.Sprintf("%d", atomic.AddInt64(&nextId, 1))
	meta = fillMetaDefaults(meta)
	pool.printf("<%v>.Create(%v, %v, %v, %v, %v)=%s...", pool.name, sbStr(parent), isLeaf, codeDir, scratchDir, meta, id)
	defer func() {
		pool.printf("...returns %v, %v", sbStr(sb), err)
	}()

	t := common.T0("Create()")
	defer t.T1()

	// block until we have enough to cover the cgroup mem limits
	t2 := t.T0("acquire-mem")
	pool.mem.adjustAvailableMB(-meta.MemLimitMB)
	t2.T1()

	t2 = t.T0("acquire-cgroup")
	// when creating a new Sandbox without a parent, we want to
	// move the cgroup memory charge (otherwise the charge will
	// exist outside any Sandbox).  But when creating a child, we
	// don't want to use this cgroup feature, because the child
	// would take the blame for ALL of the parent's allocations
	moveMemCharge := (parent == nil)
	cg := pool.cgPool.GetCg(meta.MemLimitMB, moveMemCharge)
	t2.T1()

	var cSock *SOCKContainer = &SOCKContainer{
		pool:             pool,
		id:               id,
		containerRootDir: pool.rootDirs.Make("SB-" + id),
		codeDir:          codeDir,
		scratchDir:       scratchDir,
		cg:               cg,
		cgRefCount:       1,
		children:         make(map[string]Sandbox),
		meta:             meta,
	}
	var c Sandbox = cSock

	defer func() {
		if err != nil {
			c.Destroy()
		}
	}()

	// root file system
	if isLeaf && cSock.codeDir == "" {
		return nil, fmt.Errorf("leaf sandboxes must have codeDir set")
	}

	t2 = t.T0("make-root-fs")
	if err := cSock.populateRoot(); err != nil {
		return nil, fmt.Errorf("failed to create root FS: %v", err)
	}
	t2.T1()

	// add installed packages to the path, and import the modules we'll need
	var pyCode []string

	for _, pkg := range meta.Installs {
		pyCode = append(pyCode, "sys.path.append('/packages/"+pkg+"/files')")
	}

	for _, mod := range meta.Imports {
		pyCode = append(pyCode, "import "+mod)
	}

	// handler or Zygote?
	if isLeaf {
		pyCode = append(pyCode, "web_server()")
	} else {
		pyCode = append(pyCode, "fork_server()")
	}

	path := filepath.Join(scratchDir, "bootstrap.py")
	code := []byte(strings.Join(pyCode, "\n"))
	if err := ioutil.WriteFile(path, code, 0600); err != nil {
		return nil, err
	}

	// wrap to make thread-safe and handle container death.
	// after this line, two things happen:
	// 1. listeners (e.g., evictors) become aware of the Sandbox
	// 2. if this function fails, .Destroy() cleanup will be through the safeSandbox layer
	c = newSafeSandbox(c, pool.eventHandlers)

	// create new process in container (fresh, or forked from parent)
	if parent != nil {
		t2 := t.T0("fork-proc")
		if err := parent.fork(c); err != nil {
			pool.printf("parent.fork returned %v", err)
			return nil, FORK_FAILED
		}
		cSock.parent = parent
		t2.T1()
	} else {
		t2 := t.T0("fresh-proc")
		if err := cSock.freshProc(); err != nil {
			return nil, err
		}
		t2.T1()
	}

	return c, nil
}

func (pool *SOCKPool) printf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Printf("%s [SOCK POOL %s]", strings.TrimRight(msg, "\n"), pool.name)
}

// handler(...) will be called everytime a sandbox-related event occurs,
// such as Create, Destroy, etc.
//
// the events are sent after the actions complete
//
// TODO: eventually make this part of SandboxPool API, and support in Docker?
func (pool *SOCKPool) AddListener(handler SandboxEventFunc) {
	pool.eventHandlers = append(pool.eventHandlers, handler)
}

func (pool *SOCKPool) Cleanup() {
	// user is required to kill all containers before they call
	// this.  If they did, the memory pool should be full.
	pool.printf("make sure all memory is free")
	pool.mem.adjustAvailableMB(-pool.mem.totalMB)
	pool.printf("memory pool emptied")

	pool.cgPool.Destroy()
	if err := pool.rootDirs.Cleanup(); err != nil {
		panic(err)
	}
}

func (pool *SOCKPool) DebugString() string {
	return pool.debugger.Dump()
}

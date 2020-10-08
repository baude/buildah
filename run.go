package buildah

import (
	"fmt"
	"io"

	"github.com/containers/buildah/define"
	"github.com/opencontainers/runtime-spec/specs-go"
)

const (
	// runUsingRuntimeCommand is a command we use as a key for reexec
	runUsingRuntimeCommand = Package + "-oci-runtime"
)

// TerminalPolicy takes the value DefaultTerminal, WithoutTerminal, or WithTerminal.
type TerminalPolicy int

const (
	// DefaultTerminal indicates that this Run invocation should be
	// connected to a pseudoterminal if we're connected to a terminal.
	DefaultTerminal TerminalPolicy = iota
	// WithoutTerminal indicates that this Run invocation should NOT be
	// connected to a pseudoterminal.
	WithoutTerminal
	// WithTerminal indicates that this Run invocation should be connected
	// to a pseudoterminal.
	WithTerminal
)

// String converts a TerminalPoliicy into a string.
func (t TerminalPolicy) String() string {
	switch t {
	case DefaultTerminal:
		return "DefaultTerminal"
	case WithoutTerminal:
		return "WithoutTerminal"
	case WithTerminal:
		return "WithTerminal"
	}
	return fmt.Sprintf("unrecognized terminal setting %d", t)
}

type NameSpaceOption = define.NamespaceOption

// RunOptions can be used to alter how a command is run in the container.
type RunOptions struct {
	// Hostname is the hostname we set for the running container.
	Hostname string
	// Isolation is either IsolationDefault, IsolationOCI, IsolationChroot, or IsolationOCIRootless.
	Isolation define.Isolation
	// Runtime is the name of the runtime to run.  It should accept the
	// same arguments that runc does, and produce similar output.
	Runtime string
	// Args adds global arguments for the runtime.
	Args []string
	// NoPivot adds the --no-pivot runtime flag.
	NoPivot bool
	// Mounts are additional mount points which we want to provide.
	Mounts []specs.Mount
	// Env is additional environment variables to set.
	Env []string
	// User is the user as whom to run the command.
	User string
	// WorkingDir is an override for the working directory.
	WorkingDir string
	// Shell is default shell to run in a container.
	Shell string
	// Cmd is an override for the configured default command.
	Cmd []string
	// Entrypoint is an override for the configured entry point.
	Entrypoint []string
	// NamespaceOptions controls how we set up the namespaces for the process.
	NamespaceOptions define.NamespaceOptions
	// ConfigureNetwork controls whether or not network interfaces and
	// routing are configured for a new network namespace (i.e., when not
	// joining another's namespace and not just using the host's
	// namespace), effectively deciding whether or not the process has a
	// usable network.
	ConfigureNetwork define.NetworkConfigurationPolicy
	// CNIPluginPath is the location of CNI plugin helpers, if they should be
	// run from a location other than the default location.
	CNIPluginPath string
	// CNIConfigDir is the location of CNI configuration files, if the files in
	// the default configuration directory shouldn't be used.
	CNIConfigDir string
	// Terminal provides a way to specify whether or not the command should
	// be run with a pseudoterminal.  By default (DefaultTerminal), a
	// terminal is used if os.Stdout is connected to a terminal, but that
	// decision can be overridden by specifying either WithTerminal or
	// WithoutTerminal.
	Terminal TerminalPolicy
	// TerminalSize provides a way to set the number of rows and columns in
	// a pseudo-terminal, if we create one, and Stdin/Stdout/Stderr aren't
	// connected to a terminal.
	TerminalSize *specs.Box
	// The stdin/stdout/stderr descriptors to use.  If set to nil, the
	// corresponding files in the "os" package are used as defaults.
	Stdin  io.Reader `json:"-"`
	Stdout io.Writer `json:"-"`
	Stderr io.Writer `json:"-"`
	// Quiet tells the run to turn off output to stdout.
	Quiet bool
	// AddCapabilities is a list of capabilities to add to the default set.
	AddCapabilities []string
	// DropCapabilities is a list of capabilities to remove from the default set,
	// after processing the AddCapabilities set.  If a capability appears in both
	// lists, it will be dropped.
	DropCapabilities []string
	// Devices are the additional devices to add to the containers
	Devices define.ContainerDevices
}

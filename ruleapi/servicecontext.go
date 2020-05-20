package ruleapi

import (
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/support/trace"
)

// activity init context
type initContext struct {
	settings      map[string]interface{}
	mapperFactory mapper.Factory
	logger        log.Logger
}

func newInitContext(name string, settings map[string]interface{}, l log.Logger) *initContext {
	return &initContext{
		settings:      settings,
		mapperFactory: mapper.NewFactory(resolve.GetBasicResolver()),
		logger:        log.ChildLogger(l, name),
	}
}

func (i *initContext) Settings() map[string]interface{} {
	return i.settings
}

func (i *initContext) MapperFactory() mapper.Factory {
	return i.mapperFactory
}

func (i *initContext) Logger() log.Logger {
	return i.logger
}

// ServiceContext context
type ServiceContext struct {
	TaskName     string
	activityHost activity.Host

	metadata *activity.Metadata
	settings map[string]interface{}
	inputs   map[string]interface{}
	outputs  map[string]interface{}

	shared map[string]interface{}
}

func newServiceContext(md *activity.Metadata) *ServiceContext {
	input := map[string]data.TypedValue{"Input1": data.NewTypedValue(data.TypeString, "")}
	output := map[string]data.TypedValue{"Output1": data.NewTypedValue(data.TypeString, "")}

	// TBD: rule action's details (like: metadata, name, etc) to be used here
	sHost := &ServiceHost{
		HostId:     "1",
		HostRef:    "github.com/project-flogo/rules",
		IoMetadata: &metadata.IOMetadata{Input: input, Output: output},
		HostData:   data.NewSimpleScope(nil, nil),
	}
	sContext := &ServiceContext{
		metadata:     md,
		activityHost: sHost,
		TaskName:     "Rule action service",
		inputs:       make(map[string]interface{}, len(md.Input)),
		outputs:      make(map[string]interface{}, len(md.Output)),
		settings:     make(map[string]interface{}, len(md.Settings)),
	}

	for name, tv := range md.Input {
		sContext.inputs[name] = tv.Value()
	}
	for name, tv := range md.Output {
		sContext.outputs[name] = tv.Value()
	}

	return sContext
}

// ActivityHost gets the "host" under with the activity is executing
func (sc *ServiceContext) ActivityHost() activity.Host {
	return sc.activityHost
}

//Name the name of the activity that is currently executing
func (sc *ServiceContext) Name() string {
	return sc.TaskName
}

// GetInput gets the value of the specified input attribute
func (sc *ServiceContext) GetInput(name string) interface{} {
	val, found := sc.inputs[name]
	if found {
		return val
	}
	return nil
}

// SetOutput sets the value of the specified output attribute
func (sc *ServiceContext) SetOutput(name string, value interface{}) error {
	sc.outputs[name] = value
	return nil
}

// GetInputObject gets all the activity input as the specified object.
func (sc *ServiceContext) GetInputObject(input data.StructValue) error {
	err := input.FromMap(sc.inputs)
	return err
}

// SetOutputObject sets the activity output as the specified object.
func (sc *ServiceContext) SetOutputObject(output data.StructValue) error {
	sc.outputs = output.ToMap()
	return nil
}

// GetSharedTempData get shared temporary data for activity, lifespan
// of the data dependent on the activity host implementation
func (sc *ServiceContext) GetSharedTempData() map[string]interface{} {
	if sc.shared == nil {
		sc.shared = make(map[string]interface{})
	}
	return sc.shared
}

// Logger the logger for the activity
func (sc *ServiceContext) Logger() log.Logger {
	return logger
}

// SetInput sets input
func (sc *ServiceContext) SetInput(name string, value interface{}) {
	sc.inputs[name] = value
}

// GetTracingContext returns tracing context
func (sc *ServiceContext) GetTracingContext() trace.TracingContext {
	return nil
}

// ServiceHost hosts service
type ServiceHost struct {
	HostId  string
	HostRef string

	IoMetadata *metadata.IOMetadata
	HostData   data.Scope
	ReplyData  map[string]interface{}
	ReplyErr   error
}

// ID returns the ID of the Activity Host
func (ac *ServiceHost) ID() string {
	return ac.HostId
}

// Name the name of the Activity Host
func (ac *ServiceHost) Name() string {
	return ""
}

// IOMetadata get the input/output metadata of the activity host
func (ac *ServiceHost) IOMetadata() *metadata.IOMetadata {
	return ac.IoMetadata
}

// Reply is used to reply to the activity Host with the results of the execution
func (ac *ServiceHost) Reply(replyData map[string]interface{}, err error) {
	ac.ReplyData = replyData
	ac.ReplyErr = err
}

// Return is used to indicate to the activity Host that it should complete and return the results of the execution
func (ac *ServiceHost) Return(returnData map[string]interface{}, err error) {
	ac.ReplyData = returnData
	ac.ReplyErr = err
}

// Scope returns the scope for the Host's data
func (ac *ServiceHost) Scope() data.Scope {
	return ac.HostData
}

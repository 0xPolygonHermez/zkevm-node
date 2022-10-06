package jsonrpc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unicode"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

const (
	requiredReturnParamsPerFn = 2
)

type serviceData struct {
	sv      reflect.Value
	funcMap map[string]*funcData
}

type funcData struct {
	inNum int
	reqt  []reflect.Type
	fv    reflect.Value
	isDyn bool
}

func (f *funcData) numParams() int {
	return f.inNum - 1
}

// Handler handles jsonrpc requests
type Handler struct {
	serviceMap map[string]*serviceData
}

func newJSONRpcHandler() *Handler {
	handler := &Handler{
		serviceMap: map[string]*serviceData{},
	}
	return handler
}

var connectionCounter = 0
var connectionCounterMutex sync.Mutex

// Handle is the function that knows which and how a function should
// be executed when a JSON RPC request is received
func (d *Handler) Handle(req Request) Response {
	connectionCounterMutex.Lock()
	connectionCounter++
	connectionCounterMutex.Unlock()
	defer func() {
		connectionCounterMutex.Lock()
		connectionCounter--
		connectionCounterMutex.Unlock()
		log.Debugf("Current open connections %d", connectionCounter)
	}()
	log.Debugf("Current open connections %d", connectionCounter)
	log.Debugf("request method %s id %v params %v", req.Method, req.ID, string(req.Params))

	service, fd, err := d.getFnHandler(req)
	if err != nil {
		return NewResponse(req, nil, err)
	}

	inArgs := make([]reflect.Value, fd.inNum)
	inArgs[0] = service.sv

	inputs := make([]interface{}, fd.numParams())
	for i := 0; i < fd.inNum-1; i++ {
		val := reflect.New(fd.reqt[i+1])
		inputs[i] = val.Interface()
		inArgs[i+1] = val.Elem()
	}

	if fd.numParams() > 0 {
		if err := json.Unmarshal(req.Params, &inputs); err != nil {
			return NewResponse(req, nil, newRPCError(invalidParamsErrorCode, "Invalid Params"))
		}
	}

	output := fd.fv.Call(inArgs)
	if err := getError(output[1]); err != nil {
		log.Errorf("failed to call method %s: [%v]%v. Params: %v", req.Method, err.ErrorCode(), err.Error(), string(req.Params))
		return NewResponse(req, nil, err)
	}

	var data *[]byte
	res := output[0].Interface()
	if res != nil {
		d, _ := json.Marshal(res)
		data = &d
	}

	return NewResponse(req, data, nil)
}

func (d *Handler) registerService(serviceName string, service interface{}) {
	st := reflect.TypeOf(service)
	if st.Kind() == reflect.Struct {
		panic(fmt.Sprintf("jsonrpc: service '%s' must be a pointer to struct", serviceName))
	}

	funcMap := make(map[string]*funcData)
	for i := 0; i < st.NumMethod(); i++ {
		mv := st.Method(i)
		if mv.PkgPath != "" {
			// skip unexported methods
			continue
		}

		name := lowerCaseFirst(mv.Name)
		funcName := serviceName + "_" + name
		fd := &funcData{
			fv: mv.Func,
		}
		var err error
		if fd.inNum, fd.reqt, err = validateFunc(funcName, fd.fv, true); err != nil {
			panic(fmt.Sprintf("jsonrpc: %s", err))
		}
		// check if last item is a pointer
		if fd.numParams() != 0 {
			last := fd.reqt[fd.numParams()]
			if last.Kind() == reflect.Ptr {
				fd.isDyn = true
			}
		}
		funcMap[name] = fd
	}

	d.serviceMap[serviceName] = &serviceData{
		sv:      reflect.ValueOf(service),
		funcMap: funcMap,
	}
}

func (d *Handler) getFnHandler(req Request) (*serviceData, *funcData, rpcError) {
	methodNotFoundErrorMessage := fmt.Sprintf("the method %s does not exist/is not available", req.Method)

	callName := strings.SplitN(req.Method, "_", 2) //nolint:gomnd
	if len(callName) != 2 {                        //nolint:gomnd
		return nil, nil, newRPCError(notFoundErrorCode, methodNotFoundErrorMessage)
	}

	serviceName, funcName := callName[0], callName[1]

	service, ok := d.serviceMap[serviceName]
	if !ok {
		log.Infof("Method %s not found", req.Method)
		return nil, nil, newRPCError(notFoundErrorCode, methodNotFoundErrorMessage)
	}
	fd, ok := service.funcMap[funcName]
	if !ok {
		return nil, nil, newRPCError(notFoundErrorCode, methodNotFoundErrorMessage)
	}
	return service, fd, nil
}

func validateFunc(funcName string, fv reflect.Value, isMethod bool) (inNum int, reqt []reflect.Type, err error) {
	if funcName == "" {
		err = fmt.Errorf("funcName cannot be empty")
		return
	}

	ft := fv.Type()
	if ft.Kind() != reflect.Func {
		err = fmt.Errorf("function '%s' must be a function instead of %s", funcName, ft)
		return
	}

	inNum = ft.NumIn()
	outNum := ft.NumOut()

	if outNum != requiredReturnParamsPerFn {
		err = fmt.Errorf("unexpected number of output arguments in the function '%s': %d. Expected 2", funcName, outNum)
		return
	}
	if !isRPCErrorType(ft.Out(1)) {
		err = fmt.Errorf("unexpected type for the second return value of the function '%s': '%s'. Expected '%s'", funcName, ft.Out(1), rpcErrType)
		return
	}

	reqt = make([]reflect.Type, inNum)
	for i := 0; i < inNum; i++ {
		reqt[i] = ft.In(i)
	}
	return
}

var rpcErrType = reflect.TypeOf((*rpcError)(nil)).Elem()

func isRPCErrorType(t reflect.Type) bool {
	return t.Implements(rpcErrType)
}

func getError(v reflect.Value) rpcError {
	if v.IsNil() {
		return nil
	}

	switch vt := v.Interface().(type) {
	case *RPCError:
		return vt
	default:
		return newRPCError(defaultErrorCode, "runtime error")
	}
}

func lowerCaseFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

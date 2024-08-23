package builder

/*
#cgo CFLAGS: -I/usr/include/webkitgtk-4.0
#cgo LDFLAGS: -ljavascriptcoregtk-4.0
#include <JavaScriptCore/JavaScript.h>

#include <stdlib.h>
#include <string.h>

typedef struct {
    char *name;
    int age;
} Builder;

extern JSObjectRef BuilderConstructor(JSContextRef ctx, JSObjectRef constructor, size_t argumentCount,  JSValueRef arguments[], JSValueRef* exception);
extern void BuilderDestructor(JSObjectRef object);
extern JSValueRef BuilderShow(JSContextRef ctx, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount,  JSValueRef arguments[], JSValueRef* exception);
extern JSValueRef BuilderModify(JSContextRef ctx, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount,  JSValueRef arguments[], JSValueRef* exception);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

//export BuilderConstructor
func BuilderConstructor(ctx C.JSContextRef, constructor C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSObjectRef {
	builder := (*C.Builder)(C.malloc(C.sizeof_Builder))

	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]

	if argumentCount > 0 {
		nameJS := C.JSValueToStringCopy(ctx, argumentSlice[0], exception)
		maxSize := C.JSStringGetMaximumUTF8CStringSize(nameJS)
		builder.name = (*C.char)(C.malloc(maxSize))
		C.JSStringGetUTF8CString(nameJS, builder.name, maxSize)
		C.JSStringRelease(nameJS)
	} else {
		builder.name = C.CString("Default Name")
	}

	if argumentCount > 1 {
		builder.age = C.int(C.JSValueToNumber(ctx, argumentSlice[1], exception))
	} else {
		builder.age = 30
	}

	instanceClassDefinition := C.kJSClassDefinitionEmpty
	instanceClass := C.JSClassCreate(&instanceClassDefinition)
	object := C.JSObjectMake(ctx, instanceClass, unsafe.Pointer(builder))
	C.JSClassRelease(instanceClass)

	return object
}

//export BuilderDestructor
func BuilderDestructor(object C.JSObjectRef) {
	builder := (*C.Builder)(C.JSObjectGetPrivate(object))
	if builder != nil {
		C.free(unsafe.Pointer(builder.name))
		C.free(unsafe.Pointer(builder))
	}
}

//export BuilderShow
func BuilderShow(ctx C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	fmt.Println("BuilderShow ha sido llamado.")

	builder := (*C.Builder)(C.JSObjectGetPrivate(thisObject))
	if builder != nil {
		fmt.Printf("Name: %s, Age: %d\n", C.GoString(builder.name), int(builder.age))
	} else {
		fmt.Println("Builder no encontrado.")
	}

	return C.JSValueMakeUndefined(ctx)
}

//export BuilderModify
func BuilderModify(ctx C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount C.size_t, arguments *C.JSValueRef, exception *C.JSValueRef) C.JSValueRef {
	builder := (*C.Builder)(C.JSObjectGetPrivate(thisObject))

	argumentSlice := (*[1 << 30]C.JSValueRef)(unsafe.Pointer(arguments))[:argumentCount:argumentCount]

	if builder != nil {
		if argumentCount > 0 {
			nameJS := C.JSValueToStringCopy(ctx, argumentSlice[0], exception)
			maxSize := C.JSStringGetMaximumUTF8CStringSize(nameJS)
			C.free(unsafe.Pointer(builder.name))
			builder.name = (*C.char)(C.malloc(maxSize))
			C.JSStringGetUTF8CString(nameJS, builder.name, maxSize)
			C.JSStringRelease(nameJS)
		}
		if argumentCount > 1 {
			builder.age = C.int(C.JSValueToNumber(ctx, argumentSlice[1], exception))
		}
		fmt.Printf("Propiedades modificadas: Name: %s, Age: %d\n", C.GoString(builder.name), int(builder.age))
	} else {
		fmt.Println("Builder no encontrado.")
	}

	return C.JSValueMakeUndefined(ctx)
}

func PubBuilderConstructor() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.BuilderConstructor)
}

func PubBuilderDestructor() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.BuilderDestructor)
}

func PubBuilderShow() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.BuilderShow)
}

func PubBuilderModify() C.JSObjectCallAsFunctionCallback {
	return C.JSObjectCallAsFunctionCallback(C.BuilderModify)
}

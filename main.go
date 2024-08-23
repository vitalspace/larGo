package main

/*
#cgo CFLAGS: -I/usr/include/webkitgtk-4.0
#cgo LDFLAGS: -ljavascriptcoregtk-4.0
#include <JavaScriptCore/JavaScript.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"largo/src/builder"
	"largo/src/console"
	"largo/src/math"
	"largo/src/require"
	"largo/src/utils"
	"os"
	"unsafe"
)

// createCustomFunction crea una función JavaScript personalizada y la establece como propiedad del objeto global.
func createCustomFunction(context C.JSGlobalContextRef, globalObject C.JSObjectRef, functionName string, functionCallback C.JSObjectCallAsFunctionCallback) {
	// Crear una cadena JavaScript a partir del nombre de la función en formato UTF-8.
	functionString := C.JSStringCreateWithUTF8CString(C.CString(functionName))

	// Crear un objeto de función JavaScript usando la cadena y la devolución de llamada de la función.
	functionObject := C.JSObjectMakeFunctionWithCallback(context, functionString, functionCallback)

	// Establecer la función recién creada como propiedad del objeto global.
	C.JSObjectSetProperty(context, globalObject, functionString, functionObject, C.kJSPropertyAttributeNone, nil)

	// Liberar la cadena de función creada con JSStringCreateWithUTF8CString para evitar fugas de memoria.
	C.JSStringRelease(functionString)
}

// createCustomClass crea una clase JavaScript personalizada con propiedades y métodos.
func createCustomClass(context C.JSGlobalContextRef, className string, constructor C.JSObjectCallAsConstructorCallback, finalize C.JSObjectFinalizeCallback, methods map[string]C.JSObjectCallAsFunctionCallback) {
	// Definir la clase en C.
	classDefinition := C.kJSClassDefinitionEmpty
	classDefinition.callAsConstructor = constructor
	classDefinition.finalize = finalize

	// // Configurar los getters y setters para las propiedades.
	// for propName, getter := range properties {
	// 	propNameCString := C.CString(propName)
	// 	classDefinition.getProperty = getter
	// 	// Puedes también agregar un setter si es necesario
	// 	// classDefinition.setProperty = setter
	// 	C.free(unsafe.Pointer(propNameCString))
	// }

	// Crear la clase en C.
	classRef := C.JSClassCreate(&classDefinition)

	// Crear el constructor de la clase.
	constructorObject := C.JSObjectMakeConstructor(context, classRef, constructor)

	// Configurar los métodos de la clase.
	prototype := C.JSValueToObject(context, C.JSObjectGetPrototype(context, constructorObject), nil)
	for methodName, callback := range methods {
		methodString := C.JSStringCreateWithUTF8CString(C.CString(methodName))
		methodObject := C.JSObjectMakeFunctionWithCallback(context, methodString, callback)
		C.JSObjectSetProperty(context, prototype, methodString, methodObject, C.kJSPropertyAttributeNone, nil)
		C.JSStringRelease(methodString)
	}

	// Registrar la clase en el objeto global de JavaScript.
	classString := C.JSStringCreateWithUTF8CString(C.CString(className))
	C.JSObjectSetProperty(context, C.JSContextGetGlobalObject(context), classString, constructorObject, C.kJSPropertyAttributeNone, nil)
	C.JSStringRelease(classString)

	// Liberar la referencia a la clase.
	C.JSClassRelease(classRef)
}

// Apis define las API disponibles en JavaScript.
func Apis(context C.JSGlobalContextRef, globalObject C.JSObjectRef) {

	createCustomFunction(context, globalObject, "Add", C.JSObjectCallAsFunctionCallback(math.Add()))
	createCustomFunction(context, globalObject, "Mult", C.JSObjectCallAsFunctionCallback(math.Mult()))
	createCustomFunction(context, globalObject, "require", C.JSObjectCallAsFunctionCallback(require.Require()))

	createCustomFunction(context, globalObject, "print", C.JSObjectCallAsFunctionCallback(console.Log()))
	console_str := C.CString("console")
	console_js := C.JSStringCreateWithUTF8CString(console_str)
	C.free(unsafe.Pointer(console_str))
	consoleGlobalObject := C.JSObjectMake(context, nil, nil)
	C.JSObjectSetProperty(context, globalObject, console_js, consoleGlobalObject, C.kJSPropertyAttributeNone, nil)

	createCustomFunction(context, consoleGlobalObject, "log", C.JSObjectCallAsFunctionCallback(console.Log()))
	createCustomFunction(context, consoleGlobalObject, "time", C.JSObjectCallAsFunctionCallback(console.Time()))
	createCustomFunction(context, consoleGlobalObject, "timeEnd", C.JSObjectCallAsFunctionCallback(console.TimeEnd()))
	C.JSStringRelease(console_js)

	methods := map[string]C.JSObjectCallAsFunctionCallback{
		"show":   C.JSObjectCallAsFunctionCallback(builder.PubBuilderShow()),
		"modify": C.JSObjectCallAsFunctionCallback(builder.PubBuilderModify()),
	}

	// // // Definir propiedades para la clase Builder.
	// properties := map[string]C.JSObjectGetPropertyCallback{
	// 	"name": C.JSObjectGetPropertyCallback(builder.BuilderGetName),
	// 	"age":  C.JSObjectGetPropertyCallback(builder.BuilderGetAge),
	// 	// Puedes definir más getters si es necesario
	// }

	createCustomClass(context, "Builder", C.JSObjectCallAsConstructorCallback(builder.PubBuilderConstructor()), C.JSObjectFinalizeCallback(builder.PubBuilderDestructor()), methods)

}

func main() {
	// Crear un contexto JavaScript global.
	context := C.JSGlobalContextCreate(nil)
	globalObject := C.JSContextGetGlobalObject(context)

	// builder.InitializeBuilderClass(C.JSContextRef(context))
	// classManager(C.JSContextRef(context))

	// Configurar las API en el objeto global.
	Apis(context, globalObject)

	// Verificar si hay argumentos de línea de comandos y si se proporciona el comando "run".
	if len(os.Args) > 2 && os.Args[1] == "run" {
		jsFileName := os.Args[2]

		// Leer el contenido del archivo JavaScript.
		fileContent := utils.ReadFile(jsFileName)

		// Crear una cadena JavaScript a partir del contenido del archivo.
		jsCode := C.JSStringCreateWithUTF8CString(C.CString(fileContent))
		defer C.JSStringRelease(jsCode)

		// Evaluar el script JavaScript.
		result := C.JSEvaluateScript(context, jsCode, globalObject, nil, 1, nil)

		// Convertir el resultado a una cadena de Go.
		resultStringJS := C.JSValueToStringCopy(context, result, nil)
		defer C.JSStringRelease(resultStringJS)

		// Obtener el tamaño máximo necesario para la cadena UTF-8.
		bufferSize := C.JSStringGetMaximumUTF8CStringSize(resultStringJS)
		resultCString := make([]C.char, bufferSize)
		C.JSStringGetUTF8CString(resultStringJS, &resultCString[0], bufferSize)

		// Imprimir el resultado.
		fmt.Printf("%s\n", C.GoString(&resultCString[0]))
	}

	// Liberar el contexto JavaScript global.
	C.JSGlobalContextRelease(context)
}

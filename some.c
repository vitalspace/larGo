// #define _GNU_SOURCE
// #include <JavaScriptCore/JavaScript.h>
// #include <stdio.h>
// #include <stdlib.h>
// #include <string.h>

// typedef struct {
//     char *name;
//     int age;
// } Builder;

// JSObjectRef BuilderConstructor(JSContextRef ctx, JSObjectRef constructor, size_t argumentCount, const JSValueRef arguments[], JSValueRef* exception) {
//     Builder *builder = (Builder*)malloc(sizeof(Builder));
    
//     if (argumentCount > 0) {
//         JSStringRef nameJS = JSValueToStringCopy(ctx, arguments[0], exception);
//         size_t maxSize = JSStringGetMaximumUTF8CStringSize(nameJS);
//         builder->name = (char*)malloc(maxSize);
//         JSStringGetUTF8CString(nameJS, builder->name, maxSize);
//         JSStringRelease(nameJS);
//     } else {
//         builder->name = strdup("Default Name");
//     }
//     builder->age = (argumentCount > 1) ? (int)JSValueToNumber(ctx, arguments[1], exception) : 30;
    
//     JSClassDefinition instanceClassDefinition = kJSClassDefinitionEmpty;
//     JSClassRef instanceClass = JSClassCreate(&instanceClassDefinition);
//     JSObjectRef object = JSObjectMake(ctx, instanceClass, builder);
//     JSClassRelease(instanceClass);
    
//     return object;
// }

// void BuilderDestructor(JSObjectRef object) {
//     Builder *builder = (Builder*)JSObjectGetPrivate(object);
//     if (builder) {
//         free(builder->name);
//         free(builder);
//     }
// }

// JSValueRef BuilderShow(JSContextRef ctx, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, const JSValueRef arguments[], JSValueRef* exception) {
//     printf("BuilderShow ha sido llamado.\n");
    
//     Builder *builder = (Builder*)JSObjectGetPrivate(thisObject);
//     if (builder) {
//         printf("Name: %s, Age: %d\n", builder->name, builder->age);
//     } else {
//         printf("Builder no encontrado.\n");
//     }
    
//     return JSValueMakeUndefined(ctx);
// }

// JSValueRef BuilderModify(JSContextRef ctx, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, const JSValueRef arguments[], JSValueRef* exception) {
//     Builder *builder = (Builder*)JSObjectGetPrivate(thisObject);
//     if (builder) {
//         if (argumentCount > 0) {
//             JSStringRef nameJS = JSValueToStringCopy(ctx, arguments[0], exception);
//             size_t maxSize = JSStringGetMaximumUTF8CStringSize(nameJS);
//             free(builder->name);  // Liberar la memoria del nombre anterior
//             builder->name = (char*)malloc(maxSize);
//             JSStringGetUTF8CString(nameJS, builder->name, maxSize);
//             JSStringRelease(nameJS);
//         }
//         if (argumentCount > 1) {
//             builder->age = (int)JSValueToNumber(ctx, arguments[1], exception);
//         }
//         printf("Propiedades modificadas: Name: %s, Age: %d\n", builder->name, builder->age);
//     } else {
//         printf("Builder no encontrado.\n");
//     }
    
//     return JSValueMakeUndefined(ctx);
// }

// void InitializeBuilderClass(JSContextRef ctx) {
//     JSClassDefinition classDefinition = kJSClassDefinitionEmpty;
//     classDefinition.callAsConstructor = BuilderConstructor;
//     classDefinition.finalize = BuilderDestructor;
    
//     JSClassRef builderClass = JSClassCreate(&classDefinition);
    
//     JSObjectRef constructor = JSObjectMakeConstructor(ctx, builderClass, BuilderConstructor);
    
//     JSObjectRef prototype = JSValueToObject(ctx, JSObjectGetPrototype(ctx, constructor), NULL);
    
//     JSStringRef showFuncName = JSStringCreateWithUTF8CString("show");
//     JSObjectRef showFunc = JSObjectMakeFunctionWithCallback(ctx, showFuncName, BuilderShow);
//     JSObjectSetProperty(ctx, prototype, showFuncName, showFunc, kJSPropertyAttributeNone, NULL);
//     JSStringRelease(showFuncName);
    
//     JSStringRef modifyFuncName = JSStringCreateWithUTF8CString("modify");
//     JSObjectRef modifyFunc = JSObjectMakeFunctionWithCallback(ctx, modifyFuncName, BuilderModify);
//     JSObjectSetProperty(ctx, prototype, modifyFuncName, modifyFunc, kJSPropertyAttributeNone, NULL);
//     JSStringRelease(modifyFuncName);
    
//     JSStringRef builderName = JSStringCreateWithUTF8CString("Builder");
//     JSObjectSetProperty(ctx, JSContextGetGlobalObject(ctx), builderName, constructor, kJSPropertyAttributeNone, NULL);
//     JSStringRelease(builderName);
    
//     JSClassRelease(builderClass);
// }

// int main(int argc, char* argv[]) {
//     JSGlobalContextRef ctx = JSGlobalContextCreate(NULL);
    
//     InitializeBuilderClass(ctx);
    
//     printf("Contexto y clase inicializados.\n");
    
//     JSStringRef script = JSStringCreateWithUTF8CString(
//         "var b = new Builder('Alice', 25); "
//         "b.show(); "
//         "b.modify('Bob', 30); "
//         "b.show();"
//     );
    
//     JSValueRef exception = NULL;
//     JSValueRef result = JSEvaluateScript(ctx, script, NULL, NULL, 1, &exception);
//     JSStringRelease(script);
    
//     if (exception) {
//         JSStringRef exceptionStr = JSValueToStringCopy(ctx, exception, NULL);
//         size_t exceptionSize = JSStringGetMaximumUTF8CStringSize(exceptionStr);
//         char* exceptionCStr = (char*)malloc(exceptionSize);
//         JSStringGetUTF8CString(exceptionStr, exceptionCStr, exceptionSize);
//         printf("JavaScript exception: %s\n", exceptionCStr);
//         free(exceptionCStr);
//         JSStringRelease(exceptionStr);
//     } else {
//         printf("Script ejecutado correctamente.\n");
//     }
    
//     JSGlobalContextRelease(ctx);
    
//     return 0;
// }
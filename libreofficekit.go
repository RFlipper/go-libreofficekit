package libreofficekit

// #cgo CFLAGS: -I /usr/include/LibreOfficeKit/ -D LOK_USE_UNSTABLE_API
// #cgo LDFLAGS: -ldl
// #include <stdlib.h>
// #include "LibreOfficeKit/LibreOfficeKitInit.h"
// #include "LibreOfficeKit/LibreOfficeKit.h"
/*
typedef void (*voidFunc) ();
typedef int (*intFunc) ();
typedef char* (*charFunc) ();
void destroy_bridge(voidFunc f, void* handle) {
	return f(handle);
};
LibreOfficeKitDocument* document_load_bridge(voidFunc f,
											 LibreOfficeKit* pThis,
											 const char* pURL) {
	f(pThis, pURL);
};
char* get_error_bridge(charFunc f, LibreOfficeKit* pThis) {
	return f(pThis);
};
int get_document_bridge(intFunc f, LibreOfficeKitDocument* pThis) {
	return f(pThis);
};
void get_document_size_bridge(voidFunc f,
							LibreOfficeKitDocument* pThis,
							long* pWidth,
							long* pHeight) {
	return f(pThis, pWidth, pHeight);
};
void set_document_part_bridge(voidFunc f,
							LibreOfficeKitDocument* pThis,
							int nPart) {
	return f(pThis, nPart);
};

char* get_document_part_name_bridge(charFunc f,
							LibreOfficeKitDocument* pThis,
							int nPart) {
	return f(pThis, nPart);
};

void initialize_for_rendering_bridge(voidFunc f,
									LibreOfficeKitDocument* pThis,
									const char* pArguments) {
	return f(pThis, pArguments);
};

int document_save_bridge(intFunc f,
						LibreOfficeKitDocument* pThis,
						const char* pUrl,
						const char* pFormat,
						const char* pFilterOptions) {
	return f(pThis, pUrl, pFormat, pFilterOptions);
};
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type Office struct {
	handle *C.struct__LibreOfficeKit
}

func NewOffice(path string) (*Office, error) {
	office := new(Office)

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	lokit := C.lok_init(c_path)
	if lokit == nil {
		return nil, fmt.Errorf("Failed to initialize LibreOfficeKit with path: '%s'", path)
	}

	office.handle = lokit

	return office, nil

}

func (self *Office) Close() {
	selfPointer := unsafe.Pointer(self.handle)
	C.destroy_bridge(self.handle.pClass.destroy, selfPointer)
}

func (self *Office) GetError() string {
	return C.GoString(C.get_error_bridge(self.handle.pClass.getError, self.handle))
}

func (self *Office) LoadDocument(path string) (*Document, error) {
	document := new(Document)
	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))
	handle := C.document_load_bridge(self.handle.pClass.documentLoad, self.handle, c_path)
	if handle == nil {
		return nil, fmt.Errorf("Failed to load document")
	}
	document.handle = handle
	return document, nil
}

const (
	TextDocument = iota
	SpreadsheetDocument
	PresentationDocument
	DrawingDocument
	OtherDocument
)

type Document struct {
	handle *C.struct__LibreOfficeKitDocument
}

func (self *Document) Close() {
	selfPointer := unsafe.Pointer(self.handle)
	C.destroy_bridge(self.handle.pClass.destroy, selfPointer)
}

func (self *Document) GetType() int {
	return int(C.get_document_bridge(self.handle.pClass.getDocumentType, self.handle))
}

func (self *Document) GetParts() int {
	return int(C.get_document_bridge(self.handle.pClass.getParts, self.handle))
}

func (self *Document) GetPart() int {
	return int(C.get_document_bridge(self.handle.pClass.getPart, self.handle))
}

func (self *Document) SetPart(part int) {
	c_part := C.int(part)
	C.set_document_part_bridge(self.handle.pClass.setPart, self.handle, c_part)
}

func (self *Document) GetPartName(part int) string {
	c_part := C.int(part)
	c_part_name := C.get_document_part_name_bridge(self.handle.pClass.getPartName, self.handle, c_part)
	defer C.free(unsafe.Pointer(c_part_name))
	return C.GoString(c_part_name)
}

func (self *Document) GetSize() (int, int) {
	width := C.long(0)
	heigth := C.long(0)
	C.get_document_size_bridge(self.handle.pClass.getDocumentSize, self.handle, &width, &heigth)
	return int(width), int(heigth)
}

func (self *Document) InitializeForRendering(arguments string) {
	c_arguments := C.CString(arguments)
	defer C.free(unsafe.Pointer(c_arguments))
	C.initialize_for_rendering_bridge(self.handle.pClass.initializeForRendering, self.handle, c_arguments)
}

func (self *Document) SaveAs(path string, format string, filter string) error {
	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))
	c_format := C.CString(format)
	defer C.free(unsafe.Pointer(c_format))
	c_filter := C.CString(filter)
	defer C.free(unsafe.Pointer(c_filter))
	status := C.document_save_bridge(self.handle.pClass.saveAs, self.handle, c_path, c_format, c_filter)
	if status != 0 {
		return fmt.Errorf("Failed to save document")
	} else {
		return nil
	}
}
